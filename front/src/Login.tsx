import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import "./css/registerLogin.css";

const Login: React.FC = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      const response = await fetch(
        "https://go-food-court-29eb18b0ec35.herokuapp.com/api/login",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ email, password }),
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();

      // Supposons que `data` contienne une propriété `role` avec le rôle de l'utilisateur
      console.log("Login response:", data);

      // Stocker l'email et le rôle de l'utilisateur dans le localStorage ou contexte global
      localStorage.setItem("userEmail", email);
      localStorage.setItem("userRole", data.role); // Stockez le rôle pour une utilisation future
      localStorage.setItem("userID", data.userID); // Stockez l'ID pour une utilisation future

      // Rediriger l'utilisateur en fonction de son rôle
      if (data.role === "client" || data.role === "admin") {
        navigate("/home"); // Remplacez par le chemin vers le tableau de bord client
      } else if (data.role === "restaurant") {
        navigate("/restaurant-orders"); // Remplacez par le chemin vers le tableau de bord restaurant
      } else {
        console.error("Role not recognized:", data.role);
        // Gérer le cas d'un rôle non reconnu
      }
    } catch (error) {
      console.error("Erreur lors de la connexion:", error);
      // Gérer l'erreur ici (par exemple, afficher un message d'erreur)
    }
  };

  return (
    <div className="form-container content-container">
      <form onSubmit={handleSubmit} className="form-box">
        <h2 className="titleform">Connexion</h2>
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
          Connexion
        </button>
        <div onClick={() => navigate("/")} className="switch-form">
          Pas encore de compte ? Créer un compte
        </div>
      </form>
    </div>
  );
};

export default Login;
