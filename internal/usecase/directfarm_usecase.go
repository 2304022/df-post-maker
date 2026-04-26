package usecase

import (
	"bytes"
	"df-post-maker/internal/dto"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

const directFarmAuthURL = "https://direct.farm/api/security/v2/auth/signin"
const directFarmPostAddURL = "https://direct.farm/api/resource/v2/post/add"
const directFarmUploadURL = "https://direct.farm/upload"

type DirectFarmUseCase struct{}

func NewDirectFarmUseCase() *DirectFarmUseCase {
	return &DirectFarmUseCase{}
}

type directFarmAuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

func (u *DirectFarmUseCase) Auth(req dto.AuthRequestDto) (*dto.AuthResponseDto, error) {
	authURL := os.Getenv("DIRECT_FARM_AUTH_URL")
	if authURL == "" {
		authURL = directFarmAuthURL
	}

	body, err := json.Marshal(map[string]string{
		"login":    req.Login,
		"password": req.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(authURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to call direct.farm auth: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("direct.farm auth failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var authResp directFarmAuthResponse
	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return nil, fmt.Errorf("failed to parse auth response: %w", err)
	}

	return &dto.AuthResponseDto{
		Token:        authResp.Token,
		RefreshToken: authResp.RefreshToken,
	}, nil
}

func (u *DirectFarmUseCase) CreatePost(req dto.CreatePostRequestDto, authHeader string) (*dto.CreatePostResponseDto, error) {
	postURL := os.Getenv("DIRECT_FARM_POST_ADD_URL")
	if postURL == "" {
		postURL = directFarmPostAddURL
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, postURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call direct.farm post add: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("direct.farm post add failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var postResp dto.CreatePostResponseDto
	if err := json.Unmarshal(respBody, &postResp); err != nil {
		return nil, fmt.Errorf("failed to parse post response: %w", err)
	}

	return &postResp, nil
}

func (u *DirectFarmUseCase) Upload(file multipart.File, header *multipart.FileHeader, authHeader string) (*dto.UploadResponseDto, error) {
	uploadURL := os.Getenv("DIRECT_FARM_UPLOAD_URL")
	if uploadURL == "" {
		uploadURL = directFarmUploadURL
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err = io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}
	if err = writer.WriteField("reltype", "post"); err != nil {
		return nil, fmt.Errorf("failed to write reltype field: %w", err)
	}
	writer.Close()

	httpReq, err := http.NewRequest(http.MethodPost, uploadURL, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	httpReq.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call direct.farm upload: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("direct.farm upload failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var uploadResp dto.UploadResponseDto
	if err := json.Unmarshal(respBody, &uploadResp); err != nil {
		return nil, fmt.Errorf("failed to parse upload response: %w", err)
	}

	return &uploadResp, nil
}
