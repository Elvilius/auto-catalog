package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Elvilius/auto-catalog/internal/service"
)

type httpHandler struct {
	Service *service.Service
}

type CreateCars struct {
	RegNums []string `json:"reg_nums"`
}

type Success struct {
	Success bool `json:"success"`
}

type CarId struct {
	ID int `json:"id"`
}

type CarFilter struct {
	RegNum          string
	Mark            string
	Model           string
	YearFrom        int
	YearTo          int
	OwnerName       string
	OwnerSurname    string
	OwnerPatronymic string
	Page            int
	PageSize        int
}

type UpdateCar struct {
	ID              int    `json:"id"`
	RegNum          string `json:"reg_num"`
	Mark            string `json:"mark"`
	Model           string `json:"model"`
	Year            int    `json:"year"`
	OwnerName       string `json:"owner_name"`
	OwnerSurname    string `json:"owner_surname"`
	OwnerPatronymic string `json:"owner_patronymic"`
}

func Register(server *http.ServeMux, service *service.Service) {
	h := httpHandler{Service: service}
	server.HandleFunc("/auto-catalog/create", h.Create)
	server.HandleFunc("/auto-catalog/list", h.List)
	server.HandleFunc("/auto-catalog/delete", h.Delete)
	server.HandleFunc("/auto-catalog/update", h.Update)
}

func (h *httpHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody CreateCars
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Internal error", http.StatusBadRequest)
		return
	}

	if len(requestBody.RegNums) == 0 {
		http.Error(w, "reg_nums field is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err := h.Service.Create(ctx, requestBody.RegNums)
	if err != nil {
		http.Error(w, "Internal error", http.StatusBadRequest)
		return
	}

	res := Success{Success: true}
	jsonResp, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonResp)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

func (h *httpHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()
	queryParams := r.URL.Query()

	reqNum := queryParams.Get("reg_num")
	mark := queryParams.Get("mark")
	model := queryParams.Get("model")
	ownerName := queryParams.Get("owner_name")
	ownerSurname := queryParams.Get("owner_surname")
	ownerPatronymic := queryParams.Get("owner_patronymic")

	var err error

	yearToNum := 0
	yearTo := queryParams.Get("year_to")

	if yearTo != "" {
		yearToNum, err = strconv.Atoi(yearTo)
		if err != nil {
			http.Error(w, "year_to must be int", http.StatusBadRequest)
			return
		}
	}

	yearFromNum := 0
	yearFrom := queryParams.Get("year_from")
	if yearFrom != "" {
		yearFromNum, err = strconv.Atoi(yearFrom)
		if err != nil {
			http.Error(w, "year_from must be int", http.StatusBadRequest)
			return
		}
	}

	pageNum := 1
	page := queryParams.Get("page")
	if page != "" {
		pageNum, err = strconv.Atoi(page)
		if err != nil {
			http.Error(w, "page must be int", http.StatusBadRequest)
			return
		}
	}

	pageSizeNum := 10
	pageSize := queryParams.Get("page_size")
	if pageSize != "" {
		pageSizeNum, err = strconv.Atoi(pageSize)
		if err != nil {
			http.Error(w, "page_size must be int", http.StatusBadRequest)
			return
		}
	}

	filter := service.CarFilter(CarFilter{
		Page:            pageNum,
		PageSize:        pageSizeNum,
		YearFrom:        yearFromNum,
		YearTo:          yearToNum,
		OwnerName:       ownerName,
		OwnerSurname:    ownerSurname,
		Model:           model,
		Mark:            mark,
		RegNum:          reqNum,
		OwnerPatronymic: ownerPatronymic,
	})
	cars, err := h.Service.List(ctx, filter)
	if err != nil {
		http.Error(w, "Internal error", http.StatusBadRequest)
		return
	}

	jsonResp, err := json.Marshal(cars)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonResp)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

func (h *httpHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody CarId
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Internal error", http.StatusBadRequest)
		return
	}

	if requestBody.ID == 0 {
		http.Error(w, "ID field is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err := h.Service.Delete(ctx, requestBody.ID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusBadRequest)
		return
	}

	res := Success{Success: true}
	jsonResp, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonResp)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

func (h *httpHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody UpdateCar
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Internal error", http.StatusBadRequest)
		return
	}

	if requestBody.ID == 0 {
		http.Error(w, "ID field is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err := h.Service.Update(ctx, service.UpdateCar(requestBody))
	if err != nil {
		http.Error(w, "Internal error", http.StatusBadRequest)
		return
	}

	res := Success{Success: true}
	jsonResp, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonResp)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}
