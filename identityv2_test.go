package restclient_test

import (
	"github.com/racker/go-restclient"
	"log"
)

func ExampleIdentityV2Authenticator() {
	authenticator, err := restclient.IdentityV2Authenticator(
		"https://identity.api.rackspacecloud.com",
		"username", "", "apikey")
	if err != nil {
		log.Fatal(err)
	}

	client := restclient.New()
	client.AddInterceptor(authenticator)

	// calls to client.Exchange will get x-auth-token auto populated by interceptor

	// Output:
	//
}
