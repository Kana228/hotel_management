package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
	"hotelm/repository"
	"hotelm/models"
	"hotelm/service"
)

// Parse templates once at startup.
var (
	customerDashboardTmpl = template.Must(template.ParseFiles("templates/customer_dashboard.html"))
	availableRoomsTmpl    = template.Must(template.ParseFiles("templates/available_rooms.html"))
	bookingFormTmpl       = template.Must(template.ParseFiles("templates/booking_form.html"))
	myBookingsTmpl        = template.Must(template.ParseFiles("templates/my_bookings.html"))
)

// CustomerDashboardHandler renders the customer dashboard with options.
func CustomerDashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Render the customer dashboard page.
	if err := customerDashboardTmpl.Execute(w, nil); err != nil {
		http.Error(w, "Error rendering customer dashboard", http.StatusInternalServerError)
		return
	}
}

// AvailableRoomsHandler retrieves available rooms and renders the list.
func AvailableRoomsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	rooms, err := service.GetAvailableRooms()
	if err != nil {
		http.Error(w, "Error retrieving available rooms: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the available rooms template, passing the rooms slice.
	if err := availableRoomsTmpl.Execute(w, rooms); err != nil {
		http.Error(w, "Error rendering available rooms", http.StatusInternalServerError)
		return
	}
}

// NewBookingPageHandler renders the form for creating a new booking.
// The room ID is expected as a query parameter "room_id".
func NewBookingPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get room_id from query parameters.
	roomIDStr := r.URL.Query().Get("room_id")
	if roomIDStr == "" {
		http.Error(w, "Missing room_id parameter", http.StatusBadRequest)
		return
	}
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		http.Error(w, "Invalid room_id", http.StatusBadRequest)
		return
	}

	// Pass the room ID to the booking form template.
	data := struct {
		RoomID int
	}{
		RoomID: roomID,
	}

	if err := bookingFormTmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering booking form", http.StatusInternalServerError)
		return
	}
}

func CreateBookingHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse form data.
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        return
    }

    // Get room ID, booking dates, and payment method from form.
    roomIDStr := r.FormValue("room_id")
    checkinStr := r.FormValue("checkin_date")
    checkoutStr := r.FormValue("checkout_date")
    paymentMethod := r.FormValue("payment_method")

    roomID, err := strconv.Atoi(roomIDStr)
    if err != nil {
        http.Error(w, "Invalid room ID", http.StatusBadRequest)
        return
    }

    // Expect dates in YYYY-MM-DD format.
    checkinDate, err := time.Parse("2006-01-02", checkinStr)
    if err != nil {
        http.Error(w, "Invalid check-in date", http.StatusBadRequest)
        return
    }
    checkoutDate, err := time.Parse("2006-01-02", checkoutStr)
    if err != nil {
        http.Error(w, "Invalid check-out date", http.StatusBadRequest)
        return
    }

    booking := models.Booking{
        BookingDate:   time.Now(),
        CheckinDate:   checkinDate,
        CheckoutDate:  checkoutDate,
        PaymentStatus: "Pending", // Will be updated when payment is created.
        RoomID:        roomID,
        // CustomerID will be set in the service layer.
    }

    // Create the booking using the service layer and capture the bookingID.
    bookingID, err := service.CreateBookingForCustomer(booking)
    if err != nil {
        http.Error(w, "Error creating booking: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Retrieve the room details to obtain the price.
    room, err := repository.GetRoomByID(roomID)
    if err != nil {
        http.Error(w, "Error retrieving room details: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Create a payment record with the specified payment method and room's price.
    payment := models.Payment{
        PaymentMethod:   paymentMethod,
        PaymentStatus:   "Completed",
        TransactionDate: time.Now(),
        Amount:          room.Price,
        BookingID:       bookingID, // Use the returned bookingID.
    }

    _, err = repository.CreatePayment(payment)
    if err != nil {
        http.Error(w, "Error creating payment: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // On success, redirect to My Bookings page.
    http.Redirect(w, r, "/customer/bookings", http.StatusSeeOther)
}


// MyBookingsHandler retrieves the logged-in customer's bookings and renders them.
func MyBookingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	bookings, err := service.GetMyBookings()
	if err != nil {
		http.Error(w, "Error retrieving bookings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the my bookings template, passing the bookings slice.
	if err := myBookingsTmpl.Execute(w, bookings); err != nil {
		http.Error(w, "Error rendering my bookings", http.StatusInternalServerError)
		return
	}
}

// DeleteBookingHandler processes the deletion of a booking.
// The booking ID should be passed as a form value.
func DeleteBookingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	bookingIDStr := r.FormValue("booking_id")
	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	if err := service.DeleteBookingForCustomer(bookingID); err != nil {
		http.Error(w, "Error deleting booking: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// On success, redirect back to the My Bookings page.
	http.Redirect(w, r, "/customer/bookings", http.StatusSeeOther)
}
