package books

import (
	"encoding/json"
	"net/http"
	"fmt"
)


func (c *BookController) GetCatalog(w http.ResponseWriter, r *http.Request) {
	books, err := c.catalogHandler.Handle(r.Context())
	if err != nil {

		fmt.Println("ПИЗДЕЦ В КАТАЛОГЕ:", err)
		c.sendError(w, http.StatusInternalServerError, "failed to load catalog")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}