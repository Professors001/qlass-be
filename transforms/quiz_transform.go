package transforms

import (
	"qlass-be/domain/entities"
	"qlass-be/dtos"
)

func SaveQuizDtoToQuizEntity(dto dtos.SaveQuizDto, classID uint) *entities.Quiz {
	return &entities.Quiz{
		Title:                  dto.Title,
		ClassID:                classID,
		Description:            dto.Description,
		DefaultTimePerQuestion: dto.DefaultTimePerQuestion,
	}
}

func SaveQuizQuestionDtoToQuizQuestionEntity(dto dtos.SaveQuizQuestionDto, quizID uint) *entities.QuizQuestion {
	return &entities.QuizQuestion{
		QuizID:           quizID,
		QuestionText:     dto.QuestionText,
		PointsMultiplier: dto.PointsMultiplier,
		TimeLimitSeconds: dto.TimeLimitSeconds,
		OrderIndex:       dto.OrderIndex,
	}
}

func SaveQuizOptionDtoToQuizOptionEntity(dto dtos.SaveQuizOption, questionID uint) *entities.QuizOption {
	return &entities.QuizOption{
		QuestionID: questionID,
		OptionText: dto.OptionText,
		IsCorrect:  dto.IsCorrect,
		OrderIndex: dto.OrderIndex,
	}

}
