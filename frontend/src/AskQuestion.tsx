import { useState } from 'react';
import { BACKEND_URL } from "./config";

export default function AskQuestion() {
  const [question, setQuestion] = useState('');
  const [filename, setFilename] = useState('');
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
        body: JSON.stringify({ question, filename }),
        credentials: 'include',
      });
      if (!res.ok) {
        const errorResponse = await res.json();
        throw new Error(errorResponse.error || 'Failed to get answer');
      }
      const data = await res.json();
      setAnswer(data.answer);
    } catch (e) {
      setError(e instanceof Error ? e.message : "An unknown error occurred");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div style={{ marginTop: "2rem", textAlign: "center", width: "100%", paddingInline: "1rem" }}>
      <h2 style={{ fontSize: "1.8rem", marginBottom: "1rem" }}>Ask a question</h2>
      
        <div style={{ display: "flex", flexDirection: "column", gap: "1rem", alignItems: "center" }}>
          <input
            type="text"
            value={filename}
            onChange={e => setFilename(e.target.value)}
            placeholder="Enter filename (e.g., Introduction_to_AI)"
            style={{ 
              padding: "0.75rem 1rem", 
              borderRadius: "8px", 
              border: "none", 
              fontSize: "1rem", 
              width: "100%", 
              maxWidth: "700px" 
            }}
          />
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
              width: "100%", 
              maxWidth: "700px" 
            }}
          />




        <button
          onClick={handleAsk}
          disabled={loading || !question || !filename}
          style={{
            padding: "0.75rem 1.5rem",
            borderRadius: "8px",
            border: "none",
            backgroundColor: "white",
            color: "black",
            fontSize: "1rem",
            cursor: loading || !question || !filename ? "not-allowed" : "pointer",
            opacity: loading || !question || !filename ? 0.6 : 1,
          }}
        >
          {loading ? "Asking..." : "Ask"}
        </button>
      </div>

      {answer && (
        <div style={{ display: "flex", justifyContent: "center" }}>
          <div style={{
            marginTop: "1rem",
            fontSize: "1.1rem",
            maxWidth: "800px",
            wordWrap: "break-word",
            whiteSpace: "pre-wrap",
            textAlign: "left"
          }}>
            <strong>Answer:</strong> {answer}
          </div>
        </div>
      )}

      {error && <p style={{ color: "red", marginTop: "1rem" }}>{error}</p>}
    </div>
  );
}
