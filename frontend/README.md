# Frontend — AI Resume Screening Dashboard

Next.js (App Router) + TypeScript dashboard for the resume screening platform.

## Pages

| Route              | Purpose                                                      |
| ------------------ | ----------------------------------------------------------- |
| `/`                | Candidate pool — searchable list with skills preview        |
| `/upload`          | Upload a resume (PDF/TXT file or pasted text)               |
| `/candidates/[id]` | Candidate profile + AI summary + job-match history          |
| `/jobs`            | Create job postings and browse open jobs                    |
| `/jobs/[id]`       | Job detail + one-click AI matching with ranked candidates   |

## Setup

```bash
cp .env.local.example .env.local   # NEXT_PUBLIC_API_URL=http://localhost:8080
npm install
npm run dev
```

- `lib/api.ts` — typed client for the backend REST API.
- `components/ui.tsx` — small shared UI (score badge, skill chips, alerts, spinner).
- `app/globals.css` — dark theme design system.

The frontend is a thin, fully client-rendered layer over the Go API; the backend URL
is configured via `NEXT_PUBLIC_API_URL`.
