package apis

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/theghostmac/donorwallet/internal/models"
	"github.com/theghostmac/donorwallet/internal/userops"
	"go.uber.org/zap"
)

// InitRouter initializes the Gin router with all routes.
func InitRouter() *gin.Engine {
	router := gin.Default()

	// User routes.
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.GET("/users/:id", getUser)
    router.POST("/login", loginUser)
    router.POST("/donate", makeDonation)
    router.GET("/all-donations", listUserDonations)
    router.GET("/get-donation/:donation_id", getDonation)

	return router
}

func createUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := userops.CreateUser(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

// updateUser handles updating an existing user
func updateUser(c *gin.Context) {
    userID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var updatedUser models.User
    if err := c.ShouldBindJSON(&updatedUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    updatedUser.UserID = userID
    err = userops.UpdateUser(&updatedUser)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, updatedUser)
}

// getUser handles retrieving a user by ID
func getUser(c *gin.Context) {
    userID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    user, err := userops.GetUser(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, user)
}

func loginUser(c *gin.Context) {
    var loginInfo struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&loginInfo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    token, err := userops.LoginUser(loginInfo.Username, loginInfo.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}

func makeDonation(c *gin.Context) {
    var donationInfo struct {
        BeneficiaryID uuid.UUID `json:"beneficiary_id"`
        Amount float64 `json:"amount"`
        Message string `json:"message"`
    }

    if err := c.ShouldBindJSON(&donationInfo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Retrieve the donor ID from JWT token
    donorID, err := getUserIDFromToken(c)
    if err != nil {
        logger.Error("failed to extract user ID from token", zap.Error(err))
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    logger.Info("Donor ID extracted", zap.String("donorID", donorID.String()))

    donationID, err := userops.Donate(donorID, donationInfo.BeneficiaryID, donationInfo.Amount, donationInfo.Message)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Donation successful", "donation_id": donationID})
}

func getDonation(c *gin.Context) {
    donationID, _ := strconv.ParseUint(c.Param("donation_id"), 10, 32)
    donation, err := userops.GetDonation(uint(donationID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Donation not found"})
        return
    }

    c.JSON(http.StatusOK, donation)
}

// === Pagination ===

func listUserDonations(c *gin.Context) {
   // Retrieve user ID frwom JWT token.
   userID, err := getUserIDFromToken(c)
   if err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
    return
   }

   page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
   limit, _ := strconv.Atoi(c.DefaultQuery("limit",  "10"))

   // Parse start and end dates from query parameters.
   startDateStr := c.DefaultQuery("start_date", "")
   endDateStr := c.DefaultQuery("end_date", "")

   // Default to a wide range if not specified.
   startDate, _ := time.Parse("2006-01-02", startDateStr)
   if startDateStr == "" {
    startDate = time.Time{} // zero time is the most distant past time.
   }

   endDate, _ := time.Parse("2006-01-02", endDateStr)
   if endDateStr == "" {
    endDate = time.Now() // current time, if no end date is provided.
   }

    donations, err := userops.GetUserDonations(userID, page, limit, startDate, endDate) // TODO: fix the time for End date.
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data": donations,
        "page": page,
        "limit": limit,
    })
}