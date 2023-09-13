package services

import (
	"database/sql"
	"log"
	"time"
)

const NotFound = "Item not found"

type Item struct {
	Id             	int64   	`json:"id"`
	UserId		   	int64		`json:"userId"`
	Label          	string		`json:"label"`
	Description    	string 	 	`json:"description"`
	ContainerNumber int32		`json:"containerNumber"`
	ExpirationDate 	time.Time 	`json:"expirationDate"`
	Image          	[]byte    	`json:"image"`
}

type ItemService struct {
	Db *sql.DB
}

func (service *ItemService) Add(newItem *Item, username string) (int64, error) {
	tx, txErr := service.Db.Begin()
	if txErr != nil {
		return 0, txErr
	}

	userRow := tx.QueryRow("SELECT id FROM users WHERE username = ?", username)
	userRow.Scan(&newItem.UserId)

	stmt, stmtErr := tx.Prepare(`
		INSERT INTO items (user_id, label, description, container_number, expiration_date, image) VALUES (?,?,?,?,?,?)`)
	if stmtErr != nil {
		return 0, stmtErr
	}

	defer stmt.Close()

	res, err := stmt.Exec(
		newItem.UserId,
		newItem.Label,
		newItem.Description,
		newItem.ContainerNumber,
		newItem.ExpirationDate.Format(time.RFC3339),
		newItem.Image)

	if err != nil {
		return 0, err
	}

	tx.Commit()

	id, _ := res.LastInsertId()
	return id, nil
}

func (service *ItemService) Update(item *Item, username string) error {
	tx, txErr := service.Db.Begin()
	if txErr != nil {
		return txErr
	}

	stmt, stmtErr := tx.Prepare(`
		UPDATE items 
		SET label = ?, 
			description = ?,
			container_number = ?,
			expiration_date = ?,
			image = ?
		WHERE id = ? AND user_id = ?
	`)
	if stmtErr != nil {
		return stmtErr
	}

	defer stmt.Close()

	res, err := stmt.Exec(
		item.Label,
		item.Description,
		item.ContainerNumber,
		item.ExpirationDate.Format(time.RFC3339),
		item.Image,
		item.Id,
		item.UserId)
	if err != nil {
		return err
	}

	tx.Commit()

	numRows, _ := res.RowsAffected()
	if numRows == 0 {
		return &ServiceError{Message: NotFound}
	}

	return nil
}

func (service *ItemService) Delete(id int64, username string) error {
	tx, txErr := service.Db.Begin()
	if txErr != nil {
		return txErr
	}

	stmt, stmtErr := tx.Prepare(`
		DELETE FROM items 
		WHERE id IN (
			SELECT i.id FROM items i
			INNER JOIN users u ON i.user_id = u.id
			WHERE i.id = ? AND u.username = ?
		);
	`)
	if stmtErr != nil {
		return stmtErr
	}

	defer stmt.Close()

	res, err := stmt.Exec(id, username)
	if err != nil {
		return err
	}

	tx.Commit()

	numRows, _ := res.RowsAffected()
	if numRows == 0 {
		return &ServiceError{Message: NotFound}
	}

	return nil
}

func (service *ItemService) GetUserItems(username string) ([]Item, error) {
	items := make([]Item, 0)

	stmt, stmtErr := service.Db.Prepare(`
		SELECT i.id, i.label, i.description, i.container_number, i.expiration_date, i.image
		FROM users AS u
		INNER JOIN items AS i ON i.user_id = u.id
		WHERE username = ?`)
	if stmtErr != nil {
		return nil, stmtErr
	}
	
	rows, err := stmt.Query(username)
	if err != nil {
		return nil, err
	}
	
	defer rows.Close()

	for rows.Next() {
		item := Item{}
		var expirationDate string
		rowErr := rows.Scan(&item.Id, &item.Label, &item.Description, &item.ContainerNumber, &expirationDate, &item.Image)

		if rowErr != nil {
			return nil, rowErr
		}

		item.ExpirationDate, err = time.Parse(time.RFC3339, expirationDate)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}
