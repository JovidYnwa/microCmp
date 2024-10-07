package api

import (
	"fmt"
	"net/http"

	"github.com/JovidYnwa/microCmp/db"
)

type CompanyFilterHandler struct {
	filterStore db.CompanyFilterStore
}

func NewCompanyFilterHandler(companyStore db.CompanyFilterStore) *CompanyFilterHandler {
	return &CompanyFilterHandler{
		filterStore: companyStore,
	}
}

func (h *CompanyFilterHandler) HandleListTrpls(w http.ResponseWriter, r *http.Request) {
	result, err := h.filterStore.GetTrpls()
	if err != nil {
		fmt.Printf("Trpls %s", err)
		WriteJSON(w, 401, "smth bad happaned")
	}
	WriteJSON(w, 200, result)
}

func (h *CompanyFilterHandler) HandleRgionsrpls(w http.ResponseWriter, r *http.Request) {
	result, err := h.filterStore.GetRegions()
	if err != nil {
		fmt.Printf("Trpls %s", err)
		WriteJSON(w, 401, "smth bad happaned")
	}
	WriteJSON(w, 200, result)
}

func (h *CompanyFilterHandler) HandleSubscriberStatus(w http.ResponseWriter, r *http.Request) {
	result, err := h.filterStore.GetSubsStatuses()
	if err != nil {
		fmt.Printf("Trpls %s", err)
		WriteJSON(w, 401, "smth bad happaned")
	}
	WriteJSON(w, 200, result)
}
