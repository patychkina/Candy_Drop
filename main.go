package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("ghfisg1hdifyerhvsu7dfueo"))

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
		AllowedOrigins:   []string{"http://127.0.0.1:5500"}, // разрешить доступ с frontend
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Set-Cookie"},
		AllowCredentials: true,
	})

	mux := http.NewServeMux()
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
		protectedHandler(w, r, db)
	})

	handler := c.Handler(mux)
	log.Print("Success")
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
		http.Error(w, `{"message": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	log.Printf("Attempting to login user with login: %s", creds.Login)

	var user User
	err = db.QueryRow("SELECT userid, password, privilege FROM \"User\" WHERE Login = $1", creds.Login).
		Scan(&user.UserID, &user.Password, &user.Privilege)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No user found with the provided login")
			http.Error(w, `{"message": "Такого пользователя не существует"}`, http.StatusUnauthorized)
		} else {
			log.Println("Error querying database:", err)
			http.Error(w, `{"message": "Database error"}`, http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Found user with ID: %d", user.UserID)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		log.Println("Error comparing passwords:", err)
		http.Error(w, `{"message": "Пароль неверный"}`, http.StatusUnauthorized)
		return
	}

	session, _ := store.Get(r, "session-name")
	session.Values["user_id"] = user.UserID
	session.Values["privilege"] = user.Privilege
	err = session.Save(r, w)
	if err != nil {
		log.Println("Error saving session:", err)
		http.Error(w, `{"message": "Failed to save session"}`, http.StatusInternalServerError)
		return
	}

	log.Println("Login successful, session set")

	// Устанавливаем cookie для сессии с правильными аттрибутами
	http.SetCookie(w, &http.Cookie{
		Name:     "session-name",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,                 // Защищает от XSS атак
		Secure:   false,                // Поставьте true, если используете HTTPS
		SameSite: http.SameSiteLaxMode, // Может быть Strict или None для кросс-доменных запросов
	})

	w.Header().Set("Content-Type", "application/json")
	// Возвращаем и UserID, и Privilege
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":   user.UserID,
		"privilege": user.Privilege,
	})
}

func protectedHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Получаем сессию пользователя
	session, _ := store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		log.Println("No session found")
		http.Error(w, `{"message": "No session found"}`, http.StatusUnauthorized)
		return
	}

	// Получаем привилегию пользователя из базы данных
	var privilege string
	err := db.QueryRow("SELECT privilege FROM \"User\" WHERE userid = $1", userID).Scan(&privilege)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No user found with the provided user ID")
			http.Error(w, `{"message": "User not found"}`, http.StatusUnauthorized)
		} else {
			log.Println("Error querying database:", err)
			http.Error(w, `{"message": "Database error"}`, http.StatusInternalServerError)
		}
		return
	}

	// Возвращаем данные в формате JSON
	response := map[string]interface{}{
		"user_id":   userID,
		"privilege": privilege,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
