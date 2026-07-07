// Typed client for the Go backend API.

export interface Candidate {
  id: string;
  name: string;
  email: string;
  phone: string;
  location: string;
  title: string;
  summary: string;
  skills: string[];
  experience: string[];
  education: string[];
  years_experience: number;
  source_file: string;
  created_at: string;
  matches?: Match[];
}

export interface Job {
  id: string;
  title: string;
  company: string;
  location: string;
  description: string;
  required_skills: string[];
  min_years_experience: number;
  created_at: string;
}

export interface Match {
  id: string;
  candidate_id: string;
  job_id: string;
  score: number;
  reasoning: string;
  matched_skills: string[];
  missing_skills: string[];
  candidate?: Candidate;
  job?: Job;
}

export interface MatchResponse {
  job: Job;
  matches: Match[];
  failures: string[];
}

const BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    ...init,
    headers: {
      ...(init?.body && !(init.body instanceof FormData)
        ? { "Content-Type": "application/json" }
        : {}),
      ...init?.headers,
    },
  });

  if (!res.ok) {
    let message = `Request failed (${res.status})`;
    try {
      const data = await res.json();
      if (data?.error) message = data.error;
    } catch {
      /* non-JSON error body */
    }
    throw new Error(message);
  }
  return res.json() as Promise<T>;
}

export const api = {
  // Resumes / candidates
  uploadResumeFile: (file: File) => {
    const form = new FormData();
    form.append("resume", file);
    return request<Candidate>("/api/resumes/upload", { method: "POST", body: form });
  },
  uploadResumeText: (text: string, sourceFile = "pasted-resume.txt") =>
    request<Candidate>("/api/resumes/upload", {
      method: "POST",
      body: JSON.stringify({ text, source_file: sourceFile }),
    }),
  listCandidates: (q = "") =>
    request<Candidate[]>(`/api/candidates${q ? `?q=${encodeURIComponent(q)}` : ""}`),
  getCandidate: (id: string) => request<Candidate>(`/api/candidates/${id}`),
  deleteCandidate: (id: string) =>
    request<{ deleted: string }>(`/api/candidates/${id}`, { method: "DELETE" }),

  // Jobs
  listJobs: () => request<Job[]>("/api/jobs"),
  getJob: (id: string) => request<Job>(`/api/jobs/${id}`),
  createJob: (job: Partial<Job>) =>
    request<Job>("/api/jobs", { method: "POST", body: JSON.stringify(job) }),
  deleteJob: (id: string) =>
    request<{ deleted: string }>(`/api/jobs/${id}`, { method: "DELETE" }),
  matchJob: (id: string) =>
    request<MatchResponse>(`/api/jobs/${id}/match`, { method: "POST" }),
};
