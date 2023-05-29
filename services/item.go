package services

import (
	"database/sql"
	"time"
)

const NotFound = "Item not found"

type Item struct {
	Id             int64     `json:"id"`
	Label          string    `json:"label"`
	Description    string    `json:"description"`
	ExpirationDate time.Time `json:"expirationDate"`
	Image          []byte    `json:"image"`
}

type ItemService struct {
	Db *sql.DB
}

func (service *ItemService) Add(newItem *Item) (int64, error) {
	tx, txErr := service.Db.Begin()
	if txErr != nil {
		return 0, txErr
	}

	stmt, stmtErr := tx.Prepare(`
		INSERT INTO items (label, description, expiration_date, image) VALUES (?,?,?,?)`)
	if stmtErr != nil {
		return 0, stmtErr
	}

	defer stmt.Close()

	res, err := stmt.Exec(
		newItem.Label,
		newItem.Description,
		newItem.ExpirationDate.Format(time.RFC3339),
		newItem.Image)

	if err != nil {
		return 0, err
	}

	tx.Commit()

	id, _ := res.LastInsertId()
	return id, nil
}

func (service *ItemService) Update(item *Item) error {
	tx, txErr := service.Db.Begin()
	if txErr != nil {
		return txErr
	}

	stmt, stmtErr := tx.Prepare(`
		UPDATE items 
		SET label = ?, 
			 description = ?, 
			 expiration_date = ?,
			 image = ?
		WHERE id = ?
	`)
	if stmtErr != nil {
		return stmtErr
	}

	defer stmt.Close()

	res, err := stmt.Exec(
		item.Label,
		item.Description,
		item.ExpirationDate.Format(time.RFC3339),
		item.Image,
		item.Id)

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

func (service *ItemService) Delete(id int64) error {
	tx, txErr := service.Db.Begin()
	if txErr != nil {
		return txErr
	}

	stmt, stmtErr := tx.Prepare(`DELETE FROM items WHERE id = ?`)
	if stmtErr != nil {
		return stmtErr
	}

	defer stmt.Close()

	res, err := stmt.Exec(id)
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

func (service *ItemService) GetAll() ([]Item, error) {
	items := make([]Item, 0)
	rows, err := service.Db.Query("SELECT id, label, description, expiration_date, image FROM items")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		item := Item{}
		var expirationDate string
		rowErr := rows.Scan(&item.Id, &item.Label, &item.Description, &expirationDate, &item.Image)

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
