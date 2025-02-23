package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"hotelm/db"
	"hotelm/models"
)

// CreatePayment inserts a new payment into the database
func CreatePayment(payment models.Payment) (int, error) {
	query := `INSERT INTO payment (payment_method, payment_status, transaction_date, amount, booking_id) 
		VALUES ($1, $2, $3, $4, $5) RETURNING payment_id`
	var id int
	err := db.DB.QueryRow(query, payment.PaymentMethod, payment.PaymentStatus, payment.TransactionDate, payment.Amount, payment.BookingID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create payment: %v", err)
	}
	return id, nil
}

// GetPaymentByID retrieves a payment by ID
func GetPaymentByID(paymentID int) (*models.Payment, error) {
	query := `SELECT payment_id, payment_method, payment_status, transaction_date, amount, booking_id FROM payment WHERE payment_id = $1`
	var payment models.Payment

	err := db.DB.QueryRow(query, paymentID).Scan(&payment.PaymentID, &payment.PaymentMethod, &payment.PaymentStatus, &payment.TransactionDate, &payment.Amount, &payment.BookingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("error retrieving payment: %v", err)
	}
	return &payment, nil
}

// UpdatePayment updates an existing payment
func UpdatePayment(payment models.Payment) error {
	query := `UPDATE payment SET payment_method = $1, payment_status = $2, transaction_date = $3, amount = $4, booking_id = $5 WHERE payment_id = $6`
	result, err := db.DB.Exec(query, payment.PaymentMethod, payment.PaymentStatus, payment.TransactionDate, payment.Amount, payment.BookingID, payment.PaymentID)
	if err != nil {
		return fmt.Errorf("failed to update payment: %v", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("payment not found")
	}
	return nil
}

// DeletePayment removes a payment by ID
func DeletePayment(paymentID int) error {
	query := `DELETE FROM payment WHERE payment_id = $1`
	result, err := db.DB.Exec(query, paymentID)
	if err != nil {
		return fmt.Errorf("failed to delete payment: %v", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("payment not found")
	}
	return nil
}

// GetPaymentsByBookingID retrieves all payments associated with a specific booking
func GetPaymentsByBookingID(bookingID int) ([]models.Payment, error) {
	query := `SELECT payment_id, payment_method, payment_status, transaction_date, amount, booking_id FROM payment WHERE booking_id = $1`
	rows, err := db.DB.Query(query, bookingID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve payments: %v", err)
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		if err := rows.Scan(&payment.PaymentID, &payment.PaymentMethod, &payment.PaymentStatus, &payment.TransactionDate, &payment.Amount, &payment.BookingID); err != nil {
			return nil, fmt.Errorf("error scanning payment: %v", err)
		}
		payments = append(payments, payment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading payments: %v", err)
	}
	return payments, nil
}
