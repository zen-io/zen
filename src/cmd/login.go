package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/spf13/cobra"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// Time before MFA step times out
const MFA_TIMEOUT = 30

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to AWS",
	Long:  `Execute an SSO login to aws using the provided profile [login]`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		profile, _ := cmd.Flags().GetString("profile")
		url, err := GetUrl(profile)
		if err != nil {
			eng.Errorf("getting login url: %w", err)
			return
		}

		ssoLogin(url)
	},
}

func GetUrl(profile string) (string, error) {
	cfg, err := config.LoadSharedConfigProfile(context.TODO(), profile)
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	ssooidcClient := ssooidc.New(ssooidc.Options{Region: cfg.SSORegion})
	if err != nil {
		return "", fmt.Errorf("creating ssooidc client: %w", err)
	}

	// register your client which is triggering the login flow
	register, err := ssooidcClient.RegisterClient(context.TODO(), &ssooidc.RegisterClientInput{
		ClientName: aws.String("sample-client-name"),
		ClientType: aws.String("public"),
		Scopes:     []string{"sso-portal:*"},
	})
	if err != nil {
		return "", fmt.Errorf("registering ssooidc client: %w", err)
	}

	// authorize your device using the client registration response
	deviceAuth, err := ssooidcClient.StartDeviceAuthorization(context.TODO(), &ssooidc.StartDeviceAuthorizationInput{
		ClientId:     register.ClientId,
		ClientSecret: register.ClientSecret,
		StartUrl:     aws.String(cfg.SSOStartURL),
	})
	if err != nil {
		return "", fmt.Errorf("creating device auth: %w", err)
	}

	return aws.ToString(deviceAuth.VerificationUriComplete), nil
}

// // login with hardware MFA
func ssoLogin(url string) {
	u := launcher.NewUserMode().Headless(false).Bin("/Applications/Brave Browser.app/Contents/MacOS/Brave Browser").MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	// loadCookies(*browser)
	defer browser.MustClose()

	err := rod.Try(func() {
		page := browser.MustPage(url)

		page.MustElement("li.JDAKTe").MustClick()

		// allow request
		unauthorized := true
		for unauthorized {

			txt := page.Timeout(MFA_TIMEOUT * time.Second).MustElement(".awsui-util-mb-s").MustWaitLoad().MustText()
			if txt == "Request approved" {
				unauthorized = false
			} else {
				exists, _, _ := page.HasR("button", "Allow")
				if exists {
					page.MustWaitLoad().MustElementR("button", "Allow").MustClick()
				}

				time.Sleep(500 * time.Millisecond)
			}
		}

		// saveCookies(*browser)
	})

	if errors.Is(err, context.DeadlineExceeded) {
		panic("Timed out waiting for MFA")
	} else if err != nil {
		panic(err.Error())
	}
}
