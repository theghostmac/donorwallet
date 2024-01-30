package userops

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	sendgrid "github.com/theghostmac/donorwallet/external"
	"github.com/theghostmac/donorwallet/internal/database"
	"github.com/theghostmac/donorwallet/internal/jwtauth"
	"github.com/theghostmac/donorwallet/internal/models"
	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()

// CreateUser inserts a new user into the database.
func CreateUser(user *models.User) error {
	// Generate a new UUID for the user.
	user.UserID = uuid.New()

	hashedPassword, err := HashPassword(user.PasswordHash)
	if err != nil {
		logger.Fatal("Failed to hash password.", zap.Error(err))
		return err
	}

	user.PasswordHash = hashedPassword

	// First, create the user record in the database.
    if err := database.DB.Create(user).Error; err != nil {
        logger.Error("error creating user record.", zap.Error(err))
        return err
    }

	// Now that the user exists, create a wallet for them.
    wallet := models.Wallet{
        UserID: user.UserID, // Use the UUID generated for the user
        Balance: 100.0,        // Initial balance as 100 because of testing /donate endpoint.
        CreatedAt: time.Now(),
    }

	// Set the wallet balance in the user response
    user.WalletBalance = wallet.Balance

	if err := database.DB.Create(&wallet).Error; err != nil {
        logger.Error("error creating a wallet for the user.", zap.Error(err))
        return err
    }

	return nil
}

// GetUser retrieves a user by their ID.
func GetUser(userID uuid.UUID) (*models.User, error) {
	var user models.User
	err := database.DB.Where("user_id =?", userID).First(&user).Error
	if err!= nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user in the database.
func UpdateUser(user *models.User) error {
	return database.DB.Save(user).Error
}

// DeleteUser removes a user from the database.
func DeleteUser(userID uuid.UUID) error {
	return database.DB.Where("user_id =?", userID).Delete(&models.User{}).Error
}

// LoginUser authenticates a user and returns a JWT token.
func LoginUser(username, password string) (string, error) {
	var user models.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		logger.Error("failed to find this user in the database: ", zap.Error(err))
		return "", err
	}

	// Check if the password is correct.
	if !CheckPasswordHash(password, user.PasswordHash) {
		logger.Error("error with the provided and saved passwords: ", zap.Error(err))
		return "", errors.New("invalid password")
	}

	// Generate JWT token
	token, err := jwtauth.GenerateToken(user.UserID)
	if err != nil {
		logger.Error("error generating JWT token for the user: ", zap.Error(err))
		return "", err
	}

	return token, nil
}

// Donate allows a user to donate money from their wallet to another user.
func Donate(donorID, beneficiaryID uuid.UUID, amount float64, message string) (uint, error) {
	// Ensure donor is not donating to themselves.
	if donorID == beneficiaryID {
		return 0, errors.New("cannot donate to yourself")
	}

	donation := models.Donation{
        DonorID: donorID,
        BeneficiaryID: beneficiaryID,
        Amount: amount,
        Message: message,
        CreatedAt: time.Now(),
    }

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Get donor's wallet
		var donorWallet models.Wallet
		if err := tx.Where("user_id = ?", donorID).First(&donorWallet).Error; err != nil {
			logger.Error("error fetching the donor's wallet: ", zap.Error(err))
            return err
        }

		// Check the balance.
		if donorWallet.Balance < amount {
			return errors.New("insufficient balance in wallet")
		}

		// Deduct amount from donor's wallet
		donorWallet.Balance -= amount
		if err := tx.Save(&donorWallet).Error; err != nil {
			logger.Error("error deducting from the donor's wallet: ", zap.Error(err))
			return err
		}
	
		// Get beneficiary's wallet
		var beneficiaryWallet models.Wallet
		if err := tx.Where("user_id = ?", beneficiaryID).First(&beneficiaryWallet).Error; err != nil {
			logger.Error("error fetching the beneficiary's wallet: ", zap.Error(err))
			return err
		}

		// Add amount to beneficiary's wallet.
		beneficiaryWallet.Balance += amount
		if err := tx.Save(&beneficiaryWallet).Error; err != nil {
			logger.Error("error funding the the beneficiary's wallet: ", zap.Error(err))
			return err
		}

		// Record donation.
		if err := tx.Create(&donation).Error; err != nil {
			logger.Error("error saving the donation credentials: ", zap.Error(err))
			return err
		}

		logger.Info("Donation created", zap.Uint("donationID", donation.DonationID))

		return nil
	})

	return donation.DonationID, err
}

// CountUserDonations counts the number of donations a user has done.
func CountUserDonations(userID uuid.UUID) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Donation{}).Where("donor_i =  ?", userID).Count(&count).Error
	return count, err
}

// SendThankYouMessage sends a thank you message if user equals 2 or more donations count.
func SendThankYouMessage(userID uuid.UUID) {
	count, err := CountUserDonations(userID)
	if err != nil || count < 2 {
		return
	}

	userEmail, err := GetUserEmail(userID)
	if err != nil {
		logger.Error("error fetching user's email: ", zap.Error(err))
		return
	}

	subject := "Thank You for Your Donations!"
	message := "We really appreciate your support. Thank you for your generous donations!"

	err = sendgrid.SendEmail(userEmail, subject, message)
	if err != nil {
		logger.Error("failed to send thank you email", zap.Error(err))
	} else {
		logger.Info("Thank you email sent to user", zap.String("userID", userID.String()))
	}
}

// GetUserDonations accepts start and end dates and views all user's donations.
func GetUserDonations(userID uuid.UUID, page, limit int, startDate, endDate time.Time) ([]models.Donation, error) {
    var donations []models.Donation
    offset := (page - 1) * limit
    err := database.DB.Where("donor_id = ? AND created_at >= ? AND created_at <= ?", userID, startDate, endDate).Offset(offset).Limit(limit).Find(&donations).Error
    return donations, err
}

// GetDonation retrieves a single donation.
func GetDonation(donationID uint) (*models.Donation, error) {
	var donation models.Donation
	err := database.DB.Where("donation_id = ?", donationID).First(&donation).Error
    return &donation, err
}

func GetUserEmail(userID uuid.UUID) (string, error) {
	var user models.User
	err := database.DB.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		logger.Error("failed to find this user in the database: ", zap.Error(err))
		return "", err
	}
	return user.Email, nil
}