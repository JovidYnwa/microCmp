package api

import (
	"fmt"
	"net/http"

	"github.com/JovidYnwa/microCmp/db"
)

type CompanyHandler struct {
	filterStore db.CompanyFilterStore
}

func NewCompanyHandler(companyStore db.CompanyFilterStore) *CompanyHandler {
	return &CompanyHandler{
		filterStore: companyStore,
	}
}

func (h *CompanyHandler) HandleListTrpls(w http.ResponseWriter, r *http.Request) {
	result, err := h.filterStore.GetTrpls()
	if err != nil {
		fmt.Printf("Trpls %s", err)
		WriteJSON(w, 401, "smth bad happaned")
	}
	WriteJSON(w, 200, result)
}

func (h *CompanyHandler) HandleRgionsrpls(w http.ResponseWriter, r *http.Request) {
	result, err := h.filterStore.GetRegions()
	if err != nil {
		fmt.Printf("Trpls %s", err)
		WriteJSON(w, 401, "smth bad happaned")
	}
	WriteJSON(w, 200, result)
}

func (h *CompanyHandler) HandleSubscriberStatus(w http.ResponseWriter, r *http.Request) {
	result, err := h.filterStore.GetSubsStatuses()
	if err != nil {
		fmt.Printf("Trpls %s", err)
		WriteJSON(w, 401, "smth bad happaned")
	}
	WriteJSON(w, 200, result)
}
