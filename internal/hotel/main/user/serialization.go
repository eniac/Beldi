package user

type Request struct {
	Username string
	Password string
}

type Result struct {
	Correct bool
}
