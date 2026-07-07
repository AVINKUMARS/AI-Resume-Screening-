"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api";
import { ErrorAlert, Spinner, SuccessAlert } from "@/components/ui";

export default function UploadPage() {
  const router = useRouter();
  const [tab, setTab] = useState<"file" | "text">("file");
  const [file, setFile] = useState<File | null>(null);
  const [text, setText] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setBusy(true);
    setError(null);
    setSuccess(null);
    try {
      const candidate =
        tab === "file"
          ? await api.uploadResumeFile(file!)
          : await api.uploadResumeText(text);
      setSuccess(`Parsed "${candidate.name || "candidate"}" successfully. Redirecting…`);
      setTimeout(() => router.push(`/candidates/${candidate.id}`), 900);
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setBusy(false);
    }
  }

  const canSubmit = tab === "file" ? !!file : text.trim().length > 30;

  return (
    <div style={{ maxWidth: 680, margin: "0 auto" }}>
      <h2>Upload a resume</h2>
      <p className="muted small">
        The resume is parsed with Google Gemini into structured fields, skills, and a
        recruiter-facing summary.
      </p>

      <div className="row" style={{ margin: "16px 0" }}>
        <button
          className={`btn ${tab === "file" ? "" : "ghost"}`}
          onClick={() => setTab("file")}
          type="button"
        >
          Upload file
        </button>
        <button
          className={`btn ${tab === "text" ? "" : "ghost"}`}
          onClick={() => setTab("text")}
          type="button"
        >
          Paste text
        </button>
      </div>

      <form onSubmit={submit} className="card">
        {tab === "file" ? (
          <div>
            <label>Resume file (PDF or TXT)</label>
            <input
              type="file"
              accept=".pdf,.txt"
              onChange={(e) => setFile(e.target.files?.[0] ?? null)}
            />
            {file && <p className="muted small">Selected: {file.name}</p>}
          </div>
        ) : (
          <div>
            <label>Resume text</label>
            <textarea
              placeholder="Paste the full resume text here…"
              value={text}
              onChange={(e) => setText(e.target.value)}
              style={{ minHeight: 240 }}
            />
          </div>
        )}

        <div style={{ marginTop: 16 }}>
          <button className="btn" type="submit" disabled={!canSubmit || busy}>
            {busy ? <Spinner /> : null}
            {busy ? "Parsing with AI…" : "Parse & save"}
          </button>
        </div>
      </form>

      <div style={{ marginTop: 16 }}>
        <ErrorAlert message={error} />
        <SuccessAlert message={success} />
      </div>
    </div>
  );
}
