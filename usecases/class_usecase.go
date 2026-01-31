package usecases

import (
	"context"
	"errors"
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
	"qlass-be/transform"
	"qlass-be/utils"
	"strconv"
)

type ClassUseCase interface {
	CreateClass(ctx context.Context, req *dtos.CreateClassRequestDto, ownerID uint) (*dtos.CreateClassResponseDto, error)
	GetClassDetailsByID(ctx context.Context, classID uint) (*dtos.ClassDetailsDto, error)
	GetClassDetailsByInviteCode(ctx context.Context, inviteCode string) (*dtos.ClassDetailsDto, error)
	GetAllMyClasses(ctx context.Context, userID uint) ([]dtos.ClassDetailsDto, error)
	EnrollStudent(ctx context.Context, inviteCode string, studentID uint) error
	GetEnrolledStudentsByClassID(ctx context.Context, classID uint) (*dtos.SummaryEnrolledStudentsDto, error)
}

type classUseCase struct {
	classRepo  repositories.ClassRepository
	enrollRepo repositories.EnrollRepository
}

func NewClassUseCase(classRepo repositories.ClassRepository, enrollRepo repositories.EnrollRepository) ClassUseCase {
	return &classUseCase{
		classRepo:  classRepo,
		enrollRepo: enrollRepo,
	}
}

func (c *classUseCase) CreateClass(ctx context.Context, req *dtos.CreateClassRequestDto, ownerID uint) (*dtos.CreateClassResponseDto, error) {
	// 1. Generate Invite Code
	inviteCode, err := c.generateUniqueInviteCode(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Map DTO to Entity
	class := &entities.Class{
		Name:        req.Name,
		Description: req.Description,
		Section:     req.Section,
		Term:        req.Term,
		Room:        req.Room,
		InviteCode:  inviteCode,
		OwnerID:     ownerID,
	}

	// 3. Persist to Repository
	if err := c.classRepo.Create(class); err != nil {
		return nil, err
	}

	// 4. Enroll the owner as a teacher
	if err := c.enrollRepo.EnrollWithRole(class.ID, ownerID, "teacher"); err != nil {
		return nil, err
	}

	// 5. Fetch the latest details
	classDetails, err := c.GetClassDetailsByID(ctx, uint(class.ID))
	if err != nil {
		return nil, err
	}

	// 6. Build Response following Return Message & Data requirement
	return &dtos.CreateClassResponseDto{
		Message: "Class created successfully",
		Data:    *classDetails,
	}, nil
}

func (c *classUseCase) GetClassDetailsByID(ctx context.Context, classID uint) (*dtos.ClassDetailsDto, error) {
	class, err := c.classRepo.GetByID(classID)
	if err != nil {
		return nil, err
	}

	classDetailsDto := transform.EntityToClassDetailsDto(*class)

	return &classDetailsDto, nil
}

func (c *classUseCase) GetClassDetailsByInviteCode(ctx context.Context, inviteCode string) (*dtos.ClassDetailsDto, error) {
	class, err := c.classRepo.GetByInviteCode(inviteCode)
	if err != nil {
		return nil, err
	}

	classDetailsDto := transform.EntityToClassDetailsDto(*class)

	return &classDetailsDto, nil
}

func (c *classUseCase) GetAllMyClasses(ctx context.Context, userID uint) ([]dtos.ClassDetailsDto, error) {
	enrollments, err := c.classRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var result []dtos.ClassDetailsDto
	for _, enrollment := range enrollments {
		dto := transform.EntityToClassDetailsDto(enrollment.Class)
		dto.Role = enrollment.Role
		result = append(result, dto)
	}

	if result == nil {
		result = []dtos.ClassDetailsDto{}
	}

	return result, nil
}

func (c *classUseCase) EnrollStudent(ctx context.Context, inviteCode string, studentID uint) error {
	class, err := c.classRepo.GetByInviteCode(inviteCode)
	if err != nil {
		return err
	}

	isEnrolled, err := c.enrollRepo.IsEnrolled(class.ID, studentID)
	if err != nil {
		return err
	}
	if isEnrolled {
		return errors.New("user is already enrolled in this class")
	}

	return c.enrollRepo.EnrollStudent(class.ID, studentID)
}

func (c *classUseCase) GetEnrolledStudentsByClassID(ctx context.Context, classID uint) (*dtos.SummaryEnrolledStudentsDto, error) {
	enrollments, err := c.enrollRepo.GetEnrolledStudentsByClassID(classID)
	if err != nil {
		return nil, err
	}

	var teachers []dtos.StudentDetailsDto
	var students []dtos.StudentDetailsDto

	for _, enrollment := range enrollments {
		dto := transform.EntityToStudentDetailsDto(enrollment)
		if enrollment.Role == "teacher" || enrollment.Role == "ta" {
			teachers = append(teachers, dto)
		} else {
			students = append(students, dto)
		}
	}

	return &dtos.SummaryEnrolledStudentsDto{
		ClassID:      strconv.FormatUint(uint64(classID), 10),
		StudentCount: len(students),
		Teachers:     teachers,
		Students:     students,
	}, nil
}

func (c *classUseCase) generateUniqueInviteCode(_ context.Context) (string, error) {
	const maxRetries = 10
	const codeLength = 6

	for i := 0; i < maxRetries; i++ {
		code := utils.GenerateRandomString(codeLength)

		_, err := c.classRepo.GetByInviteCode(code)
		if err != nil {
			return code, nil
		}

	}

	return "", errors.New("failed to generate unique invite code: maximum retries reached")
}
