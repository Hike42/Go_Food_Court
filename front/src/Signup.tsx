import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import "./css/registerLogin.css";

const Signup: React.FC = () => {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  const navigateToLogin = () => {
    navigate("/login");
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      const response = await fetch(
        "https://go-food-court-29eb18b0ec35.herokuapp.com/api/signup",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ name, email, password }),
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      console.log(data);

      // Affichez un message de succès ou redirigez l'utilisateur ici
    } catch (error) {
      console.error("Erreur lors de la création du compte:", error);
      // Affichez un message d'erreur ici
    }
  };

  return (
    <div className="form-container">
      <form onSubmit={handleSubmit} className="form-box">
        <h2>Créer un compte</h2>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Name"
          className="input-field"
        />
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="Email"
          className="input-field"
        />
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="Password"
          className="input-field"
        />
        <button type="submit" className="button">
          Créer un compte
        </button>
        <div onClick={navigateToLogin} className="switch-form">
          Déjà un compte ? Connecte toi
        </div>
      </form>
    </div>
  );
};

export default Signup;
