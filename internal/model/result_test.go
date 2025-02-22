package model_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/luckyAkbar/atec/internal/model"
)

func TestAnswerDetailScannerAndValuer(t *testing.T) {
	t.Run("Value valid field value should return marshaled JSON", func(t *testing.T) {
		answerDetail := model.AnswerDetail{
			1: {1: 1, 2: 2},
			2: {1: 3, 2: 4},
		}

		fieldValue := answerDetail

		expectedJSON, err := json.Marshal(fieldValue)
		if err != nil {
			t.Fatalf("failed to marshal field value: %v", err)
		}

		value, err := answerDetail.Value(context.Background(), nil, reflect.Value{}, fieldValue)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(value, expectedJSON) {
			t.Errorf("expected %s, got %s", expectedJSON, value)
		}
	})

	t.Run("Value invalid field value should return error", func(t *testing.T) {
		answerDetail := model.AnswerDetail{}
		invalidFieldValue := make(chan int) // invalid type for JSON marshaling

		_, err := answerDetail.Value(context.Background(), nil, reflect.Value{}, invalidFieldValue)
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})

	t.Run("Scan nil dbValue should return nil error", func(t *testing.T) {
		answerDetail := &model.AnswerDetail{}

		err := answerDetail.Scan(context.Background(), nil, reflect.Value{}, nil)
		if err != nil {
			t.Errorf("expected nil error, but got %v", err)
		}

		if !reflect.DeepEqual(*answerDetail, model.AnswerDetail{}) {
			t.Errorf("expected empty AnswerDetail, but got %v", *answerDetail)
		}
	})

	t.Run("Scan valid JSON byte slice should unmarshal correctly", func(t *testing.T) {
		answerDetail := &model.AnswerDetail{}
		dbValue := []byte(`{
            "1": {"1": 1, "2": 2},
            "2": {"1": 3, "2": 4}
        }`)

		err := answerDetail.Scan(context.Background(), nil, reflect.Value{}, dbValue)
		if err != nil {
			t.Errorf("expected nil error, but got %v", err)
		}

		expected := model.AnswerDetail{
			1: {1: 1, 2: 2},
			2: {1: 3, 2: 4},
		}

		if !reflect.DeepEqual(*answerDetail, expected) {
			t.Errorf("expected %v, but got %v", expected, *answerDetail)
		}
	})

	t.Run("Scan valid JSON string should unmarshal correctly", func(t *testing.T) {
		answerDetail := &model.AnswerDetail{}
		dbValue := `{
            "1": {"1": 1, "2": 2},
            "2": {"1": 3, "2": 4}
        }`

		err := answerDetail.Scan(context.Background(), nil, reflect.Value{}, dbValue)
		if err != nil {
			t.Errorf("expected nil error, but got %v", err)
		}

		expected := model.AnswerDetail{
			1: {1: 1, 2: 2},
			2: {1: 3, 2: 4},
		}

		if !reflect.DeepEqual(*answerDetail, expected) {
			t.Errorf("expected %v, but got %v", expected, *answerDetail)
		}
	})

	t.Run("Scan invalid type should return error", func(t *testing.T) {
		answerDetail := &model.AnswerDetail{}
		dbValue := 12345 // invalid type

		err := answerDetail.Scan(context.Background(), nil, reflect.Value{}, dbValue)
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})

	t.Run("Scan invalid JSON should return error", func(t *testing.T) {
		answerDetail := &model.AnswerDetail{}
		dbValue := []byte(`invalid json`)

		err := answerDetail.Scan(context.Background(), nil, reflect.Value{}, dbValue)
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})
}

