package rest

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/atec/internal/config"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/usecase"
)

// @Summary		Create a new account to access user only resources within this system
// @Description	Use this API endpoint to create a new account
// @Tags			Authentication
// @Accept			json
// @Param			signup_input	body	SignupInput	true	"Login Credentials"
// @Produce		json
// @Success		200	{object}	StandardSuccessResponse{data=SignupOutput}	"Successful response"
// @Failure		400	{object}	StandardErrorResponse						"Bad request / validation error"
// @Router			/v1/auth/signup [post]
func (s *Service) HandleSignUp() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &SignupInput{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		output, err := s.authUsecase.HandleSignup(c.Request().Context(), usecase.SignupInput{
			Email:    input.Email,
			Password: input.Password,
			Username: input.Username,
		})

		if err != nil {
			return UsecaseErrorToRESTResponse(c, err)
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data: SignupOutput{
				Message: output.Message,
			},
		})
	}
}

// @Summary		Resend email for account verification
// @Description	If something happen during account verification that cause the email not received
// @Description	or the email is lost, use this API to resend the verification email
// @Tags			Authentication
// @Accept			json
// @Param			resend_signup_input	body	ResendVerificationInput	true	"resend signup input"
// @Produce		json
// @Success		200	{object}	StandardSuccessResponse{data=ResendVerificationOutput}	"Successful response"
// @Failure		400	{object}	StandardErrorResponse									"Bad request / validation error"
// @Failure		429	{object}	StandardErrorResponse									"Too many requests"
// @Failure		500	{object}	StandardErrorResponse									"Internal Error"
// @Router			/v1/auth/signup/resend [post]
func (s *Service) HandleResendSignupVerification() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &ResendVerificationInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		output, err := s.authUsecase.HandleResendSignupVerification(c.Request().Context(), usecase.ResendSignupVerificationInput{
			Email: input.Email,
		})

		if err != nil {
			return UsecaseErrorToRESTResponse(c, err)
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data: ResendVerificationOutput{
				Message: output.Message,
			},
		})
	}
}

// @Summary		Validate account after signup
// @Description	Confirmation method to ensure the email used when signup is active and owned by requester.
// @Description	If the confirmation token is valid, the account will be activated
// @Description	and will be able to be used on login. Otherwise, the opposite will happen.
// @Tags			Authentication
// @Accept			json
// @Param			validation_token	query	string	true	"validation token received from email after signup in form of a jwt token"
// @Produce		json
// @Success		200	{object}	StandardSuccessResponse{data=VerifyAccountOutput}	"Successful response"
// @Failure		400	{object}	StandardErrorResponse								"Bad request / validation error"
// @Failure		401	{object}	StandardErrorResponse								"invalid or expired verification token. details will be given in the error message"
// @Router			/v1/auth/verify [get]
func (s *Service) HandleVerifyAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &VerifyAccountInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		_, err := s.authUsecase.HandleAccountVerification(c.Request().Context(), usecase.AccountVerificationInput{
			VerificationToken: input.ValidationToken,
		})

		if err != nil {
			return UsecaseErrorToRESTResponse(c, err)
		}

		return c.HTML(http.StatusOK, accountVerifiedWebpage())
	}
}

// @Summary		Gain access to the system by authenticating using a registered account
// @Description	Use this endpoint to login with your username and password
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Param			login_input	body		LoginInput									true	"account detail such as email and password to log in"
// @Success		200			{object}	StandardSuccessResponse{data=LoginOutput}	"Successful response"
// @Failure		400			{object}	StandardErrorResponse						"Bad request"
// @Failure		401			{object}	StandardErrorResponse						"Authentication Failed"
// @Failure		500			{object}	StandardErrorResponse						"Internal Error"
// @Router			/v1/auth/login [post]
func (s *Service) HandleLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &LoginInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		output, err := s.authUsecase.HandleLogin(c.Request().Context(), usecase.LoginInput{
			Email:    input.Email,
			Password: input.Password,
		})

		if err != nil {
			return UsecaseErrorToRESTResponse(c, err)
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data: LoginOutput{
				Token: output.Token,
			},
		})
	}
}

