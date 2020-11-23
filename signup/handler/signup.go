package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"

	eproto "github.com/embedscript/backend/emails/proto"
	aproto "github.com/m3o/services/alert/proto/alert"
	cproto "github.com/m3o/services/customers/proto"
	inviteproto "github.com/m3o/services/invite/proto"
	nproto "github.com/m3o/services/namespaces/proto"
	pproto "github.com/m3o/services/payments/proto"
	signup "github.com/m3o/services/signup/proto/signup"
	sproto "github.com/m3o/services/subscriptions/proto"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/client"
	mconfig "github.com/micro/micro/v3/service/config"
	merrors "github.com/micro/micro/v3/service/errors"
	logger "github.com/micro/micro/v3/service/logger"
	mstore "github.com/micro/micro/v3/service/store"
)

const (
	namespace          = "backend"
	internalErrorMsg   = "An error occurred during signup. Contact #m3o-support at slack.m3o.com if the issue persists"
	notInvitedErrorMsg = "You have not been invited to the service. Please request an invite on m3o.com"
)

const (
	expiryDuration      = 5 * time.Minute
	prefixPaymentMethod = "payment-method/"

	signupTopic = "signup"
)

type tokenToEmail struct {
	Email      string `json:"email"`
	Token      string `json:"token"`
	Created    int64  `json:"created"`
	CustomerID string `json:"customerID"`
}

type Signup struct {
	inviteService       inviteproto.InviteService
	customerService     cproto.CustomersService
	namespaceService    nproto.NamespacesService
	subscriptionService sproto.SubscriptionsService
	alertService        aproto.AlertService
	paymentService      pproto.ProviderService
	emailService        eproto.EmailsService
	auth                auth.Auth
	config              conf
	cache               *cache.Cache
}

var (
	// TODO: move this message to a better location
	// Message is a predefined message returned during signup
	Message = "Please complete signup at https://m3o.com/subscribe?email=%s. This command will now wait for you to finish."
)

type sendgridConf struct {
	TemplateID         string `json:"template_id"`
	RecoveryTemplateID string `json:"recovery_template_id"`
}

type conf struct {
	TestMode       bool         `json:"test_env"`
	PaymentMessage string       `json:"message"`
	Sendgrid       sendgridConf `json:"sendgrid"`
	// using a negative "nopayment" rather than "paymentrequired" because it will default to having to pay if not set
	NoPayment bool `json:"no_payment"`
}

func NewSignup(srv *service.Service, auth auth.Auth) *Signup {
	c := conf{}
	val, err := mconfig.Get("micro.signup")
	if err != nil {
		logger.Warnf("Error getting config: %v", err)
	}
	err = val.Scan(&c)
	if err != nil {
		logger.Warnf("Error scanning config: %v", err)
	}

	if len(strings.TrimSpace(c.PaymentMessage)) == 0 {
		c.PaymentMessage = Message
	}
	if !c.TestMode && len(c.Sendgrid.TemplateID) == 0 {
		logger.Fatalf("No sendgrid template ID provided")
	}

	s := &Signup{
		inviteService:       inviteproto.NewInviteService("invite", srv.Client()),
		customerService:     cproto.NewCustomersService("customers", srv.Client()),
		namespaceService:    nproto.NewNamespacesService("namespaces", srv.Client()),
		subscriptionService: sproto.NewSubscriptionsService("subscriptions", srv.Client()),
		paymentService:      pproto.NewProviderService("payments", srv.Client()),
		emailService:        eproto.NewEmailsService("emails", srv.Client()),
		auth:                auth,
		config:              c,
		cache:               cache.New(1*time.Minute, 5*time.Minute),
		alertService:        aproto.NewAlertService("alert", srv.Client()),
	}
	return s
}

// taken from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func randStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// SendVerificationEmail is the first step in the signup flow.SendVerificationEmail
// A stripe customer and a verification token will be created and an email sent.
func (e *Signup) SendVerificationEmail(ctx context.Context,
	req *signup.SendVerificationEmailRequest,
	rsp *signup.SendVerificationEmailResponse) error {
	err := e.sendVerificationEmail(ctx, req, rsp)
	if err != nil {
		logger.Warnf("Error during sending verification email: %v", err)
	}
	return err
}

