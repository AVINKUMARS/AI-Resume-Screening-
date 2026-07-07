"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { api, Job, Match } from "@/lib/api";
import { Chips, ErrorAlert, ScoreBadge, Spinner } from "@/components/ui";

export default function JobDetailPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const [job, setJob] = useState<Job | null>(null);
  const [matches, setMatches] = useState<Match[]>([]);
  const [failures, setFailures] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [matching, setMatching] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    (async () => {
      try {
        setJob(await api.getJob(id));
      } catch (e) {
        setError((e as Error).message);
      } finally {
        setLoading(false);
      }
    })();
  }, [id]);

  async function runMatch() {
    setMatching(true);
    setError(null);
    try {
      const res = await api.matchJob(id);
      setMatches(res.matches ?? []);
      setFailures(res.failures ?? []);
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setMatching(false);
    }
  }

  async function remove() {
    if (!confirm("Delete this job?")) return;
    await api.deleteJob(id);
    router.push("/jobs");
  }

  if (loading)
    return (
      <div className="row">
        <Spinner /> <span className="muted">Loading…</span>
      </div>
    );
  if (error && !job) return <ErrorAlert message={error} />;
  if (!job) return null;

  return (
    <div>
      <Link href="/jobs" className="small muted">
        ← Back to jobs
      </Link>

      <div className="row spread" style={{ margin: "12px 0" }}>
        <div>
          <h2 style={{ margin: 0 }}>{job.title}</h2>
          <div className="muted">
            {job.company}
            {job.location ? ` · ${job.location}` : ""}
            {job.min_years_experience ? ` · ${job.min_years_experience}+ yrs` : ""}
          </div>
        </div>
        <button className="btn danger" onClick={remove}>
          Delete
        </button>
      </div>

      <div className="card">
        <h3>Required skills</h3>
        <Chips items={job.required_skills} />
        {job.description && (
          <>
            <h3 style={{ marginTop: 16 }}>Description</h3>
            <p className="small">{job.description}</p>
          </>
        )}
      </div>

      <div className="row spread" style={{ margin: "24px 0 12px" }}>
        <h3 style={{ margin: 0 }}>Candidate ranking</h3>
        <button className="btn" onClick={runMatch} disabled={matching}>
          {matching ? <Spinner /> : null}
          {matching ? "Scoring candidates…" : "Run AI matching"}
        </button>
      </div>

      <ErrorAlert message={error} />

      {failures.length > 0 && (
        <div className="alert error">
          {failures.length} candidate(s) could not be scored: {failures.join("; ")}
        </div>
      )}

      {matches.length === 0 ? (
        <p className="muted small">
          Click <strong>Run AI matching</strong> to score every candidate in the pool against
          this job.
        </p>
      ) : (
        matches.map((m, idx) => (
          <div key={m.id} className="card" style={{ marginBottom: 10 }}>
            <div className="row spread">
              <div className="row">
                <span className="muted" style={{ width: 24 }}>
                  #{idx + 1}
                </span>
                <ScoreBadge score={m.score} />
                <Link href={`/candidates/${m.candidate_id}`}>
                  <strong>{m.candidate?.name ?? "Candidate"}</strong>
                  <span className="muted small"> · {m.candidate?.title}</span>
                </Link>
              </div>
            </div>
            <p className="small" style={{ margin: "10px 0 8px" }}>
              {m.reasoning}
            </p>
            <div className="small muted">Matched</div>
            <Chips items={m.matched_skills} variant="good" />
            <div className="small muted" style={{ marginTop: 6 }}>
              Missing
            </div>
            <Chips items={m.missing_skills} variant="bad" />
          </div>
        ))
      )}
    </div>
  );
}
