package itemendpoints

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/durid-ah/item-tracker/dto"
	"github.com/durid-ah/item-tracker/helpers"
	"github.com/durid-ah/item-tracker/services"
)

func DeleteItemHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			var itemDto dto.DeleteitemDto
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

			if err := json.Unmarshal(body, &itemDto); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				println(err)
			}

			deleteErr := itemSvc.Delete(itemDto.Id, username.(string))
			if deleteErr != nil && deleteErr.Error() == services.NotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			} else if deleteErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(deleteErr)
				return
			}

			w.WriteHeader(http.StatusOK)
		})
}