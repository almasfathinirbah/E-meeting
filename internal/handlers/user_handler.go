package handlers

import (
	"e_meeting/internal/models"
	"e_meeting/internal/services"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// UserHandler handles HTTP requests related to user management
type UserHandler struct {
	userService services.UserService // Service layer for user-related operations
}

// NewUserHandler creates a new instance of UserHandler
// Parameters:
//   - userService: Service implementation for user operations
//
// Returns:
//   - A pointer to the new UserHandler instance
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register handles user registration requests
// It creates a new user account with the provided registration details
// Parameters:
//   - c: The Fiber context containing the HTTP request and response
//
// Returns:
//   - error: Any error that occurred during registration
func (h *UserHandler) Register(c *fiber.Ctx) error {
	req := c.Locals("request").(models.RegisterRequest)

	user, err := h.userService.Register(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to register user")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to register user " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.RegisterResponse{
		Message: "User registered successfully",
		UserID:  user.ID,
	})
}

// Login handles user authentication requests
// It validates credentials and returns a JWT token upon successful authentication
// Parameters:
//   - c: The Fiber context containing the HTTP request and response
//
// Returns:
//   - error: Any error that occurred during login
func (h *UserHandler) Login(c *fiber.Ctx) error {
	req := c.Locals("request").(models.LoginRequest)

	token, userId, err := h.userService.Login(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to login")
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(models.LoginResponse{
		UserID: userId,
		Token:  token,
	})
}

// GetProfile handles requests to retrieve a user's profile
// It enforces access control to ensure users can only access their own profile or have admin rights
// Parameters:
//   - c: The Fiber context containing the HTTP request and response
//
// Returns:
//   - error: Any error that occurred while fetching the profile
//   - Returns specific status codes for different error conditions:
//   - 403 Forbidden: When trying to access another user's profile without admin rights
//   - 404 Not Found: When the requested profile doesn't exist
//   - 400 Bad Request: When the user ID format is invalid
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	// Get authenticated user ID from context
	authUserID, _ := c.Locals("userID").(string)
	requestedID := c.Params("id")

	// Optional: Check if user is requesting their own profile or has admin rights
	isAdmin, _ := c.Locals("isAdmin").(bool)
	if authUserID != requestedID && !isAdmin {
		return c.Status(http.StatusForbidden).JSON(models.ErrorResponse{
			Error: "Forbidden, you can only access your own profile",
		})
	}

	profile, err := h.userService.GetProfile(requestedID)
	if err != nil {
		switch err.Error() {
		case "user not found":
			return c.Status(http.StatusNotFound).JSON(models.ErrorResponse{
				Error: "user not found",
			})
		case "invalid user ID format":
			return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{
				Error: "invalid user ID format",
			})
		default:
			fmt.Printf("Error fetching user profile: %v\n", err)
			return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{
				Error: "Failed to fetch user profile",
			})
		}
	}

	return c.Status(http.StatusOK).JSON(profile)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	// Get authenticated user ID from context
	authUserID, _ := c.Locals("userID").(string)
	requestedID := c.Params("id")

	// Optional: Check if user is requesting their own profile or has admin rights
	if authUserID != requestedID {
		return c.Status(http.StatusForbidden).JSON(models.ErrorResponse{
			Error: "Forbidden",
		})
	}

	req := c.Locals("request").(models.UpdateProfileRequest)

	profile, err := h.userService.UpdateProfile(requestedID, &req)
	if err != nil {
		switch err.Error() {
		case "user not found":
			return c.Status(http.StatusNotFound).JSON(models.ErrorResponse{
				Error: "user not found",
			})
		case "invalid user ID format":
			return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{
				Error: "invalid user ID format",
			})
		default:
			fmt.Printf("Error updating user profile: %v\n", err)
			return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{
				Error: "Failed to update user profile",
			})
		}
	}

	return c.Status(http.StatusOK).JSON(profile)
}
