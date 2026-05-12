package auth

type Login struct {
	Nama     string `json:"nama" validate:"required"`
	Password string `json:"password" validate:"required"`
}
