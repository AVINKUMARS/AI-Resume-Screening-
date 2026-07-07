package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ledongthuc/pdf"
	"github.com/scalingwolf/ai-resume-screening/backend/internal/gemini"
	"github.com/scalingwolf/ai-resume-screening/backend/internal/models"
	"gorm.io/gorm"
)

// Handler wires HTTP handlers to the database and the Gemini client.
type Handler struct {
	db *gorm.DB
	ai *gemini.Client
}

func New(db *gorm.DB, ai *gemini.Client) *Handler {
	return &Handler{db: db, ai: ai}
}

// Health is a simple liveness probe.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// UploadResume accepts either a multipart file ("resume": .pdf/.txt) or a JSON
// body {"text": "..."}, extracts text, parses it with Gemini, and stores the
// resulting candidate.
func (h *Handler) UploadResume(c *gin.Context) {
	var (
		rawText    string
		sourceFile string
	)

	// Prefer a multipart upload; fall back to a JSON text body.
	if fileHeader, err := c.FormFile("resume"); err == nil {
		sourceFile = fileHeader.Filename
		f, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot open uploaded file"})
			return
		}
		defer f.Close()

		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(f); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot read uploaded file"})
			return
		}

		if strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".pdf") {
			text, err := extractPDFText(buf.Bytes())
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": fmt.Sprintf("failed to read PDF: %v", err)})
				return
			}
			rawText = text
		} else {
			rawText = buf.String()
		}
	} else {
		var body struct {
			Text       string `json:"text"`
			SourceFile string `json:"source_file"`
		}
		if err := c.ShouldBindJSON(&body); err != nil || strings.TrimSpace(body.Text) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "provide a 'resume' file or a non-empty 'text' field"})
			return
		}
		rawText = body.Text
		sourceFile = body.SourceFile
	}

	if strings.TrimSpace(rawText) == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "no text could be extracted from the resume"})
		return
	}

	parsed, err := h.ai.ParseResume(c.Request.Context(), rawText)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("AI parsing failed: %v", err)})
		return
	}

	candidate := models.Candidate{
		Name:       parsed.Name,
		Email:      parsed.Email,
		Phone:      parsed.Phone,
		Location:   parsed.Location,
		Title:      parsed.Title,
		Summary:    parsed.Summary,
		Skills:     parsed.Skills,
		Experience: parsed.Experience,
		Education:  parsed.Education,
		YearsExp:   parsed.YearsExp,
		RawText:    rawText,
		SourceFile: sourceFile,
	}
	if err := h.db.Create(&candidate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save candidate"})
		return
	}

	c.JSON(http.StatusCreated, candidate)
}

// ListCandidates returns candidates, optionally filtered by a free-text `q`
// query that matches name, title, or skills.
func (h *Handler) ListCandidates(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	query := h.db.Order("created_at desc")

	if q != "" {
		like := "%" + strings.ToLower(q) + "%"
		query = query.Where(
			"LOWER(name) LIKE ? OR LOWER(title) LIKE ? OR LOWER(skills::text) LIKE ? OR LOWER(summary) LIKE ?",
			like, like, like, like,
		)
	}

	var candidates []models.Candidate
	if err := query.Omit("raw_text").Find(&candidates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list candidates"})
		return
	}
	c.JSON(http.StatusOK, candidates)
}

// GetCandidate returns a single candidate with their match history.
func (h *Handler) GetCandidate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid candidate id"})
		return
	}

	var candidate models.Candidate
	if err := h.db.Preload("Matches", func(db *gorm.DB) *gorm.DB {
		return db.Order("score desc")
	}).Preload("Matches.Job").First(&candidate, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "candidate not found"})
		return
	}
	c.JSON(http.StatusOK, candidate)
}

// DeleteCandidate removes a candidate and their matches.
func (h *Handler) DeleteCandidate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid candidate id"})
		return
	}
	h.db.Where("candidate_id = ?", id).Delete(&models.Match{})
	if err := h.db.Delete(&models.Candidate{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete candidate"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": id})
}

// --- Jobs ----------------------------------------------------------------

func (h *Handler) CreateJob(c *gin.Context) {
	var job models.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job payload"})
		return
	}
	if strings.TrimSpace(job.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "job title is required"})
		return
	}
	if err := h.db.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create job"})
		return
	}
	c.JSON(http.StatusCreated, job)
}

func (h *Handler) ListJobs(c *gin.Context) {
	var jobs []models.Job
	if err := h.db.Order("created_at desc").Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list jobs"})
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func (h *Handler) GetJob(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}
	var job models.Job
	if err := h.db.First(&job, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}
	c.JSON(http.StatusOK, job)
}

func (h *Handler) DeleteJob(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}
	h.db.Where("job_id = ?", id).Delete(&models.Match{})
	if err := h.db.Delete(&models.Job{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete job"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": id})
}

// --- Matching ------------------------------------------------------------

// MatchJob scores every candidate against a job, upserts the results, and
// returns them ranked best-first.
func (h *Handler) MatchJob(c *gin.Context) {
	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}

	var job models.Job
	if err := h.db.First(&job, "id = ?", jobID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	var candidates []models.Candidate
	if err := h.db.Find(&candidates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load candidates"})
		return
	}

	jobText := describeJob(&job)
	results := make([]models.Match, 0, len(candidates))
	var failures []string

	for i := range candidates {
		cand := &candidates[i]
		res, err := h.ai.MatchCandidate(c.Request.Context(), describeCandidate(cand), jobText)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", cand.Name, err))
			continue
		}

		match := models.Match{
			CandidateID:   cand.ID,
			JobID:         job.ID,
			Score:         res.Score,
			Reasoning:     res.Reasoning,
			MatchedSkills: res.MatchedSkills,
			MissingSkills: res.MissingSkills,
		}

		// One match row per (candidate, job): replace any prior scoring.
		h.db.Where("candidate_id = ? AND job_id = ?", cand.ID, job.ID).Delete(&models.Match{})
		if err := h.db.Create(&match).Error; err != nil {
			failures = append(failures, fmt.Sprintf("%s: save failed", cand.Name))
			continue
		}
		match.Candidate = cand
		results = append(results, match)
	}

	// Rank best-first.
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Score > results[i].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"job":      job,
		"matches":  results,
		"failures": failures,
	})
}

// --- helpers -------------------------------------------------------------

func describeCandidate(c *models.Candidate) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Name: %s\nTitle: %s\nLocation: %s\nYears of experience: %.1f\n",
		c.Name, c.Title, c.Location, c.YearsExp)
	fmt.Fprintf(&b, "Summary: %s\n", c.Summary)
	fmt.Fprintf(&b, "Skills: %s\n", strings.Join(c.Skills, ", "))
	fmt.Fprintf(&b, "Experience:\n- %s\n", strings.Join(c.Experience, "\n- "))
	fmt.Fprintf(&b, "Education:\n- %s\n", strings.Join(c.Education, "\n- "))
	return b.String()
}

func describeJob(j *models.Job) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Title: %s\nCompany: %s\nLocation: %s\nMinimum years of experience: %.1f\n",
		j.Title, j.Company, j.Location, j.MinYearsExp)
	fmt.Fprintf(&b, "Required skills: %s\n", strings.Join(j.RequiredSkills, ", "))
	fmt.Fprintf(&b, "Description:\n%s\n", j.Description)
	return b.String()
}

// extractPDFText pulls plain text out of a PDF byte slice.
func extractPDFText(data []byte) (string, error) {
	r, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	textReader, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	if _, err := buf.ReadFrom(textReader); err != nil {
		return "", err
	}
	return buf.String(), nil
}
