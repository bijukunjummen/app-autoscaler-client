package autoscaler

import (
	"fmt"

	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Client represents the behavior of App Autoscaler API client
type Client interface {
	// Get the Service Bindings of an Application Instance
	GetServiceBindings() (*ServiceInstances, error)

	// Get the Service Binding given the Binding GUID
	GetBinding(bindingGUID string) (*BindingResource, error)

	// Update the Binding
	UpdateBinding(bindingGUID string, binding *Binding) (*BindingResource, error)

	// Get the Scaling decisions for a Binding
	GetScalingDecisions(bindingGUID string) ([]ScalingDecision, error)

	// Get the Schedueled Limit changes for a Binding
	GetScheduledLimitChanges(bindingGUID string) ([]ScheduledLimitChange, error)

	// Create Scheduled Limit change for a binding
	CreateScheduledLimitChange(bindingGUID string, scheduledLimitChange *ScheduledLimitChange) (*ScheduledLimitChange, error)

	//Update Scheduled Limit changes for a binding
	UpdateScheduledLimitChange(bindingGUID string, changeGUID string, scheduledLimitChange *ScheduledLimitChange) (*ScheduledLimitChange, error)

	//Delete Scheduled Limit changes
	DeleteScheduledLimitChange(bindingGUID string, changeGUID string) error
}

// Config holds the configuration for autoscaler settings
type Config struct {
	CFConfig         *CFConfig
	AutoscalerAPIUrl string
	InstanceGUID     string
}

// DefaultClient is the default implementation of Autoscaler Client
type DefaultClient struct {
	httpClient OauthHTTPWrapper
	config     *Config
}

// NewClient is the helper for creating a new Autoscaler Client
func NewClient(autoscalerConfig *Config) (Client, error) {
	uaaConfig := autoscalerConfig.CFConfig
	oauthWrapper, err := NewUAAClient(uaaConfig)

	if err != nil {
		return nil, err
	}

	return &DefaultClient{
		httpClient: oauthWrapper,
		config:     autoscalerConfig,
	}, nil
}

// GetServiceBindings ...
func (client *DefaultClient) GetServiceBindings() (*ServiceInstances, error) {
	serviceBindingsURL := fmt.Sprintf("%s/instances/%s/bindings", client.config.AutoscalerAPIUrl, client.config.InstanceGUID)
	request, err := client.httpClient.NewRequest("GET", serviceBindingsURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Bad Response: %s", body)
	}
	var serviceInstances ServiceInstances

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&serviceInstances); err != nil {
		return nil, err
	}

	return &serviceInstances, nil

}

//GetBinding ...
func (client *DefaultClient) GetBinding(bindingGUID string) (*BindingResource, error) {
	bindingURL := fmt.Sprintf("%s/bindings/%s", client.config.AutoscalerAPIUrl, bindingGUID)
	request, err := client.httpClient.NewRequest("GET", bindingURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Bad Response: %s", body)
	}
	var binding BindingResource

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&binding); err != nil {
		return nil, err
	}
	return &binding, nil
}

//UpdateBinding ...
func (client *DefaultClient) UpdateBinding(bindingGUID string, binding *Binding) (*BindingResource, error) {
	bindingURL := fmt.Sprintf("%s/bindings/%s", client.config.AutoscalerAPIUrl, bindingGUID)

	body, err := json.Marshal(binding)
	if err != nil {
		return nil, err
	}
	request, err := client.httpClient.NewRequest("PUT", bindingURL, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Bad Response: %s", body)
	}
	var bindingUpdated BindingResource

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&bindingUpdated); err != nil {
		return nil, err
	}
	return &bindingUpdated, nil
}

// GetScalingDecisions ...
func (client *DefaultClient) GetScalingDecisions(bindingGUID string) ([]ScalingDecision, error) {
	return nil, fmt.Errorf("%s", "Not implemented...")
}

// GetScheduledLimitChanges ...
func (client *DefaultClient) GetScheduledLimitChanges(bindingGUID string) ([]ScheduledLimitChange, error) {
	schedulesForBindingURL := fmt.Sprintf("%s/bindings/%s/scheduled_limit_changes", client.config.AutoscalerAPIUrl, bindingGUID)

	request, err := client.httpClient.NewRequest("GET", schedulesForBindingURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Bad Response: %s", body)
	}

	var changesResource ScheduledLimitChangesResource

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&changesResource); err != nil {
		return nil, err
	}
	return changesResource.ScheduledLimitChanges, nil
}

//CreateScheduledLimitChange ...
func (client *DefaultClient) CreateScheduledLimitChange(bindingGUID string, scheduledLimitChange *ScheduledLimitChange) (*ScheduledLimitChange, error) {
	schedulesForBindingURL := fmt.Sprintf("%s/bindings/%s/scheduled_limit_changes", client.config.AutoscalerAPIUrl, bindingGUID)

	body, err := json.Marshal(scheduledLimitChange)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s", string(body))

	request, err := client.httpClient.NewRequest("POST", schedulesForBindingURL, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Status Code: %d, Body:%s", resp.StatusCode, body)
	}
	var scheduledLimitChangeUpdated ScheduledLimitChange

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&scheduledLimitChangeUpdated); err != nil {
		return nil, err
	}
	return &scheduledLimitChangeUpdated, nil
}

// UpdateScheduledLimitChange ...
func (client *DefaultClient) UpdateScheduledLimitChange(bindingGUID string, changeGUID string, scheduledLimitChange *ScheduledLimitChange) (*ScheduledLimitChange, error) {
	schedulesForBindingURL := fmt.Sprintf("%s/bindings/%s/scheduled_limit_changes/%s", client.config.AutoscalerAPIUrl, bindingGUID, changeGUID)

	body, err := json.Marshal(scheduledLimitChange)
	if err != nil {
		return nil, err
	}

	request, err := client.httpClient.NewRequest("PUT", schedulesForBindingURL, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Bad Response: %s", body)
	}
	var scheduledLimitChangeUpdated ScheduledLimitChange

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&scheduledLimitChangeUpdated); err != nil {
		return nil, err
	}
	return &scheduledLimitChangeUpdated, nil
}

// DeleteScheduledLimitChange ...
func (client *DefaultClient) DeleteScheduledLimitChange(bindingGUID string, changeGUID string) error {
	schedulesForBindingURL := fmt.Sprintf("%s/bindings/%s/scheduled_limit_changes/%s", client.config.AutoscalerAPIUrl, bindingGUID, changeGUID)

	request, err := client.httpClient.NewRequest("DELETE", schedulesForBindingURL, nil)

	if err != nil {
		return err
	}

	resp, err := client.httpClient.Do(request)

	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Bad Response: %s", body)
	}

	return nil
}
