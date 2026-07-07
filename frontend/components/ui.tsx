import React from "react";

export function ScoreBadge({ score }: { score: number }) {
  const tier = score >= 75 ? "high" : score >= 50 ? "mid" : "low";
  return <span className={`score ${tier}`}>{score}</span>;
}

export function Chips({ items, variant }: { items?: string[]; variant?: "good" | "bad" }) {
  if (!items || items.length === 0) return <span className="muted small">—</span>;
  return (
    <div>
      {items.map((item, i) => (
        <span key={i} className={`chip ${variant ?? ""}`}>
          {item}
        </span>
      ))}
    </div>
  );
}

export function ErrorAlert({ message }: { message?: string | null }) {
  if (!message) return null;
  return <div className="alert error">{message}</div>;
}

export function SuccessAlert({ message }: { message?: string | null }) {
  if (!message) return null;
  return <div className="alert success">{message}</div>;
}

export function Spinner() {
  return <span className="spinner" aria-label="loading" />;
}
