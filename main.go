package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", LoginPage)
	http.HandleFunc("/signup", SignupPage)
	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/welcome", WelcomePage)
	http.HandleFunc("/logout", LogOut)

	// Serve static files from the "static" directory.
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	template.ParseGlob("/templates/*.html")
	// Start the server on port 8080.
	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func SignupPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Retrieve signup form data.
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Perform signup logic here (e.g., store user data in a database).
		// For simplicity, we'll just print the data for demonstration.
		fmt.Printf("New user signup: Username - %s, Password - %s\n", username, password)

		// Redirect to a welcome or login page after signup.
		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
		return
	}

	// If not a POST request, serve the signup page template.
	tmpl, err := template.ParseFiles("templates/signup.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// LoginPage is the handler for the login page.
func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Perform authentication logic here (e.g., check against a database).
		// For simplicity, we'll just check if the username and password are both "admin".
		if username == "admin" && password == "admin" {
			// Create a new cookie.
			cookie := http.Cookie{
				Name:     "user",
				Value:    username,
				Expires:  time.Now().Add(24 * time.Hour), // Cookie expires in 24 hours.
				HttpOnly: true,
			}

			// Set the cookie in the response.
			http.SetCookie(w, &cookie)

			// Successful login, redirect to a welcome page.
			http.Redirect(w, r, "/welcome", http.StatusSeeOther)
			return
		}

		// Invalid credentials, show the login page with an error message.
		fmt.Fprintf(w, "Invalid credentials. Please try again.")
		return
	}

	// If not a POST request, serve the login page template.
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// WelcomePage is the handler for the welcome page.
func WelcomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, "Welcome, you have successfully logged in!")
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:    "user",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour), // Expire in the past to delete the cookie.
	}
	http.SetCookie(w, &cookie)

	// Redirect to the login page after logging out.
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
