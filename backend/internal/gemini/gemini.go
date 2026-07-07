package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const apiBase = "https://generativelanguage.googleapis.com/v1beta/models"

// Client is a thin wrapper over the Gemini generateContent REST endpoint.
type Client struct {
	apiKey string
	model  string
	http   *http.Client
}

func NewClient(apiKey, model string) *Client {
	return &Client{
		apiKey: apiKey,
		model:  model,
		http:   &http.Client{Timeout: 60 * time.Second},
	}
}

// --- wire types for the Gemini REST API ---------------------------------

type geminiRequest struct {
	Contents         []content         `json:"contents"`
	GenerationConfig *generationConfig `json:"generationConfig,omitempty"`
}

type content struct {
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text"`
}

type generationConfig struct {
	Temperature      float64 `json:"temperature"`
	ResponseMIMEType string  `json:"responseMimeType,omitempty"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []part `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// --- domain result types -------------------------------------------------

// ParsedResume is the structured extraction of a raw resume.
type ParsedResume struct {
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Location   string   `json:"location"`
	Title      string   `json:"title"`
	Summary    string   `json:"summary"`
	Skills     []string `json:"skills"`
	Experience []string `json:"experience"`
	Education  []string `json:"education"`
	YearsExp   float64  `json:"years_experience"`
}

// MatchResult is the AI's assessment of a candidate against a job.
type MatchResult struct {
	Score         int      `json:"score"`
	Reasoning     string   `json:"reasoning"`
	MatchedSkills []string `json:"matched_skills"`
	MissingSkills []string `json:"missing_skills"`
}

// generate performs a single generateContent call and returns the raw text of
// the first candidate part.
func (c *Client) generate(ctx context.Context, prompt string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("gemini: API key not configured")
	}

	reqBody := geminiRequest{
		Contents: []content{{Parts: []part{{Text: prompt}}}},
		GenerationConfig: &generationConfig{
			Temperature:      0.2,
			ResponseMIMEType: "application/json",
		},
	}

	buf, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("gemini: marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/%s:generateContent?key=%s", apiBase, c.model, c.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		return "", fmt.Errorf("gemini: new request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("gemini: http: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var gr geminiResponse
	if err := json.Unmarshal(body, &gr); err != nil {
		return "", fmt.Errorf("gemini: decode response (status %d): %w", resp.StatusCode, err)
	}
	if gr.Error != nil {
		return "", fmt.Errorf("gemini: api error %d: %s", gr.Error.Code, gr.Error.Message)
	}
	if len(gr.Candidates) == 0 || len(gr.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini: empty response (status %d)", resp.StatusCode)
	}

	return gr.Candidates[0].Content.Parts[0].Text, nil
}

// jsonFence strips ```json ... ``` fences some models still emit even when JSON
// output is requested, so the payload unmarshals cleanly.
var jsonFence = regexp.MustCompile("(?s)```(?:json)?\\s*(.*?)\\s*```")

func cleanJSON(s string) string {
	s = strings.TrimSpace(s)
	if m := jsonFence.FindStringSubmatch(s); m != nil {
		return strings.TrimSpace(m[1])
	}
	return s
}

// ParseResume extracts structured fields and an AI summary from raw resume text.
func (c *Client) ParseResume(ctx context.Context, rawText string) (*ParsedResume, error) {
	prompt := fmt.Sprintf(`You are an expert technical recruiter and resume parser.
Extract the candidate's details from the resume below and return ONLY a JSON object
with this exact shape:
{
  "name": string,
  "email": string,
  "phone": string,
  "location": string,
  "title": string,              // current or most recent professional title
  "summary": string,            // 2-3 sentence recruiter-facing summary of the candidate
  "skills": string[],           // technical and professional skills
  "experience": string[],       // one entry per role: "Title at Company (dates) — key achievement"
  "education": string[],        // one entry per qualification
  "years_experience": number    // total years of professional experience, estimated
}
If a field is unknown, use an empty string, empty array, or 0. Do not invent data.

RESUME:
"""
%s
"""`, truncate(rawText, 20000))

	raw, err := c.generate(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var parsed ParsedResume
	if err := json.Unmarshal([]byte(cleanJSON(raw)), &parsed); err != nil {
		return nil, fmt.Errorf("gemini: parse resume JSON: %w (raw: %.200s)", err, raw)
	}
	return &parsed, nil
}

// MatchCandidate scores how well a candidate fits a job.
func (c *Client) MatchCandidate(ctx context.Context, candidate, job string) (*MatchResult, error) {
	prompt := fmt.Sprintf(`You are an expert technical recruiter scoring a candidate against a job.
Compare the CANDIDATE PROFILE to the JOB and return ONLY a JSON object with this shape:
{
  "score": number,             // 0-100 overall fit, weighing skills, experience, and seniority
  "reasoning": string,         // 2-4 sentences explaining the score for a hiring manager
  "matched_skills": string[],  // required/valued skills the candidate clearly has
  "missing_skills": string[]   // required/valued skills the candidate appears to lack
}
Be objective and evidence-based. Do not inflate scores.

JOB:
"""
%s
"""

CANDIDATE PROFILE:
"""
%s
"""`, truncate(job, 8000), truncate(candidate, 12000))

	raw, err := c.generate(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var result MatchResult
	if err := json.Unmarshal([]byte(cleanJSON(raw)), &result); err != nil {
		return nil, fmt.Errorf("gemini: match JSON: %w (raw: %.200s)", err, raw)
	}
	if result.Score < 0 {
		result.Score = 0
	}
	if result.Score > 100 {
		result.Score = 100
	}
	return &result, nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}
