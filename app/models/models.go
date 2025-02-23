package models

import "time"

type Customer struct {
	CustomerID int    
	Name       string 
	Phone      string 
	Email      string 
	Address    string 
}

type Vendor struct {
	VendorID   int    
	Name       string 
	Email      string 
	Phone      string 
	HotelName  string 
	Address    string 
}

type Room struct {
	RoomID        int      
	Name          string   
	Description   string   
	Location      string   
	Availability  bool     
	Price         float64  
	RoomType      string   
	AverageRating float64  
	Amenities     string   // Now a single string, e.g., "WiFi,TV,Mini Bar"
	VendorID      int      
}

type Booking struct {
	BookingID     int       
	BookingDate   time.Time 
	CheckinDate   time.Time 
	CheckoutDate  time.Time 
	PaymentStatus string    
	RoomID        int       
	CustomerID    int       
}

type Payment struct {
	PaymentID      int       
	PaymentMethod  string    
	PaymentStatus  string    
	TransactionDate time.Time 
	Amount         float64   
	BookingID      int       
}

type Review struct {
	ReviewID   int       
	Comment    string    
	Rating     int       
	ReviewDate time.Time 
	BookingID  int       
	CustomerID int       
	RoomID     int       
}
