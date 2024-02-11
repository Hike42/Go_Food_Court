import React, { useEffect } from "react";

// Définir les types des props
interface NotificationProps {
  message: string;
  type?: "info" | "success" | "error";
  onClose: () => void;
}

const Notification: React.FC<NotificationProps> = ({
  message,
  type = "info",
  onClose,
}) => {
  // Définir une durée avant de fermer la notification automatiquement
  useEffect(() => {
    const timer = setTimeout(() => {
      onClose();
    }, 10000); // 10 secondes

    return () => clearTimeout(timer);
  }, [onClose]);

  // Styles de base pour la notification
  const baseStyle = {
    position: "fixed",
    top: "20px",
    right: "20px",
    padding: "10px 20px",
    borderRadius: "10px",
    color: "white",
    backgroundColor:
      type === "success" ? "green" : type === "error" ? "red" : "blue",
    zIndex: 1000,
    animation: "slideIn 0.5s ease-out",
  } as React.CSSProperties;

  // CSS pour l'animation
  const animationStyle = `
    @keyframes slideIn {
      from {
        opacity: 0;
        transform: translateX(100%);
      }
      to {
        opacity: 1;
        transform: translateX(0);
      }
    }
  `;

  return (
    <>
      <style>{animationStyle}</style>
      <div style={baseStyle}>{message}</div>
    </>
  );
};

export default Notification;
