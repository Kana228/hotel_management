package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"hotelm/db"
	"hotelm/models"
	
)

// CreateCustomer inserts a new customer into the database
func CreateCustomer(customer models.Customer) (int, error) {
	query := `INSERT INTO customer (name, phone, email, address) VALUES ($1, $2, $3, $4) RETURNING customer_id`
	var id int
	err := db.DB.QueryRow(query, customer.Name, customer.Phone, customer.Email, customer.Address).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create customer: %v", err)
	}
	return id, nil
}

// GetCustomerByID retrieves a customer by id
func GetCustomerByID(customerID int) (*models.Customer, error) {
	query := `SELECT customer_id, name, phone, email, address FROM customer WHERE customer_id = $1`
	var customer models.Customer

	err := db.DB.QueryRow(query, customerID).Scan(&customer.CustomerID, &customer.Name, &customer.Phone, &customer.Email, &customer.Address)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("error retrieving customer: %v", err)
	}
	return &customer, nil
}

// UpdateCustomer updates an existing customer
func UpdateCustomer(customer models.Customer) error {
	query := `UPDATE customer SET name = $1, phone = $2, email = $3, address = $4 WHERE customer_id = $5`
	result, err := db.DB.Exec(query, customer.Name, customer.Phone, customer.Email, customer.Address, customer.CustomerID)
	if err != nil {
		return fmt.Errorf("failed to update customer: %v", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("customer not found")
	}
	return nil
}

// DeleteCustomer removes a customer by id
func DeleteCustomer(customerID int) error {
	query := `DELETE FROM customer WHERE customer_id = $1`
	result, err := db.DB.Exec(query, customerID)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %v", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("customer not found")
	}
	return nil
}

// GetAllCustomers retrieves all customers
func GetAllCustomers() ([]models.Customer, error) {
	query := `SELECT customer_id, name, phone, email, address FROM customer`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve customers: %v", err)
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		if err := rows.Scan(&customer.CustomerID, &customer.Name, &customer.Phone, &customer.Email, &customer.Address); err != nil {
			return nil, fmt.Errorf("error scanning customer: %v", err)
		}
		customers = append(customers, customer)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading customers: %v", err)
	}
	return customers, nil
}
