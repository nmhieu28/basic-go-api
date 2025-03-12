package requests

type CreateUserRequest struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}
type LoginRequest struct {
	Email    string
	Password string
}
type VerifyEmailRequest struct {
	Token string
}
type ResetPasswordRequest struct {
	Email    string
	Password string
	Code     string
}
type ForgotPasswordRequest struct {
	Email string
}
