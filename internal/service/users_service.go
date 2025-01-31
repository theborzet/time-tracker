package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/theborzet/time-tracker/internal/models"
	"github.com/theborzet/time-tracker/internal/pagination"
)

func InputDataError(user *models.User) error {
	passportNumberInt, err := strconv.Atoi(user.PassportNumber)
	if err != nil {
		return errors.New("passport number must be a valid number")
	}
	passportSerieInt, err := strconv.Atoi(user.PassportSerie)
	if err != nil {
		return errors.New("passport serie must be a valid number")
	}
	if len(user.PassportNumber) != 6 && passportNumberInt <= 0 {
		return errors.New("incorrect passport_number data")
	}
	if len(user.PassportSerie) != 4 && passportSerieInt <= 0 {
		return errors.New("incorrect passport_serie data")
	}
	if user.Surname == "" || user.Name == "" {
		return errors.New("incorrect input data")
	}
	return nil
}

func (s *ApiService) GetUsersWithPaginate(filters map[string]string, page int) ([]*models.User, *pagination.Paginator, error) {
	if page < 0 {
		return nil, nil, errors.New("pagination error")
	}
	users, err := s.repo.GetUsers(filters)
	if err != nil {
		return nil, nil, err
	}
	pagintatedUsers, paginator, err := pagination.PaginateUser(users, page)
	if err != nil {
		return nil, nil, err
	}
	return pagintatedUsers, &paginator, nil
}

func (s *ApiService) CreateUser(passportNumber string) error {
	passportParts := strings.Split(passportNumber, " ")
	if len(passportParts) != 2 {
		return errors.New("invalid passport format. Passport number should be in format '1234 567890'")
	}

	passportSerie := passportParts[0]
	passportNum := passportParts[1]

	peopleInfo, err := s.exApi.FetchPeopleInfo(passportSerie, passportNum)
	if err != nil {
		return fmt.Errorf("failed to fetch people info: %w", err)
	}

	user := &models.User{
		PassportSerie:  passportSerie,
		PassportNumber: passportNum,
		Surname:        peopleInfo.Surname,
		Name:           peopleInfo.Name,
		Patronymic:     peopleInfo.Patronymic,
		Address:        peopleInfo.Address,
	}
	if err := InputDataError(user); err != nil {
		return err
	}

	if err := s.repo.CreateUser(user); err != nil {
		return err
	}
	return nil
}

func (s *ApiService) UpdateUser(user *models.User) error {
	if err := InputDataError(user); err != nil {
		return err
	}
	if err := s.repo.UpdateUser(user); err != nil {
		return err
	}
	return nil
}

func (s *ApiService) DeleteUser(userID int) error {
	if userID <= 0 {
		return errors.New("incorrect userId value")
	}
	if err := s.repo.DeleteUser(userID); err != nil {
		return err
	}
	return nil
}
