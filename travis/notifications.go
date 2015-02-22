package travis

type Notification struct {
	Payload Payload `json:"payload"`
}

type Payload struct {
	Status     string     `json:"status_message"`
	Commit     string     `json:"commit"`
	Branch     string     `json:"branch"`
	Message    string     `json:"message"`
	Repository Repository `json:"repository"`
}

type Repository struct {
	Name  string `json:"name"`
	Owner string `json:"owner_name"`
}

func (n Notification) Valid() bool {
	if (n.Payload.Status != "Passed") && (n.Payload.Status != "Fixed") {
		return false
	}

	return true
}
