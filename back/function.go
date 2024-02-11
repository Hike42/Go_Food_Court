package main

import (
	"database/sql"
	"log"
)

func findUserByID(db *sql.DB, id int) (*User, error) {
	var user User
	query := `SELECT id, name, email, password, role FROM Users WHERE id = ?`
	if err := db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role); err != nil {
		return nil, err
	}
	return &user, nil
}

func findAllUsers(db *sql.DB) ([]User, error) {
	var users []User
	query := `SELECT id, name, email, password, role FROM Users`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func findRestaurantByID(db *sql.DB, id int) (*Restaurant, error) {
	var restaurant Restaurant
	query := `SELECT id, name, address, user_id FROM Restaurants WHERE id = ?`
	if err := db.QueryRow(query, id).Scan(&restaurant.ID, &restaurant.Name, &restaurant.Address, &restaurant.UserID); err != nil {
		return nil, err
	}
	return &restaurant, nil
}

func findAllRestaurants(db *sql.DB) ([]Restaurant, error) {
	var restaurants []Restaurant
	query := `SELECT id, name, address, user_id FROM Restaurants`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var restaurant Restaurant
		if err := rows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Address, &restaurant.UserID); err != nil {
			return nil, err
		}
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}

func findOrdersByRestaurantID(db *sql.DB, restaurantID int) ([]Order, error) {
	query := `SELECT o.id, o.client_email, o.dish_id, o.quantity, o.status, o.date_time
		FROM Orders o
		INNER JOIN Menus m ON o.dish_id = m.id
		WHERE m.restaurant_id = ?`
	rows, err := db.Query(query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.ClientEmail, &order.DishID, &order.Quantity, &order.Status, &order.DateTime); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func findMenuByID(db *sql.DB, id int) (*Menu, error) {
	var menu Menu
	query := `SELECT id, name, price, description, restaurant_id FROM Menus WHERE id = ?`
	if err := db.QueryRow(query, id).Scan(&menu.ID, &menu.Name, &menu.Price, &menu.Description, &menu.RestaurantID); err != nil {
		return nil, err
	}
	return &menu, nil
}

func findMenusByRestaurant(db *sql.DB, restaurantID int) ([]Menu, error) {
	query := `SELECT id, name, price, description, restaurant_id FROM Menus WHERE restaurant_id = ?`
	rows, err := db.Query(query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menus []Menu
	for rows.Next() {
		var menu Menu
		if err := rows.Scan(&menu.ID, &menu.Name, &menu.Price, &menu.Description, &menu.RestaurantID); err != nil {
			return nil, err
		}
		menus = append(menus, menu)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return menus, nil
}

func findAllMenus(db *sql.DB) ([]Menu, error) {
	var menus []Menu
	query := `SELECT id, name, price, description, restaurant_id FROM Menus`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var menu Menu
		if err := rows.Scan(&menu.ID, &menu.Name, &menu.Price, &menu.Description, &menu.RestaurantID); err != nil {
			return nil, err
		}
		menus = append(menus, menu)
	}

	return menus, nil
}

func findOrderByID(db *sql.DB, id int) (*Order, error) {
	var order Order
	query := `SELECT id, client_email, dish_id, quantity, status, date_time FROM Orders WHERE id = ?`
	if err := db.QueryRow(query, id).Scan(&order.ID, &order.ClientEmail, &order.DishID, &order.Quantity, &order.Status, &order.DateTime); err != nil {
		return nil, err
	}
	return &order, nil
}

func findAllOrdersbyUserEmail(db *sql.DB, email string) ([]Order, error) {
	var orders []Order
	query := `SELECT id, client_email, dish_id, quantity, status, date_time FROM Orders WHERE client_email = ?`
	rows, err := db.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.ClientEmail, &order.DishID, &order.Quantity, &order.Status, &order.DateTime); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func findOrderStatusByID(db *sql.DB, id int) (*OrderStatus, error) {
	var orderStatus OrderStatus
	query := `SELECT order_id, status, date_time FROM OrderStatus WHERE order_id = ?`
	if err := db.QueryRow(query, id).Scan(&orderStatus.OrderID, &orderStatus.Status, &orderStatus.DateTime); err != nil {
		return nil, err
	}
	return &orderStatus, nil
}

// Mutations

func createUser(db *sql.DB, user User) (*User, error) {
	if user.Role == "" {
		user.Role = "client"
	}
	query := `INSERT INTO Users (name, email, password, role) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, user.Name, user.Email, user.Password, user.Role)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return nil, err
	}

	user.ID = int(id)
	return &user, nil
}

func updateUser(db *sql.DB, user User) (*User, error) {
	query := `UPDATE Users SET name = ?, email = ?, password = ?, role = ? WHERE id = ?`
	_, err := db.Exec(query, user.Name, user.Email, user.Password, user.Role, user.ID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func deleteUser(db *sql.DB, id int) error {
	query := `DELETE FROM Users WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

func createRestaurant(db *sql.DB, restaurant Restaurant) (*Restaurant, error) {
	query := `INSERT INTO Restaurants (name, address, user_id) VALUES (?, ?, ?)`
	_, err := db.Exec(query, restaurant.Name, restaurant.Address, restaurant.UserID)
	if err != nil {
		return nil, err
	}
	return &restaurant, nil
}

func updateRestaurant(db *sql.DB, restaurant Restaurant) (*Restaurant, error) {
	query := `UPDATE Restaurants SET name = ?, address = ?, user_id = ? WHERE id = ?`
	_, err := db.Exec(query, restaurant.Name, restaurant.Address, restaurant.UserID, restaurant.ID)
	if err != nil {
		return nil, err
	}
	return &restaurant, nil
}

func deleteRestaurant(db *sql.DB, id int) error {
	query := `DELETE FROM Restaurants WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

func createMenu(db *sql.DB, menu Menu) (*Menu, error) {
	query := `INSERT INTO Menus (name, price, description, restaurant_id) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, menu.Name, menu.Price, menu.Description, menu.RestaurantID)
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func updateMenu(db *sql.DB, menu Menu) (*Menu, error) {
	query := `UPDATE Menus SET name = ?, price = ?, description = ?, restaurant_id = ? WHERE id = ?`
	_, err := db.Exec(query, menu.Name, menu.Price, menu.Description, menu.RestaurantID, menu.ID)
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func deleteMenu(db *sql.DB, id int) error {
	query := `DELETE FROM Menus WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

func createOrder(db *sql.DB, order Order) (*Order, error) {
	query := `INSERT INTO Orders (client_email, dish_id, quantity, status, date_time) VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, order.ClientEmail, order.DishID, order.Quantity, order.Status, order.DateTime)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func updateOrder(db *sql.DB, order Order) (*Order, error) {
	query := `UPDATE Orders SET client_email = ?, dish_id = ?, quantity = ?, status = ?, date_time = ? WHERE id = ?`
	_, err := db.Exec(query, order.ClientEmail, order.DishID, order.Quantity, order.Status, order.DateTime, order.ID)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func deleteOrder(db *sql.DB, id int) error {
	query := `DELETE FROM Orders WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

func checkUserCredentials(db *sql.DB, email, password string) (bool, string, int, error) {
	var dbPassword, role string
	var userID int // Ajouté pour stocker l'ID de l'utilisateur
	// Modifiez la requête pour récupérer également l'ID de l'utilisateur
	err := db.QueryRow("SELECT id, password, role FROM Users WHERE email = ?", email).Scan(&userID, &dbPassword, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			// Aucun utilisateur trouvé avec cet email
			return false, "", 0, nil
		}
		// Une erreur s'est produite lors de la requête
		return false, "", 0, err
	}

	// Ici, comparez le mot de passe fourni avec celui en base de données.
	// Si vous stockez des mots de passe hachés, utilisez une fonction de comparaison appropriée.
	if password == dbPassword {
		return true, role, userID, nil // Les identifiants sont corrects, renvoyez également le rôle et l'ID de l'utilisateur
	}

	return false, "", 0, nil // Mauvais mot de passe
}

func updateOrderStatus(db *sql.DB, id int, status string) error {
	query := `UPDATE Orders SET status = ? WHERE id = ?`
	_, err := db.Exec(query, status, id)
	if err != nil {
		return err
	}

	// Exemple d'utilisation de la fonction sendMail
	// À adapter selon votre logique pour récupérer l'email du client associé à la commande
	clientEmail := "hikingoff@gmail.com" // Récupérer l'email réel du client ici
	subject := "Mise à jour de votre commande Go Food Court"
	body := "Votre commande a été mise à jour avec le statut : " + status

	err = sendMail(clientEmail, subject, body)
	if err != nil {
		log.Printf("Erreur lors de l'envoi de l'email : %v", err)
		// Décidez si vous voulez renvoyer l'erreur ou simplement la logger
	}

	return nil
}