package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type APIMetrics struct {
	SuccessfulHttpRequests prometheus.Counter
	ErrorHttpRequests      prometheus.Counter

	SuccessSessionVerification prometheus.Counter
	ErrorSessionVerification   prometheus.Counter

	// user
	VerifyEmailRequests prometheus.Counter

	// auth
	SignInRequests       prometheus.Counter
	SignUpRequests       prometheus.Counter
	RefreshRequest       prometheus.Counter
	LogoutRequest        prometheus.Counter
	VerifySessionRequest prometheus.Counter

	// oauth
	GoogleRedirectRequest prometheus.Counter
	GoogleCallBackRequest prometheus.Counter
	GmailAuthRequest      prometheus.Counter
}

// NewAPIMetrics creates a new instance of APIMetrics with Prometheus counters initialized.
func NewAPIMetrics() *APIMetrics {
	serviceName := "auth_api"

	return &APIMetrics{
		SuccessfulHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_http_requests", serviceName),
			Help: "The total number of successful http requests",
		}),
		ErrorHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_http_requests", serviceName),
			Help: "The total number of unsuccessful http requests",
		}),
		SuccessSessionVerification: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_session_verification_requests", serviceName),
			Help: "The total number of successful session verification http requests",
		}),
		ErrorSessionVerification: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_session_verification_requests", serviceName),
			Help: "The total number of unsuccessful session verification http requests",
		}),
		VerifyEmailRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_verify_email_requests", serviceName),
			Help: "The total number of verify email http requests",
		}),
		SignInRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_sign_in_requests", serviceName),
			Help: "The total number of sign in http requests",
		}),
		SignUpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_sign_up_requests", serviceName),
			Help: "The total number of sign up http requests",
		}),
		RefreshRequest: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_refresh_requests", serviceName),
			Help: "The total number of refresh http requests",
		}),
		LogoutRequest: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_logout_requests", serviceName),
			Help: "The total number of logout http requests",
		}),
		VerifySessionRequest: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_verify_session_requests", serviceName),
			Help: "The total number of verify session http requests",
		}),
		GoogleRedirectRequest: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_google_redirect_requests", serviceName),
			Help: "The total number of Google redirect http requests",
		}),
		GoogleCallBackRequest: prometheus.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_google_callback_requests", serviceName),
			Help: "The total number of Google callback http requests",
		}),
		GmailAuthRequest: prometheus.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_gmail_auth_requests", serviceName),
			Help: "The total number of Gmail Auth http requests",
		}),
	}
}