// @Summary		Initiate change password process for an active account
// @Description	If the user wants to change their password, use this API.
// @Description	when the request succeed, an email containing confirmation link will
// @Description	be sent to the account email.
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Param			init_change_password_input	body		InitResetPasswordInput									true	"the email of the account which password will be reset"
// @Success		200							{object}	StandardSuccessResponse{data=InitResetPasswordOutput}	"Successful response"
// @Failure		400							{object}	StandardErrorResponse									"Bad request"	"validation error"
// @Failure		500							{object}	StandardErrorResponse									"Internal Error"
// @Router			/v1/auth/password [patch]
func (s *Service) HandleInitResetPassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &InitResetPasswordInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		output, err := s.authUsecase.HandleInitesetPassword(c.Request().Context(), usecase.InitResetPasswordInput{
			Email: input.Email,
		})

		if err != nil {
			return UsecaseErrorToRESTResponse(c, err)
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data: InitResetPasswordOutput{
				Message: output.Message,
			},
		})
	}
}

// @Summary		Change user's account password
// @Description	After using init reset password method, the resulting JWT token can be used here
// @Description	to authorize the password changes.
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Param			change_password_input	body		ResetPasswordInput									true	"the email of the account which password will be reset"
// @Success		200						{object}	StandardSuccessResponse{data=ResetPasswordOutput}	"Successful response"
// @Failure		400						{object}	StandardErrorResponse								"Bad request"	"validation error"
// @Failure		500						{object}	StandardErrorResponse								"Internal Error"
// @Router			/v1/auth/password [post]
func (s *Service) HandleResetPassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &ResetPasswordInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		output, err := s.authUsecase.HandleResetPassword(c.Request().Context(), usecase.ResetPasswordInput{
			ResetPasswordToken: input.Token,
			NewPassword:        input.NewPassword,
		})

		if err != nil {
			return UsecaseErrorToRESTResponse(c, err)
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data: ResetPasswordOutput{
				Message: output.Message,
			},
		})
	}
}

// @Summary		Webpage's form to change user's account password
// @Description	After user click the link from email for reset password, he will be redirected here to change the password. This might not work in this page
// @Description	since the returned value in this API is HTML
// @Tags			Authentication
// @Accept			application/x-www-form-urlencoded
// @Produce		html
// @Success		200	{object}	StandardSuccessResponse	"Successful response"
// @Failure		400	{object}	StandardErrorResponse	"Bad request"	"validation error"
// @Failure		500	{object}	StandardErrorResponse	"Internal Error"
// @Router			/v1/auth/password [get]
func (s *Service) HandleRenderChangePasswordPage() echo.HandlerFunc {
	return func(c echo.Context) error {
		changePasswordToken := c.QueryParam(model.ChangePasswordTokenQuery)
		if changePasswordToken == "" {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "missing required query parameter",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		return c.HTML(http.StatusOK, changePasswordWebpage(changePasswordToken))
	}
}

// @Summary		Delete user's account
// @Description	Delete all user's related data from the system
// @Tags			Authentication
// @Accept			application/json
// @Security		ParentLevelAuth
// @Param			Authorization			header	string				true	"JWT Token"
// @Param			change_password_input	body	DeleteAccountInput	true	"the email of the account which password will be reset"
//
// @Produce		json
// @Success		204	{object}	StandardSuccessResponse	"Successful response"
// @Failure		400	{object}	StandardErrorResponse	"Bad request"	"validation error"
// @Failure		500	{object}	StandardErrorResponse	"Internal Error"
// @Router			/v1/auth/accounts [delete]
func (s *Service) HandleDeleteAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &DeleteAccountInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		err := s.authUsecase.HandleDeleteUserData(c.Request().Context(), usecase.DeleteUserDataInput{
			Email:    input.Email,
			Password: input.Password,
		})

		if err != nil {
			return UsecaseErrorToRESTResponse(c, err)
		}

		return c.NoContent(http.StatusNoContent)
	}
}

