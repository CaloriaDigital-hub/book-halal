package commands

type SignInCommand struct {
	Email    string
	Password string
}

type SignInResult struct {
	Token     string
	ExpiresAt string
}