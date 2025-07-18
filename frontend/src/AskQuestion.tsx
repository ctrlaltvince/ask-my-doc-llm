import React, { useState } from 'react';

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
      const res = await fetch('http://localhost:8081/ask', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ question }),
      });
      if (!res.ok) {
        throw new Error('Failed to get answer');
      }
      const data = await res.json();
      setAnswer(data.answer);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div>
      <h2>Ask a question</h2>
      <input
        type="text"
        value={question}
        onChange={e => setQuestion(e.target.value)}
        placeholder="Enter your question"
      />
      <button onClick={handleAsk} disabled={loading || !question}>
        {loading ? 'Asking...' : 'Ask'}
      </button>
      {answer && <p><strong>Answer:</strong> {answer}</p>}
      {error && <p style={{color: 'red'}}>{error}</p>}
    </div>
  );
}
