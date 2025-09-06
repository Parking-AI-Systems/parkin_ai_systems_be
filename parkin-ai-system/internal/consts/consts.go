package consts

const (
	TokenTypeAccess  = "token_access"
	TokenTypeRefresh = "token_refresh"
)
const (
	RoleAdmin = "role_admin"
	RoleUser  = "role_user"
)

var ValidRoles = []string{
	RoleAdmin,
	RoleUser,
}

const (
	VehicleTypeCar       = "car"
	VehicleTypeMotorbike = "motorbike"
	VehicleTypeTruck     = "truck"
)

var ValidVehicleTypes = []string{VehicleTypeCar, VehicleTypeMotorbike, VehicleTypeTruck}

// Slot types from parking_slots
const (
	SlotTypeStandard = "standard"
	SlotTypeDisabled = "disabled"
	SlotTypeElectric = "electric"
	SlotTypeVIP      = "VIP"
)

var ValidSlotTypes = []string{SlotTypeStandard, SlotTypeDisabled, SlotTypeElectric, SlotTypeVIP}

// Vehicle-Slot compatibility map
var VehicleSlotCompatibility = map[string][]string{
	VehicleTypeCar:       {SlotTypeStandard, SlotTypeDisabled, SlotTypeVIP},
	VehicleTypeMotorbike: {SlotTypeStandard},
	VehicleTypeTruck:     {SlotTypeStandard, SlotTypeVIP},
}

const (
	TransactionTypeDeposit    = "deposit"
	TransactionTypeWithdrawal = "withdrawal"
	TransactionTypePayment    = "payment"
	TransactionTypeRefund     = "refund"
)

var ValidTransactionTypes = []string{
	TransactionTypeDeposit,
	TransactionTypeWithdrawal,
	TransactionTypePayment,
	TransactionTypeRefund,
}
