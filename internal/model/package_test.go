package model_test

import (
	"testing"

	"github.com/luckyAkbar/atec/internal/model"
)

func TestChecklistGroupValidation(t *testing.T) {
	t.Run("empty custom name should return error", func(t *testing.T) {
		cg := model.ChecklistGroup{
			CustomName: "",
		}

		err := cg.Validate()
		if err == nil {
			t.Errorf("on empty custom name, expected error, but got nil")
		}
	})

	t.Run("empty questions should return error", func(t *testing.T) {
		cg := model.ChecklistGroup{
			CustomName: "just to not trigger the error here",
			Questions:  []string{},
		}

		err := cg.Validate()
		if err == nil {
			t.Errorf("on empty questions, expected error, but got nil")
		}
	})

	t.Run("empty options should return error", func(t *testing.T) {
		cg := model.ChecklistGroup{
			CustomName: "just to not trigger the error here",
			Questions:  []string{"just to not trigger the error here"},
			Options:    []model.AnswerOption{},
		}

		err := cg.Validate()
		if err == nil {
			t.Errorf("on empty options, expected error, but got nil")
		}
	})
}
