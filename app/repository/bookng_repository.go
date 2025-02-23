package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"hotelm/db"
	"hotelm/models"
)

// CreateBooking inserts a new booking into the database
func CreateBooking(booking models.Booking) (int, error) {
	query := `INSERT INTO booking (booking_date, checkin_date, checkout_date, payment_status, room_id, customer_id) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING booking_id`
	var id int
	err := db.DB.QueryRow(query, booking.BookingDate, booking.CheckinDate, booking.CheckoutDate, booking.PaymentStatus, booking.RoomID, booking.CustomerID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create booking: %v", err)
	}
	return id, nil
}

// GetBookingByID retrieves a booking by ID
func GetBookingByID(bookingID int) (*models.Booking, error) {
	query := `SELECT booking_id, booking_date, checkin_date, checkout_date, payment_status, room_id, customer_id FROM booking WHERE booking_id = $1`
	var booking models.Booking

	err := db.DB.QueryRow(query, bookingID).Scan(&booking.BookingID, &booking.BookingDate, &booking.CheckinDate, &booking.CheckoutDate, &booking.PaymentStatus, &booking.RoomID, &booking.CustomerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, fmt.Errorf("error retrieving booking: %v", err)
	}
	return &booking, nil
}

// UpdateBooking updates an existing booking
func UpdateBooking(booking models.Booking) error {
	query := `UPDATE booking SET booking_date = $1, checkin_date = $2, checkout_date = $3, payment_status = $4, room_id = $5, customer_id = $6 WHERE booking_id = $7`
	result, err := db.DB.Exec(query, booking.BookingDate, booking.CheckinDate, booking.CheckoutDate, booking.PaymentStatus, booking.RoomID, booking.CustomerID, booking.BookingID)
	if err != nil {
		return fmt.Errorf("failed to update booking: %v", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("booking not found")
	}
	return nil
}

// DeleteBooking removes a booking by ID
func DeleteBooking(bookingID int) error {
	query := `DELETE FROM booking WHERE booking_id = $1`
	result, err := db.DB.Exec(query, bookingID)
	if err != nil {
		return fmt.Errorf("failed to delete booking: %v", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("booking not found")
	}
	return nil
}

// GetBookingsByCustomerID retrieves all bookings made by a specific customer
func GetBookingsByCustomerID(customerID int) ([]models.Booking, error) {
	query := `SELECT booking_id, booking_date, checkin_date, checkout_date, payment_status, room_id, customer_id FROM booking WHERE customer_id = $1`
	rows, err := db.DB.Query(query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve bookings: %v", err)
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var booking models.Booking
		if err := rows.Scan(&booking.BookingID, &booking.BookingDate, &booking.CheckinDate, &booking.CheckoutDate, &booking.PaymentStatus, &booking.RoomID, &booking.CustomerID); err != nil {
			return nil, fmt.Errorf("error scanning booking: %v", err)
		}
		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading bookings: %v", err)
	}
	return bookings, nil
}