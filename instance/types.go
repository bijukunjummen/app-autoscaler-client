package instance

// MostRecentEvent -
type MostRecentEvent struct {
	GUID               string `json:"guid"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	ReadingID          int    `json:"reading_id"`
	ServiceBindingGUID string `json:"service_binding_guid"`
	ScalingFactor      int    `json:"scaling_factor"`
	Description        string `json:"description"`
}

// NextScheduledLimitChange -
type NextScheduledLimitChange struct {
	GUID               string `json:"guid"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	ExecutesAt         string `json:"executes_at"`
	MinInstances       int    `json:"min_instances"`
	MaxInstances       int    `json:"max_instances"`
	ServiceBindingGUID string `json:"service_binding_guid"`
	Recurrence         int    `json:"recurrence"`
	Enabled            bool   `json:"enabled"`
}

//Rule -
type Rule struct {
	GUID               string `json:"guid"`
	ServiceBindingGUID string `json:"service_binding_guid"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	Type               string `json:"type"`
	Enabled            bool   `json:"enabled"`
	MinThreshold       int    `json:"min_threshold"`
	MaxThreshold       int    `json:"max_threshold"`
}

//Relationships -
type Relationships struct {
	MostRecentEvent          MostRecentEvent          `json:"most_recent_event"`
	NextScheduledLimitChange NextScheduledLimitChange `json:"next_scheduled_limit_change"`
	Rules                    []Rule                   `json:"rules"`
}

//Binding -
type Binding struct {
	GUID                  string          `json:"guid"`
	CreatedAt             string          `json:"created_at"`
	UpdatedAt             string          `json:"updated_at"`
	AppName               string          `json:"app_name"`
	MinInstances          int             `json:"min_instances"`
	MaxInstances          int             `json:"max_instances"`
	ExpectedInstanceCount int             `json:"expected_instance_count"`
	Enabled               bool            `json:"enabled"`
}

// Resource -
type BindingResource struct {
	Binding
	Relationships         Relationships   `json:"relationships"`
	Links                 map[string]Link `json:"links"`
}

// Link -
type Link struct {
	Href string `json:"href"`
}

// ServiceInstances -
type ServiceInstances struct {
	BindingResources []BindingResource `json:"resources"`
}
