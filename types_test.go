package autoscaler_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"

	. "github.com/bijukunjummen/app-autoscaler-client"
)

var _ = Describe("ServiceInstances type", func() {
	sampleJSON := []byte(`
    {
    "resources": [
        {
            "guid": "540f43bc-b9cc-4126-97a4-a56b64052da4",
            "created_at": "2021-01-01T00:00:00Z",
            "updated_at": "2021-01-01T00:00:00Z",
            "app_name": "your-app-name",
            "min_instances": 5,
            "max_instances": 10,
            "expected_instance_count": 5,
            "enabled": true,
            "relationships": {
                "most_recent_event": {
                    "guid": "193809b3-35d4-4483-b5f8-3d9d57bc9f30",
                    "created_at": "2021-01-01T00:00:00Z",
                    "updated_at": "2021-02-01T00:00:00Z",
                    "reading_id": 23,
                    "service_binding_guid": "540f43bc-b9cc-4126-97a4-a56b64052da4",
                    "scaling_factor": 3,
                    "description": "Minimum instance limit of 1 reached"
                },
                "next_scheduled_limit_change": {
                    "guid": "bdb345ab-8c04-4547-9938-d1a857943f6a",
                    "created_at": "2021-01-01T00:00:00Z",
                    "updated_at": "2021-01-01T00:00:00Z",
                    "executes_at": "2021-01-01T00:00:00Z",
                    "min_instances": 2,
                    "max_instances": 5,
                    "service_binding_guid": "540f43bc-b9cc-4126-97a4-a56b64052da4",
                    "recurrence": 1,
                    "enabled": true
                },
                "rules": [
                    {
                        "guid": "59c19991-ee1d-4e09-8057-94e0a614941a",
                        "service_binding_guid": "540f43bc-b9cc-4126-97a4-a56b64052da4",
                        "created_at": "2021-01-01T00:00:00Z",
                        "updated_at": "2021-01-01T00:00:00Z",
                        "type": "cpu",
                        "enabled": true,
                        "min_threshold": 50,
                        "max_threshold": 80
                    }
                ]
            },
            "links": {
                "self": {
                    "href": "/api/bindings/540f43bc-b9cc-4126-97a4-a56b64052da4"
                },
                "scheduled_rules": {
                    "href": "/api/bindings/540f43bc-b9cc-4126-97a4-a56b64052da4/scheduled_limit_changes"
                },
                "events": {
                    "href": "/api/bindings/540f43bc-b9cc-4126-97a4-a56b64052da4/events"
                }
            }
        }
    ]}
    `)
	var serviceInstances ServiceInstances
	err := json.Unmarshal(sampleJSON, &serviceInstances)
	Context("Given a Sample raw json", func() {

		It("Should be transformed to ServiceInstances without errors", func() {
			Ω(err).Should(BeNil())
		})
		It("Should have exactly 1 resource under it", func() {
			Ω(len(serviceInstances.BindingResources)).Should(Equal(1))
		})
		It("Should have all tags under the Resource element", func() {
			resource := serviceInstances.BindingResources[0]

			Ω(resource.GUID).Should(Equal("540f43bc-b9cc-4126-97a4-a56b64052da4"))
			Ω(resource.AppName).Should(Equal("your-app-name"))
			Ω(resource.MinInstances).Should(Equal(5))
			Ω(resource.MaxInstances).Should(Equal(10))
			Ω(resource.ExpectedInstanceCount).Should(Equal(5))
			Ω(resource.Enabled).Should(Equal(true))
		})
		It("Should be able to retrieve the Binding with all sub tags", func() {
			resource := serviceInstances.BindingResources[0]
			binding := resource.Binding

			Ω(binding.GUID).Should(Equal("540f43bc-b9cc-4126-97a4-a56b64052da4"))
			Ω(binding.AppName).Should(Equal("your-app-name"))
			Ω(binding.MinInstances).Should(Equal(5))
			Ω(binding.MaxInstances).Should(Equal(10))
			Ω(binding.ExpectedInstanceCount).Should(Equal(5))
			Ω(binding.Enabled).Should(Equal(true))
		})
	})

	Context("Given the Resource tag", func() {
		resource := serviceInstances.BindingResources[0]

		Context("And the relationships tag under it", func() {
			relationships := resource.Relationships

			It("Should have a most_recent_event subtag", func() {
				Ω(relationships.MostRecentEvent).ShouldNot(BeNil())
			})
			It("Should have all tags in most_recent_event", func() {
				mre := relationships.MostRecentEvent
				Ω(mre.GUID).Should(Equal("193809b3-35d4-4483-b5f8-3d9d57bc9f30"))
				Ω(mre.ReadingID).Should(Equal(23))
				Ω(mre.ServiceBindingGUID).Should(Equal("540f43bc-b9cc-4126-97a4-a56b64052da4"))
				Ω(mre.ScalingFactor).Should(Equal(3))
				Ω(mre.Description).Should(HavePrefix("Minimum instance limit"))
				cd, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
				Ω(*mre.CreatedAt).Should(Equal(cd))
			})

			It("Should have a next_scheduled_limit_change tag", func() {
				Ω(relationships.NextScheduledLimitChange).ShouldNot(BeNil())
			})
			It("Should have all tags in next_scheduled_limit_change element", func() {
				nsl := relationships.NextScheduledLimitChange
				Ω(nsl.GUID).Should(Equal("bdb345ab-8c04-4547-9938-d1a857943f6a"))
				Ω(nsl.MinInstances).Should(Equal(2))
				Ω(nsl.MaxInstances).Should(Equal(5))
				Ω(nsl.ServiceBindingGUID).Should(Equal("540f43bc-b9cc-4126-97a4-a56b64052da4"))
				Ω(nsl.Recurrence).Should(Equal(1))
				Ω(nsl.Enabled).Should(Equal(true))
				ea, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
				Ω(*nsl.ExecutesAt).Should(Equal(ea))
			})
			It("Should have rules element with exactly 1 rule", func() {
				rules := relationships.Rules
				Ω(len(rules)).Should(Equal(1))
			})
			It("Each rule should have all expected sub elements", func() {
				rules := relationships.Rules
				rule := rules[0]
				Ω(rule.GUID).Should(Equal("59c19991-ee1d-4e09-8057-94e0a614941a"))
				Ω(rule.ServiceBindingGUID).Should(Equal("540f43bc-b9cc-4126-97a4-a56b64052da4"))
				Ω(rule.Type).Should(Equal("cpu"))
				Ω(rule.Enabled).Should(Equal(true))
				Ω(rule.MinThreshold).Should(Equal(50))
				Ω(rule.MaxThreshold).Should(Equal(80))
			})
		})
		Context("Given the Links tag", func() {
			links := serviceInstances.BindingResources[0].Links

			It("should have appropriate links", func() {
				Ω(len(links)).Should(Equal(3))
				Ω(links["self"].Href).Should(Equal("/api/bindings/540f43bc-b9cc-4126-97a4-a56b64052da4"))
				Ω(links["scheduled_rules"].Href).Should(Equal("/api/bindings/540f43bc-b9cc-4126-97a4-a56b64052da4/scheduled_limit_changes"))
				Ω(links["events"].Href).Should(Equal("/api/bindings/540f43bc-b9cc-4126-97a4-a56b64052da4/events"))

			})
		})
	})

})
