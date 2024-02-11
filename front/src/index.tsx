import { useEffect, useState } from "react";
import Notification from "./components/Notification";
import { useNavigate } from "react-router-dom";
import "./css/homepage.css";

type Restaurant = {
  id: number;
  name: string;
  address: string;
  user_id: number;
};

type Menu = {
  id: number;
  name: string;
  price: number;
  description: string;
  restaurant_id: number;
};

const Index = () => {
  const [selectedRestaurant, setSelectedRestaurant] =
    useState<Restaurant | null>(null);
  const [restaurants, setRestaurants] = useState<Restaurant[]>([]);
  const [menus, setMenus] = useState<Menu[]>([]);

  const [cart, setCart] = useState<Menu[]>([]);

  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem("userEmail");
    localStorage.removeItem("userID");
    localStorage.removeItem("role");
    console.log("Local storage removed");
    navigate("/login");
  };

  const [isCartOpen, setIsCartOpen] = useState(false);

  const [notification, setNotification] = useState({
    show: false,
    message: "",
    type: "", // 'success', 'info', 'error'
  });

  const dateTime = new Date().toISOString().replace("T", " ").replace("Z", "");

  const addToCart = (menu: Menu) => {
    setCart([...cart, menu]);
  };

  const closeCart = () => {
    setIsCartOpen(false);
  };

  const openCart = () => {
    setIsCartOpen(true);
    setSelectedRestaurant(null); // Ajoutez cette ligne pour fermer la fenêtre modale du menu
  };

  const removeFromCart = (menu: Menu) => {
    const updatedCart = cart.filter((item) => item.id !== menu.id);
    setCart(updatedCart);
  };

  interface CartIndicatorProps {
    count: number;
  }

  const CartIndicator: React.FC<CartIndicatorProps> = ({ count }) => {
    if (count === 0) {
      return null; // Ne pas afficher si le panier est vide
    }

    return <div className="cart-indicator">{count}</div>;
  };

  useEffect(() => {
    // Fonction pour charger la liste des restaurants depuis l'API GraphQL
    const fetchRestaurants = () => {
      fetch("http://localhost:8080/graphql", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          query: `
                        query {
                            getAllRestaurants {
                                id
                                name
                                address
                                user_id
                            }
                        }
                    `,
        }),
      })
        .then((response) => response.json())
        .then((data) => {
          if (data && data.data && data.data.getAllRestaurants) {
            setRestaurants(data.data.getAllRestaurants);
          } else {
            // Gérer le cas où les données ne sont pas dans le format attendu
            console.error("Format de données inattendu:", data);
          }
        })
        .catch((error) => {
          console.error("Error fetching data: ", error);
        });
    };

    const client_email = localStorage.getItem("userEmail");
    console.log("Client email:", client_email);

    // Appeler la fonction pour charger les restaurants au chargement de la page
    fetchRestaurants();
  }, []);

  // Fonction pour charger la liste des menus par restaurant
  useEffect(() => {
    if (selectedRestaurant) {
      fetch("http://localhost:8080/graphql", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          query: `
                        query GetMenusByRestaurant($restaurantId: Int!) {
                            getMenusByRestaurant(restaurant_id: $restaurantId) {
                                id
                                name
                                price
                                description
                                restaurant_id
                            }
                        }
                    `,
          variables: {
            restaurantId: selectedRestaurant.id,
          },
        }),
      })
        .then((response) => response.json())
        .then((data) => {
          if (data && data.data && data.data.getMenusByRestaurant) {
            setMenus(data.data.getMenusByRestaurant);
          } else {
            // Gérer le cas où les données ne sont pas dans le format attendu
            console.error("Format de données inattendu:", data);
          }
        })
        .catch((error) => {
          console.error("Error fetching data: ", error);
        });
    }
  }, [selectedRestaurant]);

  // Fonction pour ouvrir la fenêtre modale
  const openModal = (restaurant: Restaurant) => {
    setSelectedRestaurant(restaurant);
    setIsCartOpen(false); // Ajoutez cette ligne pour fermer le panier
  };

  // Fonction pour fermer la fenêtre modale
  const closeModal = () => {
    setSelectedRestaurant(null);
  };
  const placeOrder = () => {
    cart.forEach((menu) => {
      // Préparer les variables pour la mutation
      const variables = {
        client_email: localStorage.getItem("userEmail"),
        dish_id: menu.id,
        quantity: 1,
        status: "pending",
        date_time: dateTime,
      };

      fetch("http://localhost:8080/graphql", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          query: `
        mutation CreateOrder($client_email: String!, $dish_id: Int!, $quantity: Int!, $status: String!, $date_time: String!) {
          createOrder(client_email: $client_email, dish_id: $dish_id, quantity: $quantity, status: $status, date_time: $date_time) {
            id
            client_email
            dish_id
            quantity
            status
            date_time
          }
        }`,
          variables: variables,
        }),
      })
        .then((response) => response.json())
        .then((data) => {
          console.log("Order created:", data);
          // Afficher la notification de succès
          setNotification({
            show: true,
            message: `Votre commande pour "${menu.name}" a été passée avec succès.`,
            type: "success",
          });
          setCart([]); // Vider le panier
          setIsCartOpen(false); // Fermer le panier
        })
        .catch((error) => {
          console.error("Error creating order:", error);
          // Ici, vous pouvez également déclencher une notification d'erreur
        });
    });
  };

  return (
    <div className="container">
      <h1>Liste des Restaurants</h1>
      <button onClick={handleLogout} className="disconnect">
        Déconnexion
      </button>
      <div className="subtitle">
        <button className="button-cart" onClick={openCart}>
          Mon Panier
          <CartIndicator count={cart.length} />
        </button>

        <a href="/my-orders">
          <button>Mes commandes</button>
        </a>
      </div>
      <ul className="restaurant-list">
        {restaurants.length === 0 ? (
          <p className="no-restaurant">Aucun restaurant disponible.</p>
        ) : (
          restaurants.map((restaurant, index) => (
            <li key={index}>
              <div className="restaurant-info">
                <div className="left">
                  <h2>{restaurant.name}</h2>
                  <p>Adresse: {restaurant.address}</p>
                </div>
                <div className="right">
                  <button onClick={() => openModal(restaurant)}>
                    Voir le menu
                  </button>
                </div>
              </div>
            </li>
          ))
        )}
      </ul>

      {notification.show && (
        <Notification
          message={notification.message}
          type={notification.type as "success" | "info" | "error" | undefined}
          onClose={() => setNotification({ ...notification, show: false })}
        />
      )}

      {/* Fenêtre modale pour afficher le menu */}
      {selectedRestaurant && (
        <div className="modal">
          <div className="modal-content">
            <div className="header">
              <h2>Menu de {selectedRestaurant.name}</h2>
              <button onClick={closeModal}>Fermer</button>
            </div>
            <ul>
              {menus.map((menu) => (
                <li key={menu.id} style={{ background: "#333" }}>
                  <h3>{menu.name}</h3>
                  <p>Prix: {menu.price} €</p>
                  <p>{menu.description}</p>
                  <button onClick={() => addToCart(menu)}>
                    Ajouter au panier
                  </button>
                </li>
              ))}
            </ul>
          </div>
        </div>
      )}

      {isCartOpen && (
        <div className="modal-content">
          <h2>Votre commande</h2>
          {cart.length === 0 ? (
            <p>Votre panier est vide.</p>
          ) : (
            <ul>
              {cart.map((item) => (
                <li key={item.id}>
                  <h3>{item.name}</h3>
                  <p>Prix: {item.price} €</p>
                  <p>Description: {item.description}</p>
                  <button onClick={() => removeFromCart(item)}>Retirer</button>
                </li>
              ))}
            </ul>
          )}
          <div className="validate_button">
            <button onClick={placeOrder}>Commander</button>
            <br />
            <button onClick={closeCart}>Fermer le panier</button>
          </div>
        </div>
      )}
    </div>
  );
};

export default Index;