func TestResultDetailScannerAndValuer(t *testing.T) {
	t.Run("Value valid field value should return marshaled JSON", func(t *testing.T) {
		resultDetail := model.ResultDetail{
			1: {Name: "Subtest 1", Grade: 85},
			2: {Name: "Subtest 2", Grade: 90},
		}

		fieldValue := resultDetail

		expectedJSON, err := json.Marshal(fieldValue)
		if err != nil {
			t.Fatalf("failed to marshal field value: %v", err)
		}

		value, err := resultDetail.Value(context.Background(), nil, reflect.Value{}, fieldValue)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(value, expectedJSON) {
			t.Errorf("expected %s, got %s", expectedJSON, value)
		}
	})

	t.Run("Value invalid field value should return error", func(t *testing.T) {
		resultDetail := model.ResultDetail{}
		invalidFieldValue := make(chan int) // invalid type for JSON marshaling

		_, err := resultDetail.Value(context.Background(), nil, reflect.Value{}, invalidFieldValue)
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})

	t.Run("Scan nil dbValue should return nil error", func(t *testing.T) {
		resultDetail := &model.ResultDetail{}

		err := resultDetail.Scan(context.Background(), nil, reflect.Value{}, nil)
		if err != nil {
			t.Errorf("expected nil error, but got %v", err)
		}

		if !reflect.DeepEqual(*resultDetail, model.ResultDetail{}) {
			t.Errorf("expected empty ResultDetail, but got %v", *resultDetail)
		}
	})

	t.Run("Scan valid JSON byte slice should unmarshal correctly", func(t *testing.T) {
		resultDetail := &model.ResultDetail{}
		dbValue := []byte(`{
            "1": {"name": "Subtest 1", "grade": 85},
            "2": {"name": "Subtest 2", "grade": 90}
        }`)

		err := resultDetail.Scan(context.Background(), nil, reflect.Value{}, dbValue)
		if err != nil {
			t.Errorf("expected nil error, but got %v", err)
		}

		expected := model.ResultDetail{
			1: {Name: "Subtest 1", Grade: 85},
			2: {Name: "Subtest 2", Grade: 90},
		}

		if !reflect.DeepEqual(*resultDetail, expected) {
			t.Errorf("expected %v, but got %v", expected, *resultDetail)
		}
	})

	t.Run("Scan valid JSON string should unmarshal correctly", func(t *testing.T) {
		resultDetail := &model.ResultDetail{}
		dbValue := `{
            "1": {"name": "Subtest 1", "grade": 85},
            "2": {"name": "Subtest 2", "grade": 90}
        }`

		err := resultDetail.Scan(context.Background(), nil, reflect.Value{}, dbValue)
		if err != nil {
			t.Errorf("expected nil error, but got %v", err)
		}

		expected := model.ResultDetail{
			1: {Name: "Subtest 1", Grade: 85},
			2: {Name: "Subtest 2", Grade: 90},
		}

		if !reflect.DeepEqual(*resultDetail, expected) {
			t.Errorf("expected %v, but got %v", expected, *resultDetail)
		}
	})

	t.Run("Scan invalid type should return error", func(t *testing.T) {
		resultDetail := &model.ResultDetail{}
		dbValue := 12345 // invalid type

		err := resultDetail.Scan(context.Background(), nil, reflect.Value{}, dbValue)
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})

	t.Run("Scan invalid JSON should return error", func(t *testing.T) {
		resultDetail := &model.ResultDetail{}
		dbValue := []byte(`invalid json`)

		err := resultDetail.Scan(context.Background(), nil, reflect.Value{}, dbValue)
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})
}

func TestResultDetailCountTotalScore(t *testing.T) {
	t.Run("CountTotalScore should return total score from all subtest", func(t *testing.T) {
		resultDetail := model.ResultDetail{
			1: {Name: "Subtest 1", Grade: 85},
			2: {Name: "Subtest 2", Grade: 90},
		}

		expectedTotalScore := 175

		totalScore := resultDetail.CountTotalScore()
		if totalScore != expectedTotalScore {
			t.Errorf("expected %d, got %d", expectedTotalScore, totalScore)
		}
	})
}
