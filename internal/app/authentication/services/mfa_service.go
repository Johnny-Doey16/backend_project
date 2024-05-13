package services

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"io"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	ut "github.com/steve-mir/diivix_backend/internal/app/authentication/utils"
	"github.com/steve-mir/diivix_backend/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO: Refactor - remove unnecessary return data
// userID is the user's uid, issuer is the website (organization) issuing the otp
// 2.[]string the recovery codes
// 3. []byte the qr code incase the image will be sent or display in the backend, 4. error the error
// return 4. the secret, 2. The url (to be used in the frontend to generate the qr code)
func RegisterMFA(ctx context.Context, config utils.Config, store *sqlc.Store, userID uuid.UUID, issuer, password string) (string, []byte, string, error) {
	user, err := store.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", []byte{}, "", status.Error(codes.NotFound, "user with ID "+userID.String()+" not found")
		}
		return "", []byte{}, "", err
	}

	// Check password
	err = ut.CheckPassword(password, user.PasswordHash)
	if err != nil {
		return "", []byte{}, "", status.Error(codes.Unauthenticated, "password incorrect")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: userID.String(),
	})
	if err != nil {
		return "", []byte{}, "", status.Errorf(codes.Internal, err.Error())
	}

	secret := key.Secret()
	url := key.URL()

	// Convert TOTP key into a QR code encoded as a PNG image.
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return "", []byte{}, "", status.Errorf(codes.Internal, err.Error())
	}
	png.Encode(&buf, img)

	// Generate "recovery codes" for the user
	/*
		recoveryCodes, err := generateRecoveryCodes(10, 12) // generates 10 codes, each 12 characters long
		if err != nil {
			return "", []string{}, []byte{}, "", status.Errorf(codes.Internal, err.Error())
		}

		// Here you would securely store the secret key with the user's account in the database.
		// Make sure to encrypt the secret (and maybe the recovery codes) before storing it.

		encryptedSecret, err := encryptSecret(config, secret)
		if err != nil {
			return "", []string{}, []byte{}, "", status.Errorf(codes.Internal, err.Error())
		}

		encryptedRecoveryCodes, err := encryptRecoveryCodes(config, recoveryCodes)
		if err != nil {
			return "", []string{}, []byte{}, "", status.Errorf(codes.Internal, err.Error())
		}

		// Begin transaction
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return "", []string{}, []byte{}, "", status.Errorf(codes.Internal, err.Error())
		}
		defer tx.Rollback()

		qtx := store.WithTx(tx)

		// Save the secret to the db
		err = qtx.AddMfaSecret(ctx, sqlc.AddMfaSecretParams{
			UserID:    userID,
			SecretKey: encryptedSecret,
		})
		if err != nil {
			tx.Rollback()
			return "", []string{}, []byte{}, "", status.Errorf(codes.Internal, err.Error())
		}

		// Store the recovery codes to the db.
		err = qtx.AddRecoveryCodes(ctx, sqlc.AddRecoveryCodesParams{
			UserID:  userID,
			Column2: encryptedRecoveryCodes,
		})
		if err != nil {
			tx.Rollback()
			return "", []string{}, []byte{}, "", status.Errorf(codes.Internal, err.Error())
		}
	*/

	return secret, buf.Bytes(), url, nil
}

func ValidateMFAWorks(ctx context.Context, config utils.Config, db *sql.DB, store *sqlc.Store, userID uuid.UUID, pass, secret string) ([]string, error) {

	isValid := totp.Validate(pass, secret)
	if !isValid {
		return []string{}, status.Errorf(codes.Unauthenticated, "OTP is invalid")
	}

	recoveryCodes, err := generateRecoveryCodes(10, 12) // generates 10 codes, each 12 characters long
	if err != nil {
		return []string{}, status.Errorf(codes.Internal, err.Error())
	}

	encryptedSecret, err := encryptSecret(config, secret)
	if err != nil {
		return []string{}, status.Errorf(codes.Internal, err.Error())
	}

	encryptedRecoveryCodes, err := encryptRecoveryCodes(config, recoveryCodes)
	if err != nil {
		return []string{}, status.Errorf(codes.Internal, err.Error())
	}

	// Begin transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return []string{}, status.Errorf(codes.Internal, err.Error())
	}
	defer tx.Rollback()

	qtx := store.WithTx(tx)

	// Save the secret to the db
	err = qtx.AddMfaSecret(ctx, sqlc.AddMfaSecretParams{
		UserID:    userID,
		SecretKey: encryptedSecret,
	})
	if err != nil {
		tx.Rollback()
		return []string{}, status.Errorf(codes.Internal, err.Error())
	}

	// Store the recovery codes to the db.
	err = qtx.AddRecoveryCodes(ctx, sqlc.AddRecoveryCodesParams{
		UserID:  userID,
		Column2: encryptedRecoveryCodes,
	})
	if err != nil {
		tx.Rollback()
		return []string{}, status.Errorf(codes.Internal, err.Error())
	}

	return recoveryCodes, tx.Commit()
}

