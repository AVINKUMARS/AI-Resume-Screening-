# Backend — AI Resume Screening API

Go (Gin) REST API with GORM/PostgreSQL persistence and Google Gemini for resume
parsing and candidate–job matching.

## Layout

```
backend/
├── main.go                     # wiring + startup
└── internal/
    ├── config/     config.go   # env-driven configuration
    ├── database/   database.go # GORM connect + auto-migrate
    ├── models/     models.go   # Candidate, Job, Match (JSONB skill lists)
    ├── gemini/     gemini.go   # Gemini generateContent client
    ├── handlers/   handlers.go # HTTP handlers
    └── router/     router.go   # routes + CORS
```

## Configuration (`.env`)

| Var              | Default                              | Notes                          |
| ---------------- | ------------------------------------ | ------------------------------ |
| `PORT`           | `8080`                               | HTTP listen port               |
| `DATABASE_URL`   | `postgres://…/resume_screening`      | PostgreSQL DSN                 |
| `GEMINI_API_KEY` | —                                    | **Required** for AI features   |
| `GEMINI_MODEL`   | `gemini-2.5-flash`                   | Any generateContent model      |
| `ALLOWED_ORIGIN` | `http://localhost:3000`              | CORS origin for the frontend   |

## Run

```bash
go mod tidy
go run .
```

The schema (`candidates`, `jobs`, `matches`) is auto-migrated on startup. Skill,
experience, education, and matched/missing-skill lists are stored as JSONB.

## How the AI works

- **Parsing** (`ParseResume`) — the raw resume text is sent to Gemini with a strict
  JSON schema instruction; the response is unmarshalled into a `Candidate`.
- **Matching** (`MatchCandidate`) — a candidate profile and a job description are sent
  together; Gemini returns a 0–100 score, reasoning, and matched/missing skills.

Both calls request `application/json` output and defensively strip any markdown code
fences before decoding.

## Test the API

```bash
# health
curl localhost:8080/health

# upload a resume as text
curl -X POST localhost:8080/api/resumes/upload \
  -H 'Content-Type: application/json' \
  -d '{"text":"Jane Doe — Senior Go Engineer. 6 years building REST APIs with Gin, PostgreSQL, GORM..."}'

# create a job
curl -X POST localhost:8080/api/jobs \
  -H 'Content-Type: application/json' \
  -d '{"title":"Backend Engineer","required_skills":["Go","PostgreSQL"],"min_years_experience":4}'

# rank candidates for a job
curl -X POST localhost:8080/api/jobs/<JOB_ID>/match
```
