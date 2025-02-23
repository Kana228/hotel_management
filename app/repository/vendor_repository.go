package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"hotelm/db"
	"hotelm/models"
)

// CreateVendor inserts a new vendor into the database
func CreateVendor(vendor models.Vendor) (int, error) {
	query := `INSERT INTO vendor (name, email, phone, hotel_name, address) VALUES ($1, $2, $3, $4, $5) RETURNING vendor_id`
	var id int
	err := db.DB.QueryRow(query, vendor.Name, vendor.Email, vendor.Phone, vendor.HotelName, vendor.Address).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create vendor: %v", err)
	}
	return id, nil
}

// GetVendorByID retrieves a vendor by ID
func GetVendorByID(vendorID int) (*models.Vendor, error) {
	query := `SELECT vendor_id, name, email, phone, hotel_name, address FROM vendor WHERE vendor_id = $1`
	var vendor models.Vendor

	err := db.DB.QueryRow(query, vendorID).Scan(&vendor.VendorID, &vendor.Name, &vendor.Email, &vendor.Phone, &vendor.HotelName, &vendor.Address)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("vendor not found")
		}
		return nil, fmt.Errorf("error retrieving vendor: %v", err)
	}
	return &vendor, nil
}

// UpdateVendor updates an existing vendor
func UpdateVendor(vendor models.Vendor) error {
	query := `UPDATE vendor SET name = $1, email = $2, phone = $3, hotel_name = $4, address = $5 WHERE vendor_id = $6`
	result, err := db.DB.Exec(query, vendor.Name, vendor.Email, vendor.Phone, vendor.HotelName, vendor.Address, vendor.VendorID)
	if err != nil {
		return fmt.Errorf("failed to update vendor: %v", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("vendor not found")
	}
	return nil
}

// DeleteVendor removes a vendor by ID
func DeleteVendor(vendorID int) error {
	query := `DELETE FROM vendor WHERE vendor_id = $1`
	result, err := db.DB.Exec(query, vendorID)
	if err != nil {
		return fmt.Errorf("failed to delete vendor: %v", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("vendor not found")
	}
	return nil
}

// GetAllVendors retrieves all vendors
func GetAllVendors() ([]models.Vendor, error) {
	query := `SELECT vendor_id, name, email, phone, hotel_name, address FROM vendor`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve vendors: %v", err)
	}
	defer rows.Close()

	var vendors []models.Vendor
	for rows.Next() {
		var vendor models.Vendor
		if err := rows.Scan(&vendor.VendorID, &vendor.Name, &vendor.Email, &vendor.Phone, &vendor.HotelName, &vendor.Address); err != nil {
			return nil, fmt.Errorf("error scanning vendor: %v", err)
		}
		vendors = append(vendors, vendor)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading vendors: %v", err)
	}
	return vendors, nil
}
