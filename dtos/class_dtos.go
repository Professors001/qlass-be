package dtos

type CreateClassRequestDto struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Section     string `json:"section" binding:"required"`
	Term        string `json:"term" binding:"required"`
	Room        string `json:"room"`
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
	Role            string `json:"role"`
}

type EnrollRequestDto struct {
	InviteCode string `json:"invite_code" binding:"required"`
}

type StudentDetailsDto struct {
	UniversityID string `json:"university_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	ProfileImg   string `json:"profile_img"`
	Email        string `json:"email"`
	EnrolledRole string `json:"enrolled_role"`
	IsArchived   bool   `json:"is_archived"`
	EnrolledAt   string `json:"enrolled_at"`
	Status       string `json:"status"`
}

type SummaryEnrolledStudentsDto struct {
	ClassID      string              `json:"class_id"`
	StudentCount int                 `json:"student_count"`
	Teachers     []StudentDetailsDto `json:"teachers"`
	Students     []StudentDetailsDto `json:"students"`
}

type GetEnrolledStudentsResponseDto struct {
	Message string            `json:"message"`
	Data    StudentDetailsDto `json:"data"`
}
