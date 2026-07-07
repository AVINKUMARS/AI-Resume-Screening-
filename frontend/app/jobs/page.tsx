"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api, Job } from "@/lib/api";
import { Chips, ErrorAlert, Spinner } from "@/components/ui";

export default function JobsPage() {
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // create-form state
  const [title, setTitle] = useState("");
  const [company, setCompany] = useState("");
  const [location, setLocation] = useState("");
  const [skills, setSkills] = useState("");
  const [minExp, setMinExp] = useState("");
  const [description, setDescription] = useState("");
  const [creating, setCreating] = useState(false);

  async function load() {
    setLoading(true);
    try {
      setJobs(await api.listJobs());
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, []);

  async function create(e: React.FormEvent) {
    e.preventDefault();
    setCreating(true);
    setError(null);
    try {
      await api.createJob({
        title,
        company,
        location,
        description,
        required_skills: skills
          .split(",")
          .map((s) => s.trim())
          .filter(Boolean),
        min_years_experience: parseFloat(minExp) || 0,
      });
      setTitle("");
      setCompany("");
      setLocation("");
      setSkills("");
      setMinExp("");
      setDescription("");
      await load();
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setCreating(false);
    }
  }

  return (
    <div className="grid cols-2">
      <div>
        <h2>Post a job</h2>
        <form onSubmit={create} className="card">
          <label>Title *</label>
          <input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="Senior Backend Engineer" />
          <label>Company</label>
          <input value={company} onChange={(e) => setCompany(e.target.value)} placeholder="Scaling Wolf" />
          <label>Location</label>
          <input value={location} onChange={(e) => setLocation(e.target.value)} placeholder="Remote / Bengaluru" />
          <label>Required skills (comma-separated)</label>
          <input value={skills} onChange={(e) => setSkills(e.target.value)} placeholder="Go, PostgreSQL, REST APIs" />
          <label>Minimum years of experience</label>
          <input value={minExp} onChange={(e) => setMinExp(e.target.value)} placeholder="4" type="number" min="0" />
          <label>Description</label>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Responsibilities, tech stack, must-haves…"
          />
          <div style={{ marginTop: 14 }}>
            <button className="btn" type="submit" disabled={!title.trim() || creating}>
              {creating ? <Spinner /> : null} Create job
            </button>
          </div>
        </form>
        <div style={{ marginTop: 12 }}>
          <ErrorAlert message={error} />
        </div>
      </div>

      <div>
        <h2>Open jobs</h2>
        {loading ? (
          <div className="row">
            <Spinner /> <span className="muted">Loading…</span>
          </div>
        ) : jobs.length === 0 ? (
          <div className="card muted">No jobs yet — create one to start matching.</div>
        ) : (
          jobs.map((j) => (
            <Link key={j.id} href={`/jobs/${j.id}`} className="list-item">
              <div className="row spread">
                <strong>{j.title}</strong>
                <span className="muted small">{j.company}</span>
              </div>
              <div className="muted small">
                {j.location || "—"}
                {j.min_years_experience ? ` · ${j.min_years_experience}+ yrs` : ""}
              </div>
              <div style={{ marginTop: 6 }}>
                <Chips items={j.required_skills?.slice(0, 6)} />
              </div>
            </Link>
          ))
        )}
      </div>
    </div>
  );
}
