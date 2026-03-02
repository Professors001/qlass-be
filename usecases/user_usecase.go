package usecases

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
	"qlass-be/middleware"
	"qlass-be/transforms"
	"qlass-be/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	RegisterFirstStep(ctx context.Context, req *dtos.RegisterRequestStepOneDto) (*dtos.ResponseRegisterStepOneDto, error)
	RegisterSecondStep(ctx context.Context, req *dtos.RegisterRequestStepTwoDto) (*dtos.ResponseRegisterStepTwoDto, error)
	Login(ctx context.Context, req *dtos.LoginRequestDto) (*dtos.LoginResponseDto, error)
	CreateTeacher(ctx context.Context, req *dtos.CreateTeacherRequestDto) (*dtos.CreateTeacherResponseDto, error)
	UpdateUser(req *dtos.UpdateUserRequestDto, userID uint) (*dtos.UserDisplayData, error)
	ChangePassword(req *dtos.ChangePasswordRequestDto, userID uint) (*dtos.ChangePasswordResponseDto, error)
	ForgetPasswordStep1(ctx context.Context, req *dtos.ForgetPasswordStep1RequestDto) (*dtos.ForgetPasswordStep1ResponseDto, error)
	ForgetPasswordStep2(ctx context.Context, req *dtos.ForgetPasswordStep2RequestDto) (*dtos.ForgetPasswordStep2ResponseDto, error)
}

type userUseCase struct {
	userRepo          repositories.UserRepository
	userCacheRepo     repositories.UserCacheRepository
	jwtService        middleware.JwtService
	emailService      EmailService
	attachmentUseCase AttachmentUseCase
}

func NewUserUseCase(repo repositories.UserRepository, cacheRepo repositories.UserCacheRepository, jwtService middleware.JwtService, email EmailService, attachmentUC AttachmentUseCase) UserUseCase {
	return &userUseCase{
		userRepo:          repo,
		userCacheRepo:     cacheRepo,
		jwtService:        jwtService,
		emailService:      email,
		attachmentUseCase: attachmentUC,
	}
}

func (u *userUseCase) RegisterFirstStep(ctx context.Context, req *dtos.RegisterRequestStepOneDto) (*dtos.ResponseRegisterStepOneDto, error) {
	// Check if registration is already pending in Redis
	if _, err := u.userCacheRepo.GetRegistrationData(ctx, req.Email); err == nil {
		return nil, errors.New("registration pending, please check your email for OTP")
	}

	// Check if user already exists
	existingUserByEmail, _ := u.userRepo.GetByEmail(req.Email)
	if existingUserByEmail != nil {
		return nil, errors.New("Email already exists")
	}

	existingUserByUniID, _ := u.userRepo.GetByUniID(req.UniversityID)
	if existingUserByUniID != nil {
		return nil, errors.New("University ID already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	otp := utils.GenerateRandomString(6)

	tempData := transforms.RequestToTempRegisterDataDto(req, string(hashedPassword), otp)

	err = u.emailService.SendOTP(req.Email, otp)
	if err != nil {
		log.Println("Error sending OTP email:", err)
		return nil, err
	}

	// Store into Redis with Email as key (TTL 5 minutes)
	err = u.userCacheRepo.SetRegistrationData(ctx, req.Email, tempData, 5*time.Minute)
	if err != nil {
		log.Println("Error storing registration data in cache:", err)
		return nil, err
	}

	return &dtos.ResponseRegisterStepOneDto{
		Message:          "Please proceed to step 2",
		RefEmail:         req.Email,
		ExpiresInSeconds: 300,
	}, nil
}

func (u *userUseCase) RegisterSecondStep(ctx context.Context, req *dtos.RegisterRequestStepTwoDto) (*dtos.ResponseRegisterStepTwoDto, error) {
	// Retrieve temp data from Redis
	tempData, err := u.userCacheRepo.GetRegistrationData(ctx, req.Email)
	if err != nil {
		return nil, errors.New("no pending registration found or OTP expired")
	}

	// Validate OTP
	if tempData.OTP != req.OTP {
		return nil, errors.New("invalid OTP")
	}

	// Create user in DB
	newUser := transforms.TempRegisterDataDtoToUserEntity(tempData)

	err = u.userRepo.Create(newUser)
	if err != nil {
		log.Println("Error creating user in DB:", err)
		return nil, err
	}

	// Download and upload profile image
	if len(newUser.FirstName) > 0 && len(newUser.LastName) > 0 {
		initials := fmt.Sprintf("%c%c", newUser.FirstName[0], newUser.LastName[0])
		resp, err := http.Get("https://ui-avatars.com/api/?name=" + initials)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				filename := fmt.Sprintf("profile_%s_%s.png", newUser.FirstName, newUser.LastName)

				var b bytes.Buffer
				w := multipart.NewWriter(&b)

				fw, err := w.CreateFormFile("file", filename)
				if err != nil {
					log.Println("Error creating form file:", err)
					return nil, err
				}

				if _, err = io.Copy(fw, resp.Body); err != nil {
					log.Println("Error copying image data:", err)
					return nil, err
				}
				w.Close()

				req, err := http.NewRequest("POST", "", &b)
				if err != nil {
					log.Println("Error creating dummy request:", err)
					return nil, err
				}
				req.Header.Set("Content-Type", w.FormDataContentType())

				if err := req.ParseMultipartForm(10 << 20); err != nil {
					log.Println("Error parsing multipart form:", err)
					return nil, err
				}

				fileHeader := req.MultipartForm.File["file"][0]

				attachment, err := u.attachmentUseCase.UploadAttachment(newUser.ID, fileHeader)
				if err == nil {
					newUser.ProfileImgAttachmentID = &attachment.AttachmentID
					if err := u.userRepo.Update(newUser); err != nil {
						log.Println("Error updating user with profile image:", err)
					}
				} else {
					log.Println("Error uploading profile image:", err)
				}
			}
		} else {
			log.Println("Error downloading profile image:", err)
		}
	}

	return &dtos.ResponseRegisterStepTwoDto{
		Message: "Registration successful",
	}, nil

}

