package autoscaler

import (
	"fmt"


	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type AutoScalerClient interface {
	GetServiceBindings() (*ServiceInstances, error)
	GetBinding(bindingGuid string) (*BindingResource, error)
	UpdateBinding(bindingGuid string, binding *Binding) (*BindingResource, error)
	GetScalingDecisions(bindingGuid string) ([]ScalingDecision, error)
	GetScheduledLimitChanges(bindingGuid string) ([]ScheduledLimitChange, error)
	CreateScheduledLimitChange(bindingGuid string, scheduledLimitChange *ScheduledLimitChange) (*ScheduledLimitChange, error)
	UpdateScheduledLimitChange(bindingGuid string, changeGuid string, scheduledLimitChange *ScheduledLimitChange) (*ScheduledLimitChange, error)
	DeleteScheduledLimitChange(bindingGuid string, changeGuid string) error
}

type AutoscalerConfig struct {
	UAAConfig        *Config
	AutoscalerAPIUrl string
	InstanceGUID     string
}

type DefaultAutoScalerClient struct {
	httpClient OauthHttpWrapper
	config     *AutoscalerConfig
}

func NewAutoScalerClient(autoscalerConfig *AutoscalerConfig) (AutoScalerClient, error) {
	uaaConfig := autoscalerConfig.UAAConfig
	oauthWrapper, err := NewClient(uaaConfig)

	if err != nil {
		return nil, err
	}

	return &DefaultAutoScalerClient{
		httpClient: oauthWrapper,
		config:     autoscalerConfig,
	}, nil
}

func (autoscalerClient *DefaultAutoScalerClient) GetServiceBindings() (*ServiceInstances, error) {
	serviceBindingsUrl := fmt.Sprintf("%s/instances/%s/bindings", autoscalerClient.config.AutoscalerAPIUrl, autoscalerClient.config.InstanceGUID)
	request, err := autoscalerClient.httpClient.NewRequest("GET", serviceBindingsUrl, nil)
	if err != nil {
		return nil, err
	}
	resp, err := autoscalerClient.httpClient.Do(request)
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

func (autoscalerClient *DefaultAutoScalerClient) GetBinding(bindingGuid string) (*BindingResource, error) {
	bindingUrl := fmt.Sprintf("%s/bindings/%s", autoscalerClient.config.AutoscalerAPIUrl, bindingGuid)
	request, err := autoscalerClient.httpClient.NewRequest("GET", bindingUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := autoscalerClient.httpClient.Do(request)
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

func (autoscalerClient *DefaultAutoScalerClient) UpdateBinding(bindingGuid string, binding *Binding) (*BindingResource, error) {
	bindingUrl := fmt.Sprintf("%s/bindings/%s", autoscalerClient.config.AutoscalerAPIUrl, bindingGuid)

	body, err := json.Marshal(binding)
	if err != nil {
		return nil, err
	}
	request, err := autoscalerClient.httpClient.NewRequest("PUT", bindingUrl, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	resp, err := autoscalerClient.httpClient.Do(request)
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

func (autoscalerClient *DefaultAutoScalerClient) GetScalingDecisions(bindingGuid string) ([]ScalingDecision, error) {
	return nil, fmt.Errorf("%s", "Not implemented...")
}

func (autoscalerClient *DefaultAutoScalerClient) GetScheduledLimitChanges(bindingGuid string) ([]ScheduledLimitChange, error) {
	schedulesForBindingUrl := fmt.Sprintf("%s/bindings/%s/scheduled_limit_changes", autoscalerClient.config.AutoscalerAPIUrl, bindingGuid)

	request, err := autoscalerClient.httpClient.NewRequest("GET", schedulesForBindingUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := autoscalerClient.httpClient.Do(request)
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

func (autoscalerClient *DefaultAutoScalerClient) CreateScheduledLimitChange(bindingGuid string, scheduledLimitChange *ScheduledLimitChange) (*ScheduledLimitChange, error) {
	schedulesForBindingUrl := fmt.Sprintf("%s/bindings/%s/scheduled_limit_changes", autoscalerClient.config.AutoscalerAPIUrl, bindingGuid)

	body, err := json.Marshal(scheduledLimitChange)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s", string(body))

	request, err := autoscalerClient.httpClient.NewRequest("POST", schedulesForBindingUrl, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	resp, err := autoscalerClient.httpClient.Do(request)
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

func (autoscalerClient *DefaultAutoScalerClient) UpdateScheduledLimitChange(bindingGuid string, changeGuid string, scheduledLimitChange *ScheduledLimitChange) (*ScheduledLimitChange, error) {
	schedulesForBindingUrl := fmt.Sprintf("%s/bindings/%s/scheduled_limit_changes/%s", autoscalerClient.config.AutoscalerAPIUrl, bindingGuid, changeGuid)

	body, err := json.Marshal(scheduledLimitChange)
	if err != nil {
		return nil, err
	}

	request, err := autoscalerClient.httpClient.NewRequest("PUT", schedulesForBindingUrl, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	resp, err := autoscalerClient.httpClient.Do(request)
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
func (autoscalerClient *DefaultAutoScalerClient) DeleteScheduledLimitChange(bindingGuid string, changeGuid string) error {
	schedulesForBindingUrl := fmt.Sprintf("%s/bindings/%s/scheduled_limit_changes/%s", autoscalerClient.config.AutoscalerAPIUrl, bindingGuid, changeGuid)

	request, err := autoscalerClient.httpClient.NewRequest("DELETE", schedulesForBindingUrl, nil)

	if err != nil {
		return err
	}

	resp, err := autoscalerClient.httpClient.Do(request)

	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Bad Response: %s", body)
	}

	return nil
}
