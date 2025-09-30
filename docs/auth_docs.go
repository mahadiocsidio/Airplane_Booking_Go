package docs

// ===================== Register =====================
//create doc for register user
// @Summary Register a new user
// @Description Register a new account with name, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param register body RegisterRequest true "User registration request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func RegisterUserDoc() {}

// ===================== Login =====================

// @Summary Login user
// @Description Login with email and password, returns JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "User login request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
func LoginUserDoc() {}
