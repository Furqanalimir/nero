package forms

type Authenticate struct {
	Phone    int    `json:"phone" validate:"required"`
	Password string `json:"password" validate:"password"`
}
type data struct{}
type ReqResSwagger struct {
	Data  data
	Error string
}
