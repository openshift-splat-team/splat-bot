package util

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/vmware/govmomi/ssoadmin"
	"github.com/vmware/govmomi/ssoadmin/types"
	"github.com/vmware/govmomi/sts"
	"log"
	"math/big"
	"net/url"
	"os"
	"time"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
)

// ClientLogout is empty function that logs out of vSphere clients
type ClientLogout func()

// CreateVSphereClients creates the SOAP and REST client to access
// different portions of the vSphere API
// e.g. tags are only available in REST
func CreateVSphereClients(ctx context.Context, vcenter, username, password string) (*vim25.Client, *rest.Client, ClientLogout, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	u, err := soap.ParseURL(vcenter)
	if err != nil {
		return nil, nil, nil, err
	}
	u.User = url.UserPassword(username, password)
	c, err := govmomi.NewClient(ctx, u, true)

	if err != nil {
		return nil, nil, nil, err
	}

	restClient := rest.NewClient(c.Client)
	err = restClient.Login(ctx, u.User)
	if err != nil {
		logoutErr := c.Logout(context.TODO())
		if logoutErr != nil {
			err = logoutErr
		}
		return nil, nil, nil, err
	}

	return c.Client, restClient, func() {
		//nolint:errcheck
		c.Logout(context.TODO())
		//nolint:errcheck
		restClient.Logout(context.TODO())
	}, nil
}

func getSsoAdminClient(ctx context.Context, user *url.Userinfo, vc *vim25.Client) (*ssoadmin.Client, error) {
	c, err := sts.NewClient(ctx, vc)
	if err != nil {
		return nil, fmt.Errorf("failed to create sts client: %v", err)
	}

	req := sts.TokenRequest{
		Certificate: c.Certificate(),
		Userinfo:    user,
		Renewable:   true,
		Delegatable: true,
	}

	issue := c.Issue

	s, err := issue(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to issue token: %v", err)
	}

	if req.Token != "" {
		duration := s.Lifetime.Expires.Sub(s.Lifetime.Created)
		if duration < req.Lifetime {
			// The granted lifetime is that of the bearer token, which is 5min max.
			// Extend the lifetime via Renew.
			req.Token = s.Token
			if s, err = c.Renew(ctx, req); err != nil {
				return nil, fmt.Errorf("failed to renew token: %v", err)
			}
		}
	}

	admin, err := ssoadmin.NewClient(ctx, vc)
	if err != nil {
		return nil, fmt.Errorf("unable to create vSphere admin client: %v", err)
	}

	header := soap.Header{
		Security: &sts.Signer{
			Certificate: vc.Certificate(),
			Token:       s.Token,
		},
	}

	if err = admin.Login(c.WithHeader(ctx, header)); err != nil {
		return nil, fmt.Errorf("unable to login: %v", err)
	}
	return admin, nil
}

func GetRandomIdentifier(length int) (string, error) {
	const (
		lowerLetters = "abcdefghijkmnopqrstuvwxyz"
		upperLetters = "ABCDEFGHIJKLMNPQRSTUVWXYZ"
		digits       = "23456789"
		all          = lowerLetters + upperLetters + digits
	)
	var password string
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(all))))
		if err != nil {
			return "", err
		}
		newchar := string(all[n.Int64()])
		if password == "" {
			password = newchar
		}
		if i < length-1 {
			n, err = rand.Int(rand.Reader, big.NewInt(int64(len(password)+1)))
			if err != nil {
				return "", err
			}
			j := n.Int64()
			password = password[0:j] + newchar + password[j:]
		}
	}
	pw := []rune(password)
	for _, replace := range []int{5, 11, 17} {
		pw[replace] = '-'
	}
	return string(pw), nil
}

func DeleteUserAccount(ctx context.Context, vcenterUrl, principalUser, vCenterUser, vCenterPass string) error {
	log.Printf("deleting account %s in vcenter %s", principalUser, vcenterUrl)
	vim25Client, _, logout, err := CreateVSphereClients(ctx, vcenterUrl, vCenterUser, vCenterPass)

	if err != nil {
		return fmt.Errorf("unable to create client: %v", err)
	}

	defer logout()

	userInfo := url.UserPassword(vCenterUser, vCenterPass)

	ssoAdminClient, err := getSsoAdminClient(ctx, userInfo, vim25Client)
	if err != nil {
		return fmt.Errorf("unable to create vSphere admin client: %v", err)
	}
	//nolint:errcheck
	defer ssoAdminClient.Logout(ctx)

	err = ssoAdminClient.DeletePrincipal(ctx, principalUser)
	if err != nil {
		return fmt.Errorf("unable to delete vSphere admin client: %v", err)
	}

	return nil
}

func CreateUserAccount(ctx context.Context,
	vcenterUrl,
	domain,
	vCenterUser,
	vCenterPass,
	newUser,
	newPassword,
	group string) error {
	vim25Client, _, logout, err := CreateVSphereClients(ctx, vcenterUrl, vCenterUser, vCenterPass)

	if err != nil {
		fmt.Printf("unable to create client: %v", err)
		os.Exit(1)
	}

	defer logout()

	userInfo := url.UserPassword(vCenterUser, vCenterPass)

	ssoAdminClient, err := getSsoAdminClient(ctx, userInfo, vim25Client)
	if err != nil {
		return fmt.Errorf("unable to create vSphere admin client: %v", err)
	}
	//nolint:errcheck
	defer ssoAdminClient.Logout(ctx)

	adminUsers, err := ssoAdminClient.FindPersonUser(ctx, vCenterUser)
	if err != nil {
		return fmt.Errorf("unable to get admin users: %v", err)
	}

	if domain == "" {
		domain = "vsphere.local"
	}

	user, err := ssoAdminClient.FindPersonUser(ctx, newUser)
	if err != nil {
		log.Printf("unable to get user: %v", err)
	}
	if user != nil {
		log.Printf("user %s already exists", newUser)
		return nil
	}

	err = ssoAdminClient.CreatePersonUser(ctx, newUser, adminUsers.Details, newPassword)
	if err != nil {
		return fmt.Errorf("unable to create user: %v", err)

	}
	err = ssoAdminClient.AddUsersToGroup(ctx, group, types.PrincipalId{
		Name:   newUser,
		Domain: domain,
	})
	if err != nil {
		return fmt.Errorf("unable to add user to group: %v", err)
	}
	return nil
}
