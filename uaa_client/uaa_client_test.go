package uaa_client_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/bijukunjummen/app-autoscaler-client/uaa_client"
	"net/http"
	"github.com/onsi/gomega/ghttp"
)

//AccessToken - represents a authenticated token from UAA
type AccessToken struct {
	Token     string `json:"access_token"`
	Type      string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
	Scope     string
	JTI       string
}

var _ = Describe("UAA Client", func() {


	Context("Given a sample UAA Client", func() {
		var server *ghttp.Server
		var config Config

		BeforeEach(func() {
			server = ghttp.NewServer()

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/info"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, Endpoint {
						AuthorizationEndpoint: server.URL(),
						TokenEndpoint: server.URL(),
					}),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/oauth/token"),
					ghttp.VerifyBasicAuth("cf", ""),
					ghttp.RespondWithJSONEncoded(http.StatusOK, AccessToken{
						Token: "test-token",
					}),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/organizations"),

					ghttp.VerifyHeader(http.Header{
						"Authorization": []string{"Bearer test-token"},
					}),
					ghttp.RespondWithJSONEncoded(http.StatusOK, AccessToken{
						Token: "test-token",
					}),
				),
			)
			config = Config{
				CCApiUrl: server.URL() ,
				Username: "admin",
				Password: "admin",
				SkipSslValidation: true,
			}

		})

		AfterEach(func() {
			server.Close()
		})

		It("Should be able to get the token for the user", func() {

			client, err := NewClient(&config)

			Ω(err).Should(BeNil())

			request, err := client.NewCCRequest("GET", "/v2/organizations", nil)
			Ω(err).Should(BeNil())

			_, err = client.Do(request)

			Ω(err).Should(BeNil())

		})


	})
})