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

      console.log("Login response:", data);

      localStorage.setItem("userEmail", email);
      localStorage.setItem("userRole", data.role);
      localStorage.setItem("userID", data.userID);

      if (data.role === "client" || data.role === "admin") {
        navigate("/home");
      } else if (data.role === "restaurant") {
        navigate("/restaurant-orders");
      } else {
        console.error("Role not recognized:", data.role);
      }
    } catch (error) {
      console.error("Erreur lors de la connexion:", error);
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
