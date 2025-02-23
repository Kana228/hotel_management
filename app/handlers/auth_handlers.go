package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"hotelm/service"
	"hotelm/session"
)

// loginTemplate is the HTML template for the login page.
// Ensure that the file "templates/login.html" exists and contains the proper form.
var loginTemplate = template.Must(template.ParseFiles("templates/login.html"))

// LoginPageHandler serves the login page.
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// Render the login page template.
	if err := loginTemplate.Execute(w, nil); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// LoginPostHandler processes the login form submission.
func LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Retrieve form values.
	role := r.FormValue("role") // Expected values: "vendor" or "customer"
	idStr := r.FormValue("id")
	name := r.FormValue("name")

	// Convert id from string to integer.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Authenticate based on role.
	if role == "customer" {
		err = service.LoginCustomer(id, name)
	} else if role == "vendor" {
		err = service.LoginVendor(id, name)
	} else {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	if err != nil {
		// On authentication failure, return an error.
		http.Error(w, "Login failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// On success, redirect to the appropriate dashboard.
	if role == "customer" {
		http.Redirect(w, r, "/customer", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/vendor", http.StatusSeeOther)
	}
}

// LogoutHandler clears the current user session and redirects to the login page.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session.ClearCurrentUser()
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
