package auth

type RegisterRequestDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponseDto struct {
	Message string `json:"message"`
}

type LoginRequestDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseDto struct {
	AuthToken string `json:"authToken"`
}
