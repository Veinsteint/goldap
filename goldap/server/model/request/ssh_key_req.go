package request

// SSHKeyAddReq add SSH public key request
type SSHKeyAddReq struct {
	Title string `json:"title" validate:"required,min=1,max=100"`
	Key   string `json:"key" validate:"required,min=50"`
}

// SSHKeyListReq list SSH keys request
type SSHKeyListReq struct{}

// SSHKeyDeleteReq delete SSH key request
type SSHKeyDeleteReq struct {
	ID uint `uri:"id" json:"id" validate:"required"`
}
