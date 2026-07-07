"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { api, Candidate } from "@/lib/api";
import { Chips, ErrorAlert, ScoreBadge, Spinner } from "@/components/ui";

export default function CandidateDetailPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const [candidate, setCandidate] = useState<Candidate | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    (async () => {
      try {
        setCandidate(await api.getCandidate(id));
      } catch (e) {
        setError((e as Error).message);
      } finally {
        setLoading(false);
      }
    })();
  }, [id]);

  async function remove() {
    if (!confirm("Delete this candidate?")) return;
    await api.deleteCandidate(id);
    router.push("/");
  }

  if (loading)
    return (
      <div className="row">
        <Spinner /> <span className="muted">Loading…</span>
      </div>
    );
  if (error) return <ErrorAlert message={error} />;
  if (!candidate) return null;

  return (
    <div>
      <Link href="/" className="small muted">
        ← Back to candidates
      </Link>

      <div className="row spread" style={{ margin: "12px 0" }}>
        <div>
          <h2 style={{ margin: 0 }}>{candidate.name || "Unnamed candidate"}</h2>
          <div className="muted">
            {candidate.title}
            {candidate.location ? ` · ${candidate.location}` : ""}
            {candidate.years_experience ? ` · ${candidate.years_experience} yrs exp` : ""}
          </div>
        </div>
        <button className="btn danger" onClick={remove}>
          Delete
        </button>
      </div>

      <div className="grid cols-2">
        <div className="card">
          <h3>Summary</h3>
          <p className="small">{candidate.summary || "—"}</p>

          <h3 style={{ marginTop: 18 }}>Contact</h3>
          <p className="small muted" style={{ margin: 0 }}>
            {candidate.email || "no email"} · {candidate.phone || "no phone"}
          </p>

          <h3 style={{ marginTop: 18 }}>Skills</h3>
          <Chips items={candidate.skills} />
        </div>

        <div className="card">
          <h3>Experience</h3>
          <ul className="small">
            {candidate.experience?.length
              ? candidate.experience.map((x, i) => <li key={i}>{x}</li>)
              : "—"}
          </ul>

          <h3 style={{ marginTop: 18 }}>Education</h3>
          <ul className="small">
            {candidate.education?.length
              ? candidate.education.map((x, i) => <li key={i}>{x}</li>)
              : "—"}
          </ul>
        </div>
      </div>

      <h3 style={{ marginTop: 24 }}>Job match history</h3>
      {candidate.matches && candidate.matches.length > 0 ? (
        candidate.matches.map((m) => (
          <div key={m.id} className="card" style={{ marginBottom: 10 }}>
            <div className="row spread">
              <strong>{m.job?.title ?? "Job"}</strong>
              <ScoreBadge score={m.score} />
            </div>
            <p className="small" style={{ marginBottom: 8 }}>
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
      ) : (
        <p className="muted small">
          No matches yet. Run matching from the <Link href="/jobs">Jobs</Link> page.
        </p>
      )}
    </div>
  );
}
