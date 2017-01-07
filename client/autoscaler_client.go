package client

import (
	"fmt"

	"github.com/bijukunjummen/app-autoscaler-client/types"
	"github.com/bijukunjummen/app-autoscaler-client/uaa_client"

	"bytes"
	"encoding/json"
)

type AutoScalerClient interface {
	GetServiceBindings() (*types.ServiceInstances, error)
	GetBinding(bindingGuid string) (*types.BindingResource, error)
	UpdateBinding(bindingGuid string, binding *types.Binding) (*types.BindingResource, error)
	GetScalingDecisions(bindingGuid string) ([]types.ScalingDecision, error)
	GetScheduledLimitChanges(bindingGuid string) ([]types.ScheduledLimitChange, error)
	CreateScheduledLimitChange(bindingGuid string, scheduledLimitChange *types.ScheduledLimitChange) (*types.ScheduledLimitChange, error)
	UpdateScheduledLimitChange(bindingGuid string, changeGuid string, scheduledLimitChange *types.ScheduledLimitChange) (*types.ScheduledLimitChange, error)
	DeleteScheduledLimitChange(bindingGuid string, changeGuid string, scheduledLimitChange *types.ScheduledLimitChange) error
}

type AutoscalerConfig struct {
	UAAConfig        *uaa_client.Config
	AutoscalerAPIUrl string
	InstanceGUID     string
}

type DefaultAutoScalerClient struct {
	httpClient uaa_client.OauthHttpWrapper
	config     *AutoscalerConfig
}

func NewAutoScalerClient(autoscalerConfig *AutoscalerConfig) (AutoScalerClient, error) {
	uaaConfig := autoscalerConfig.UAAConfig
	oauthWrapper, err := uaa_client.NewClient(uaaConfig)

	if err != nil {
		return nil, err
	}

	return &DefaultAutoScalerClient{
		httpClient: oauthWrapper,
		config:     autoscalerConfig,
	}, nil
}

func (autoscalerClient *DefaultAutoScalerClient) GetServiceBindings() (*types.ServiceInstances, error) {
	serviceBindingsUrl := fmt.Sprintf("%s/instances/%s/bindings", autoscalerClient.config.AutoscalerAPIUrl, autoscalerClient.config.InstanceGUID)
	request, err := autoscalerClient.httpClient.NewRequest("GET", serviceBindingsUrl, nil)
	if err != nil {
		return nil, err
	}
	resp, err := autoscalerClient.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	var serviceInstances types.ServiceInstances

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&serviceInstances); err != nil {
		return nil, err
	}

	return &serviceInstances, nil

}

func (autoscalerClient *DefaultAutoScalerClient) GetBinding(bindingGuid string) (*types.BindingResource, error) {
	bindingUrl := fmt.Sprintf("%s/bindings/%s", autoscalerClient.config.AutoscalerAPIUrl, bindingGuid)
	request, err := autoscalerClient.httpClient.NewRequest("GET", bindingUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := autoscalerClient.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	var binding types.BindingResource

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&binding); err != nil {
		return nil, err
	}
	return &binding, nil
}

func (autoscalerClient *DefaultAutoScalerClient) UpdateBinding(bindingGuid string, binding *types.Binding) (*types.BindingResource, error) {
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
	var bindingUpdated types.BindingResource

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&bindingUpdated); err != nil {
		return nil, err
	}
	return &bindingUpdated, nil
}

func (autoscalerClient *DefaultAutoScalerClient) GetScalingDecisions(bindingGuid string) ([]types.ScalingDecision, error) {
	return nil, fmt.Errorf("Not implemented yet!")
}

func (autoscalerClient *DefaultAutoScalerClient) GetScheduledLimitChanges(bindingGuid string) ([]types.ScheduledLimitChange, error) {
	schedulesForBindingUrl := fmt.Sprintf("%s/bindings/%s/scheduled_limit_changes", autoscalerClient.config.AutoscalerAPIUrl, bindingGuid)
	request, err := autoscalerClient.httpClient.NewRequest("GET", schedulesForBindingUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := autoscalerClient.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	var scheduledLimitChanges []types.ScheduledLimitChange

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&scheduledLimitChanges); err != nil {
		return nil, err
	}
	return scheduledLimitChanges, nil
}

func (autoscalerClient *DefaultAutoScalerClient) CreateScheduledLimitChange(bindingGuid string, scheduledLimitChange *types.ScheduledLimitChange) (*types.ScheduledLimitChange, error) {
	return nil, fmt.Errorf("Not implemented yet!")
}
func (autoscalerClient *DefaultAutoScalerClient) UpdateScheduledLimitChange(bindingGuid string, changeGuid string, scheduledLimitChange *types.ScheduledLimitChange) (*types.ScheduledLimitChange, error) {
	return nil, fmt.Errorf("Not implemented yet!")
}
func (autoscalerClient *DefaultAutoScalerClient) DeleteScheduledLimitChange(bindingGuid string, changeGuid string, scheduledLimitChange *types.ScheduledLimitChange) error {
	return fmt.Errorf("Not implemented yet!")
}
