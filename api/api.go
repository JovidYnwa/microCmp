package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/JovidYnwa/microCmp/db"
	"github.com/JovidYnwa/microCmp/types"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type CompanyHandler struct {
	store db.CompanyStore
}

// Constructor function for CompanyHandler
func NewCompanyHandler(store db.CompanyStore) *CompanyHandler {
	return &CompanyHandler{
		store: store,
	}
}

// Run method to start the server
func (s *CompanyHandler) Run() {
	// router := mux.NewRouter()

	// router.HandleFunc("/companies", makeHTTPHandleFunc(s.handleGetCompanies))
	// router.HandleFunc("/company", makeHTTPHandleFunc(s.handleCreateCompany))

	// router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	// router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleGetAccountByID), s.store))
	// router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer))

	// Enable CORS for all routes
	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"*"}, // Allow all origins
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowedHeaders:   []string{"Authorization", "Content-Type"},
	// 	AllowCredentials: true,
	// })

	// handler := c.Handler(router)

	// log.Println("Json API server running on port: ", s.listenAddr)
	// http.ListenAndServe(s.listenAddr, handler)
}

// Handler methods (all returning error)
func (s *CompanyHandler) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	} else if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	return fmt.Errorf("method is not allowed %s", r.Method)
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}

func withJWTAuth(handlerFunc http.HandlerFunc, s db.CompanyStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")

		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}
		userID, err := getID(r)
		if err != nil {
			permissionDenied(w)
			return
		}
		account, err := s.GetAccountByID(userID)
		if err != nil {
			permissionDenied(w)
		}

		claims := token.Claims.(jwt.MapClaims)
		if account.Number != int64(claims["accountNumber"].(float64)) {
			permissionDenied(w)
			return
		}
		handlerFunc(w, r)
	}
}

const jwtSecret = "forTest"

func validateJWT(tokenString string) (*jwt.Token, error) {

	return jwt.ParseWithClaims(tokenString, nil, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unxpected singign method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

}

func (s *CompanyHandler) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		idStr := mux.Vars(r)["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("invalid id given %s", idStr)
		}

		account, err := s.store.GetAccountByID(id)
		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, account)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

// Get /account
func (s *CompanyHandler) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *CompanyHandler) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountRequest := new(types.CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountRequest); err != nil {
		return err
	}

	account := types.NewAccount(createAccountRequest.FirstName, createAccountRequest.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}
	token, err := createJWT(account)
	if err != nil {
		return err
	}
	fmt.Println("jwt token= ", token)

	return WriteJSON(w, http.StatusOK, account)
}

func (s *CompanyHandler) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id given %s", idStr)
	}
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *CompanyHandler) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(types.TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()
	return WriteJSON(w, http.StatusOK, transferReq)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return json.NewEncoder(w).Encode(v)
}

func createJWT(account *types.Account) (string, error) {
	// Create the Claims
	claims := &jwt.MapClaims{
		"ExpiresAt":     jwt.NewNumericDate(time.Unix(1516239022, 0)),
		"Issuer":        "test",
		"AccountNumber": account.Number,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))

}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}

func (h *CompanyHandler) HandleCreateCompany(w http.ResponseWriter, r *http.Request) {
	createCompanyRequest := new(types.CreateCompanyReq)
	if err := json.NewDecoder(r.Body).Decode(createCompanyRequest); err != nil {
		// return err
		WriteJSON(w, http.StatusBadRequest, "bad bad ")
	}

	cmpID, err := h.store.SetCompany(createCompanyRequest.Company)
	if err != nil {
		// return err
		WriteJSON(w, http.StatusBadRequest, "bad bad bad ")

	}
	createCompanyRequest.CompanyInfo.CompanyID = *cmpID

	if err := h.store.SetCompanyInfo(createCompanyRequest.CompanyInfo); err != nil {
		// return err
		WriteJSON(w, http.StatusBadRequest, "bad bad 1")
	}

	WriteJSON(w, http.StatusOK, createCompanyRequest)
}

// Get /account /companies?page=1&pageSize=10
func (h *CompanyHandler) HandleGetCompanies(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for pagination
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10 // Default page size
	}

	paginatedResponse, err := h.store.GetCompanies(page, pageSize)
	if err != nil {
		fmt.Println(err)
		WriteJSON(w, http.StatusOK, paginatedResponse)
	}
	WriteJSON(w, http.StatusOK, paginatedResponse)
}

// Get /companies?id=1
func (h *CompanyHandler) HandleGetCompanyDetail(w http.ResponseWriter, r *http.Request) {
	companyID := mux.Vars(r)["id"]

	// s, err := h.store.GetCompanies(page, pageSize)
	// if err != nil {
	// 	fmt.Println(err)
	// 	WriteJSON(w, http.StatusOK, paginatedResponse)
	// }
	WriteJSON(w, http.StatusOK, "dsalkjfhakjdshf"+companyID)
}
