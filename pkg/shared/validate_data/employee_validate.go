package validate_data

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/helper"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateEmployee(input *domain.Employee) error {
	if input.ID == primitive.NilObjectID {
		return errors.New("the employee's information cannot be empty")
	}

	if input.OrganizationID == primitive.NilObjectID {
		return errors.New("the employee's information cannot be empty")
	}

	if input.Email == "" {
		return errors.New("the employee's information cannot be empty")
	}
	if !helper.EmailValid(input.Email) {
		return errors.New("email Invalid ")
	}

	if input.FirstName == "" {
		return errors.New("the employee's information cannot be empty")
	}

	if input.LastName == "" {
		return errors.New("the employee's information cannot be empty")
	}

	if input.JobTitle == "" {
		return errors.New("the employee's information cannot be empty")
	}
	return nil
}

func ValidateEmployeeInput(input *domain.EmployeeInput) error {
	if input.OrganizationID == primitive.NilObjectID {
		return errors.New("the employee's information cannot be empty")
	}

	if input.Email == "" {
		return errors.New("the employee's information cannot be empty")
	}
	if !helper.EmailValid(input.Email) {
		return errors.New("email Invalid ")
	}

	if input.FirstName == "" {
		return errors.New("the employee's information cannot be empty")
	}

	if input.LastName == "" {
		return errors.New("the employee's information cannot be empty")
	}

	if input.JobTitle == "" {
		return errors.New("the employee's information cannot be empty")
	}
	return nil
}
