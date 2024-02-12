package main

import "github.com/graphql-go/graphql"

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id":       &graphql.Field{Type: graphql.Int},
		"name":     &graphql.Field{Type: graphql.String},
		"email":    &graphql.Field{Type: graphql.String},
		"password": &graphql.Field{Type: graphql.String},
		"role":     &graphql.Field{Type: graphql.String},
	},
})

var restaurantType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Restaurant",
	Fields: graphql.Fields{
		"id":      &graphql.Field{Type: graphql.Int},
		"name":    &graphql.Field{Type: graphql.String},
		"address": &graphql.Field{Type: graphql.String},
		"user_id": &graphql.Field{Type: graphql.Int},
	},
})

var menuType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Menu",
	Fields: graphql.Fields{
		"id":            &graphql.Field{Type: graphql.Int},
		"name":          &graphql.Field{Type: graphql.String},
		"price":         &graphql.Field{Type: graphql.Float},
		"description":   &graphql.Field{Type: graphql.String},
		"restaurant_id": &graphql.Field{Type: graphql.Int},
	},
})

var orderType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Order",
	Fields: graphql.Fields{
		"id":           &graphql.Field{Type: graphql.Int},
		"client_email": &graphql.Field{Type: graphql.String},
		"dish_id":      &graphql.Field{Type: graphql.Int},
		"quantity":     &graphql.Field{Type: graphql.Int},
		"status":       &graphql.Field{Type: graphql.String},
		"date_time":    &graphql.Field{Type: graphql.String},
	},
})

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{

		"getUserByID": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(int)
				db := initDB()
				defer db.Close()
				return findUserByID(db, id)
			},
		},

		"getAllUsers": &graphql.Field{
			Type: graphql.NewList(userType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				db := initDB()
				defer db.Close()
				return findAllUsers(db)
			},
		},

		"getRestaurantByID": &graphql.Field{
			Type: restaurantType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(int)
				db := initDB()
				defer db.Close()
				return findRestaurantByID(db, id)
			},
		},
		"getAllRestaurants": &graphql.Field{
			Type: graphql.NewList(restaurantType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				db := initDB()
				defer db.Close()
				return findAllRestaurants(db)
			},
		},
		"getMenuByID": &graphql.Field{
			Type: menuType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(int)
				db := initDB()
				defer db.Close()
				return findMenuByID(db, id)
			},
		},
		"getMenusByRestaurant": &graphql.Field{
			Type: graphql.NewList(menuType),
			Args: graphql.FieldConfigArgument{
				"restaurant_id": &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				restaurantID, _ := p.Args["restaurant_id"].(int)
				db := initDB()
				defer db.Close()
				return findMenusByRestaurant(db, restaurantID)
			},
		},
		"getAllMenus": &graphql.Field{
			Type: graphql.NewList(menuType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				db := initDB()
				defer db.Close()
				return findAllMenus(db)
			},
		},
		"getOrderByID": &graphql.Field{
			Type: orderType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(int)
				db := initDB()
				defer db.Close()
				return findOrderByID(db, id)
			},
		},
		"getALlOrdersByEmail": &graphql.Field{
			Type: graphql.NewList(orderType),
			Args: graphql.FieldConfigArgument{
				"client_email": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				email, _ := p.Args["client_email"].(string)
				db := initDB()
				defer db.Close()
				return findAllOrdersbyUserEmail(db, email)
			},
		},
		"getOrdersbyRestaurant": &graphql.Field{
			Type: graphql.NewList(orderType),
			Args: graphql.FieldConfigArgument{
				"restaurant_id": &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				restaurantID, _ := p.Args["restaurant_id"].(int)
				db := initDB()
				defer db.Close()
				return findOrdersByRestaurantID(db, restaurantID)
			},
		},
	},
})

var RootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"createUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"name":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"role":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Extraire les paramètres de la requête GraphQL
				name := p.Args["name"].(string)
				email := p.Args["email"].(string)
				password := p.Args["password"].(string)
				role := p.Args["role"].(string)

				// Appeler la fonction createUser avec les paramètres extraits
				db := initDB()
				defer db.Close()
				return createUser(db, User{Name: name, Email: email, Password: password, Role: role})
			},
		},

		"updateUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"name":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"role":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				user := User{
					ID:       p.Args["id"].(int),
					Name:     p.Args["name"].(string),
					Email:    p.Args["email"].(string),
					Password: p.Args["password"].(string),
					Role:     p.Args["role"].(string),
				}
				db := initDB()
				defer db.Close()
				return updateUser(db, user)
			},
		},
		"deleteUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(int)
				db := initDB()
				defer db.Close()
				return deleteUser(db, id), nil
			},
		},
		"createRestaurant": &graphql.Field{
			Type: restaurantType,
			Args: graphql.FieldConfigArgument{
				"name":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"address": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"user_id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				restaurant := Restaurant{
					Name:    p.Args["name"].(string),
					Address: p.Args["address"].(string),
					UserID:  p.Args["user_id"].(int),
				}
				db := initDB()
				defer db.Close()
				return createRestaurant(db, restaurant)
			},
		},
		"updateRestaurant": &graphql.Field{
			Type: restaurantType,
			Args: graphql.FieldConfigArgument{
				"id":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"name":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"address": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"user_id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				restaurant := Restaurant{
					ID:      p.Args["id"].(int),
					Name:    p.Args["name"].(string),
					Address: p.Args["address"].(string),
					UserID:  p.Args["user_id"].(int),
				}
				db := initDB()
				defer db.Close()
				return updateRestaurant(db, restaurant)
			},
		},
		"deleteRestaurant": &graphql.Field{
			Type: restaurantType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(int)
				db := initDB()
				defer db.Close()
				return deleteRestaurant(db, id), nil
			},
		},
		"createMenu": &graphql.Field{
			Type: menuType,
			Args: graphql.FieldConfigArgument{
				"name":          &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"price":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
				"description":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"restaurant_id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				menu := Menu{
					Name:         p.Args["name"].(string),
					Price:        p.Args["price"].(float64),
					Description:  p.Args["description"].(string),
					RestaurantID: p.Args["restaurant_id"].(int),
				}
				db := initDB()
				defer db.Close()
				return createMenu(db, menu)
			},
		},
		"updateMenu": &graphql.Field{
			Type: menuType,
			Args: graphql.FieldConfigArgument{
				"id":            &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"name":          &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"price":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
				"description":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"restaurant_id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				menu := Menu{
					ID:           p.Args["id"].(int),
					Name:         p.Args["name"].(string),
					Price:        p.Args["price"].(float64),
					Description:  p.Args["description"].(string),
					RestaurantID: p.Args["restaurant_id"].(int),
				}
				db := initDB()
				defer db.Close()
				return updateMenu(db, menu)
			},
		},
		"deleteMenu": &graphql.Field{
			Type: menuType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(int)
				db := initDB()
				defer db.Close()
				return deleteMenu(db, id), nil
			},
		},
		"createOrder": &graphql.Field{
			Type: orderType,
			Args: graphql.FieldConfigArgument{
				"client_email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"dish_id":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"quantity":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"status":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"date_time":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				order := Order{
					ClientEmail: p.Args["client_email"].(string),
					DishID:      p.Args["dish_id"].(int),
					Quantity:    p.Args["quantity"].(int),
					Status:      p.Args["status"].(string),
					DateTime:    p.Args["date_time"].(string),
				}
				db := initDB()
				defer db.Close()
				return createOrder(db, order)
			},
		},
		"updateOrder": &graphql.Field{
			Type: orderType,
			Args: graphql.FieldConfigArgument{
				"id":           &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"client_email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"dish_id":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"quantity":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"status":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"date_time":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				order := Order{
					ID:          p.Args["id"].(int),
					ClientEmail: p.Args["client_email"].(string),
					DishID:      p.Args["dish_id"].(int),
					Quantity:    p.Args["quantity"].(int),
					Status:      p.Args["status"].(string),
					DateTime:    p.Args["date_time"].(string),
				}
				db := initDB()
				defer db.Close()
				return updateOrder(db, order)
			},
		},
		"updateOrderStatus": &graphql.Field{
			Type: orderType,
			Args: graphql.FieldConfigArgument{
				"id":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"status": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				orderStatus := Order{
					ID: p.Args["id"].(int),
					Status:  p.Args["status"].(string),
				}
				db := initDB()
				defer db.Close()
				err := updateOrderStatus(db, orderStatus.ID, orderStatus.Status)
				if err != nil {
					return nil, err
				}
				return orderStatus, nil
			},
		},
		"deleteOrder": &graphql.Field{
			Type: orderType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(int)
				db := initDB()
				defer db.Close()
				return deleteOrder(db, id), nil
			},
		},
	},
})
