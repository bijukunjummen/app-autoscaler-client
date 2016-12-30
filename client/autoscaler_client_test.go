package client_test

import (
	. "github.com/bijukunjummen/app-autoscaler-client/client"
	"github.com/bijukunjummen/app-autoscaler-client/uaa_client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("Behavior of Auto Scaler", func() {

	var server *ghttp.Server
	var config *AutoscalerConfig

	sampleServiceInstancesJson := `
	{
	  "resources": [
	    {
	      "guid": "d3156757-6e5e-40fc-8301-2ca8dec795d2",
	      "service_instance_guid": "c017ec06-cf4c-42fa-adbd-1b6a290d8d6a",
	      "app_guid": "8d7069ee-7f77-4383-85db-595cec95638d",
	      "app_name": "sample-spring-cloud-svc-ci",
	      "expected_instance_count": 2,
	      "min_instances": 2,
	      "max_instances": 5,
	      "enabled": false,
	      "created_at": "2016-12-20T13:53:23Z",
	      "updated_at": "2016-12-29T07:27:14Z",
	      "relationships": {
		"rules": [
		  {
		    "guid": "d53b1ca4-9e01-47c2-7b35-8805cd66f351",
		    "created_at": "2016-12-20T13:53:23Z",
		    "updated_at": "2016-12-21T04:08:40Z",
		    "type": "cpu",
		    "enabled": true,
		    "min_threshold": 20,
		    "max_threshold": 80,
		    "service_binding_guid": "d3156757-6e5e-40fc-8301-2ca8dec795d2",
		    "sub_type": ""
		  }
		],
		"most_recent_event": {
		  "guid": "d99630bd-da36-4002-5177-27eeb9b30775",
		  "service_binding_guid": "d3156757-6e5e-40fc-8301-2ca8dec795d2",
		  "scaling_factor": 0,
		  "description": "Cannot scale: at min limit of 2 instances.\nCurrent CPU of 0.09% is below lower threshold of 20%.",
		  "created_at": "2016-12-29T07:26:31Z",
		  "updated_at": "2016-12-29T07:26:31Z"
		}
	      },
	      "links": {
		"events": {
		  "href": "/api/bindings/d3156757-6e5e-40fc-8301-2ca8dec795d2/scaling_events"
		},
		"scheduled_limit_changes": {
		  "href": "/api/bindings/d3156757-6e5e-40fc-8301-2ca8dec795d2/scheduled_limit_changes"
		},
		"self": {
		  "href": "/api/bindings/d3156757-6e5e-40fc-8301-2ca8dec795d2"
		}
	      }
	    }
	  ]
	}
	`

	sampleBindingJson := `
	{
	  "guid": "d3156757-6e5e-40fc-8301-2ca8dec795d2",
	  "service_instance_guid": "c017ec06-cf4c-42fa-adbd-1b6a290d8d6a",
	  "app_guid": "8d7069ee-7f77-4383-85db-595cec95638d",
	  "app_name": "sample-spring-cloud-svc-ci",
	  "expected_instance_count": 2,
	  "min_instances": 2,
	  "max_instances": 5,
	  "enabled": true,
	  "created_at": "2016-12-20T13:53:23Z",
	  "updated_at": "2016-12-30T00:48:36Z",
	  "relationships": {
	    "rules": [
	      {
		"guid": "d5524ca3-5eda-43ea-688c-f5685d70b64e",
		"created_at": "2016-12-30T00:48:36Z",
		"updated_at": "2016-12-30T00:48:36Z",
		"type": "http_latency",
		"enabled": false,
		"min_threshold": 100,
		"max_threshold": 900,
		"service_binding_guid": "d3156757-6e5e-40fc-8301-2ca8dec795d2",
		"sub_type": "avg_95th"
	      },
	      {
		"guid": "d53b1ca4-9e01-47c2-7b35-8805cd66f351",
		"created_at": "2016-12-20T13:53:23Z",
		"updated_at": "2016-12-30T00:48:36Z",
		"type": "http_throughput",
		"enabled": true,
		"min_threshold": 20,
		"max_threshold": 80,
		"service_binding_guid": "d3156757-6e5e-40fc-8301-2ca8dec795d2",
		"sub_type": ""
	      }
	    ],
	    "most_recent_event": {
	      "guid": "3960c9d5-a2a5-4404-7940-8c6a5f1bef3c",
	      "service_binding_guid": "d3156757-6e5e-40fc-8301-2ca8dec795d2",
	      "scaling_factor": 0,
	      "description": "Cannot scale: at min limit of 2 instances.\nCurrent HTTP Throughput of 0.00/s/instance is below lower threshold of 20/s/instance.",
	      "created_at": "2016-12-30T01:25:00Z",
	      "updated_at": "2016-12-30T01:25:00Z"
	    }
	  },
	  "links": {
	    "events": {
	      "href": "/api/bindings/d3156757-6e5e-40fc-8301-2ca8dec795d2/scaling_events"
	    },
	    "scheduled_limit_changes": {
	      "href": "/api/bindings/d3156757-6e5e-40fc-8301-2ca8dec795d2/scheduled_limit_changes"
	    },
	    "self": {
	      "href": "/api/bindings/d3156757-6e5e-40fc-8301-2ca8dec795d2"
	    }
	  }
	}
	`

	BeforeEach(func() {

		server = ghttp.NewServer()
		config = &AutoscalerConfig{
			UAAConfig: &uaa_client.Config{
				CCApiUrl:          server.URL(),
				Username:          "user",
				Password:          "pwd",
				SkipSslValidation: true,
			},
			AutoscalerAPIUrl: server.URL() + "/api",
			InstanceGUID:     "c017ec06-cf4c-42fa-adbd-1b6a290d8d6a",
		}
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/v2/info"),
				ghttp.RespondWithJSONEncoded(http.StatusOK, uaa_client.Endpoint{
					AuthorizationEndpoint: server.URL(),
					TokenEndpoint:         server.URL(),
				}),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/oauth/token"),
				ghttp.VerifyBasicAuth("cf", ""),
				ghttp.RespondWithJSONEncoded(http.StatusOK, uaa_client.AccessToken{
					Token: "test-token",
				}),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/instances/c017ec06-cf4c-42fa-adbd-1b6a290d8d6a/bindings"),

				ghttp.VerifyHeader(http.Header{
					"Authorization": []string{"Bearer test-token"},
				}),
				ghttp.RespondWith(http.StatusOK, sampleServiceInstancesJson),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/bindings/mybinding"),

				ghttp.VerifyHeader(http.Header{
					"Authorization": []string{"Bearer test-token"},
				}),
				ghttp.RespondWith(http.StatusOK, sampleBindingJson),
			),
		)
	})

	Context("Given a Autoscaler Client", func() {

		It("Should be able to query for information", func() {
			client, err := NewAutoScalerClient(config)
			Ω(err).Should(BeNil())

			serviceInstances, err := client.GetServiceBindings()

			Ω(err).Should(BeNil())

			Ω(len(serviceInstances.BindingResources)).Should(Equal(1))

			binding, err := client.GetBinding("mybinding")

			Ω(binding.AppName).Should(Equal("sample-spring-cloud-svc-ci"))
		})
	})
})
