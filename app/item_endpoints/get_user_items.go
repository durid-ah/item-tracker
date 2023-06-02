package itemendpoints

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/durid-ah/item-tracker/helpers"
	"github.com/durid-ah/item-tracker/services"
)

func GetItemsHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			username := r.Context().Value(helpers.UserNameKey)
			itemSvc := services.ItemService{Db: db}
			items, err := itemSvc.GetUserItems(username.(string))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(items)
		})
}