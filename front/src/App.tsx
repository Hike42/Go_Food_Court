import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Signup from "./Signup";
import Login from "./Login";
import MyOrders from "./ClientOrder";
import RestaurantOrder from "./RestaurantOrder";
import Index from ".";
import "./css/App.css";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Signup />} />
        <Route path="/login" element={<Login />} />
        <Route path="/home" element={<Index />} />
        <Route path="/my-orders" element={<MyOrders />} />
        <Route path="/restaurant-orders" element={<RestaurantOrder />} />
      </Routes>
    </Router>
  );
}

export default App;
