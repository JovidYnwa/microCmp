package api

import (
	"net/http"

	"github.com/JovidYnwa/microCmp/db"
)

type CompanyHandler struct {
	store db.Storage
}

func NewCompanyHandler(companryStore db.Storage) *CompanyHandler {
	return &CompanyHandler{
		store: companryStore,
	}
}

func HandleTestFunc1(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, 200, "yo12")
}
