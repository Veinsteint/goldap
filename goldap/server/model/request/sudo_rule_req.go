package request

// SudoRuleAddReq create sudo rule request
type SudoRuleAddReq struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=255"`
	User        string `json:"user" validate:"required"`    // user or group (%group1)
	Host        string `json:"host" validate:"required"`    // ALL or CIDR
	Command     string `json:"command" validate:"required"` // ALL or command path
	RunAsUser   string `json:"runAsUser"`                   // run as user (ALL or root)
	RunAsGroup  string `json:"runAsGroup"`                  // run as group
	Options     string `json:"options"`                     // NOPASSWD, NOEXEC, etc.
}

// SudoRuleUpdateReq update sudo rule request
type SudoRuleUpdateReq struct {
	ID          uint   `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=255"`
	User        string `json:"user" validate:"required"`
	Host        string `json:"host" validate:"required"`
	Command     string `json:"command" validate:"required"`
	RunAsUser   string `json:"runAsUser"`
	RunAsGroup  string `json:"runAsGroup"`
	Options     string `json:"options"`
}

// SudoRuleDeleteReq delete sudo rule request
type SudoRuleDeleteReq struct {
	ID uint `json:"id" validate:"required"`
}