func (e *Signup) sendVerificationEmail(ctx context.Context,
	req *signup.SendVerificationEmailRequest,
	rsp *signup.SendVerificationEmailResponse) error {
	logger.Info("Received Signup.SendVerificationEmail request")

	k := randStringBytesMaskImprSrc(8)
	tok := &tokenToEmail{
		Token:   k,
		Email:   req.Email,
		Created: time.Now().Unix(),
		// @todo this is wrong
		CustomerID: req.Email,
	}

	bytes, err := json.Marshal(tok)
	if err != nil {
		logger.Error(err)
		return merrors.InternalServerError("signup.SendVerificationEmail", internalErrorMsg)
	}

	if err := mstore.Write(&mstore.Record{
		Key:   req.Email,
		Value: bytes,
	}); err != nil {
		logger.Error(err)
		return merrors.InternalServerError("signup.SendVerificationEmail", internalErrorMsg)
	}
	// HasPaymentMethod needs to resolve email from token, so we save the
	// same record under a token too
	if err := mstore.Write(&mstore.Record{
		Key:   tok.Token,
		Value: bytes,
	}); err != nil {
		logger.Error(err)
		return merrors.InternalServerError("signup.SendVerificationEmail", internalErrorMsg)
	}

	if e.config.TestMode {
		logger.Infof("Sending verification token '%v'", k)
	}

	// Send email
	// @todo send different emails based on if the account already exists
	// ie. registration vs login email.

	err = e.sendEmail(ctx, req.Email, e.config.Sendgrid.TemplateID, map[string]interface{}{
		"token": k,
	})
	if err != nil {
		logger.Errorf("Error when sending email to %v: %v", req.Email, err)
		return merrors.InternalServerError("signup.SendVerificationEmail", internalErrorMsg)
	}

	return nil
}

// Lifted  from the invite service https://github.com/m3o/services/blob/master/projects/invite/handler/invite.go#L187
// sendEmailInvite sends an email invite via the sendgrid API using the
// predesigned email template. Docs: https://bit.ly/2VYPQD1
func (e *Signup) sendEmail(ctx context.Context, email, templateID string, templateData map[string]interface{}) error {
	b, _ := json.Marshal(templateData)
	_, err := e.emailService.Send(context.TODO(), &eproto.SendRequest{To: email, TemplateId: templateID, TemplateData: b}, client.WithAuthToken())
	return err
}

func (e *Signup) Verify(ctx context.Context, req *signup.VerifyRequest, rsp *signup.VerifyResponse) error {
	err := e.verify(ctx, req, rsp)
	if err != nil {
		logger.Warnf("Error during verification: %v", err)
	}
	return err
}

func (e *Signup) verify(ctx context.Context, req *signup.VerifyRequest, rsp *signup.VerifyResponse) error {
	logger.Info("Received Signup.Verify request")

	recs, err := mstore.Read(req.Email)
	if err == mstore.ErrNotFound {
		logger.Errorf("Can't verify, record for %v is not found", req.Email)
		return merrors.InternalServerError("signup.Verify", internalErrorMsg)
	} else if err != nil {
		logger.Errorf("email verification error: %v", err)
		return merrors.InternalServerError("signup.Verify", internalErrorMsg)
	}

	tok := &tokenToEmail{}
	if err := json.Unmarshal(recs[0].Value, tok); err != nil {
		return err
	}

	if tok.Token != req.Token {
		return merrors.Forbidden("signup.Verify", "The token you provided is invalid")
	}

	if time.Since(time.Unix(tok.Created, 0)) > expiryDuration {
		return merrors.Forbidden("signup.Verify", "The token you provided has expired")
	}

	// set the response message
	rsp.Message = fmt.Sprintf(e.config.PaymentMessage, req.Email)

	rsp.Namespaces = []string{namespace}
	rsp.CustomerID = tok.CustomerID

	return nil
}

