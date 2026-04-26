package dto

type AuthRequestDto struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponseDto struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
