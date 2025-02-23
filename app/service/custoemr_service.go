package service

import (
	"fmt"

	"hotelm/models"
	"hotelm/repository"
	"hotelm/session"
)

// GetAvailableRooms returns a global list of available rooms.
func GetAvailableRooms() ([]models.Room, error) {
	rooms, err := repository.GetAvailableRooms()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve available rooms: %v", err)
	}
	return rooms, nil
}

// CreateBookingForCustomer creates a new booking for the logged-in customer.
// It sets the Booking.CustomerID to the current customer's ID.
func CreateBookingForCustomer(booking models.Booking) (int, error) {
	// Ensure a customer is logged in.
	user := session.GetCurrentUser()
	customer, ok := user.(*models.Customer)
	if !ok {
		return 0, fmt.Errorf("no customer is currently logged in")
	}
	booking.CustomerID = customer.CustomerID

	// Retrieve the room details.
	room, err := repository.GetRoomByID(booking.RoomID)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve room: %v", err)
	}
	// Check if the room is available.
	if !room.Availability {
		return 0, fmt.Errorf("room is not available")
	}

	// Create the booking.
	bookingID, err := repository.CreateBooking(booking)
	if err != nil {
		return 0, fmt.Errorf("failed to create booking: %v", err)
	}

	// Update room availability to false.
	room.Availability = false
	if err := repository.UpdateRoom(*room); err != nil {
		return 0, fmt.Errorf("booking created but failed to update room availability: %v", err)
	}

	return bookingID, nil
}

// GetMyBookings retrieves all bookings for the logged-in customer.
func GetMyBookings() ([]models.Booking, error) {
	// Ensure a customer is logged in.
	user := session.GetCurrentUser()
	customer, ok := user.(*models.Customer)
	if !ok {
		return nil, fmt.Errorf("no customer is currently logged in")
	}

	bookings, err := repository.GetBookingsByCustomerID(customer.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve customer bookings: %v", err)
	}
	return bookings, nil
}

// DeleteBookingForCustomer deletes a booking if it belongs to the logged-in customer.
func DeleteBookingForCustomer(bookingID int) error {
	// Ensure a customer is logged in.
	user := session.GetCurrentUser()
	customer, ok := user.(*models.Customer)
	if !ok {
		return fmt.Errorf("no customer is currently logged in")
	}

	// Retrieve the customer's bookings to verify ownership.
	bookings, err := repository.GetBookingsByCustomerID(customer.CustomerID)
	if err != nil {
		return fmt.Errorf("failed to retrieve customer bookings: %v", err)
	}

	var bookingFound *models.Booking
	for _, b := range bookings {
		if b.BookingID == bookingID {
			bookingFound = &b
			break
		}
	}
	if bookingFound == nil {
		return fmt.Errorf("unauthorized: booking does not belong to the logged-in customer")
	}

	// Delete the booking.
	if err := repository.DeleteBooking(bookingID); err != nil {
		return fmt.Errorf("failed to delete booking: %v", err)
	}

	// Retrieve the room associated with the booking.
	room, err := repository.GetRoomByID(bookingFound.RoomID)
	if err != nil {
		return fmt.Errorf("failed to retrieve room: %v", err)
	}

	// Mark the room as available.
	room.Availability = true
	if err := repository.UpdateRoom(*room); err != nil {
		return fmt.Errorf("booking deleted but failed to update room availability: %v", err)
	}

	return nil
}