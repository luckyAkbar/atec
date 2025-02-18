package usecase_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/usecase"
)

func TestRegisterChildInputValidate(t *testing.T) {
	t.Run("Valid RegisterChildInput", func(t *testing.T) {
		input := usecase.RegisterChildInput{
			DateOfBirth: time.Now(),
			Gender:      true,
			Name:        "John Doe",
		}

		if err := input.Validate(); err != nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, false)
		}
	})

	t.Run("Invalid RegisterChildInput - Missing DateOfBirth", func(t *testing.T) {
		input := usecase.RegisterChildInput{
			Gender: true,
			Name:   "John Doe",
		}

		if err := input.Validate(); err == nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("Invalid RegisterChildInput - Missing Name", func(t *testing.T) {
		input := usecase.RegisterChildInput{
			DateOfBirth: time.Now(),
			Gender:      true,
		}

		if err := input.Validate(); err == nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("Invalid RegisterChildInput - Missing DateOfBirth and Name", func(t *testing.T) {
		input := usecase.RegisterChildInput{
			Gender: true,
		}

		if err := input.Validate(); err == nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, true)
		}
	})
}

func TestUpdateChildInputValidate(t *testing.T) {
	name := "John Doe"

	t.Run("Valid UpdateChildInput", func(t *testing.T) {
		dateOfBirth := time.Now()
		gender := true

		input := usecase.UpdateChildInput{
			ChildID:     uuid.New(),
			DateOfBirth: &dateOfBirth,
			Gender:      &gender,
			Name:        &name,
		}

		if err := input.Validate(); err != nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, false)
		}
	})

	t.Run("Invalid UpdateChildInput - Missing ChildID", func(t *testing.T) {
		dateOfBirth := time.Now()
		gender := true
		input := usecase.UpdateChildInput{
			DateOfBirth: &dateOfBirth,
			Gender:      &gender,
			Name:        &name,
		}

		if err := input.Validate(); err == nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("Valid UpdateChildInput - Only ChildID", func(t *testing.T) {
		input := usecase.UpdateChildInput{
			ChildID: uuid.New(),
		}

		if err := input.Validate(); err != nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, false)
		}
	})

	t.Run("Valid UpdateChildInput - ChildID and DateOfBirth", func(t *testing.T) {
		dateOfBirth := time.Now()
		input := usecase.UpdateChildInput{
			ChildID:     uuid.New(),
			DateOfBirth: &dateOfBirth,
		}

		if err := input.Validate(); err != nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, false)
		}
	})

	t.Run("Valid UpdateChildInput - ChildID and Gender", func(t *testing.T) {
		gender := true
		input := usecase.UpdateChildInput{
			ChildID: uuid.New(),
			Gender:  &gender,
		}

		if err := input.Validate(); err != nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, false)
		}
	})

	t.Run("Valid UpdateChildInput - ChildID and Name", func(t *testing.T) {
		input := usecase.UpdateChildInput{
			ChildID: uuid.New(),
			Name:    &name,
		}

		if err := input.Validate(); err != nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, false)
		}
	})
}
