package usecases

import (
	"encoding/json"
	"errors"
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
	"qlass-be/transforms"
	"qlass-be/utils"
	"sort"
	"time"

	"gorm.io/datatypes"
)

type ClassMaterialUseCase interface {
	CreateClassMaterial(dto *dtos.CreateClassMaterialDto, ownerID uint) error
	CreateQuizMaterial(dto dtos.CreateQuizClassMaterialDto) error
	GetMaterialByID(id uint) (*dtos.GetClassMaterialDto, error)
	GetMaterialsByClassID(classID uint, userID uint) ([]*dtos.GetThumnailClassMaterialDto, error)
	UpdatePostClassMaterial(dto *dtos.UpdatePostClassMaterialDto, ownerID uint) error
	UpdateAssignmentClassMaterial(dto *dtos.UpdateAssignmentClassMaterialDto, ownerID uint) error
	UpdateQuizClassMaterial(dto *dtos.UpdateQuizClassMaterialDto, ownerID uint) error
	DeleteClassMaterial(classMaterialID uint, ownerID uint) error
}

type classMaterialUseCase struct {
	classMaterialRepo repositories.ClassMaterialRepository
	classRepo         repositories.ClassRepository
	attachmentRepo    repositories.AttachmentRepository
	attachmentUseCase AttachmentUseCase
	userUseCase       UserUseCase
	quizGameLogRepo   repositories.QuizGameLogRepository
	quizRepo          repositories.QuizRepository
}

func NewClassMaterialUseCase(
	classMaterialRepo repositories.ClassMaterialRepository,
	classRepo repositories.ClassRepository,
	attachmentRepo repositories.AttachmentRepository,
	attachmentUseCase AttachmentUseCase,
	quizGameLogRepo repositories.QuizGameLogRepository,
	userUseCase UserUseCase,
	quizRepo repositories.QuizRepository,
) ClassMaterialUseCase {
	return &classMaterialUseCase{
		classMaterialRepo: classMaterialRepo,
		classRepo:         classRepo,
		attachmentRepo:    attachmentRepo,
		attachmentUseCase: attachmentUseCase,
		quizGameLogRepo:   quizGameLogRepo,
		userUseCase:       userUseCase,
		quizRepo:          quizRepo,
	}
}

