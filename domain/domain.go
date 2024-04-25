package domain

type Car struct {
	ID              int    `json:"id"`
	RegNum          string `json:"reg_num"`
	Mark            string `json:"mark"`
	Model           string `json:"model"`
	Year            int    `json:"year"`
	OwnerName       string `json:"owner_name"`
	OwnerSurname    string `json:"owner_surname"`
	OwnerPatronymic string `json:"owner_patronymic"`
}
