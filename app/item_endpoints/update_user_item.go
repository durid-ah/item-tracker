package itemendpoints

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/durid-ah/item-tracker/helpers"
	"github.com/durid-ah/item-tracker/services"
)

func UpdateItemHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			var item services.Item
			itemSvc := services.ItemService{Db: db}
			username := r.Context().Value(helpers.UserNameKey)

			body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err := r.Body.Close(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				return
			}

			if err := json.Unmarshal(body, &item); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				println(err)
			}

			updateErr := itemSvc.Update(&item, username.(string))
			if updateErr != nil && updateErr.Error() == services.NotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			} else if updateErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(updateErr)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(item)
		})
}