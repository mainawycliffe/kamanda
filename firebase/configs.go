package firebase

type FirebaseProjectConfigs struct {
	Projects map[string]string `json:"projects"`
	Targets  interface{}       `json:"-"`
}