func (u *userUseCase) Login(ctx context.Context, req *dtos.LoginRequestDto) (*dtos.LoginResponseDto, error) {
	// Retrieve user by email
	user, err := u.userRepo.GetByUniID(req.UniversityID)
	if err != nil {
		return nil, errors.New("This Account is not registered")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("Incorrect password")
	}

	userDisplay := transforms.UserEntityToUserDisplayResponse(user, "")

	if user.ProfileImgAttachmentID != nil {
		attachment, err := u.attachmentUseCase.GetAttachmentByID(*user.ProfileImgAttachmentID)

		if err != nil {
			log.Println("Error getting attachment:", err)
		} else {
			log.Println("Img URL:", attachment.FileURL)
			userDisplay = transforms.UserEntityToUserDisplayResponse(user, attachment.FileURL)
		}
	}

	token, err := u.jwtService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dtos.LoginResponseDto{
		Message: "Login successful",
		Token:   token,
		User:    *userDisplay,
	}, nil
}

func (u *userUseCase) CreateTeacher(ctx context.Context, req *dtos.CreateTeacherRequestDto) (*dtos.CreateTeacherResponseDto, error) {
	// Check if user already exists
	if _, err := u.userRepo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}
	if _, err := u.userRepo.GetByUniID(req.UniversityID); err == nil {
		return nil, errors.New("university ID already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	teacher := &entities.User{
		UniversityID: req.UniversityID,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         "teacher",
		IsVerified:   false,
		IsActive:     true,
	}

	if err := u.userRepo.Create(teacher); err != nil {
		log.Println("Error creating teacher in DB:", err)
		return nil, err
	}

	return &dtos.CreateTeacherResponseDto{
		Message: "Teacher created successfully",
		UserID:  teacher.ID,
	}, nil
}

func (u *userUseCase) UpdateUser(req *dtos.UpdateUserRequestDto, userID uint) (*dtos.UserDisplayData, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}

	if req.ProfileImgAttachmentID != *user.ProfileImgAttachmentID {
		attachmentToDelete := *user.ProfileImgAttachmentID

		user.ProfileImgAttachmentID = &req.ProfileImgAttachmentID

		if err := u.attachmentUseCase.DeleteAttachment(attachmentToDelete); err != nil {
			log.Println("Error deleting attachment:", err)
			return nil, err
		}
	}

	if err := u.userRepo.Update(user); err != nil {
		log.Println("Error updating user:", err)
		return nil, err
	}

	var profileImgURL string
	if user.ProfileImgAttachmentID != nil {
		if attachment, err := u.attachmentUseCase.GetAttachmentByID(*user.ProfileImgAttachmentID); err == nil {
			profileImgURL = attachment.FileURL
		}
	}

	return transforms.UserEntityToUserDisplayResponse(user, profileImgURL), nil
}

func (u *userUseCase) ChangePassword(req *dtos.ChangePasswordRequestDto, userID uint) (*dtos.ChangePasswordResponseDto, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("This Account is not registered")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword))
	if err != nil {
		return nil, errors.New("Incorrect password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = string(hashedPassword)

	if err := u.userRepo.Update(user); err != nil {
		log.Println("Error updating user:", err)
		return nil, err
	}

	return &dtos.ChangePasswordResponseDto{
		Message: "Password changed successfully",
	}, nil
}

func (u *userUseCase) ForgetPasswordStep1(ctx context.Context, req *dtos.ForgetPasswordStep1RequestDto) (*dtos.ForgetPasswordStep1ResponseDto, error) {
	user, err := u.userRepo.GetByUniID(req.UniversityID)
	if err != nil {
		return nil, errors.New("This Account is not registered")
	}

	otp := utils.GenerateRandomString(6)

	err = u.emailService.SendOTP(user.Email, otp)

	tempData := &dtos.TempForgetPasswordData{
		UniversityID: user.UniversityID,
		OTP:          otp,
		User:         *user,
	}

	err = u.userCacheRepo.SetForgetPasswordData(ctx, user.UniversityID, tempData, 5*time.Minute)
	if err != nil {
		log.Println("Error storing registration data in cache:", err)
		return nil, err
	}

	return &dtos.ForgetPasswordStep1ResponseDto{
		Message: "Please proceed to step 2",
		Email:   user.Email,
	}, nil
}

func (u *userUseCase) ForgetPasswordStep2(ctx context.Context, req *dtos.ForgetPasswordStep2RequestDto) (*dtos.ForgetPasswordStep2ResponseDto, error) {
	tempData, err := u.userCacheRepo.GetForgetPasswordData(ctx, req.UniversityID)
	if err != nil {
		return nil, errors.New("no pending registration found or OTP expired")
	}

	// Validate OTP
	if tempData.OTP != req.OTP {
		return nil, errors.New("invalid OTP")
	}

	user := tempData.User

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = string(hashedPassword)

	if err := u.userRepo.Update(&user); err != nil {
		log.Println("Error updating user:", err)
		return nil, err
	}

	return &dtos.ForgetPasswordStep2ResponseDto{
		Message: "Password changed successfully",
	}, nil
}
