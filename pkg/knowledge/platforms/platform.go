package platforms

import (
	"fmt"
	"strings"

	"github.com/openshift-splat-team/splat-bot/data"
)

func GetPlatformTermsVSphere() data.TokenMatch {
	return data.TokenMatch{
		Tokens: []string{"vsphere", "vmware", "vcenter"},
		Type:   "or",
	}
}

func GetPlatformTermsAWS() data.TokenMatch {
	return data.TokenMatch{
		Tokens: []string{"aws", "ec2"},
		Type:   "or",
	}
}

func GetInstallTerms() data.TokenMatch {
	return data.TokenMatch{
		Tokens: []string{"install", "installation", "ipi", "upi", "install-config"},
		Type:   "or",
	}
}

// GetPathContextExpr returns the platform expressions for a given path
// if unknown, it returns nil
func GetPathContextExpr(path string) string {
	var additionalTerms []data.TokenMatch
	if strings.Contains(path, "vmware") {
		additionalTerms = append(additionalTerms, GetPlatformTermsVSphere())
	} else if strings.Contains(path, "aws") {
		additionalTerms = append(additionalTerms, GetPlatformTermsAWS())
	}
	if strings.Contains(path, "install") {
		additionalTerms = append(additionalTerms, GetInstallTerms())
	}

	expressions := []string{}
	for _, term := range additionalTerms {
		wrapped := []string{}
		for _, term := range term.Tokens {
			wrapped = append(wrapped, fmt.Sprintf("\"%s\"", term))
		}
		expressions = append(expressions, fmt.Sprintf("containsAny(tokens, [%s])", strings.Join(wrapped, ",")))
	}

	return strings.Join(expressions, " and ")
}

// GetPathContextTerms returns the platform terms for a given path
// if unknown, it returns nil
func GetPathContextTerms(path string) []data.TokenMatch {
	var additionalTerms []data.TokenMatch
	if strings.Contains(path, "vmware") {
		additionalTerms = append(additionalTerms, GetPlatformTermsVSphere())
	} else if strings.Contains(path, "aws") {
		additionalTerms = append(additionalTerms, GetPlatformTermsAWS())
	}
	if strings.Contains(path, "install") {
		additionalTerms = append(additionalTerms, GetInstallTerms())
	}
	return additionalTerms
}
