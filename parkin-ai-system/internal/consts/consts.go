package consts

const (
	RoleUser  = "user"
	RoleAdmin = "admin"

	TokenTypeAccess  = "token_access"
	TokenTypeRefresh = "token_refresh"
)

var ValidSlotTypes = []string{
	"compact",
	"large",
	"handicapped",
	"electric",
}
