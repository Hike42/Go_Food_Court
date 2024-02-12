package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
	"github.com/joho/godotenv"
)

var db *sql.DB

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Restaurant struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	UserID  int    `json:"user_id"`
}

type Menu struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Description  string  `json:"description"`
	RestaurantID int     `json:"restaurant_id"`
}

type Order struct {
	ID          int    `json:"id"`
	ClientEmail string `json:"client_email"`
	DishID      int    `json:"dish_id"`
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
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true&interpolateParams=true", username, password, host, dbName)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

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

	// Config SMTP Gmail
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
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	db = initDB()
	r := gin.Default()

	config := cors.Config{
		AllowOrigins: []string{
			"https://master--gofoodcourt.netlify.app",
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  false,
	}

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
			c.JSON(http.StatusOK, gin.H{"message": "Login successful", "role": role, "userID": userID})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