func (e *Signup) CompleteSignup(ctx context.Context, req *signup.CompleteSignupRequest, rsp *signup.CompleteSignupResponse) error {
	err := e.completeSignup(ctx, req, rsp)
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (e *Signup) completeSignup(ctx context.Context, req *signup.CompleteSignupRequest, rsp *signup.CompleteSignupResponse) error {
	logger.Info("Received Signup.CompleteSignup request")

	recs, err := mstore.Read(req.Email)
	if err == mstore.ErrNotFound {
		logger.Errorf("Can't verify record for %v: record not found", req.Email)
		return merrors.InternalServerError("signup.CompleteSignup", internalErrorMsg)
	} else if err != nil {
		logger.Errorf("Error reading store: err")
		return merrors.InternalServerError("signup.CompleteSignup", internalErrorMsg)
	}

	tok := &tokenToEmail{}
	if err := json.Unmarshal(recs[0].Value, tok); err != nil {
		logger.Errorf("Error when unmarshaling stored token object for %v: %v", req.Email, err)
		return merrors.InternalServerError("signup.CompleteSignup", internalErrorMsg)
	}
	if tok.Token != req.Token { // not checking expiry here because we've already checked it during Verify() step
		return merrors.Forbidden("signup.CompleteSignup.invalid_token", "The token you provided is incorrect")
	}

	rsp.Namespace = namespace

	// take secret from the request
	secret := req.Secret

	// generate a random secret
	if len(req.Secret) == 0 {
		secret = uuid.New().String()
	}
	_, err = e.auth.Generate(tok.CustomerID, auth.WithSecret(secret), auth.WithIssuer(namespace), auth.WithName(req.Email))
	if err != nil {
		logger.Errorf("Error generating token for %v: %v", tok.CustomerID, err)
		return merrors.InternalServerError("signup.CompleteSignup", internalErrorMsg)
	}

	t, err := e.auth.Token(auth.WithCredentials(tok.CustomerID, secret), auth.WithTokenIssuer(namespace))
	if err != nil {
		logger.Errorf("Can't get token for %v: %v", tok.CustomerID, err)
		return merrors.InternalServerError("signup.CompleteSignup", internalErrorMsg)
	}
	rsp.AuthToken = &signup.AuthToken{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expiry.Unix(),
		Created:      t.Created.Unix(),
	}

	return nil
}

func (e *Signup) Recover(ctx context.Context, req *signup.RecoverRequest, rsp *signup.RecoverResponse) error {
	logger.Info("Received Signup.Recover request")
	_, found := e.cache.Get(req.Email)
	if found {
		return merrors.BadRequest("signup.recover", "We have issued a recovery email recently. Please check that.")
	}

	custResp, err := e.customerService.Read(ctx, &cproto.ReadRequest{Email: req.Email}, client.WithAuthToken())
	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == 404 { // not found
			return merrors.NotFound("signup.recover", "Could not find an account with that email address")
		}
		return merrors.InternalServerError("signup.recover", "Error looking up account")
	}

	listRsp, err := e.namespaceService.List(ctx, &nproto.ListRequest{
		User: custResp.Customer.Id,
	}, client.WithAuthToken())
	if err != nil {
		return merrors.InternalServerError("signup.recover", "Error calling namespace service: %v", err)
	}
	if len(listRsp.Namespaces) == 0 {
		return merrors.BadRequest("signup.recover", "We don't recognize this account")
	}

	// Sendgrid wants objects in a list not string
	namespaces := []map[string]string{}
	for _, v := range listRsp.Namespaces {
		namespaces = append(namespaces, map[string]string{
			"id": v.Id,
		})
	}

	logger.Infof("Sending email with data %v", namespaces)
	err = e.sendEmail(ctx, req.Email, e.config.Sendgrid.RecoveryTemplateID, map[string]interface{}{
		"namespaces": namespaces,
	})
	if err == nil {
		e.cache.Set(req.Email, true, cache.DefaultExpiration)
	}

	return err
}
