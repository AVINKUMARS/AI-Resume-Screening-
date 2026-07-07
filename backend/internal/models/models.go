package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StringSlice is a []string that persists to a JSONB column. Skills, education
// lines, etc. are variable-length lists we want to store and query as JSON.
type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}

func (s *StringSlice) Scan(src interface{}) error {
	if src == nil {
		*s = StringSlice{}
		return nil
	}
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		return errors.New("models: unsupported type for StringSlice scan")
	}
}

// Candidate is a parsed resume plus the AI-generated summary.
type Candidate struct {
	ID         uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	Name       string      `gorm:"index" json:"name"`
	Email      string      `gorm:"index" json:"email"`
	Phone      string      `json:"phone"`
	Location   string      `json:"location"`
	Title      string      `json:"title"`
	Summary    string      `gorm:"type:text" json:"summary"`
	Skills     StringSlice `gorm:"type:jsonb" json:"skills"`
	Experience StringSlice `gorm:"type:jsonb" json:"experience"`
	Education  StringSlice `gorm:"type:jsonb" json:"education"`
	YearsExp   float64     `json:"years_experience"`
	RawText    string      `gorm:"type:text" json:"raw_text,omitempty"`
	SourceFile string      `json:"source_file"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`

	Matches []Match `gorm:"foreignKey:CandidateID" json:"matches,omitempty"`
}

func (c *Candidate) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// Job is a role a recruiter wants to match candidates against.
type Job struct {
	ID             uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	Title          string      `gorm:"index" json:"title"`
	Company        string      `json:"company"`
	Location       string      `json:"location"`
	Description    string      `gorm:"type:text" json:"description"`
	RequiredSkills StringSlice `gorm:"type:jsonb" json:"required_skills"`
	MinYearsExp    float64     `json:"min_years_experience"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`

	Matches []Match `gorm:"foreignKey:JobID" json:"matches,omitempty"`
}

func (j *Job) BeforeCreate(tx *gorm.DB) error {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return nil
}

// Match is an AI-scored fit between a candidate and a job.
type Match struct {
	ID            uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	CandidateID   uuid.UUID   `gorm:"type:uuid;index" json:"candidate_id"`
	JobID         uuid.UUID   `gorm:"type:uuid;index" json:"job_id"`
	Score         int         `json:"score"` // 0-100
	Reasoning     string      `gorm:"type:text" json:"reasoning"`
	MatchedSkills StringSlice `gorm:"type:jsonb" json:"matched_skills"`
	MissingSkills StringSlice `gorm:"type:jsonb" json:"missing_skills"`
	CreatedAt     time.Time   `json:"created_at"`

	Candidate *Candidate `gorm:"foreignKey:CandidateID" json:"candidate,omitempty"`
	Job       *Job       `gorm:"foreignKey:JobID" json:"job,omitempty"`
}

func (m *Match) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
