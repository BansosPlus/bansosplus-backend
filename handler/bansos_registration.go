package handler

import (
	"net/http"
	"strconv"
	"strings"
	"fmt"
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/shopspring/decimal"

	"github.com/BansosPlus/bansosplus-backend.git/model"
	"github.com/BansosPlus/bansosplus-backend.git/repository"
)

type BansosRegistrationHandler struct {
	bansosRegistrationRepository repository.BansosRegistrationRepository
	apiMlUrl string
}

func NewBansosRegistrationHandler(bansosRegistrationRepository repository.BansosRegistrationRepository, apiMlUrl string) *BansosRegistrationHandler {
	return &BansosRegistrationHandler{
		bansosRegistrationRepository: bansosRegistrationRepository,
		apiMlUrl: apiMlUrl,
	}
}

func valueToIndex(value string, array []string) int {
	for i, v := range array {
		if strings.EqualFold(value, v) { // Case-insensitive comparison
			return i
		}
	}
	return -1 // Return -1 if the value is not found in the array
}

func generateQueryParams(bansosRegistrationID int, penghasilan, jumlahMakan, berobat, tanggungan, bahanBakar, jumlahAset, luasLantai, jenisDinding, pendidikan string) string {
	penghasilanArray := []string{"<500 ribu", "500 ribu-1 juta", "1 juta-1.5 juta", ">1.5 juta"}
	luasLantaiArray := []string{"Diatas 8m²", "Dibawah 8m²"}
	kualitasDindingArray := []string{"Buruk", "Normal", "Bagus"}
	jumlahMakanArray := []string{"0", "1", "2", "3"}
	bahanBakarArray := []string{"Kayu/Arang", "Gas/LPG"}
	pendidikanArray := []string{"SD", "SMP", "SMA", "Sarjana"}
	asetArray := []string{"<500 ribu", "500 ribu-1 juta", "1 juta-1.5 juta", ">1.5 juta"}
	berobatArray := []string{"Mampu", "Tidak Mampu"}
	tanggunganArray := []string{"0", "1", "2", ">2"}

	penghasilanIndex := valueToIndex(penghasilan, penghasilanArray)
	jumlahMakanIndex := valueToIndex(jumlahMakan, jumlahMakanArray)
	berobatIndex := valueToIndex(berobat, berobatArray)
	tanggunganIndex := valueToIndex(tanggungan, tanggunganArray)
	bahanBakarIndex := valueToIndex(bahanBakar, bahanBakarArray)
	jumlahAsetIndex := valueToIndex(jumlahAset, asetArray)
	luasLantaiIndex := valueToIndex(luasLantai, luasLantaiArray)
	jenisDindingIndex := valueToIndex(jenisDinding, kualitasDindingArray)
	pendidikanIndex := valueToIndex(pendidikan, pendidikanArray)

	queryParams := fmt.Sprintf("bansos_registration_id=%d&penghasilan=%d&jumlah_makan=%d&berobat=%d&tanggungan=%d&bahan_bakar=%d&jumlah_aset=%d&luas_lantai=%d&jenis_dinding=%d&pendidikan=%d",
		bansosRegistrationID, penghasilanIndex, jumlahMakanIndex, berobatIndex, tanggunganIndex, bahanBakarIndex, jumlahAsetIndex, luasLantaiIndex, jenisDindingIndex, pendidikanIndex)

	return queryParams
}

func processBansosRegistrationToMLAPI(url string) {
	response, err := http.Post(url, "application/json", nil)
	if err != nil {
		fmt.Println("Failed to make request to ML API:", err)
		return
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		fmt.Println("Failed to get a successful response from ML API. Status code:", response.StatusCode)

		// Read the error message from the response body
		errorBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Failed to read the error response body:", err)
			return
		}

		fmt.Println("Error message:", string(errorBody))
		return
	}

	fmt.Println("ML request succeeded!")
}

