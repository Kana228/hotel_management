package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"hotelm/models"
	"hotelm/service"
)

// Parse templates once at startup.
// Ensure these template files exist in the "templates" directory.
var (
	vendorDashboardTmpl = template.Must(template.ParseFiles("templates/vendor_dashboard.html"))
	vendorRoomsTmpl     = template.Must(template.ParseFiles("templates/vendor_rooms.html"))
	newRoomTmpl         = template.Must(template.ParseFiles("templates/new_room.html"))
	editRoomTmpl        = template.Must(template.ParseFiles("templates/edit_room.html"))
	vendorPaymentsTmpl  = template.Must(template.ParseFiles("templates/vendor_payments.html"))
)

// VendorDashboardHandler renders the vendor dashboard page.
func VendorDashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := vendorDashboardTmpl.Execute(w, nil); err != nil {
		http.Error(w, "Error rendering vendor dashboard", http.StatusInternalServerError)
	}
}

// VendorRoomsHandler displays the list of rooms for the logged-in vendor.
func VendorRoomsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	rooms, err := service.GetVendorRooms()
	if err != nil {
		http.Error(w, "Error retrieving vendor rooms: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := vendorRoomsTmpl.Execute(w, rooms); err != nil {
		http.Error(w, "Error rendering vendor rooms", http.StatusInternalServerError)
		return
	}
}

// NewRoomPageHandler renders the form to create a new room.
func NewRoomPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := newRoomTmpl.Execute(w, nil); err != nil {
		http.Error(w, "Error rendering new room page", http.StatusInternalServerError)
		return
	}
}

// CreateRoomHandler processes the form submission to create a new room.
func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Extract form values.
	name := r.FormValue("name")
	description := r.FormValue("description")
	location := r.FormValue("location")
	availabilityStr := r.FormValue("availability") // expect "true" or "false"
	priceStr := r.FormValue("price")
	roomType := r.FormValue("room_type")
	averageRatingStr := r.FormValue("average_rating")
	amenities := r.FormValue("amenities") // simple comma-separated string

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}
	avail, err := strconv.ParseBool(availabilityStr)
	if err != nil {
		avail = true // default to true if parsing fails
	}
	averageRating, err := strconv.ParseFloat(averageRatingStr, 64)
	if err != nil {
		averageRating = 0.0
	}

	room := models.Room{
		Name:          name,
		Description:   description,
		Location:      location,
		Availability:  avail,
		Price:         price,
		RoomType:      roomType,
		AverageRating: averageRating,
		Amenities:     amenities,
		// VendorID will be set in the service layer.
	}

	_, err = service.CreateRoomForVendor(room)
	if err != nil {
		http.Error(w, "Error creating room: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// On success, redirect to the vendor rooms list.
	http.Redirect(w, r, "/vendor/rooms", http.StatusSeeOther)
}

// EditRoomPageHandler renders the form to edit an existing room.
// It expects a query parameter "room_id".
func EditRoomPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
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

	// Retrieve the room details, ensuring the logged-in vendor owns it.
	room, err := service.GetRoomByIDForVendor(roomID)
	if err != nil {
		http.Error(w, "Error retrieving room: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := editRoomTmpl.Execute(w, room); err != nil {
		http.Error(w, "Error rendering edit room page", http.StatusInternalServerError)
		return
	}
}

// UpdateRoomHandler processes the form submission to update a room.
func UpdateRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	roomIDStr := r.FormValue("room_id")
	name := r.FormValue("name")
	description := r.FormValue("description")
	location := r.FormValue("location")
	availabilityStr := r.FormValue("availability")
	priceStr := r.FormValue("price")
	roomType := r.FormValue("room_type")
	averageRatingStr := r.FormValue("average_rating")
	amenities := r.FormValue("amenities")

	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}
	avail, err := strconv.ParseBool(availabilityStr)
	if err != nil {
		avail = true
	}
	averageRating, err := strconv.ParseFloat(averageRatingStr, 64)
	if err != nil {
		averageRating = 0.0
	}

	room := models.Room{
		RoomID:        roomID,
		Name:          name,
		Description:   description,
		Location:      location,
		Availability:  avail,
		Price:         price,
		RoomType:      roomType,
		AverageRating: averageRating,
		Amenities:     amenities,
		// VendorID will be set in the service layer.
	}

	err = service.UpdateRoomForVendor(room)
	if err != nil {
		http.Error(w, "Error updating room: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// On success, redirect back to the vendor rooms list.
	http.Redirect(w, r, "/vendor/rooms", http.StatusSeeOther)
}

// DeleteRoomHandler processes deletion of a room.
// Expects a POST request with a form value "room_id".
func DeleteRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	roomIDStr := r.FormValue("room_id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}
	err = service.DeleteRoomForVendor(roomID)
	if err != nil {
		http.Error(w, "Error deleting room: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Redirect back to the vendor rooms list.
	http.Redirect(w, r, "/vendor/rooms", http.StatusSeeOther)
}

// VendorPaymentsHandler displays the payments received for the vendor's rooms.
func VendorPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	payments, err := service.GetVendorPayments()
	if err != nil {
		http.Error(w, "Error retrieving vendor payments: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := vendorPaymentsTmpl.Execute(w, payments); err != nil {
		http.Error(w, "Error rendering vendor payments", http.StatusInternalServerError)
		return
	}
}
