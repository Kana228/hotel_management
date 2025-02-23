package service

import (
	"fmt"

	"hotelm/db"
	"hotelm/models"
	"hotelm/repository"
	"hotelm/session"
)

// GetVendorRooms retrieves all rooms belonging to the currently logged-in vendor.
func GetVendorRooms() ([]models.Room, error) {
	// Ensure we have a vendor logged in.
	user := session.GetCurrentUser()
	vendor, ok := user.(*models.Vendor)
	if !ok {
		return nil, fmt.Errorf("no vendor is currently logged in")
	}

	// Query rooms where vendor_id matches the current vendor.
	query := `
		SELECT room_id, name, description, location, availability, price, room_type, average_rating, amenities, vendor_id 
		FROM room
		WHERE vendor_id = $1
	`
	rows, err := db.DB.Query(query, vendor.VendorID)
	if err != nil {
		return nil, fmt.Errorf("error querying vendor rooms: %v", err)
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		// Since amenities is now stored as a simple TEXT field (per our earlier decision),
		// we can scan it directly into room.Amenities.
		if err := rows.Scan(
			&room.RoomID,
			&room.Name,
			&room.Description,
			&room.Location,
			&room.Availability,
			&room.Price,
			&room.RoomType,
			&room.AverageRating,
			&room.Amenities,
			&room.VendorID,
		); err != nil {
			return nil, fmt.Errorf("error scanning room: %v", err)
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading vendor rooms: %v", err)
	}

	return rooms, nil
}

// CreateRoomForVendor creates a new room for the logged-in vendor.
// It sets the VendorID in the room to that of the logged-in vendor.
func CreateRoomForVendor(room models.Room) (int, error) {
	// Ensure we have a vendor logged in.
	user := session.GetCurrentUser()
	vendor, ok := user.(*models.Vendor)
	if !ok {
		return 0, fmt.Errorf("no vendor is currently logged in")
	}

	// Set the room's VendorID to the current vendor.
	room.VendorID = vendor.VendorID

	// Call repository function to create the room.
	id, err := repository.CreateRoom(room)
	if err != nil {
		return 0, fmt.Errorf("failed to create room: %v", err)
	}
	return id, nil
}

// UpdateRoomForVendor updates an existing room if it belongs to the logged-in vendor.
func UpdateRoomForVendor(room models.Room) error {
	// Ensure we have a vendor logged in.
	user := session.GetCurrentUser()
	vendor, ok := user.(*models.Vendor)
	if !ok {
		return fmt.Errorf("no vendor is currently logged in")
	}

	// Retrieve the current room to verify ownership.
	existingRoom, err := repository.GetRoomByID(room.RoomID)
	if err != nil {
		return fmt.Errorf("failed to retrieve room: %v", err)
	}
	if existingRoom.VendorID != vendor.VendorID {
		return fmt.Errorf("unauthorized: this room does not belong to the logged-in vendor")
	}

	// Ensure that the VendorID is correct.
	room.VendorID = vendor.VendorID

	// Call repository function to update the room.
	if err := repository.UpdateRoom(room); err != nil {
		return fmt.Errorf("failed to update room: %v", err)
	}
	return nil
}

// DeleteRoomForVendor deletes a room if it belongs to the logged-in vendor.
func DeleteRoomForVendor(roomID int) error {
	// Ensure we have a vendor logged in.
	user := session.GetCurrentUser()
	vendor, ok := user.(*models.Vendor)
	if !ok {
		return fmt.Errorf("no vendor is currently logged in")
	}

	// Retrieve the room to verify ownership.
	room, err := repository.GetRoomByID(roomID)
	if err != nil {
		return fmt.Errorf("failed to retrieve room: %v", err)
	}
	if room.VendorID != vendor.VendorID {
		return fmt.Errorf("unauthorized: this room does not belong to the logged-in vendor")
	}

	// Call repository function to delete the room.
	if err := repository.DeleteRoom(roomID); err != nil {
		return fmt.Errorf("failed to delete room: %v", err)
	}
	return nil
}

// GetVendorPayments retrieves all payments for bookings on rooms belonging to the logged-in vendor.
func GetVendorPayments() ([]models.Payment, error) {
	// Ensure we have a vendor logged in.
	user := session.GetCurrentUser()
	vendor, ok := user.(*models.Vendor)
	if !ok {
		return nil, fmt.Errorf("no vendor is currently logged in")
	}

	// This query joins payments, bookings, and rooms to retrieve payments for the vendor's rooms.
	query := `
		SELECT p.payment_id, p.payment_method, p.payment_status, p.transaction_date, p.amount, p.booking_id
		FROM payment p
		JOIN booking b ON p.booking_id = b.booking_id
		JOIN room r ON b.room_id = r.room_id
		WHERE r.vendor_id = $1
	`
	rows, err := db.DB.Query(query, vendor.VendorID)
	if err != nil {
		return nil, fmt.Errorf("error querying vendor payments: %v", err)
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		if err := rows.Scan(
			&payment.PaymentID,
			&payment.PaymentMethod,
			&payment.PaymentStatus,
			&payment.TransactionDate,
			&payment.Amount,
			&payment.BookingID,
		); err != nil {
			return nil, fmt.Errorf("error scanning payment: %v", err)
		}
		payments = append(payments, payment)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading vendor payments: %v", err)
	}

	return payments, nil
}





// GetRoomByIDForVendor retrieves a room by its ID and checks that it belongs to the logged-in vendor.
func GetRoomByIDForVendor(roomID int) (*models.Room, error) {
    room, err := repository.GetRoomByID(roomID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve room: %v", err)
    }
    // Retrieve the logged-in vendor from session
    user := session.GetCurrentUser()
    vendor, ok := user.(*models.Vendor)
    if !ok {
        return nil, fmt.Errorf("no vendor is currently logged in")
    }
    if room.VendorID != vendor.VendorID {
        return nil, fmt.Errorf("unauthorized: this room does not belong to the logged-in vendor")
    }
    return room, nil
}