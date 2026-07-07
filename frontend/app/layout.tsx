import "./globals.css";
import type { Metadata } from "next";
import Link from "next/link";

export const metadata: Metadata = {
  title: "AI Resume Screening",
  description: "AI-powered resume parsing and candidate matching",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <header className="header">
          <div className="brand">
            <span style={{ fontSize: 22 }}>🧭</span>
            <h1>AI Resume Screening</h1>
          </div>
          <nav>
            <Link href="/">Candidates</Link>
            <Link href="/upload">Upload</Link>
            <Link href="/jobs">Jobs &amp; Matching</Link>
          </nav>
        </header>
        <main className="container">{children}</main>
      </body>
    </html>
  );
}
