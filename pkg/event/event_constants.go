package event

const (
	// Event types
	TypeNormal  = "Normal"
	TypeWarning = "Warning"

	// Event Reasons
	//nolint
	ReasonErrFetchingCredstashSecret = "ErrCredstash"
	ReasonErrGeneric                 = "ErrGeneric"

	ReasonErrCreateSecret     = "ErrCreateSecret"
	ReasonSuccessCreateSecret = "SuccessCreateSecret"

	//nolint
	ReasonErrDeleteOldSecret     = "ErrDeleteOldSecret"
	ReasonSuccessDeleteOldSecret = "SuccessDeleteOldSecret"

	ReasonErrUpdateSecret     = "ErrUpdateSecret"
	ReasonSuccessUpdateSecret = "SuccessUpdateSecret"

	// Event Messages
	MessageFailedFetchingCredstashSecret = "Failed fetching credstash secret. Key: %s. Version: %s. Table: %s. Context: %s. Error %s"

	MessageFailedCreatingSecret  = "Failed creating secret. Name: %s. Namespace: %s. Error %s"
	MessageSuccessCreatingSecret = "Successfully created secret. Name: %s. Namespace: %s"

	//nolint
	MessageFailedDeletingOldSecret  = "Failed deleting old secret. Name: %s. Namespace: %s. Error %s"
	MessageSuccessDeletingOldSecret = "Successfully deleted old secret. Name: %s. Namespace: %s"

	MessageFailedUpdatingSecret  = "Failed updating secret. Name: %s. Namespace: %s. Error %s"
	MessageSuccessUpdatingSecret = "Successfully updated secret. Name: %s. Namespace: %s"
)
