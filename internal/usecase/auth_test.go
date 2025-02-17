package usecase_test

import (
	"testing"

	"github.com/luckyAkbar/atec/internal/usecase"
)

func TestLoginInput(t *testing.T) {
	t.Run("email is required and in valid format", func(t *testing.T) {
		li := usecase.LoginInput{
			Email:    "just filling randomly",
			Password: "thisShouldBeVal1dPass!",
		}

		err := li.Validate()
		if err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("empty email is unacceptable", func(t *testing.T) {
		li := usecase.LoginInput{
			Email:    "",
			Password: "thisShouldBeVal1dPass!",
		}

		err := li.Validate()
		if err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("password must not empty", func(t *testing.T) {
		li := usecase.LoginInput{
			Email:    "valid@mail.test",
			Password: "",
		}

		err := li.Validate()
		if err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("password must be atleast 8 chars", func(t *testing.T) {
		li := usecase.LoginInput{
			Email:    "valid@mail.test",
			Password: "only7ch",
		}

		err := li.Validate()
		if err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("ok", func(t *testing.T) {
		li := usecase.LoginInput{
			Email:    "valid@mail.test",
			Password: "only8char",
		}

		err := li.Validate()
		if err != nil {
			t.Errorf("expecting nil err but got %v", err)
		}
	})
}

func TestSignupInput(t *testing.T) {
	t.Run("email is required and in valid format", func(t *testing.T) {
		si := usecase.SignupInput{
			Email:    "just filling randomly",
			Password: "thisShouldBeVal1dPass!",
			Username: "valid",
		}

		err := si.Validate()
		if err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("empty email is unacceptable", func(t *testing.T) {
		si := usecase.SignupInput{
			Email:    "",
			Password: "thisShouldBeVal1dPass!",
			Username: "valid",
		}

		err := si.Validate()
		if err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("password must not empty", func(t *testing.T) {
		si := usecase.SignupInput{
			Email:    "valid@mail.test",
			Password: "",
			Username: "valid",
		}

		err := si.Validate()
		if err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("password must be atleast 8 chars", func(t *testing.T) {
		si := usecase.SignupInput{
			Email:    "valid@mail.test",
			Password: "only7ch",
			Username: "valid",
		}

		err := si.Validate()
		if err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("username must not empty", func(t *testing.T) {
		si := usecase.SignupInput{
			Email:    "valid@mail.test",
			Password: "only8char",
			Username: "",
		}

		err := si.Validate()
		if err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("ok", func(t *testing.T) {
		si := usecase.SignupInput{
			Email:    "valid@mail.test",
			Password: "only8char",
			Username: "valid",
		}

		err := si.Validate()
		if err != nil {
			t.Errorf("expecting nil err but got %v", err)
		}
	})
}

func TestAccountVerificationInput(t *testing.T) {
	t.Run("token is required", func(t *testing.T) {
		vi := usecase.AccountVerificationInput{
			VerificationToken: "",
		}

		if err := vi.Validate(); err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("ok", func(t *testing.T) {
		vi := usecase.AccountVerificationInput{
			VerificationToken: "ok",
		}

		if err := vi.Validate(); err != nil {
			t.Errorf("expecting nil but got %v", err)
		}
	})
}

func TestInitResetPasswordInput(t *testing.T) {
	t.Run("email is required", func(t *testing.T) {
		irpi := usecase.InitResetPasswordInput{
			Email: "",
		}

		if err := irpi.Validate(); err == nil {
			t.Error("expecting error but got nil")
		}
	})

	t.Run("email must also in valid format", func(t *testing.T) {
		irpi := usecase.InitResetPasswordInput{
			Email: "invalid format",
		}

		if err := irpi.Validate(); err == nil {
			t.Error("expecting error but got nil")
		}
	})

	t.Run("ok", func(t *testing.T) {
		irpi := usecase.InitResetPasswordInput{
			Email: "valid@email.test",
		}

		if err := irpi.Validate(); err != nil {
			t.Errorf("expecting nil but got %v", err)
		}
	})
}

func TestResetPasswordInput(t *testing.T) {
	t.Run("reset password token is required", func(t *testing.T) {
		rpi := usecase.ResetPasswordInput{
			ResetPasswordToken: "",
			NewPassword:        "validPass!",
		}

		if err := rpi.Validate(); err == nil {
			t.Error("expecting error but got nil")
		}
	})

	t.Run("new password is required", func(t *testing.T) {
		rpi := usecase.ResetPasswordInput{
			ResetPasswordToken: "validToken",
			NewPassword:        "",
		}

		if err := rpi.Validate(); err == nil {
			t.Error("expecting error but got nil")
		}
	})

	t.Run("new password must be atleast 8 chars", func(t *testing.T) {
		rpi := usecase.ResetPasswordInput{
			ResetPasswordToken: "validToken",
			NewPassword:        "7chars!",
		}

		if err := rpi.Validate(); err == nil {
			t.Error("expecting error but got nil")
		}
	})

	t.Run("ok", func(t *testing.T) {
		rpi := usecase.ResetPasswordInput{
			ResetPasswordToken: "validToken",
			NewPassword:        "validPass!",
		}

		if err := rpi.Validate(); err != nil {
			t.Errorf("expecting nil but got %v", err)
		}
	})
}
