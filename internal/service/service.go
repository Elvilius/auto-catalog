package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/Elvilius/auto-catalog/domain"
	"github.com/Elvilius/auto-catalog/internal/config"
	"github.com/Elvilius/auto-catalog/internal/repo"

	"github.com/alitto/pond"
)

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
	ID              int
	RegNum          string
	Mark            string
	Model           string
	Year            int
	OwnerName       string
	OwnerSurname    string
	OwnerPatronymic string
}

func NewService(repo repo.RepoInterface, cfg *config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

type Service struct {
	repo repo.RepoInterface
	cfg  *config.Config
}

func (s *Service) Create(ctx context.Context, regNums []string) error {
	cars, err := s.getCarsInfo(regNums)
	if err != nil {
		return err
	}
	return s.repo.CreateCars(ctx, cars)
}

func (s *Service) getCarsInfo(regNums []string) ([]domain.Car, error) {
	carsInfo := make([]domain.Car, len(regNums))

	pool := pond.New(s.cfg.MaxWorkers, s.cfg.MaxWorkers*10)
	defer pool.StopAndWait()

	group, ctx := pool.GroupContext(context.Background())

	for i, regNum := range regNums {
		regNum := regNum
		group.Submit(func() error {
			params := url.Values{}
			params.Add("regNum", regNum)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.cfg.CarInfoUrl, nil)
			if err != nil {
				return err
			}
			req.URL.RawQuery = params.Encode()
			resp, err := http.DefaultClient.Do(req)
			if err == nil {
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				var data map[string]any
				err = json.Unmarshal(body, &data)
				if err != nil {
					return err
				}
				car := s.getFormatCar(data)
				carsInfo[i] = car

			}
			return err
		})
	}

	err := group.Wait()
	return carsInfo, err
}

func (s *Service) getFormatCar(data map[string]any) domain.Car {
	yearFloat := data["year"].(float64)
	year := int(yearFloat)
	return domain.Car{
		RegNum:          data["regNum"].(string),
		Mark:            data["mark"].(string),
		Model:           data["model"].(string),
		Year:            year,
		OwnerName:       data["owner"].(map[string]interface{})["name"].(string),
		OwnerPatronymic: data["owner"].(map[string]interface{})["patronymic"].(string),
		OwnerSurname:    data["owner"].(map[string]interface{})["surname"].(string),
	}
}

func (s *Service) List(ctx context.Context, filter CarFilter) ([]domain.Car, error) {
	return s.repo.GetCars(ctx, repo.CarFilter(filter))
}

func (s *Service) Delete(ctx context.Context, ID int) error {
	return s.repo.DeleteCar(ctx, ID)
}

func (s *Service) Update(ctx context.Context, updateCar UpdateCar) error {
	return s.repo.UpdateCar(ctx, repo.UpdateCar(updateCar))
}
