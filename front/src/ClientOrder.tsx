import React, { useState, useEffect } from "react";
import "./css/RestaurantOrder.css";
import { format } from "date-fns";

type Menu = {
  id: number;
  name: string;
  price: number;
  description: string;
  restaurant_id: number;
};

type Order = {
  id: number;
  dish_id: number;
  quantity: number;
  status: string;
  date_time: string;
  menu?: Menu;
};

const ClientOrders: React.FC = () => {
  const [orders, setOrders] = useState<Order[]>([]);
  const email = localStorage.getItem("userEmail") || "";

  useEffect(() => {
    const fetchOrders = async () => {
      const response = await fetch(
        "https://go-food-court-29eb18b0ec35.herokuapp.com/graphql",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            query: `
            query GetALlOrdersByEmail($clientEmail: String!) {
              getALlOrdersByEmail(client_email: $clientEmail) {
                id
                dish_id
                quantity
                status
                date_time
              }
            }
          `,
            variables: { clientEmail: email },
          }),
        }
      );

      const responseData = await response.json();
      const fetchedOrders: Order[] = responseData.data.getALlOrdersByEmail;

      for (const order of fetchedOrders) {
        const menuResponse = await fetch(
          "https://go-food-court-29eb18b0ec35.herokuapp.com/graphql",
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              query: `
              query GetMenuByID($id: Int!) {
                getMenuByID(id: $id) {
                  id
                  name
                  price
                  description
                  restaurant_id
                }
              }
            `,
              variables: { id: order.dish_id },
            }),
          }
        );

        const menuData = await menuResponse.json();
        order.menu = menuData.data.getMenuByID;
      }

      setOrders(fetchedOrders);
    };

    if (email) {
      fetchOrders();
    }
  }, [email]);

  return (
    <div>
      <h2>My Orders</h2>
      <a href="/home">
        <button className="back-button">Retour à la carte</button>
      </a>
      {orders.length > 0 ? (
        <div className="table-container">
          <table>
            <thead>
              <tr>
                <th>ID Commande</th>
                <th>Nom du Plat</th>
                <th>Quantité</th>
                <th>Status</th>
                <th>Date/Heure</th>
                <th>Prix</th>
              </tr>
            </thead>
            <tbody>
              {orders.map((order) => (
                <tr key={order.id}>
                  <td>{order.id}</td>
                  <td>{order.menu?.name || "Inconnu"}</td>
                  <td>{order.quantity}</td>
                  <td>{order.status}</td>
                  <td>{format(new Date(order.date_time), "Pp")}</td>
                  <td>{`${order.menu?.price || 0} €`}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        <p>No orders found.</p>
      )}
    </div>
  );
};

export default ClientOrders;
