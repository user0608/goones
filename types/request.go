package types

type DeleteRequestUUIDs struct {
	Values UUIDArray `json:"code"`
}

type DeleteRequest struct {
	Values StrArray `json:"code"`
}
