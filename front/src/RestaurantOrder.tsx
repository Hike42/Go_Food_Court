import { useState, useEffect } from "react";
import { format } from "date-fns";
import "./css/RestaurantOrder.css"; // Assurez-vous d'avoir ce fichier CSS
import { useNavigate } from "react-router-dom";

type Order = {
  id: number;
  dish_id: number;
  quantity: number;
  status: string;
  date_time: string;
  client_email: string;
  menuName?: string;
};

const ActionDialog = ({
  order,
  onClose,
  onUpdateOrderStatus,
  onDeleteOrder,
}: {
  order: Order;
  onClose: () => void;
  onUpdateOrderStatus: (orderId: number, newStatus: string) => void;
  onDeleteOrder: (orderId: number) => void;
}) => {
  const [selectedStatus, setSelectedStatus] = useState(order.status);

  return (
    <div className="dialog-overlay">
      <div className="dialog">
        <h3>Modifier la commande {order.id}</h3>
        <select
          value={selectedStatus}
          onChange={(e) => setSelectedStatus(e.target.value)}
        >
          <option value="pending">En attente</option>
          <option value="cooking">En préparation</option>
          <option value="ready">Prête</option>
        </select>
        <div className="dialog-actions">
          <button onClick={() => onUpdateOrderStatus(order.id, selectedStatus)}>
            Confirmer
          </button>
          <button onClick={() => onDeleteOrder(order.id)}>Supprimer</button>
          <button onClick={onClose}>Annuler</button>
        </div>
      </div>
    </div>
  );
};

export default function RestaurantOrder() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [currentOrder, setCurrentOrder] = useState<Order | null>(null);
  const restaurantID = localStorage.getItem("userID"); // Supposons que l'ID du restaurant est stocké dans le localStorage
  const navigate = useNavigate();

  useEffect(() => {
    const fetchOrders = async () => {
      const response = await fetch("http://localhost:8080/graphql", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          query: `
          query GetOrdersbyRestaurant($restaurant_id: Int!) {
            getOrdersbyRestaurant(restaurant_id: $restaurant_id) {
              id
              dish_id
              quantity
              status
              date_time
              client_email
            }
          }
        `,
          variables: { restaurant_id: parseInt(restaurantID || "0") },
        }),
      });
      const { data } = await response.json();
      if (data && data.getOrdersbyRestaurant) {
        const ordersWithNames = await Promise.all(
          data.getOrdersbyRestaurant.map(async (order: any) => {
            const menuResponse = await fetch("http://localhost:8080/graphql", {
              method: "POST",
              headers: { "Content-Type": "application/json" },
              body: JSON.stringify({
                query: `
              query GetMenuByID($id: Int!) {
                getMenuByID(id: $id) {
                  name
                }
              }
            `,
                variables: { id: order.dish_id },
              }),
            });
            const menuData = await menuResponse.json();
            return { ...order, menuName: menuData.data.getMenuByID.name };
          })
        );
        setOrders(ordersWithNames);
      }
    };

    if (restaurantID) {
      fetchOrders();
    }
  }, [restaurantID]);

  const handleUpdateOrderStatus = async (
    orderId: number,
    newStatus: string
  ) => {
    const response = await fetch("http://localhost:8080/graphql", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        query: `
        mutation UpdateOrderStatus($id: Int!, $status: String!) {
          updateOrderStatus(id: $id, status: $status) {
            id
            status
          }
        }
      `,
        variables: { id: orderId, status: newStatus },
      }),
    });

    const { data, errors } = await response.json();

    if (errors) {
      console.error("Failed to update order status:", errors);
    } else if (data) {
      // Mise à jour réussie, mettez à jour l'état local si nécessaire
      setOrders((prevOrders) =>
        prevOrders.map((order) =>
          order.id === orderId ? { ...order, status: newStatus } : order
        )
      );
    }

    setCurrentOrder(null); // Fermer le dialogue après l'action
  };

  const handleDeleteOrder = async (orderId: number) => {
    const response = await fetch("http://localhost:8080/graphql", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        query: `
        mutation DeleteOrder($id: Int!) {
          deleteOrder(id: $id) {
            id
          }
        }
      `,
        variables: { id: orderId },
      }),
    });

    const { data, errors } = await response.json();

    if (errors) {
      console.error("Failed to delete order:", errors);
    } else if (data) {
      // Suppression réussie, mettez à jour l'état local
      setOrders((prevOrders) =>
        prevOrders.filter((order) => order.id !== orderId)
      );
    }

    setCurrentOrder(null); // Fermer le dialogue après l'action
  };

  const handleLogout = () => {
    localStorage.removeItem("userEmail");
    localStorage.removeItem("userID");
    localStorage.removeItem("role");
    console.log("Local storage removed");
    navigate("/login");
  };

  return (
    <div className="orders-container">
      <button onClick={handleLogout} className="disconnect">
        Déconnexion
      </button>
      <h2>Commandes du Restaurant</h2>
      <table>
        <thead>
          <tr>
            <th>ID Commande</th>
            <th>Email Client</th>
            <th>Nom du Plat</th>
            <th>Quantité</th>
            <th>Statut</th>
            <th>Date/Heure</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {orders.map((order) => (
            <tr key={order.id}>
              <td>{order.id}</td>
              <td>{order.client_email}</td>
              <td>{order.menuName || "Inconnu"}</td>
              <td>{order.quantity}</td>
              <td>{order.status}</td>
              <td>{format(new Date(order.date_time), "Pp")}</td>
              <td>
                <button onClick={() => setCurrentOrder(order)}>Actions</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      {currentOrder && (
        <ActionDialog
          order={currentOrder}
          onClose={() => setCurrentOrder(null)}
          onUpdateOrderStatus={handleUpdateOrderStatus}
          onDeleteOrder={handleDeleteOrder}
        />
      )}
    </div>
  );
}
