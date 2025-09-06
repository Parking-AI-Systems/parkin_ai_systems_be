package consts

const (
	RoleUser  = "user"
	RoleAdmin = "admin"

	TokenTypeAccess  = "token_access"
	TokenTypeRefresh = "token_refresh"
)

const (
	SlotTypeStandard = "standard"
	SlotTypeDisabled = "disabled"
	SlotTypeElectric = "electric"
	SlotTypeVIP      = "VIP"
)

var ValidSlotTypes = []string{SlotTypeStandard, SlotTypeDisabled, SlotTypeElectric, SlotTypeVIP}

const (
	VehicleTypeCar       = "car"
	VehicleTypeMotorbike = "motorbike"
	VehicleTypeTruck     = "truck"
)

var ValidVehicleTypes = []string{VehicleTypeCar, VehicleTypeMotorbike, VehicleTypeTruck}
