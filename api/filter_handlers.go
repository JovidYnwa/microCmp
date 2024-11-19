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

func (h *CompanyFilterHandler) HandleServList(w http.ResponseWriter, r *http.Request) {
	result, err := h.filterStore.GetServs()
	if err != nil {
		fmt.Printf("Servs %s", err)
		WriteJSON(w, 401, err)
	}
	WriteJSON(w, 200, result)
}

func (h *CompanyFilterHandler) HandleSimStatus(w http.ResponseWriter, r *http.Request) {
	result, err := h.filterStore.GetSimTypes()
	if err != nil {
		fmt.Printf("Trpls %s", err)
		WriteJSON(w, 401, "smth bad happaned")
	}
	WriteJSON(w, 200, result)
}

func (h *CompanyFilterHandler) HandleDivceTypes(w http.ResponseWriter, r *http.Request) {
	result := []map[string]any{
		{
			"id": 1,
			"name":"android",
		},
		{
			"id": 2,
			"name":"ios",
		},
	}
	WriteJSON(w, 200, result)
}

func (h *CompanyFilterHandler) HandlePrizeList(w http.ResponseWriter, r *http.Request) {
	result := []map[string]any{
		{
			"id": 1,
			"name":"1000 gb",
		},
		{
			"id": 2,
			"name":"2000 gb",
		},
	}
	WriteJSON(w, 200, result)
}


func (h *CompanyFilterHandler) HandleActionCmp(w http.ResponseWriter, r *http.Request) {
	result := []map[string]any{
		{
			"id": 1,
			"name":"Пополнение баланса",
		},
		{
			"id": 2,
			"name":"Что то еще",
		},
	}
	WriteJSON(w, 200, result)
}

