"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api, Candidate } from "@/lib/api";
import { Chips, ErrorAlert, Spinner } from "@/components/ui";

export default function CandidatesPage() {
  const [candidates, setCandidates] = useState<Candidate[]>([]);
  const [q, setQ] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  async function load(query = "") {
    setLoading(true);
    setError(null);
    try {
      setCandidates(await api.listCandidates(query));
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, []);

  function onSearch(e: React.FormEvent) {
    e.preventDefault();
    load(q);
  }

  return (
    <div>
      <div className="row spread" style={{ marginBottom: 16 }}>
        <div>
          <h2 style={{ margin: 0 }}>Candidates</h2>
          <p className="muted small" style={{ margin: "4px 0 0" }}>
            {candidates.length} candidate{candidates.length === 1 ? "" : "s"} in the talent pool
          </p>
        </div>
        <Link className="btn" href="/upload">
          + Upload resume
        </Link>
      </div>

      <form onSubmit={onSearch} className="row" style={{ marginBottom: 16 }}>
        <input
          placeholder="Search by name, title, or skill…"
          value={q}
          onChange={(e) => setQ(e.target.value)}
        />
        <button className="btn ghost" type="submit">
          Search
        </button>
      </form>

      <ErrorAlert message={error} />

      {loading ? (
        <div className="row">
          <Spinner /> <span className="muted">Loading candidates…</span>
        </div>
      ) : candidates.length === 0 ? (
        <div className="card muted">
          No candidates yet. <Link href="/upload">Upload a resume</Link> to get started.
        </div>
      ) : (
        candidates.map((c) => (
          <Link key={c.id} href={`/candidates/${c.id}`} className="list-item">
            <div className="row spread">
              <div>
                <strong style={{ fontSize: 16 }}>{c.name || "Unnamed candidate"}</strong>
                <div className="muted small">
                  {c.title || "—"}
                  {c.location ? ` · ${c.location}` : ""}
                  {c.years_experience ? ` · ${c.years_experience} yrs` : ""}
                </div>
              </div>
              <span className="muted small">{c.email}</span>
            </div>
            <p className="small" style={{ margin: "10px 0 6px" }}>
              {c.summary}
            </p>
            <Chips items={c.skills?.slice(0, 8)} />
          </Link>
        ))
      )}
    </div>
  );
}