// Secret is from db while pass is the 6 digit code from the authy app
func ValidateMFA(ctx context.Context, config utils.Config, store *sqlc.Store, userID uuid.UUID, pass string) error {
	// Retrieve the secret key associated with the user's account from the database.
	secret, err := store.GetMfaSecret(ctx, userID)
	if err != nil {
		// fmt.Errorf("MFA validation error %s", err.Error())
		return status.Errorf(codes.Unauthenticated, "OTP not valid")
	}

	// decrypt secret
	decryptedSecret, err := decryptSecret(config, secret)
	if err != nil {
		// fmt.Errorf("MFA decryption error %s", err.Error())
		return status.Errorf(codes.Unauthenticated, "OTP not valid")
	}

	isValid := totp.Validate(pass, decryptedSecret)
	if !isValid {
		return status.Errorf(codes.Unauthenticated, "OTP is invalid")
	}

	return nil
}

// By pass mfa in case they lose their device
func BypassMFA(ctx context.Context, config utils.Config, store *sqlc.Store, userID uuid.UUID, code string) error {
	// Get table id and codes of user that are true(false) and store in codes var
	codesAndID, err := store.GetRecoveryCodes(ctx, userID)
	if err != nil {
		fmt.Println(err.Error())
		return status.Errorf(codes.Unauthenticated, "The code is either incorrect or has been used")
	}

	_, ID, err := validateRecoveryCode(config, code, codesAndID)
	if err != nil {
		fmt.Println(err.Error())
		return status.Errorf(codes.Unauthenticated, "The code is either incorrect or has been used")
	}

	// Using the uid check if the code exists in the db
	// _, err = store.UpdatedRecoveryCodeToUsed(ctx, sqlc.UpdatedRecoveryCodeToUsedParams{
	// 	UserID: userID,
	// 	Code:   code,
	// })
	_, err = store.UpdatedByIdRecoveryCodeToUsed(ctx, sqlc.UpdatedByIdRecoveryCodeToUsedParams{
		ID:     ID,
		UserID: userID,
	})
	if err != nil {
		fmt.Println(err.Error())
		return status.Errorf(codes.Unauthenticated, "The code is either incorrect or has been used")
	}
	// If the code exists make it as used then log the user in
	return nil
}

// GenerateRecoveryCodes generates a specified number of recovery codes.
func generateRecoveryCodes(count int, length int) ([]string, error) {
	if count <= 0 || length <= 0 {
		return nil, errors.New("count and length must be positive")
	}

	codes := make([]string, count)
	for i := 0; i < count; i++ {
		bytes := make([]byte, length)
		if _, err := rand.Read(bytes); err != nil {
			return nil, err
		}
		// Use base32 encoding which is URL safe and doesn't require case sensitivity
		codes[i] = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes)[:length]
	}

	return codes, nil
}

func encryptSecret(config utils.Config, secret string) (string, error) {
	block, err := aes.NewCipher([]byte(config.MfaSecretSymmetricKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(secret), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptSecret(config utils.Config, encryptedSecret string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedSecret)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(config.MfaSecretSymmetricKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// encryptRecoveryCodes takes a slice of recovery codes and a configuration,
// encrypts each code with AES-GCM, and returns a slice of base64-encoded encrypted strings.
func encryptRecoveryCodes(config utils.Config, recoveryCodes []string) ([]string, error) {
	block, err := aes.NewCipher([]byte(config.MfaBackupSymmetricKey))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	encryptedCodes := make([]string, len(recoveryCodes))
	for i, code := range recoveryCodes {
		nonce := make([]byte, gcm.NonceSize())
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, err
		}

		// Encrypt the code using GCM
		encryptedCode := gcm.Seal(nonce, nonce, []byte(code), nil)
		// Encode the encrypted code to base64
		encryptedCodes[i] = base64.StdEncoding.EncodeToString(encryptedCode)
	}

	return encryptedCodes, nil
}

func decryptRecoveryCodes(config utils.Config, encryptedRecoveryCodes string) (string, error) {
	// decrypt encryptedSecret
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedRecoveryCodes)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(config.MfaBackupSymmetricKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func validateRecoveryCode(config utils.Config, submittedCode string, encryptedCodesAndIdInDB []sqlc.GetRecoveryCodesRow) (bool, int32, error) {
	for _, encryptedCode := range encryptedCodesAndIdInDB {
		// Decrypt each stored code
		decryptedCode, err := decryptRecoveryCodes(config, encryptedCode.Code)
		if err != nil {
			return false, 0, err // handle the error properly
		}
		// Check if the decrypted code matches the submitted code
		if decryptedCode == submittedCode {
			fmt.Println("Found", encryptedCode.ID)
			return true, encryptedCode.ID, nil // the submitted code is valid
		}
	}
	// no match found, the submitted code is invalid
	return false, 0, nil
}
