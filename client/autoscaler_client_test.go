package client_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/bijukunjummen/app-autoscaler-client/client"
	"github.com/bijukunjummen/app-autoscaler-client/uaa_client"
	"fmt"

	"github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("Behavior of Auto Scaler", func() {

	var server *ghttp.Server
	var config *AutoscalerConfig

	BeforeEach(func() {

		server = ghttp.NewServer()
		config = &AutoscalerConfig {
			UAAConfig: &uaa_client.Config {
				CCApiUrl:          server.URL(),
				Username:          "user",
				Password:          "pwd",
				SkipSslValidation: true,
			},
			AutoscalerAPIUrl: server.URL() + "/api",
			InstanceGUID: "c017ec06-cf4c-42fa-adbd-1b6a290d8d6a",
		}
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/v2/info"),
				ghttp.RespondWithJSONEncoded(http.StatusOK, uaa_client.Endpoint {
					AuthorizationEndpoint: server.URL(),
					TokenEndpoint: server.URL(),
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
				ghttp.RespondWithJSONEncoded(http.StatusOK, uaa_client.AccessToken{
					Token: "test-token",
				}),
			),
		)
	})

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

