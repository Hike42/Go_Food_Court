import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import "./css/registerLogin.css";
import Notification from "./components/Notification";

const Signup: React.FC = () => {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();
  const [notification, setNotification] = useState({
    show: false,
    message: "",
    type: "info", // "success", "error", "info"
  });

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

      setNotification({
        show: true,
        message: "Compte créé avec succès. Redirection...",
        type: "success",
      });

      setTimeout(() => navigate("/login"), 3000); // 3 secondes
    } catch (error) {
      console.error("Erreur lors de la création du compte:", error);
      setNotification({
        show: true,
        message: "Erreur lors de la création du compte.",
        type: "error",
      });
    }
  };

  return (
    <div className="form-container">
      {notification.show && (
        <Notification
          message={notification.message}
          type={notification.type as "info" | "success" | "error" | undefined}
          onClose={() => setNotification({ ...notification, show: false })}
        />
      )}

      <form onSubmit={handleSubmit} className="form-box">
        <h2 className="titleform">Créer un compte</h2>
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
