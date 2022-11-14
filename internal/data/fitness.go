//Filename: internal/data/fitness.go

package data

import (
	"database/sql"
	"time"
	"fmt"
	"context"

	"fitness.zioncastillo.net/internal/validator"

)

type Fitness struct {
	ID	    int     `json:"id"`
	User_id  int     `json:"user_id"`
	Steps   int     `json:"steps"`
	Cups    int     `json:"cups"`
	Date 	time.Time `json:"date"`
}

func ValidateItem(v *validator.Validator, fitness *Fitness) {
	// Use the Check() method to execute our validation checks
	//v.Check(fitness.Steps < 0, "Steps", "cannot be less than 0")
	//v.Check(fitness.Cups < 0, "Cups", "cannot be less than 0")
}

 //Define a FitnessModel which wraps a sql.DB connection pool
type FitnessModel struct {
	DB *sql.DB
}

//Insert function that will insert the users fitness tracked for the day
func (m FitnessModel) Insert(fitness * Fitness) error {
	
	query := `
		INSERT INTO dailyfitness (user_id, steps, cups)
		VALUES ($1, $2, $3)
		RETURNING id, date
	`
	args := []interface{}{
		fitness.User_id,
		fitness.Steps,
		fitness.Cups,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&fitness.ID, &fitness.Date)
}

//Get all function that will list all the records stored
func (m FitnessModel) GetAll(id int, user_id int, steps int, cups int, date time.Time, filters Filters ) ([]*Fitness, Metadata, error) {
	
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, user_id, steps, cups, date
		FROM dailyfitness
		WHERE (to_tsvector('simple', steps) @@ plainto_tsquery('simple',$1) OR $1 = '')
		AND (to_tsvector('simple', cups) @@ plainto_tsquery('simple', $2) OR $2 = '')
		ORDER BY %s %s, id DESC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortOrder())

		// Create a 3-second-timout context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Execute the query
	args := []interface{}{steps, cups, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	// Close the resultset
	defer rows.Close()

	totalRecords := 0

	// Initialize an empty slice to hold the School data
	lists := []*Fitness{}

	// Iterate over the rows in the resultset
	for rows.Next() {
		var fitness Fitness
		// Scan the values from the row into school
		err := rows.Scan(
			&totalRecords,
			&fitness.ID,
			&fitness.User_id,
			&fitness.Steps,
			&fitness.Cups,
			&fitness.Date,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		// Add the School to our slice
		lists = append(lists, &fitness)
	}
	// Check for errors after looping through the resultset
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// Return the slice of Schools
	return lists, metadata, nil
}

// Delete() removes a specific record *only beingn used for testing*
func (m FitnessModel) Delete(id int64) error {
	// Ensure that there is a valid id
	if id < 1 {
		return ErrRecordNotFound
	}
	// Create the delete query
	query := `
		DELETE FROM dailyfitness
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	// Execute the query
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	// Check how many rows were affected by the delete operation. We
	// call the RowsAffected() method on the result variable
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// Check if no rows were affected
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
