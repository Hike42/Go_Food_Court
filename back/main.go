package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/smtp"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
)

var db *sql.DB

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"` // Assurez-vous de gérer les mots de passe de manière sécurisée
	Role     string `json:"role"`
}

type Restaurant struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	UserID  int    `json:"user_id"` // Correspond à user_id, sans contrainte de clé étrangère
}

type Menu struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Description  string  `json:"description"`
	RestaurantID int     `json:"restaurant_id"` // Correspond à restaurant_id
}

type Order struct {
	ID          int    `json:"id"`
	ClientEmail string `json:"client_email"` // Correspond à client_email
	DishID      int    `json:"dish_id"`      // Correspond à dish_id
	Quantity    int    `json:"quantity"`
	Status      string `json:"status"`
	DateTime    string `json:"date_time"`
}

type OrderStatus struct {
	OrderID  int    `json:"order_id"`
	Status   string `json:"status"`
	DateTime string `json:"date_time"`
}

func initDB() *sql.DB {
	db, err := sql.Open("mysql", "7x51ws5f6twn882m4ofl:pscale_pw_hzU7p7Pn4v6mF8DhvCwG8ZWarurr7vb7znUko4pklig@tcp(aws.connect.psdb.cloud)/gfc-db?tls=true&interpolateParams=true")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("Connected to the database successfully.")
	return db
}

func graphqlHandler(schema graphql.Schema) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params struct {
			Query         string                 `json:"query"`
			OperationName string                 `json:"operationName"`
			Variables     map[string]interface{} `json:"variables"`
		}
		if err := c.BindJSON(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  params.Query,
			VariableValues: params.Variables,
			OperationName:  params.OperationName,
		})

		if len(result.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, result.Errors)
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}

func sendMail(to, subject, body string) error {
	from := "baptiste.verd@gmail.com"
	password := "cfsr ydwm ynby wmfq"

	// Configurer le serveur SMTP de Gmail
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message
	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	// Authentification
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Envoi de l'email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		log.Printf("sendMail error: %s", err)
		return err
	}

	log.Println("Email sent successfully")
	return nil
}


func main() {
	db = initDB()
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowMethods("OPTIONS")
	r.Use(cors.New(config))

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    RootQuery,
		Mutation: RootMutation,
	})
	if err != nil {
		panic(err)
	}

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Bienvenue sur mon serveur Go!")
	})
	r.POST("/graphql", graphqlHandler(schema))
	r.POST("/api/signup", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newUser, err := createUser(db, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": newUser})
	})

	r.POST("/api/login", func(c *gin.Context) {
		var credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		isValid, role, userID, err := checkUserCredentials(db, credentials.Email, credentials.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if isValid {
			// Utilisez l'ID de l'utilisateur (userID) ici pour récupérer des informations supplémentaires si nécessaire
			c.JSON(http.StatusOK, gin.H{"message": "Login successful", "role": role, "userID": userID})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		}
	})

	r.Run(":8080")
}
