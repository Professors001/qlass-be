package usecases

import (
	"context"
	"errors"
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
	"qlass-be/infrastructure/utils"
	"qlass-be/transform"
)

type ClassUseCase interface {
	CreateClass(ctx context.Context, req *dtos.CreateClassRequestDto, ownerID uint) (*dtos.CreateClassResponseDto, error)
	GetClassDetailsByID(ctx context.Context, classID uint) (*dtos.ClassDetailsDto, error)
	GetClassDetailsByInviteCode(ctx context.Context, inviteCode string) (*dtos.ClassDetailsDto, error)
}

type classUseCase struct {
	classRepo repositories.ClassRepository
}

func NewClassUseCase(repo repositories.ClassRepository) ClassUseCase {
	return &classUseCase{
		classRepo: repo,
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

	// 4. Fetch the latest details
	classDetails, err := c.GetClassDetailsByID(ctx, uint(class.ID))
	if err != nil {
		return nil, err
	}

	// 5. Build Response following Return Message & Data requirement
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
