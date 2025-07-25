import { useState } from 'react';
import { BACKEND_URL } from "./config";

export default function AskQuestion() {
  const [question, setQuestion] = useState('');
  const [answer, setAnswer] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  async function handleAsk() {
    setLoading(true);
    setError('');
    setAnswer('');
    try {
      const token = sessionStorage.getItem("id_token"); 
      const res = await fetch(`${BACKEND_URL}/ask`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ question }),
        credentials: 'include', // REQUIRED for CORS + Authorization header
      });
      if (!res.ok) {
        throw new Error('Failed to get answer');
      }
      const data = await res.json();
      setAnswer(data.answer);
    } catch (e) {
      if (e instanceof Error) {
        setError(e.message);
      } else {
        setError("An unknown error occurred");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <div style={{ marginTop: "2rem", textAlign: "center", maxWidth: "600px", marginInline: "auto" }}>
      <h2 style={{ fontSize: "1.8rem", marginBottom: "1rem" }}>Ask a question</h2>
      <div style={{ display: "flex", justifyContent: "center", gap: "1rem", flexWrap: "wrap" }}>
        <input
          type="text"
          value={question}
          onChange={e => setQuestion(e.target.value)}
          placeholder="Enter your question"
          style={{
            padding: "0.75rem 1rem",
            borderRadius: "8px",
            border: "none",
            fontSize: "1rem",
            minWidth: "300px",
            maxWidth: "100%",
          }}
        />
        <button
          onClick={handleAsk}
          disabled={loading || !question}
          style={{
            padding: "0.75rem 1.5rem",
            borderRadius: "8px",
            border: "none",
            backgroundColor: "white",
            color: "black",
            fontSize: "1rem",
            cursor: loading || !question ? "not-allowed" : "pointer",
            opacity: loading || !question ? 0.6 : 1,
          }}
        >
          {loading ? "Asking..." : "Ask"}
        </button>
      </div>
      {answer && <p style={{ marginTop: "1rem", fontSize: "1.1rem" }}><strong>Answer:</strong> {answer}</p>}
      {error && <p style={{ color: "red", marginTop: "1rem" }}>{error}</p>}
    </div>
  );
}