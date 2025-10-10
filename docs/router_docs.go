package docs

// ============================================================================
// Swagger (Swaggo) annotations for Omni Channel API
// Place this file anywhere in your project that gets compiled.
// Then run:  go install github.com/swaggo/swag/cmd/swag@latest && swag init
// In main, remember to import generated docs:  _ "your/module/path/docs"
// And expose UI (Fiber example): app.Get("/swagger/*", fiberSwagger.WrapHandler)
// ============================================================================

// @title           Omni Channel API
// @version         1.0
// @description     REST API for the Omni Channel service.
// @BasePath        /api/v1
// @schemes         http https
//
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization

// =========================
// AUTH
// =========================

// Login godoc
// @Summary      Login
// @Description  Authenticate with email and password to obtain a JWT access token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      model.AuthRequest  true  "Login payload"
// @Success      200      {object}  model.GlobalResponse          "Login success; token returned in data"
// @Failure      400      {object}  model.ErrorValidationResponse "Validation error"
// @Failure      401      {object}  model.GlobalResponse          "Invalid credentials"
// @Failure      500      {object}  model.GlobalResponse          "Server error"
// @Router       /login [post]
func docLogin() {}

// =========================
// USERS (Protected)
// =========================

// CreateUser godoc
// @Summary      Create user
// @Description  Create a new user. Requires a valid JWT and permission **C**.
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      model.UserRequest  true  "User payload"
// @Success      201      {object}  model.GlobalResponse          "User created"
// @Failure      400      {object}  model.ErrorValidationResponse "Validation error"
// @Failure      401      {object}  model.GlobalResponse          "Unauthorized (missing/invalid JWT)"
// @Failure      403      {object}  model.GlobalResponse          "Forbidden (insufficient permission: requires C)"
// @Failure      500      {object}  model.GlobalResponse          "Server error"
// @Router       /users/ [post]
func docCreateUser() {}

// CreateUserMatrix godoc
// @Summary      Create user matrix
// @Description  Create a user permission/role matrix. Requires a valid JWT and permission **A**.
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      model.UserMatrix  true  "UserMatrix payload"
// @Success      201      {object}  model.GlobalResponse
// @Failure      400      {object}  model.ErrorValidationResponse
// @Failure      401      {object}  model.GlobalResponse          "Unauthorized (missing/invalid JWT)"
// @Failure      403      {object}  model.GlobalResponse          "Forbidden (insufficient permission: requires A)"
// @Failure      500      {object}  model.GlobalResponse
// @Router       /users/:id/matrix [post]
func docCreateUserMatrix() {}
