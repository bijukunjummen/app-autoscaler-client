package client

import (
	"fmt"
	"github.com/bijukunjummen/app-autoscaler-client/instance"
	"github.com/bijukunjummen/app-autoscaler-client/uaa_client"

	"encoding/json"
)

type AutoScalerClient interface {
	GetServiceBindings() (*instance.ServiceInstances, error)
	GetBinding(bindingGuid string) (*instance.BindingResource, error)
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

func (autoscalerClient *DefaultAutoScalerClient) GetServiceBindings() (*instance.ServiceInstances, error) {
	serviceBindingsUrl := fmt.Sprintf("%s/instances/%s/bindings", autoscalerClient.config.AutoscalerAPIUrl, autoscalerClient.config.InstanceGUID)
	request, err := autoscalerClient.httpClient.NewRequest("GET", serviceBindingsUrl, nil)
	if err != nil {
		return nil, err
	}
	resp, err := autoscalerClient.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	var serviceInstances instance.ServiceInstances

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&serviceInstances); err != nil {
		return nil, err
	}

	return &serviceInstances, nil

}

func (autoscalerClient *DefaultAutoScalerClient) GetBinding(bindingGuid string) (*instance.BindingResource, error) {
	bindingUrl := fmt.Sprintf("%s/bindings/%s", autoscalerClient.config.AutoscalerAPIUrl, bindingGuid)
	request, err := autoscalerClient.httpClient.NewRequest("GET", bindingUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := autoscalerClient.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	var binding instance.BindingResource

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&binding); err != nil {
		return nil, err
	}
	return &binding, nil
}
