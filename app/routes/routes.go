package routes

import (
	"net/http"

	"hotelm/handlers"
)

// SetupRoutes registers all the HTTP routes for the application.
func SetupRoutes() {
	// Auth routes
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.LoginPageHandler(w, r)
		} else if r.Method == http.MethodPost {
			handlers.LoginPostHandler(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/logout", handlers.LogoutHandler)

	// Customer routes
	http.HandleFunc("/customer", handlers.CustomerDashboardHandler)
	http.HandleFunc("/customer/rooms", handlers.AvailableRoomsHandler)         // List available rooms
	http.HandleFunc("/customer/booking/new", handlers.NewBookingPageHandler)     // Show booking form
	http.HandleFunc("/customer/booking", func(w http.ResponseWriter, r *http.Request) { // Create booking (POST)
		if r.Method == http.MethodPost {
			handlers.CreateBookingHandler(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/customer/bookings", handlers.MyBookingsHandler)
	http.HandleFunc("/customer/booking/delete", handlers.DeleteBookingHandler)

	// Vendor routes
	http.HandleFunc("/vendor", handlers.VendorDashboardHandler)
	http.HandleFunc("/vendor/rooms", func(w http.ResponseWriter, r *http.Request) {
		// This route is used to list rooms.
		if r.Method == http.MethodGet {
			handlers.VendorRoomsHandler(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/vendor/rooms/new", func(w http.ResponseWriter, r *http.Request) {
		// Route to display the new room form (GET) and create a new room (POST).
		if r.Method == http.MethodGet {
			handlers.NewRoomPageHandler(w, r)
		} else if r.Method == http.MethodPost {
			handlers.CreateRoomHandler(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/vendor/rooms/edit", func(w http.ResponseWriter, r *http.Request) {
		// Route to display the edit room form (GET) and update room details (POST).
		if r.Method == http.MethodGet {
			handlers.EditRoomPageHandler(w, r)
		} else if r.Method == http.MethodPost {
			handlers.UpdateRoomHandler(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/vendor/rooms/delete", func(w http.ResponseWriter, r *http.Request) {
		// Route to delete a room (POST only).
		if r.Method == http.MethodPost {
			handlers.DeleteRoomHandler(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})



	// Registration routes
http.HandleFunc("/register/customer", func(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handlers.RegistrationCustomerPageHandler(w, r)
	} else if r.Method == http.MethodPost {
		handlers.RegistrationCustomerPostHandler(w, r)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
})
http.HandleFunc("/register/vendor", func(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handlers.RegistrationVendorPageHandler(w, r)
	} else if r.Method == http.MethodPost {
		handlers.RegistrationVendorPostHandler(w, r)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
})

http.HandleFunc("/vendor/payments", func(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        handlers.VendorPaymentsHandler(w, r)
    } else {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    }
})


	
}
