package domain

import (
	"fmt"
	"time"

	"github.com/todennus/shared/errordef"
	"github.com/todennus/x/scope"
	"github.com/xybor-x/snowflake"
)

type OAuth2ConsentResult struct {
	UserID    snowflake.ID
	ClientID  snowflake.ID
	Accepted  bool
	Scope     scope.Scopes
	ExpiresAt time.Time
}

type OAuth2Consent struct {
	UserID    snowflake.ID
	ClientID  snowflake.ID
	Scope     scope.Scopes
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

type OAuth2ConsentDomain struct {
	ConsentSessionExpiration time.Duration
	ConsentExpiration        time.Duration
}

func NewOAuth2ConsentDomain(
	failureExpiration, consentExpiration time.Duration,
) *OAuth2ConsentDomain {
	return &OAuth2ConsentDomain{
		ConsentSessionExpiration: failureExpiration,
		ConsentExpiration:        consentExpiration,
	}
}

func (domain *OAuth2ConsentDomain) NewConsentDeniedResult(userID, clientID snowflake.ID) *OAuth2ConsentResult {
	return &OAuth2ConsentResult{
		Accepted:  false,
		UserID:    userID,
		ClientID:  clientID,
		ExpiresAt: time.Now().Add(domain.ConsentSessionExpiration),
	}
}

func (domain *OAuth2ConsentDomain) NewConsentAcceptedResult(userID, clientID snowflake.ID, scope scope.Scopes) *OAuth2ConsentResult {
	return &OAuth2ConsentResult{
		Accepted:  true,
		UserID:    userID,
		ClientID:  clientID,
		Scope:     scope,
		ExpiresAt: time.Now().Add(domain.ConsentSessionExpiration),
	}
}

func (domain *OAuth2ConsentDomain) NewConsent(
	userID, clientID snowflake.ID,
	scope scope.Scopes,
) *OAuth2Consent {
	return &OAuth2Consent{
		UserID:    userID,
		ClientID:  clientID,
		Scope:     scope,
		ExpiresAt: time.Now().Add(domain.ConsentExpiration),
	}
}

func (domain *OAuth2ConsentDomain) ValidateConsent(consent *OAuth2Consent, requestedScope scope.Scopes) error {
	if consent.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("%wconsent expired", errordef.ErrDomainKnown)
	}

	if !consent.Scope.Contains(requestedScope...) {
		return fmt.Errorf("%wrequested scope exceeds the consented scope", errordef.ErrDomainKnown)
	}

	return nil
}