func (h *BansosRegistrationHandler) RegisterBansosHandler(c echo.Context) error {
	var bansosRegistration model.BansosRegistration

	token, ok := c.Get("token").(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"code":    http.StatusUnauthorized,
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	userID, ok := token["id"].(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"code":    http.StatusUnauthorized,
			"status":  "error",
			"message": "Invalid token format",
		})
	}

	userRole, ok := token["role"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"code":    http.StatusUnauthorized,
			"status":  "error",
			"message": "Invalid token format",
		})
	}

	bansosRegistration.UserID = int(userID)

	if err := c.Bind(&bansosRegistration); err != nil || bansosRegistration.BansosID == 0 || bansosRegistration.Name == "" || bansosRegistration.Nik == "" || bansosRegistration.NoKK == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"status":  "error",
			"message": "Invalid request payload",
		})
	}

	if err := h.bansosRegistrationRepository.RegisterBansos(&bansosRegistration); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Failed to do bansos registeration",
		})
	}

	if userRole == "admin" {
		if err := h.bansosRegistrationRepository.AcceptBansosRegis(&bansosRegistration, decimal.New(1, 0)); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"code":    http.StatusInternalServerError,
				"status":  "error",
				"message": "Failed to accept registration",
			})
		}
	} else {
		// Make a request to ML API with parameters asynchronously
    queryParams := generateQueryParams(
			bansosRegistration.ID, 
			bansosRegistration.Income,
			bansosRegistration.NumberOfMeals,
			bansosRegistration.Treatment,
			bansosRegistration.NumberOfDependents,
			bansosRegistration.Fuel,
			bansosRegistration.TotalAsset,
			bansosRegistration.FloorArea,
			bansosRegistration.WallQuality,
			bansosRegistration.Education,
		)

    apiMlUrl := fmt.Sprintf("%s?%s", h.apiMlUrl, queryParams)
		fmt.Println("Sending request to ML API:", apiMlUrl)

    go processBansosRegistrationToMLAPI(apiMlUrl)
	}

	// Success
	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"status":  "success",
		"message": "Bansos registration added successfully",
		"data": echo.Map{
			"id":        bansosRegistration.ID,
			"user_id":   bansosRegistration.UserID,
			"bansos_id": bansosRegistration.BansosID,
			"status":    bansosRegistration.Status,
		},
	})
}

func (h *BansosRegistrationHandler) AcceptBansosRegisHandler(c echo.Context) error {

	bansosRegisIDStr := c.QueryParam("bansos_registration_id")
	pointStr := c.QueryParam("point")

	bansosRegisID, err := strconv.Atoi(bansosRegisIDStr)
	point, err := decimal.NewFromString(pointStr)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"status":  "error",
			"message": "Invalid bansos_id parameter",
		})
	}

	bansosRegistration, err := h.bansosRegistrationRepository.GetBansosRegisByID(bansosRegisID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Registration Not Found",
		})
	}

	if err := h.bansosRegistrationRepository.AcceptBansosRegis(bansosRegistration, point); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Failed to accept registration",
		})
	}
	// Success
	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"status":  "success",
		"message": "Registration accepted successfully",
		"data": echo.Map{
			"bansos_registration_id": bansosRegistration.ID,
			"user_id":                bansosRegistration.UserID,
			"bansos_id":              bansosRegistration.BansosID,
			"status":                 bansosRegistration.Status,
		},
	})
}

func (h *BansosRegistrationHandler) RejectBansosRegisHandler(c echo.Context) error {

	bansosRegisIDStr := c.QueryParam("bansos_registration_id")
	pointStr := c.QueryParam("point")

	bansosRegisID, err := strconv.Atoi(bansosRegisIDStr)
	point, err := decimal.NewFromString(pointStr)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"status":  "error",
			"message": "Invalid bansos_id parameter",
		})
	}

	bansosRegistration, err := h.bansosRegistrationRepository.GetBansosRegisByID(bansosRegisID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Registration Not Found",
		})
	}

	if err := h.bansosRegistrationRepository.RejectBansosRegis(bansosRegistration, point); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Failed to reject registration",
		})
	}
	// Success
	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"status":  "success",
		"message": "Registration rejected successfully",
		"data": echo.Map{
			"bansos_registration_id": bansosRegistration.ID,
			"user_id":                bansosRegistration.UserID,
			"bansos_id":              bansosRegistration.BansosID,
			"status":                 bansosRegistration.Status,
		},
	})
}

