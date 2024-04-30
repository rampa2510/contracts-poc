package contracts

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/rampa2510/contracts-poc/config"
	"github.com/rampa2510/contracts-poc/internal/api/middleware"
	"github.com/rampa2510/contracts-poc/internal/utils"
)

type ContractsController struct {
	storage   *ContractsStorage
	awsClient *utils.AwsClient
	env       *config.EnvVars
}

func NewContractsController(storage *ContractsStorage, awsClient *utils.AwsClient, env *config.EnvVars) *ContractsController {
	slog.Info("Initalised Contracts controller")
	return &ContractsController{storage: storage, awsClient: awsClient, env: env}
}

type createContractsRequest struct {
	S3Key  string `json:"s3Key" validate:"required"`
	UserId int    `json:"userId" validate:"required"`
}

func (cc *ContractsController) Create(w http.ResponseWriter, r *http.Request) {
	var req createContractsRequest

	err, validationErrors := utils.ValidateRequest(r, &req)

	if err != nil {
		fmt.Println(err)
		panic(&middleware.APIError{
			Message: "Invalid request body",
			Status:  http.StatusBadRequest,
		})
	}

	if len(validationErrors) > 0 {
		utils.SendResponse(w, http.StatusBadRequest, validationErrors)
		return
	}

	url := cc.awsClient.GetPresignedUrl(cc.env.AWS_BUCKET_NAME, req.S3Key)
	insertedId, err := cc.storage.CreateNewContract(req.S3Key, req.UserId)
	if err != nil {
		fmt.Println("Here")
		panic(&middleware.APIError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		})
	}

	response := map[string]string{
		"insertedId": insertedId,
		"url":        url,
	}

	utils.SendResponse(w, http.StatusCreated, response)
}

type replaceSignatureRequest struct {
	ContractKey string `json:"contractKey" validate:"required"`
	ImageKey    string `json:"imageKey" validate:"required"`
}

func (cc *ContractsController) Replace(w http.ResponseWriter, r *http.Request) {
	var req replaceSignatureRequest

	err, validationErrors := utils.ValidateRequest(r, &req)

	if err != nil {
		fmt.Println(err)
		panic(&middleware.APIError{
			Message: "Invalid request body",
			Status:  http.StatusBadRequest,
		})
	}

	if len(validationErrors) > 0 {
		utils.SendResponse(w, http.StatusBadRequest, validationErrors)
		return
	}

	contractS3Resp, err := cc.awsClient.S3Download(cc.env.AWS_BUCKET_NAME, req.ContractKey)
	if err != nil {
		fmt.Println(err)
		panic(&middleware.APIError{
			Message: "Internal Server Error",
			Status:  http.StatusInternalServerError,
		})
	}

	contractInBytes, err := io.ReadAll(contractS3Resp.Body)
	if err != nil {
		fmt.Println(err)
		panic(&middleware.APIError{
			Message: "Internal Server Error",
			Status:  http.StatusInternalServerError,
		})
	}

	uploadErr := cc.awsClient.S3Upload(cc.env.AWS_BUCKET_NAME, "updated-contract.pdf", contractInBytes)
	if uploadErr != nil {
		fmt.Println(uploadErr)
		panic(&middleware.APIError{
			Message: "Internal Server Error",
			Status:  http.StatusInternalServerError,
		})
	}

	response := map[string]string{
		"success": "Ok",
	}
	utils.SendResponse(w, http.StatusCreated, response)
}
