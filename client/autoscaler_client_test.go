package client_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/bijukunjummen/app-autoscaler-client/client"
	"github.com/bijukunjummen/app-autoscaler-client/uaa_client"
	"fmt"
)

var _ = Describe("Behavior of Auto Scaler", func() {

	config := AutoscalerConfig {
		Config: uaa_client.Config {
			CCApiUrl:          "https://api.run.pez.pivotal.io",
			Username:          "",
			Password:          "",
			SkipSslValidation: true,
		},
		AutoscalerAPIUrl: "https://autoscale.run.pez.pivotal.io/api",
		InstanceGUID: "c017ec06-cf4c-42fa-adbd-1b6a290d8d6a",
	}
	Context("Given a Autoscaler Client", func() {

		It("Should be able to query for information", func() {
			client, err := NewAutoScalerClient(config)
			Ω(err).Should(BeNil())

			bindings, err := client.GetServiceBindings()
			Ω(err).Should(BeNil())
			fmt.Printf("%+v", bindings)
		})
	})
})

