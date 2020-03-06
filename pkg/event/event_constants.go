package credstashsecret

const (
	// Event types
	TypeNormal  = "Normal"
	TypeWarning = "Warning"

	// Event Reasons
	ReasonErrFetchingCredstashSecret = "ErrCredstash"

	// Event Messages
	MessageFailedFetchingCredstashSecret = "Failed fetching credstash secret. Key: %s. Table: %s. Version: %s. Error %s"
)
