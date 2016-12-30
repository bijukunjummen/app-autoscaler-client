package instance

// MostRecentEvent -
type MostRecentEvent struct {
	GUID               string `json:"guid,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
	UpdatedAt          string `json:"updated_at,omitempty"`
	ReadingID          int    `json:"reading_id"`
	ServiceBindingGUID string `json:"service_binding_guid"`
	ScalingFactor      int    `json:"scaling_factor"`
	Description        string `json:"description"`
}

// NextScheduledLimitChange -
type NextScheduledLimitChange struct {
	GUID               string `json:"guid,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
	UpdatedAt          string `json:"updated_at,omitempty"`
	ExecutesAt         string `json:"executes_at"`
	MinInstances       int    `json:"min_instances"`
	MaxInstances       int    `json:"max_instances"`
	ServiceBindingGUID string `json:"service_binding_guid,omitempty"`
	Recurrence         int    `json:"recurrence"`
	Enabled            bool   `json:"enabled"`
}

//Rule -
type Rule struct {
	GUID               string `json:"guid,omitempty"`
	ServiceBindingGUID string `json:"service_binding_guid,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
	UpdatedAt          string `json:"updated_at,omitempty"`
	Type               string `json:"type"`
	Enabled            bool   `json:"enabled"`
	MinThreshold       int    `json:"min_threshold"`
	MaxThreshold       int    `json:"max_threshold"`
}

//Relationships -
type Relationships struct {
	MostRecentEvent          MostRecentEvent          `json:"most_recent_event,omitempty"`
	NextScheduledLimitChange NextScheduledLimitChange `json:"next_scheduled_limit_change,omitempty"`
	Rules                    []Rule                   `json:"rules"`
}

//Binding -
type Binding struct {
	GUID                  string        `json:"guid,omitempty"`
	CreatedAt             string        `json:"created_at,omitempty"`
	UpdatedAt             string        `json:"updated_at,omitempty"`
	AppName               string        `json:"app_name,omitempty"`
	MinInstances          int           `json:"min_instances,omitempty"`
	MaxInstances          int           `json:"max_instances,omitempty"`
	ExpectedInstanceCount int           `json:"expected_instance_count,omitempty"`
	Enabled               bool          `json:"enabled,omitempty"`
	Relationships         Relationships `json:"relationships,omitempty"`
}

// Resource -
type BindingResource struct {
	Binding

	Links map[string]Link `json:"links,omitempty"`
}

// Link -
type Link struct {
	Href string `json:"href"`
}

// ServiceInstances -
type ServiceInstances struct {
	BindingResources []BindingResource `json:"resources"`
}
