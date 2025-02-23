package repository

import (
	"hotelm/db"
	"hotelm/models"
	"errors"
)

// GetReviewByID retrieves a review by ID (CRUD: Read)
func GetReviewByID(reviewID int) (*models.Review, error) {
	var r models.Review
	err := db.DB.QueryRow("SELECT review_id, comment, rating, review_date, booking_id, customer_id, room_id FROM review WHERE review_id = $1", reviewID).
		Scan(&r.ReviewID, &r.Comment, &r.Rating, &r.ReviewDate, &r.BookingID, &r.CustomerID, &r.RoomID)

	if err != nil {
		return nil, err
	}
	return &r, nil
}

// CreateReview inserts a new review, ensuring rating is between 1 and 5
func CreateReview(r models.Review) (int, error) {
	if r.Rating < 1 || r.Rating > 5 {
		return 0, errors.New("rating must be between 1 and 5")
	}

	var newID int
	err := db.DB.QueryRow(
		"INSERT INTO review (comment, rating, review_date, booking_id, customer_id, room_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING review_id",
		r.Comment, r.Rating, r.ReviewDate, r.BookingID, r.CustomerID, r.RoomID,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}
	return newID, nil
}

// UpdateReview modifies an existing review, ensuring rating is between 1 and 5
func UpdateReview(r models.Review) error {
	if r.Rating < 1 || r.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	_, err := db.DB.Exec(
		"UPDATE review SET comment=$1, rating=$2, review_date=$3, booking_id=$4, customer_id=$5, room_id=$6 WHERE review_id=$7",
		r.Comment, r.Rating, r.ReviewDate, r.BookingID, r.CustomerID, r.RoomID, r.ReviewID,
	)

	return err
}

// DeleteReview removes a review by ID (CRUD: Delete)
func DeleteReview(reviewID int) error {
	_, err := db.DB.Exec("DELETE FROM review WHERE review_id=$1", reviewID)
	return err
}

// GetAllReviews retrieves all reviews.
func GetAllReviews() ([]models.Review, error) {
	rows, err := db.DB.Query("SELECT review_id, comment, rating, review_date, booking_id, customer_id, room_id FROM review")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		if err := rows.Scan(&review.ReviewID, &review.Comment, &review.Rating, &review.ReviewDate, &review.BookingID, &review.CustomerID, &review.RoomID); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}
