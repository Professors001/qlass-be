package dtos

type CreateClassRequestDto struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Section     string `json:"section" binding:"required"`
	Term        string `json:"term" binding:"required"`
	Room        string `json:"room" binding:"required"`
}

type CreateClassResponseDto struct {
	Message string          `json:"message"`
	Data    ClassDetailsDto `json:"data"`
}

type ClassDetailsDto struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Section         string `json:"section"`
	Term            string `json:"term"`
	Room            string `json:"room"`
	InviteCode      string `json:"invite_code"`
	IsArchived      bool   `json:"is_archived"`
	OwnerID         string `json:"owner_id"`
	OwnerFirstName  string `json:"owner_first_name"`
	OwnerLastName   string `json:"owner_last_name"`
	OwnerProfileImg string `json:"owner_profile_img"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}
