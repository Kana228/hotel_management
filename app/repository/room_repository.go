package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"hotelm/db"
	"hotelm/models"
)

// CreateRoom inserts a new room into the database
func CreateRoom(room models.Room) (int, error) {
	query := `INSERT INTO room (name, description, location, availability, price, room_type, average_rating, amenities, vendor_id) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING room_id`
	var id int
	err := db.DB.QueryRow(query, room.Name, room.Description, room.Location, room.Availability, room.Price, room.RoomType, room.AverageRating, room.Amenities, room.VendorID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create room: %v", err)
	}
	return id, nil
}


func GetRoomByID(roomID int) (*models.Room, error) {
	query := `SELECT room_id, name, description, location, availability, price, room_type, average_rating, amenities, vendor_id FROM room WHERE room_id = $1`
	var room models.Room

	err := db.DB.QueryRow(query, roomID).Scan(
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
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("room not found")
		}
		return nil, fmt.Errorf("error retrieving room: %v", err)
	}
	return &room, nil
}


// UpdateRoom updates an existing room
func UpdateRoom(room models.Room) error {
	query := `UPDATE room SET name = $1, description = $2, location = $3, availability = $4, price = $5, room_type = $6, average_rating = $7, amenities = $8, vendor_id = $9 WHERE room_id = $10`
	result, err := db.DB.Exec(query, room.Name, room.Description, room.Location, room.Availability, room.Price, room.RoomType, room.AverageRating, room.Amenities, room.VendorID, room.RoomID)
	if err != nil {
		return fmt.Errorf("failed to update room: %v", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("room not found")
	}
	return nil
}


// DeleteRoom removes a room by ID
func DeleteRoom(roomID int) error {
	query := `DELETE FROM room WHERE room_id = $1`
	result, err := db.DB.Exec(query, roomID)
	if err != nil {
		return fmt.Errorf("failed to delete room: %v", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("room not found")
	}
	return nil
}

// GetAvailableRooms retrieves all available rooms
func GetAvailableRooms() ([]models.Room, error) {
	query := `SELECT room_id, name, description, location, availability, price, room_type, average_rating, amenities, vendor_id FROM room WHERE availability = TRUE`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve available rooms: %v", err)
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		if err := rows.Scan(&room.RoomID, &room.Name, &room.Description, &room.Location, &room.Availability, &room.Price, &room.RoomType, &room.AverageRating, &room.Amenities, &room.VendorID); err != nil {
			return nil, fmt.Errorf("error scanning room: %v", err)
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rooms: %v", err)
	}
	return rooms, nil
}