func (u *classMaterialUseCase) CreateClassMaterial(dto *dtos.CreateClassMaterialDto, ownerID uint) error {

	classId := dto.ClassID

	class, err := u.classRepo.GetByID(classId)
	if err != nil {
		return err
	}

	if class.OwnerID != ownerID {
		return errors.New("only class owner can create class material")
	}

	classMaterial := transforms.CreateToEntity(dto)

	if dto.Action == "publish" {
		classMaterial.IsPublished = true
		now := time.Now()
		classMaterial.PublishedAt = &now
	}

	err = u.classMaterialRepo.Create(classMaterial)
	if err != nil {
		return err
	}

	for _, attachmentID := range dto.AttachmentIds {
		attachment, err := u.attachmentRepo.GetByID(attachmentID)
		if err != nil {
			return err
		}

		attachment.OwnerType = utils.Ptr("class_material")
		attachment.OwnerID = &classMaterial.ID

		err = u.attachmentRepo.Update(attachment)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *classMaterialUseCase) CreateQuizMaterial(dto dtos.CreateQuizClassMaterialDto) error {
	// 1. Fetch the Quiz to create a snapshot
	quiz, err := u.quizRepo.GetByID(dto.QuizID)
	if err != nil {
		return err
	}

	// 2. Create the ClassMaterial
	material := &entities.ClassMaterial{
		ClassID:     dto.ClassID,
		Type:        "quiz", // Enforce type
		Title:       dto.Title,
		Description: dto.Description,
		Points:      dto.Points,
		DueAt:       dto.DueAt,
		IsPublished: false, // Default to draft; will be updated if Action is "publish"
	}

	if dto.Action == "publish" {
		now := time.Now()
		material.PublishedAt = &now
		material.IsPublished = true
	}

	if err := u.classMaterialRepo.Create(material); err != nil {
		return err
	}

	// 3. Create the QuizGameLog
	// Serialize quiz to JSON for the snapshot
	var questionDtos []dtos.GetQuizQuestionResponse
	for _, q := range quiz.Questions {
		var attachmentDto *dtos.GetAttachmentResponseDto
		atts, err := u.attachmentUseCase.GetAttachmentsByOwner("quiz_question", q.ID)
		if err == nil && len(atts) > 0 {
			attachmentDto = atts[0]
		}

		var optionDtos []dtos.GetQuizOptionResponse
		for _, o := range q.Options {
			optionDtos = append(optionDtos, dtos.GetQuizOptionResponse{
				ID:         o.ID,
				OptionText: o.OptionText,
				IsCorrect:  o.IsCorrect,
				OrderIndex: o.OrderIndex,
			})
		}

		questionDtos = append(questionDtos, dtos.GetQuizQuestionResponse{
			ID:               q.ID,
			QuestionText:     q.QuestionText,
			MediaAttachment:  attachmentDto,
			PointsMultiplier: q.PointsMultiplier,
			TimeLimitSeconds: q.TimeLimitSeconds,
			OrderIndex:       q.OrderIndex,
			Options:          optionDtos,
		})
	}

	snapshotDto := dtos.GetQuizResponseDto{
		ID:                     quiz.ID,
		ClassID:                quiz.ClassID,
		Title:                  quiz.Title,
		Description:            quiz.Description,
		DefaultTimePerQuestion: quiz.DefaultTimePerQuestion,
		Questions:              questionDtos,
	}

	quizSnapshot, err := json.Marshal(snapshotDto)
	if err != nil {
		return err
	}

	gameLog := &entities.QuizGameLog{
		ClassMaterialID: material.ID,
		QuizPin:         "",            // Initially empty/null
		Status:          "not_started", // Initial status
		QuizSnapshot:    datatypes.JSON(quizSnapshot),
	}

	if err := u.quizGameLogRepo.Create(gameLog); err != nil {
		return err
	}

	return nil
}

func (u *classMaterialUseCase) GetMaterialByID(id uint) (*dtos.GetClassMaterialDto, error) {
	material, err := u.classMaterialRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	attachmentDtos, err := u.attachmentUseCase.GetAttachmentsByOwner("class_material", material.ID)
	if err != nil {
		return nil, err
	}

	res := transforms.EntityToGetClassMaterialDtoWithAttachments(material, attachmentDtos)

	res.CreatedBy = dtos.CreatedByDto{
		ID:       material.Class.OwnerID,
		FullName: material.Class.Owner.FirstName + " " + material.Class.Owner.LastName,
		ImgURL:   u.userUseCase.GetProfileImgUrlByUserID(material.Class.OwnerID),
	}

	if material.Type == "quiz" {
		logs, err := u.quizGameLogRepo.GetByClassMaterialID(material.ID)
		if err == nil && len(logs) > 0 {
			log := logs[0]
			res.QuizGameLog = &dtos.QuizGameLogDto{
				ID:              log.ID,
				ClassMaterialID: log.ClassMaterialID,
				QuizPin:         log.QuizPin,
				Status:          log.Status,
				StartedAt:       log.StartedAt,
				FinishedAt:      log.FinishedAt,
				QuizSnapshot:    log.QuizSnapshot,
			}
		}
	}

	return res, nil
}

func (u *classMaterialUseCase) GetMaterialsByClassID(classID uint, userID uint) ([]*dtos.GetThumnailClassMaterialDto, error) {
	// 1. Fetch the class to check ownership
	class, err := u.classRepo.GetByID(classID)
	if err != nil {
		return nil, err
	}

	// 2. Fetch all materials for this class
	materials, err := u.classMaterialRepo.GetByClassID(classID)
	if err != nil {
		return nil, err
	}

	isOwner := class.OwnerID == userID
	var filteredMaterials []*entities.ClassMaterial

	// 3. Filtering Logic
	if isOwner {
		// Teacher sees everything (Drafts + Published)
		filteredMaterials = materials
	} else {
		// Students only see Published items
		for _, m := range materials {
			if m.IsPublished {
				filteredMaterials = append(filteredMaterials, m)
			}
		}
	}

	// 4. Sorting Logic: "Today First" (Latest PublishedAt first)
	sort.Slice(filteredMaterials, func(i, j int) bool {
		// If someone doesn't have a PublishedAt (like a draft), move it to the end
		if filteredMaterials[i].PublishedAt == nil {
			return false
		}
		if filteredMaterials[j].PublishedAt == nil {
			return true
		}

		// Sort by time descending (Today/Newest first)
		return filteredMaterials[i].PublishedAt.After(*filteredMaterials[j].PublishedAt)
	})

	// 5. Transform to DTO
	response := make([]*dtos.GetThumnailClassMaterialDto, 0, len(filteredMaterials))
	for _, material := range filteredMaterials {
		response = append(response, transforms.EntityToGetThumnailClassMaterialDto(material))
	}

	return response, nil
}

func (u *classMaterialUseCase) UpdatePostClassMaterial(dto *dtos.UpdatePostClassMaterialDto, ownerID uint) error {
	material, err := u.classMaterialRepo.GetByID(dto.ClassMaterialID)
	if err != nil {
		return err
	}

	class, err := u.classRepo.GetByID(material.ClassID)
	if err != nil {
		return err
	}

	if class.OwnerID != ownerID {
		return errors.New("only class owner can update class material")
	}

	// Validate that we are updating the correct type
	if material.Type != "lecture" {
		return errors.New("material type mismatch: expected lecture")
	}

	material.Title = dto.Title
	material.Description = dto.Description

	u.handlePublishedState(material, dto.Published)
	u.handleAttachments(dto.ClassMaterialID, dto.AttachmentIds)

	return u.classMaterialRepo.Update(material)
}

func (u *classMaterialUseCase) UpdateAssignmentClassMaterial(dto *dtos.UpdateAssignmentClassMaterialDto, ownerID uint) error {
	material, err := u.classMaterialRepo.GetByID(dto.ClassMaterialID)
	if err != nil {
		return err
	}

	class, err := u.classRepo.GetByID(material.ClassID)
	if err != nil {
		return err
	}

	if class.OwnerID != ownerID {
		return errors.New("only class owner can update class material")
	}

	if material.Type != "assignment" {
		return errors.New("material type mismatch: expected assignment")
	}

	material.Title = dto.Title
	material.Description = dto.Description

	// Only update optional fields if they are provided (not nil)
	if dto.Points != nil {
		material.Points = dto.Points
	}
	if dto.DueAt != nil {
		material.DueAt = dto.DueAt
	}

	u.handlePublishedState(material, dto.Published)
	u.handleAttachments(dto.ClassMaterialID, dto.AttachmentIds)

	return u.classMaterialRepo.Update(material)
}

func (u *classMaterialUseCase) UpdateQuizClassMaterial(dto *dtos.UpdateQuizClassMaterialDto, ownerID uint) error {
	material, err := u.classMaterialRepo.GetByID(dto.ClassMaterialID)
	if err != nil {
		return err
	}

	class, err := u.classRepo.GetByID(material.ClassID)
	if err != nil {
		return err
	}

	if class.OwnerID != ownerID {
		return errors.New("only class owner can update class material")
	}

	if material.Type != "quiz" {
		return errors.New("material type mismatch: expected quiz")
	}

	material.Title = dto.Title
	material.Description = dto.Description

	if dto.Points != nil {
		material.Points = dto.Points
	}

	u.handlePublishedState(material, dto.Published)

	return u.classMaterialRepo.Update(material)
}

func (u *classMaterialUseCase) DeleteClassMaterial(classMaterialID uint, ownerID uint) error {
	material, err := u.classMaterialRepo.GetByID(classMaterialID)
	if err != nil {
		return err
	}

	class, err := u.classRepo.GetByID(material.ClassID)
	if err != nil {
		return err
	}

	if class.OwnerID != ownerID {
		return errors.New("only class owner can delete class material")
	}

	if err := u.handleAttachments(classMaterialID, []uint{}); err != nil {
		return err
	}

	return u.classMaterialRepo.Delete(classMaterialID)
}

// Helper to toggle PublishedAt based on the boolean flag
func (u *classMaterialUseCase) handlePublishedState(material *entities.ClassMaterial, published bool) {
	if published {
		// Only set PublishedAt if it wasn't already published to preserve the original publish date
		if material.PublishedAt == nil {
			material.IsPublished = true
			now := time.Now()
			material.PublishedAt = &now
		}
	} else {
		material.IsPublished = false
		material.PublishedAt = nil
	}
}

func (u *classMaterialUseCase) handleAttachments(materialID uint, attachmentIds []uint) error {
	attachments, err := u.attachmentRepo.GetByOwnerTypeAndOwnerID("class_material", materialID)
	if err == nil {
		for _, att := range attachments {
			att.OwnerID = nil
			att.OwnerType = nil
			if err := u.attachmentRepo.Update(att); err != nil {
				return err
			}
		}
	}

	for _, attachmentID := range attachmentIds {
		attachment, err := u.attachmentRepo.GetByID(attachmentID)
		if err != nil {
			return err
		}

		attachment.OwnerType = utils.Ptr("class_material")
		attachment.OwnerID = &materialID

		err = u.attachmentRepo.Update(attachment)
		if err != nil {
			return err
		}
	}

	return nil
}
