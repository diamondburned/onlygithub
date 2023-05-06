package onlygithub

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Config describes the global configuration of the service.
type Config struct {
	// TierWhitelist controls whether or not tiers must be explicitly
	// whitelisted in order to access sponsored content.
	TierWhitelist *TierWhitelist `json:"tierWhitelist,omitempty"`
	// CustomCSS is custom CSS that will be injected into the page.
	// Each page will be identified by a unique ID that can be used to
	// scope the CSS to that page.
	CustomCSS string `json:"customCSS,omitempty"`
}

// TierWhitelist is a list of tiers that are allowed to access sponsored
// content.
type TierWhitelist struct {
	// Requires is a list of requirements that a tier must satisfy in order to
	// be allowed to access sponsored content.
	Requires []TierRequirement
	// AllowCustom, if true, means that custom amounts are allowed to access
	// sponsored content as long as they are greater than any of the tiers
	// specified in Requires.
	AllowCustom bool
}

// TierRequirement is a description of a tier (or what a tier should have) that
// the author desires. At least one of these fields must be set.
//
// Some fields can have values surrounded by slashes to indicate regular
// expressions. For example, if the name is "/^Tier [0-9]+$/", then it will
// match any tier name that starts with "Tier " and ends with a number.
//
// The following fields can have regular expressions:
// - Name
// - Description
//
type TierRequirement struct {
	Name        string `json:"name,omitempty"`
	Price       Cents  `json:"price,omitempty"`
	Description string `json:"description,omitempty"`
	// AllowHigher means that tiers with a higher price than the one specified
	// are allowed to access sponsored content.
	AllowHigher bool `json:"allowHigher,omitempty"`
}

// Validate validates the TierRequirement.
func (r TierRequirement) Validate() error {
	if r.Name == "" && r.Price == 0 && r.Description == "" {
		return fmt.Errorf("at least one of name, price, or description must be set")
	}
	if isRegex(r.Name) {
		_, err := regexp.Compile(trimRegex(r.Name))
		if err != nil {
			return errors.Wrap(err, "invalid name regex")
		}
	}
	if isRegex(r.Description) {
		_, err := regexp.Compile(trimRegex(r.Name))
		if err != nil {
			return errors.Wrap(err, "invalid description regex")
		}
	}
	return nil
}

// Matches returns true if the given tier matches the requirement.
func (r TierRequirement) Matches(t Tier) bool {
	return false ||
		matches(r.Name, t.Name) ||
		matches(r.Description, t.Description) ||
		(r.Price != 0 && t.Price == r.Price) ||
		(r.Price != 0 && r.AllowHigher && r.Price >= r.Price)
}

func isRegex(str string) bool {
	return strings.HasPrefix(str, "/") && strings.HasSuffix(str, "/")
}

func trimRegex(str string) string {
	return strings.TrimPrefix(strings.TrimSuffix(str, "/"), "/")
}

func matches(maybeRegex, str string) bool {
	if maybeRegex == "" {
		return false
	}
	if isRegex(maybeRegex) {
		re := regexp.MustCompile(trimRegex(maybeRegex))
		return re.MatchString(str)
	}
	return maybeRegex == str
}
