package rundeck

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	multierror "github.com/hashicorp/go-multierror"
	responses "github.com/lusis/go-rundeck/pkg/rundeck.v21/responses"
)

// ACLPolicies represents ACL Policies
type ACLPolicies responses.ACLResponse

// GetACLPolicies gets the system ACL Policies
func (c *Client) GetACLPolicies() (*ACLPolicies, error) {
	data := &ACLPolicies{}
	res, err := c.httpGet("system/acl/", requestJSON(), requestExpects(200))
	if err != nil {
		return nil, err
	}
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}
	return data, nil
}

// GetACLPolicy returns the named acl policy
func (c *Client) GetACLPolicy(policy string) ([]byte, error) {
	url := fmt.Sprintf("system/acl/%s.aclpolicy", policy)
	res, err := c.httpGet(url, accept("application/yaml"), requestExpects(200))
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CreateACLPolicy creates a system acl policy
func (c *Client) CreateACLPolicy(name string, contents io.Reader) error {
	url := fmt.Sprintf("system/acl/%s.aclpolicy", name)
	res, err := c.httpPost(url, withBody(contents), accept("application/json"), contentType("application/yaml"), requestExpects(201))
	if err != nil {
		jsonError := &responses.FailedACLValidationResponse{}
		jsonErr := json.Unmarshal(res, jsonError)
		if jsonErr != nil {
			// just return the original error
			return err
		}
		var finalErr error
		for _, v := range jsonError.Policies {
			line := fmt.Sprintf("%s: %s", v.Policy, strings.Join(v.Errors, ","))
			finalErr = multierror.Append(finalErr, fmt.Errorf("%s", line))
		}
		return &PolicyValidationError{msg: finalErr.Error()}
	}
	return nil
}

// UpdateACLPolicy creates a system acl policy
func (c *Client) UpdateACLPolicy(name string, contents io.Reader) error {
	url := fmt.Sprintf("system/acl/%s.aclpolicy", name)
	res, err := c.httpPut(url, withBody(contents), accept("application/json"), contentType("application/yaml"), requestExpects(201))
	if err != nil {
		jsonError := &responses.FailedACLValidationResponse{}
		jsonErr := json.Unmarshal(res, jsonError)
		if jsonErr != nil {
			// just return the original error
			return err
		}
		var finalErr error
		for _, v := range jsonError.Policies {
			line := fmt.Sprintf("%s: %s", v.Policy, strings.Join(v.Errors, ","))
			finalErr = multierror.Append(finalErr, fmt.Errorf("%s", line))
		}
		return &PolicyValidationError{msg: finalErr.Error()}
	}
	return nil
}
