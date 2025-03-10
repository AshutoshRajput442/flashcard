import { useEffect, useState } from "react";
import axios from "axios";
import "./App.css"; 

export default function App() {
  const [flashcards, setFlashcards] = useState([]);
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [selectedFlashcard, setSelectedFlashcard] = useState(null);

  useEffect(() => {
    fetchFlashcards();
  }, []);

  const fetchFlashcards = async () => {
    setLoading(true);
    setError("");
    try {
      const response = await axios.get("http://localhost:8080/flashcards");
      setFlashcards(response.data || []);
    } catch (error) {
      console.error("Error fetching flashcards:", error);
      setError("Failed to load flashcards. Try again!");
    } finally {
      setLoading(false);
    }
  };

  const createFlashcard = async () => {
    if (!title.trim() || !content.trim()) {
      setError("Title and Content cannot be empty!");
      return;
    }
    setError("");
    try {
      await axios.post("http://localhost:8080/flashcards", { title, content });
      setTitle("");
      setContent("");
      fetchFlashcards();
    } catch (error) {
      console.error("Error creating flashcard:", error);
      setError("Failed to create flashcard.");
    }
  };

  const deleteFlashcard = async (id) => {
    setError("");
    try {
      await axios.delete(`http://localhost:8080/flashcards/${id}`);
      fetchFlashcards();
    } catch (error) {
      console.error("Error deleting flashcard:", error);
      setError("Failed to delete flashcard.");
    }
  };

  return (
    <div className="container">
      <h1>📚 Flashcards & Quick Notes</h1>

      {error && <p className="error">{error}</p>}

      <div className="form">
        <input type="text" placeholder="Enter Title" value={title} onChange={(e) => setTitle(e.target.value)} />
        <textarea placeholder="Enter Content" value={content} onChange={(e) => setContent(e.target.value)} />
        <button onClick={createFlashcard} className="add-btn">➕ Add Flashcard</button>
      </div>

      {loading ? <p>Loading flashcards...</p> : (
        <ul className="flashcard-list">
          {flashcards.map((fc) => (
            <li key={fc.id} className="flashcard">
              <h3>{fc.title}</h3>
              <div className="flashcard-buttons">
                <button className="view-btn" onClick={() => setSelectedFlashcard(fc)}>👁 View</button>
                <button className="delete-btn" onClick={() => deleteFlashcard(fc.id)}>❌ Delete</button>
              </div>
            </li>
          ))}
        </ul>
      )}

      {selectedFlashcard && (
        <div className="full-content">
          <div className="full-content-box">
            <h2>{selectedFlashcard.title}</h2>
            <p>{selectedFlashcard.content}</p>
            <button className="close-btn" onClick={() => setSelectedFlashcard(null)}>❌ Close</button>
          </div>
        </div>
      )}
    </div>
  );
}