func (h *BansosRegistrationHandler) ValidateBansosRegisHandler(c echo.Context) error {
	_, ok := c.Get("token").(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"code":    http.StatusUnauthorized,
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	bansosRegisIDStr := c.QueryParam("bansos_registration_id")

	bansosRegisID, err := strconv.Atoi(bansosRegisIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"status":  "error",
			"message": "Invalid bansos_id parameter",
		})
	}

	bansosRegistration, err := h.bansosRegistrationRepository.GetBansosRegisByID(bansosRegisID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Registration Not Found",
		})
	}

	if err := h.bansosRegistrationRepository.ValidateBansosRegis(bansosRegistration); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Failed to validate registration",
		})
	}
	// Success
	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"status":  "success",
		"message": "Registration rejected successfully",
		"data": echo.Map{
			"bansos_registration_id": bansosRegistration.ID,
			"user_id":                bansosRegistration.UserID,
			"bansos_id":              bansosRegistration.BansosID,
			"status":                 "TAKEN",
		},
	})
}

func (h *BansosRegistrationHandler) GetOnProgressBansosRegisHandler(c echo.Context) error {
	token, ok := c.Get("token").(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"code":    http.StatusUnauthorized,
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	userRole, ok := token["role"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"code":    http.StatusUnauthorized,
			"status":  "error",
			"message": "Invalid token format",
		})
	}

	if userRole == "admin" {
		registrations, err := h.bansosRegistrationRepository.GetBansosRegisByStatus("ON_PROGRESS")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"code":    http.StatusInternalServerError,
				"status":  "error",
				"message": "Failed to retrieve bansos registrations",
			})
		}

		// Success
		return c.JSON(http.StatusOK, echo.Map{
			"code":    http.StatusOK,
			"status":  "success",
			"message": "Bansos registrations retrieved successfully",
			"data":    registrations,
		})
	}

	return c.JSON(http.StatusUnauthorized, echo.Map{
		"code":    http.StatusUnauthorized,
		"status":  "error",
		"message": "Unauthorized",
	})
}

func (h *BansosRegistrationHandler) GetBansosRegisByUserIDHandler(c echo.Context) error {

	token, ok := c.Get("token").(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"code":    http.StatusUnauthorized,
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	userID, ok := token["id"].(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"code":    http.StatusUnauthorized,
			"status":  "error",
			"message": "Invalid token format",
		})
	}

	statusStr := c.QueryParam("status")
	statusValues := strings.Split(statusStr, ",")

	bansosRegistrations, err := h.bansosRegistrationRepository.GetBansosRegisByUserID(int(userID), statusValues)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Failed to retrieve bansos registration",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"status":  "success",
		"message": "Bansos registration retrieved successfully",
		"data":    bansosRegistrations,
	})
}

func (h *BansosRegistrationHandler) GetBansosRegisByBansosIDHandler(c echo.Context) error {

	bansosIDStr := c.QueryParam("bansos_id")

	bansosID, err := strconv.Atoi(bansosIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"status":  "error",
			"message": "Invalid bansos_id parameter",
		})
	}

	token, ok := c.Get("token").(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"code":    http.StatusUnauthorized,
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	userRole, ok := token["role"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"code":    http.StatusUnauthorized,
			"status":  "error",
			"message": "Invalid token format",
		})
	}

	if userRole == "admin" {
		bansosRegistrations, err := h.bansosRegistrationRepository.GetBansosRegisByBansosID(bansosID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"code":    http.StatusInternalServerError,
				"status":  "error",
				"message": "Failed to retrieve bansos registrations",
			})
		}

		// Success
		return c.JSON(http.StatusOK, echo.Map{
			"code":    http.StatusOK,
			"status":  "success",
			"message": "Bansos registrations retrieved successfully",
			"data":    bansosRegistrations,
		})
	}

	return c.JSON(http.StatusUnauthorized, echo.Map{
		"code":    http.StatusUnauthorized,
		"status":  "error",
		"message": "Unauthorized",
	})
}

func (h *BansosRegistrationHandler) GetBansosRegisByIDHandler(c echo.Context) error {

	bansosRegisIDStr := c.QueryParam("bansos_registration_id")

	bansosRegisID, err := strconv.Atoi(bansosRegisIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"status":  "error",
			"message": "Invalid bansos_id parameter",
		})
	}

	_, ok := c.Get("token").(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"code":    http.StatusUnauthorized,
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	bansosRegistrations, err := h.bansosRegistrationRepository.GetBansosRegisByID(bansosRegisID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Failed to retrieve bansos registrations",
		})
	}

	// Success
	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"status":  "success",
		"message": "Bansos registrations retrieved successfully",
		"data":    bansosRegistrations,
	})
}