//nolint:funlen
func changePasswordWebpage(token string) string {
	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Formulir Penggantian Kata Sandi</title>
		<style>
			/* General styles */
			body {
				font-family: Arial, sans-serif;
				margin: 0;
				padding: 0;
				background-color: #f6f6f6;
				color: #333;
			}
			
			.email-container {
				max-width: 600px;
				margin: 20px auto;
				background-color: #ffffff;
				border: 1px solid #ddd;
				border-radius: 8px;
				overflow: hidden;
			}
			
			.header {
				background-color: #4CAF50;
				color: white;
				padding: 20px;
				text-align: center;
			}
			
			.content {
				padding: 20px;
			}
			
			.content p {
				margin: 0 0 15px;
				line-height: 1.6;
			}
			
			.form-group {
				margin-bottom: 20px;
				position: relative;
			}
			
			.form-group label {
				display: block;
				margin-bottom: 8px;
				font-weight: bold;
			}
			
			.form-group input {
				width: 100%%;
				padding: 10px;
				border: 1px solid #ddd;
				border-radius: 4px;
				box-sizing: border-box;
				font-size: 16px;
			}
			
			.password-toggle {
				position: absolute;
				right: 10px;
				top: 38px;
				cursor: pointer;
				color: #666;
				background: none;
				border: none;
			}
			
			.btn-container {
				text-align: center;
				margin: 20px 0;
			}
			
			.btn {
				display: inline-block;
				background-color: #4CAF50;
				color: white;
				text-decoration: none;
				padding: 10px 20px;
				font-size: 16px;
				border-radius: 5px;
				border: none;
				cursor: pointer;
			}
			
			.btn:hover {
				background-color: #45a049;
			}
			
			.footer {
				background-color: #f1f1f1;
				text-align: center;
				padding: 10px;
				font-size: 12px;
				color: #666;
			}

			/* Modal styles */
			.modal {
				display: none;
				position: fixed;
				top: 0;
				left: 0;
				width: 100%%;
				height: 100%%;
				background-color: rgba(0, 0, 0, 0.5);
				justify-content: center;
				align-items: center;
			}
			.modal-content {
				background-color: white;
				padding: 20px;
				border-radius: 8px;
				width: 300px;
				text-align: center;
				box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
			}
			.modal-title {
				font-size: 20px;
				font-weight: bold;
				color: #4CAF50;
				margin-bottom: 15px;
			}
			.modal-body {
				font-size: 14px;
				color: #333;
				margin-bottom: 20px;
			}
			.modal-close-btn {
				background-color: #4CAF50;
				color: white;
				border: none;
				padding: 8px 16px;
				border-radius: 4px;
				cursor: pointer;
			}
			.modal-close-btn:hover {
				background-color: #45a049;
			}
		</style>
	</head>
	<body>
		<div class="email-container">
			<div class="header">
				<h2>Atur Ulang Kata Sandi</h2>
			</div>
			
			<div class="content">
				<p>Mohon masukan kata sandi baru anda di formulir di bawah ini.</p>
				
				<form id="passwordForm">
					<div class="form-group">
						<label for="password">Kata Sandi</label>
						<input type="password" id="password" name="password" required>
						<button type="button" class="password-toggle" onclick="togglePassword('password')">
							👁️
						</button>
					</div>
					
					<div class="form-group">
						<label for="confirmPassword">Konfirmasi Kata Sandi</label>
						<input type="password" id="confirmPassword" name="confirmPassword" required>
						<button type="button" class="password-toggle" onclick="togglePassword('confirmPassword')">
							👁️
						</button>
					</div>
					
					<div class="btn-container">
						<button type="button" class="btn" onclick="changePassword()">Ganti Password</button>
					</div>
				</form>
			</div>
			
			<div class="footer">
							<p>&copy; 2025. Hak Cipta Dilindungi.</p>
			</div>
		</div>

		<!-- Modal HTML -->
		<div id="successModal" class="modal">
			<div class="modal-content">
				<div class="modal-title">Kata Sandi Berhasil Diubah</div>
				<div class="modal-body">Silahkan <i>login</i> kembali di aplikasi dengan menggunakan akun Anda</div>
				<button class="modal-close-btn" onclick="closeModal()">Tutup</button>
			</div>
		</div>

		<script>
			function togglePassword(inputId) {
				const passwordInput = document.getElementById(inputId);
				if (passwordInput.type === "password") {
					passwordInput.type = "text";
				} else {
					passwordInput.type = "password";
				}
			}
			
			document.addEventListener('DOMContentLoaded', function() {
				const changePasswordBtn = document.querySelector('.btn');
				
				changePasswordBtn.addEventListener('click', async function() {
					// Disable the button
					changePasswordBtn.disabled = true;
					changePasswordBtn.textContent = 'Memproses...';
					changePasswordBtn.style.opacity = '0.7';
					changePasswordBtn.style.cursor = 'not-allowed';

					try {
						const password = document.getElementById('password').value;
						const confirmPassword = document.getElementById('confirmPassword').value;
						
						// Validation
						if (password === '' || confirmPassword === '') {
							throw new Error('Harap isi kedua kolom kata sandi!');
						}
						
						if (password !== confirmPassword) {
							throw new Error('kata sandi tidak cocok!');
						}

						if (password.length < 8 || confirmPassword.length < 8) {
							throw new Error('kata sandi minimal 8 karakter!');
						}

						await resetPasswordAPI(password);
						
						showModal();
						
						document.getElementById('passwordForm').reset();
					} catch (error) {
						alert(error.message);
					} finally {
						changePasswordBtn.disabled = false;
						changePasswordBtn.textContent = 'Ganti Kata Sandi';
						changePasswordBtn.style.opacity = '1';
						changePasswordBtn.style.cursor = 'pointer';
					}
				});
			});

			// Get modal element
			const modal = document.getElementById('successModal');

			function showModal() {
				modal.style.display = 'flex';
			}

			function closeModal() {
				modal.style.display = 'none';
			}


			async function resetPasswordAPI(new_password) {
				const baseURL = '%s';
				const body = {
					'token': '%s',
					'new_password': new_password,
				};

				headers = {
					'Content-Type': 'application/json',
				}

				const response = await fetch(baseURL, {
					method: 'POST',
					headers: headers,
					body: JSON.stringify(body),
				});

				const data = await response.json();

				if (response.status !== 200) {
					throw new Error(data.error_message);
				}
			}
		</script>
		</body>
	</html>`, config.ServerResetPasswordBaseURL(), token)
}

func accountVerifiedWebpage() string {
	return `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Aktifasi Akun</title>
				<style>
					/* General styles */
					body {
						font-family: Arial, sans-serif;
						margin: 0;
						padding: 0;
						background-color: #f6f6f6;
						color: #333;
					}
					.email-container {
						max-width: 600px;
						margin: 20px auto;
						background-color: #ffffff;
						border: 1px solid #ddd;
						border-radius: 8px;
						overflow: hidden;
					}
					.header {
						background-color: #4CAF50;
						color: white;
						padding: 20px;
						text-align: center;
					}
					.content {
						padding: 20px;
					}
					.content p {
						margin: 0 0 15px;
						line-height: 1.6;
					}
					.btn-container {
						text-align: center;
						margin: 20px 0;
					}
					.btn {
						display: inline-block;
						background-color: #4CAF50;
						color: white;
						text-decoration: none;
						padding: 10px 20px;
						font-size: 16px;
						border-radius: 5px;
					}
					.btn:hover {
						background-color: #45a049;
					}
					.footer {
						background-color: #f1f1f1;
						text-align: center;
						padding: 10px;
						font-size: 12px;
						color: #666;
					}
				</style>
			</head>
			<body>
				<div class="email-container">
					<div class="header">
						<h2>Verifikasi Akun</h2>
					</div>
					<div class="content">
						<p>Akun Anda telah diaktifasi dan sudah bisa digunakan.</p>
						<p>Silakan <i>Log In </i> kembali dengan akun yang telah dibuat.</p>
					</div>
					<div class="footer">
						<p>&copy; 2025. Hak Cipta Dilindungi.</p>
					</div>
				</div>
			</body>
			</html>`
}
