package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("my_secret_key"))

type User struct {
	UserID    int    `json:"user_id"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	Privilege string `json:"privilege"`
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func sayhello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Привет!")
}

func main() {
	connStr := "user=postgres dbname=candy_drop sslmode=disable password=1710 port=5555 host=127.0.0.1"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Настройка CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:5500"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Set-Cookie"},
		AllowCredentials: true,
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/", sayhello) // Устанавливаем роутер
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received register request")
		registerHandler(w, r, db)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received login request")
		loginHandler(w, r, db)
	})

	mux.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received protected request")
		protectedHandler(w, r)
	})

	handler := c.Handler(mux)
	err = http.ListenAndServe(":8080", handler) // устанавливаем порт веб-сервера

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error generating password hash:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO \"User\" (Login, Password, Privilege) VALUES ($1, $2, $3)", user.Login, string(hashedPassword), user.Privilege)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func loginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Attempting to login user with login: %s", creds.Login)

	var user User
	err = db.QueryRow("SELECT userid, password FROM \"User\" WHERE Login = $1", creds.Login).Scan(&user.UserID, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No user found with the provided login")
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			log.Println("Error querying database:", err)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
		return
	}

	log.Printf("Found user with ID: %d", user.UserID)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		log.Println("Error comparing passwords:", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	session, _ := store.Get(r, "session-name")
	session.Values["user_id"] = user.UserID
	session.Save(r, w)

	log.Println("Login successful, session set")
	w.Write([]byte("Login successful"))
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		log.Println("No session found")
		http.Error(w, "No session found", http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("Protected content for user ID: %d", userID)))
}
