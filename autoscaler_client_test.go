package autoscaler_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"
	"net/http"

	"github.com/bijukunjummen/app-autoscaler-client"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Behavior of Auto Scaler", func() {

	var server *ghttp.Server
	var config *autoscaler.Config

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
		config = &autoscaler.Config{
			CFConfig: &autoscaler.CFConfig{
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
				ghttp.RespondWithJSONEncoded(http.StatusOK, autoscaler.Endpoint{
					AuthorizationEndpoint: server.URL(),
					TokenEndpoint:         server.URL(),
				}),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/oauth/token"),
				ghttp.VerifyBasicAuth("cf", ""),
				ghttp.RespondWithJSONEncoded(http.StatusOK, autoscaler.AccessToken{
					Token: "test-token",
				}),
			),
		)
	})

	Context("Given a Autoscaler Client", func() {

		It("Should be able to Get Service Bindings for a App Scaler instance", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/instances/c017ec06-cf4c-42fa-adbd-1b6a290d8d6a/bindings"),

					ghttp.VerifyHeader(http.Header{
						"Authorization": []string{"Bearer test-token"},
					}),
					ghttp.RespondWith(http.StatusOK, sampleServiceInstancesJson),
				),
			)
			client, err := autoscaler.NewClient(config)
			Ω(err).Should(BeNil())

			serviceInstances, err := client.GetServiceBindings()

			Ω(err).Should(BeNil())

			Ω(len(serviceInstances.BindingResources)).Should(Equal(1))
		})

		It("Should be able to Get Details of a binding given a binding Id", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/bindings/mybinding"),

					ghttp.VerifyHeader(http.Header{
						"Authorization": []string{"Bearer test-token"},
					}),
					ghttp.RespondWith(http.StatusOK, sampleBindingJson),
				),
			)
			client, err := autoscaler.NewClient(config)
			Ω(err).Should(BeNil())

			binding, err := client.GetBinding("mybinding")
			Ω(err).Should(BeNil())

			Ω(binding.AppName).Should(Equal("sample-spring-cloud-svc-ci"))
		})
		It("Should be able to Update Details of a binding given a binding Id", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/api/bindings/mybinding"),

					ghttp.VerifyHeader(http.Header{
						"Authorization": []string{"Bearer test-token"},
					}),
					ghttp.RespondWith(http.StatusOK, sampleBindingJson),
				),
			)
			client, err := autoscaler.NewClient(config)
			Ω(err).Should(BeNil())

			var binding autoscaler.Binding
			bindingJsonBytes := []byte(sampleBindingJson)
			json.Unmarshal(bindingJsonBytes, &binding)

			bindingResource, err := client.UpdateBinding("mybinding", &binding)
			Ω(err).Should(BeNil())

			Ω(bindingResource.AppName).Should(Equal("sample-spring-cloud-svc-ci"))
		})
	})

	scheduledLimitChangesResource :=
		`{"resources": [
	  {
	    "guid": "de743e10-e427-4b02-b3b4-7bf560c201ab",
	    "created_at": "2021-01-01T00:00:00Z",
	    "updated_at": "2021-01-01T00:00:00Z",
	    "executes_at": "2014-11-22T16:00:00Z",
	    "min_instances": 2,
	    "max_instances": 3,
	    "service_binding_guid": "540f43bc-b9cc-4126-97a4-a56b64052da4",
	    "recurrence": 20,
	    "enabled": true,
	    "links":{"self":{"href":"/api/bindings/d3156757-6e5e-40fc-8301-2ca8dec795d2/scheduled_limit_changes/0aebba6f-d389-4d64-6bb6-49c46530945c"},"service_binding":{"href":"/api/bindings/d3156757-6e5e-40fc-8301-2ca8dec795d2"}}
	  },
	  {
	    "guid": "cb9e6e41-ffef-4e68-a585-eb004d9bb122",
	    "created_at": "2021-01-01T00:00:00Z",
	    "updated_at": "2021-01-01T00:00:00Z",
	    "executes_at": "2014-11-22T16:00:00Z",
	    "min_instances": 4,
	    "max_instances": 5,
	    "service_binding_guid": "540f43bc-b9cc-4126-97a4-a56b64052da4",
	    "recurrence": 60,
	    "enabled": true
	  }
	]}`

	scheduledLimitChange :=
		`
		{
		  "guid": "de743e10-e427-4b02-b3b4-7bf560c201ab",
		  "created_at": "2021-01-01T00:00:00Z",
		  "updated_at": "2021-01-01T00:00:00Z",
		  "executes_at": "2014-11-22T16:00:00Z",
		  "min_instances": 2,
		  "max_instances": 3,
		  "service_binding_guid": "540f43bc-b9cc-4126-97a4-a56b64052da4",
		  "recurrence": 20,
		  "enabled": true
		}
		`

	It("Should be able to Retrieve Schedules for a binding given a Binding Id", func() {
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/bindings/mybinding/scheduled_limit_changes"),

				ghttp.VerifyHeader(http.Header{
					"Authorization": []string{"Bearer test-token"},
				}),
				ghttp.RespondWith(http.StatusOK, scheduledLimitChangesResource),
			),
		)
		client, _ := autoscaler.NewClient(config)
		schedules, err := client.GetScheduledLimitChanges("mybinding")
		Ω(err).Should(BeNil())
		Ω(len(schedules)).Should(Equal(2))
		change1 := schedules[0]
		Ω(change1.GUID).Should(Equal("de743e10-e427-4b02-b3b4-7bf560c201ab"))
		Ω(change1.MaxInstances).Should(Equal(3))
		Ω(change1.MinInstances).Should(Equal(2))
		auditDate, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
		eDate, _ := time.Parse(time.RFC3339, "2014-11-22T16:00:00Z")
		Ω(*change1.CreatedAt).Should(Equal(auditDate))
		Ω(*change1.UpdatedAt).Should(Equal(auditDate))
		Ω(*change1.ExecutesAt).Should(Equal(eDate))
		Ω(change1.ServiceBindingGUID).Should(Equal("540f43bc-b9cc-4126-97a4-a56b64052da4"))
		Ω(change1.Recurrence).Should(Equal(20))
		Ω(change1.Enabled).Should(Equal(true))
	})

	It("Should be able to update a schedule", func() {
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/api/bindings/mybinding/scheduled_limit_changes/changeid"),
				ghttp.VerifyHeader(http.Header{
					"Authorization": []string{"Bearer test-token"},
				}),
				ghttp.RespondWith(http.StatusOK, scheduledLimitChange),
			),
		)

		client, _ := autoscaler.NewClient(config)
		var scheduledLimitChangeObj autoscaler.ScheduledLimitChange
		err := json.Unmarshal([]byte(scheduledLimitChange), &scheduledLimitChangeObj)
		Ω(err).Should(BeNil())
		scheduleUpdated, err := client.UpdateScheduledLimitChange("mybinding", "changeid", &scheduledLimitChangeObj)
		Ω(scheduleUpdated.Enabled).Should(BeTrue())
		Ω(scheduleUpdated.Recurrence).Should(Equal(20))
		Ω(scheduleUpdated.ServiceBindingGUID).Should(Equal("540f43bc-b9cc-4126-97a4-a56b64052da4"))
		Ω(scheduleUpdated.MinInstances).Should(Equal(2))
		Ω(scheduleUpdated.MaxInstances).Should(Equal(3))
		Ω(scheduleUpdated.GUID).Should(Equal("de743e10-e427-4b02-b3b4-7bf560c201ab"))
	})

	It("Should be able to create a schedule", func() {
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/bindings/mybinding/scheduled_limit_changes"),
				ghttp.VerifyHeader(http.Header{
					"Authorization": []string{"Bearer test-token"},
				}),
				ghttp.RespondWith(http.StatusCreated, scheduledLimitChange),
			),
		)

		client, _ := autoscaler.NewClient(config)
		var scheduledLimitChangeObj autoscaler.ScheduledLimitChange
		err := json.Unmarshal([]byte(scheduledLimitChange), &scheduledLimitChangeObj)
		Ω(err).Should(BeNil())
		scheduleUpdated, err := client.CreateScheduledLimitChange("mybinding", &scheduledLimitChangeObj)
		Ω(scheduleUpdated.Enabled).Should(BeTrue())
		Ω(scheduleUpdated.Recurrence).Should(Equal(20))
		Ω(scheduleUpdated.ServiceBindingGUID).Should(Equal("540f43bc-b9cc-4126-97a4-a56b64052da4"))
		Ω(scheduleUpdated.MinInstances).Should(Equal(2))
		Ω(scheduleUpdated.MaxInstances).Should(Equal(3))
		Ω(scheduleUpdated.GUID).Should(Equal("de743e10-e427-4b02-b3b4-7bf560c201ab"))
	})

	It("Should be able to delete a schedule", func() {
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/api/bindings/mybinding/scheduled_limit_changes/changeid"),
				ghttp.VerifyHeader(http.Header{
					"Authorization": []string{"Bearer test-token"},
				}),
				ghttp.RespondWith(http.StatusOK, nil),
			),
		)

		client, _ := autoscaler.NewClient(config)
		var scheduledLimitChangeObj autoscaler.ScheduledLimitChange
		err := json.Unmarshal([]byte(scheduledLimitChange), &scheduledLimitChangeObj)
		Ω(err).Should(BeNil())
		err = client.DeleteScheduledLimitChange("mybinding", "changeid")
		Ω(err).Should(BeNil())
	})
})
