<h1 align="center">🧭 AI Resume Screening & Candidate Matching System</h1>

<p align="center">
  An AI-powered recruitment platform that parses resumes, manages a candidate pool,
  and ranks applicants against job openings with explainable match scores.
</p>

<p align="center">
  <img alt="Next.js"    src="https://img.shields.io/badge/Next.js-14-black?logo=next.js">
  <img alt="TypeScript" src="https://img.shields.io/badge/TypeScript-5-3178C6?logo=typescript&logoColor=white">
  <img alt="Go"         src="https://img.shields.io/badge/Go-Gin-00ADD8?logo=go&logoColor=white">
  <img alt="Postgres"   src="https://img.shields.io/badge/PostgreSQL-GORM-4169E1?logo=postgresql&logoColor=white">
  <img alt="Gemini"     src="https://img.shields.io/badge/Google-Gemini_API-8E75FF?logo=google&logoColor=white">
</p>

---

## ✨ Overview

Recruiters spend hours reading resumes and guessing who fits a role. This platform
automates the first pass: upload a resume and **Google Gemini** extracts structured
candidate data and a recruiter-facing summary; create a job and the same AI scores every
candidate **0–100** with a written rationale and matched/missing skills.

**Tech stack:** Next.js · TypeScript · Go (Gin) · PostgreSQL · GORM · Google Gemini API

```
ai-resume-screening/
├── backend/     # Go + Gin REST API · GORM/PostgreSQL · Gemini integration
└── frontend/    # Next.js (App Router) + TypeScript dashboard
```

## 🚀 Features

- **AI resume parsing** — upload a PDF/TXT resume (or paste text); Gemini extracts name,
  contact, title, skills, experience, education, and years of experience, plus a concise
  summary.
- **Candidate management** — searchable talent pool (by name, title, skill, or summary)
  with full candidate profiles.
- **Job postings** — define roles with required skills, minimum experience, and a
  description.
- **AI candidate matching** — score every candidate against a job with reasoning and
  matched/missing skills, ranked best-first.
- **Secure REST API** — clean endpoints with CORS scoped to the frontend origin; secrets
  kept in environment files, never committed.

## 🏗️ Architecture

```
┌──────────────────────┐        REST / JSON        ┌──────────────────────┐
│   Next.js Frontend   │  ───────────────────────▶ │    Go (Gin) API      │
│  (App Router + TS)   │  ◀─────────────────────── │   handlers/router    │
└──────────────────────┘                           └───────────┬──────────┘
                                                                │ GORM
                                          ┌─────────────────────┼─────────────────────┐
                                          ▼                     ▼                     ▼
                                    ┌───────────┐        ┌─────────────┐      ┌───────────────┐
                                    │ PostgreSQL│        │  Gemini API │      │ candidates /  │
                                    │ (GORM)    │        │ parse+match │      │ jobs / matches│
                                    └───────────┘        └─────────────┘      └───────────────┘
```

## ⚡ Quick start

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- A Google Gemini API key — <https://aistudio.google.com/app/apikey>

### 1. Database

Create an empty database (the backend auto-migrates the schema on startup):

```sql
CREATE DATABASE resume_screening;
```

### 2. Backend

```bash
cd backend
cp .env.example .env          # set DATABASE_URL and GEMINI_API_KEY
go mod tidy
go run .                      # http://localhost:8080
```

### 3. Frontend

```bash
cd frontend
cp .env.local.example .env.local   # NEXT_PUBLIC_API_URL=http://localhost:8080
npm install
npm run dev                        # http://localhost:3000
```

Open **http://localhost:3000** → **Upload** a resume → **Jobs** → **Run AI matching**.

## 🔌 API reference

| Method | Path                    | Description                                        |
| ------ | ----------------------- | -------------------------------------------------- |
| GET    | `/health`               | Liveness probe                                     |
| POST   | `/api/resumes/upload`   | Upload resume (multipart `resume` or JSON `text`); parses + stores candidate |
| GET    | `/api/candidates?q=`    | List / search candidates                           |
| GET    | `/api/candidates/:id`   | Candidate detail with match history                |
| DELETE | `/api/candidates/:id`   | Delete candidate                                   |
| POST   | `/api/jobs`             | Create a job posting                               |
| GET    | `/api/jobs`             | List jobs                                          |
| GET    | `/api/jobs/:id`         | Job detail                                         |
| DELETE | `/api/jobs/:id`         | Delete job                                         |
| POST   | `/api/jobs/:id/match`   | Score all candidates against the job (ranked)      |

<details>
<summary>Example: parse a resume from raw text</summary>

```bash
curl -X POST localhost:8080/api/resumes/upload \
  -H 'Content-Type: application/json' \
  -d '{"text":"Jane Doe — Senior Go Engineer. 6 years building REST APIs with Gin, PostgreSQL, GORM..."}'
```
</details>

## 🔐 Configuration

| Variable            | Where           | Description                              |
| ------------------- | --------------- | ---------------------------------------- |
| `DATABASE_URL`      | `backend/.env`  | PostgreSQL connection string             |
| `GEMINI_API_KEY`    | `backend/.env`  | Google Gemini API key (**required**)     |
| `GEMINI_MODEL`      | `backend/.env`  | Model id (default `gemini-2.5-flash`)    |
| `ALLOWED_ORIGIN`    | `backend/.env`  | CORS origin for the frontend             |
| `NEXT_PUBLIC_API_URL` | `frontend/.env.local` | Backend base URL                 |

> `.env` files are git-ignored — copy the provided `.env.example` templates and fill in
> your own values.

## 📦 Project structure

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

frontend/
├── app/                        # App Router pages
│   ├── page.tsx                # candidate pool
│   ├── upload/                 # resume upload
│   ├── candidates/[id]/        # candidate profile + matches
│   └── jobs/                   # jobs + matching
├── components/ui.tsx           # shared UI
└── lib/api.ts                  # typed API client
```

See [backend/README.md](backend/README.md) and [frontend/README.md](frontend/README.md)
for module-level details.

## 📄 License

Released for educational and portfolio use.
