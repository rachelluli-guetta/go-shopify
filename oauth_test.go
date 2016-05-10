package goshopify

import (
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestAppAuthorizeUrl(t *testing.T) {
	setup()
	defer teardown()

	cases := []struct {
		shopName string
		nonce    string
		expected string
	}{
		{"fooshop", "thenonce", "https://fooshop.myshopify.com/admin/oauth/authorize?client_id=apikey&redirect_uri=https%3A%2F%2Fexample.com%2Fcallback&scope=read_products&state=thenonce"},
	}

	for _, c := range cases {
		actual := app.AuthorizeUrl(c.shopName, c.nonce)
		if actual != c.expected {
			t.Errorf("App.AuthorizeUrl(): expected %s, actual %s", c.expected, actual)
		}
	}
}

func TestAppGetAccessToken(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token",
		httpmock.NewStringResponder(200, `{"access_token":"footoken"}`))

	token, err := app.GetAccessToken("fooshop", "foocode")

	if err != nil {
		t.Fatalf("App.GetAccessToken(): %v", err)
	}

	expected := "footoken"
	if token != expected {
		t.Errorf("Token = %v, expected %v", token, expected)
	}
}

func TestAppVerifyAuthorization(t *testing.T) {
	// These credentials are from the Shopify example page:
	// https://help.shopify.com/api/guides/authentication/oauth#verification
	code := "0907a61c0c8d55e99db179b68161bc00"
	shop := "some-shop.myshopify.com"
	timestamp := "1337178173"

	cases := []struct {
		messageHmac string
		expected    bool
	}{
		{"4712bf92ffc2917d15a2f5a273e39f0116667419aa4b6ac0b3baaf26fa3c4d20", true},
		{"notOK", false},
	}

	for _, c := range cases {
		actual := app.VerifyAuthorization(shop, code, timestamp, c.messageHmac)
		if actual != c.expected {
			t.Errorf("App.VerifyAuthorization(..., %s): expected %s, actual %s", c.messageHmac, c.expected, actual)
		}
	}
}

func TestAppVerifyAuthorizationURL(t *testing.T) {
	// These credentials are from the Shopify example page:
	// https://help.shopify.com/api/guides/authentication/oauth#verification
	urlOk, _ := url.Parse("http://example.com/callback?code=0907a61c0c8d55e99db179b68161bc00&hmac=4712bf92ffc2917d15a2f5a273e39f0116667419aa4b6ac0b3baaf26fa3c4d20&shop=some-shop.myshopify.com&signature=11813d1e7bbf4629edcda0628a3f7a20&timestamp=1337178173")
	urlNotOk, _ := url.Parse("http://example.com/callback?code=0907a61c0c8d55e99db179b68161bc00&hmac=4712bf92ffc2917d15a2f5a273e39f0116667419aa4b6ac0b3baaf26fa3c4d20&shop=some-shop.myshopify.com&signature=11813d1e7bbf4629edcda0628a3f7a20&timestamp=133717817")

	cases := []struct {
		u        *url.URL
		expected bool
	}{
		{urlOk, true},
		{urlNotOk, false},
	}

	for _, c := range cases {
		actual := app.VerifyAuthorizationURL(c.u)
		if actual != c.expected {
			t.Errorf("App.VerifyAuthorizationURL(..., %s): expected %s, actual %s", c.u, c.expected, actual)
		}
	}
}