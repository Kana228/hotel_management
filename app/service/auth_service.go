package service

import (
	"fmt"

	"hotelm/repository"
	"hotelm/session"
)

// LoginCustomer checks if a customer exists with the given id and name.
// If successful, it sets the global session pointer to the customer.
func LoginCustomer(customerID int, name string) error {
	customer, err := repository.GetCustomerByID(customerID)
	if err != nil {
		return fmt.Errorf("customer login failed: %v", err)
	}
	// Check if the provided name matches the retrieved customer
	if customer.Name != name {
		return fmt.Errorf("customer login failed: name does not match")
	}

	// Set the global session pointer for the logged-in customer.
	session.SetCurrentUser(customer)
	return nil
}

// LoginVendor checks if a vendor exists with the given id and name.
// If successful, it sets the global session pointer to the vendor.
func LoginVendor(vendorID int, name string) error {
	vendor, err := repository.GetVendorByID(vendorID)
	if err != nil {
		return fmt.Errorf("vendor login failed: %v", err)
	}
	// Check if the provided name matches the retrieved vendor
	if vendor.Name != name {
		return fmt.Errorf("vendor login failed: name does not match")
	}

	// Set the global session pointer for the logged-in vendor.
	session.SetCurrentUser(vendor)
	return nil
}
