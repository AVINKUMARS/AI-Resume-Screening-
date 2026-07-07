# AI Resume Screening & Candidate Matching System

An AI-powered recruitment platform for resume parsing, candidate management, and
job-fit scoring. Resumes are parsed with the **Google Gemini API** into structured
profiles, and candidates are ranked against job postings with explainable match scores.

**Stack:** Next.js · TypeScript · Go (Gin) · PostgreSQL · GORM · Google Gemini API

```
ai-resume-screening/
├── backend/     # Go + Gin REST API, GORM/PostgreSQL, Gemini integration
└── frontend/    # Next.js (App Router) + TypeScript dashboard
```

## Features

- **AI resume parsing** — upload a PDF/TXT resume (or paste text); Gemini extracts
  name, contact, title, skills, experience, education, years of experience, and a
  recruiter-facing summary.
- **Candidate management** — searchable talent pool (by name, title, skill, or summary)
  with full candidate profiles.
- **Job postings** — create roles with required skills, minimum experience, and a
  description.
- **AI candidate matching** — score every candidate against a job (0–100) with a written
  rationale plus matched/missing skills, ranked best-first.
- **Secure REST API** — clean, versioned endpoints with CORS scoped to the frontend origin.

## Quick start

### 1. Backend

```bash
cd backend
cp .env.example .env          # set DATABASE_URL and GEMINI_API_KEY
go mod tidy
go run .                      # serves on :8080
```

Requires a running PostgreSQL database — GORM auto-migrates the schema on startup.

### 2. Frontend

```bash
cd frontend
cp .env.local.example .env.local   # NEXT_PUBLIC_API_URL=http://localhost:8080
npm install
npm run dev                        # serves on :3000
```

Open http://localhost:3000.

## API reference

| Method | Path                    | Description                                   |
| ------ | ----------------------- | --------------------------------------------- |
| GET    | `/health`               | Liveness probe                                |
| POST   | `/api/resumes/upload`   | Upload resume (multipart `resume` or JSON `text`); parses + stores candidate |
| GET    | `/api/candidates?q=`    | List/search candidates                        |
| GET    | `/api/candidates/:id`   | Candidate detail with match history           |
| DELETE | `/api/candidates/:id`   | Delete candidate                              |
| POST   | `/api/jobs`             | Create a job posting                          |
| GET    | `/api/jobs`             | List jobs                                     |
| GET    | `/api/jobs/:id`         | Job detail                                    |
| DELETE | `/api/jobs/:id`         | Delete job                                    |
| POST   | `/api/jobs/:id/match`   | Score all candidates against the job (ranked) |

See [backend/README.md](backend/README.md) and [frontend/README.md](frontend/README.md)
for details.
