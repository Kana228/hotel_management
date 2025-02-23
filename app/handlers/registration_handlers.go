package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"hotelm/models"
	"hotelm/repository"
)

// Parse templates for registration.
var (
	customerRegTmpl  = template.Must(template.ParseFiles("templates/registration_customer.html"))
	vendorRegTmpl    = template.Must(template.ParseFiles("templates/registration_vendor.html"))
	regSuccessTmpl   = template.Must(template.ParseFiles("templates/registration_success.html"))
)

// RegistrationCustomerPageHandler renders the customer registration form.
func RegistrationCustomerPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// Render customer registration template without any error initially.
	customerRegTmpl.Execute(w, nil)
}

// RegistrationCustomerPostHandler processes the customer registration form.
func RegistrationCustomerPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Basic field validation.
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	email := r.FormValue("email")
	address := r.FormValue("address")

	if name == "" || phone == "" || email == "" || address == "" {
		customerRegTmpl.Execute(w, map[string]string{"Error": "All fields are required."})
		return
	}

	// Create the customer using the repository.
	customerID, err := repository.CreateCustomer(models.Customer{
		Name:    name,
		Phone:   phone,
		Email:   email,
		Address: address,
	})
	if err != nil {
		// Friendly message for duplicate email or other errors.
		customerRegTmpl.Execute(w, map[string]string{"Error": "Registration failed: " + err.Error()})
		return
	}

	// Render success page with the new user ID.
	regSuccessTmpl.Execute(w, map[string]string{"UserID": strconv.Itoa(customerID)})
}

// RegistrationVendorPageHandler renders the vendor registration form.
func RegistrationVendorPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	vendorRegTmpl.Execute(w, nil)
}

// RegistrationVendorPostHandler processes the vendor registration form.
func RegistrationVendorPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Basic field validation.
	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	hotelName := r.FormValue("hotel_name")
	address := r.FormValue("address")

	if name == "" || email == "" || phone == "" || hotelName == "" || address == "" {
		vendorRegTmpl.Execute(w, map[string]string{"Error": "All fields are required."})
		return
	}

	// Create the vendor using the repository.
	vendorID, err := repository.CreateVendor(models.Vendor{
		Name:      name,
		Email:     email,
		Phone:     phone,
		HotelName: hotelName,
		Address:   address,
	})
	if err != nil {
		vendorRegTmpl.Execute(w, map[string]string{"Error": "Registration failed: " + err.Error()})
		return
	}

	// Render success page with the new vendor ID.
	regSuccessTmpl.Execute(w, map[string]string{"UserID": strconv.Itoa(vendorID)})
}
