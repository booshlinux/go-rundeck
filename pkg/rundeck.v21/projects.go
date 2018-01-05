package rundeck

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"

	requests "github.com/lusis/go-rundeck/pkg/rundeck.v21/requests"
	responses "github.com/lusis/go-rundeck/pkg/rundeck.v21/responses"
)

// Project is a rundeck project
type Project struct {
	URL         string
	Name        string
	Description string
	Properties  map[string]string
}

// Projects is a collection of `Project`
type Projects []*Project

// NewProject represents a new project
type NewProject struct {
	XMLName     xml.Name         `xml:"project"`
	Name        string           `xml:"name"`
	Description string           `xml:"description"`
	Config      []ConfigProperty `xml:"config>property,omitempty"`
}

// ConfigProperty is a configuration property
type ConfigProperty struct {
	XMLName xml.Name `xml:"property"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:"value,attr"`
}

// GetProject gets a project by name
func (c *Client) GetProject(name string) (*Project, error) {
	p := &responses.ProjectInfoResponse{}
	res, err := c.httpGet("project/"+name, requestJSON(), requestExpects(200))
	if err != nil {
		return nil, err
	}
	if jsonErr := json.Unmarshal(res, &p); jsonErr != nil {
		return nil, errDecoding
	}

	project := &Project{
		URL:         p.URL,
		Name:        p.Name,
		Description: p.Description,
		Properties:  *p.Config,
	}
	return project, nil
}

// ListProjects lists all projects
func (c *Client) ListProjects() (*Projects, error) {
	data := &responses.ListProjectsResponse{}
	res, err := c.httpGet("projects", requestJSON(), requestExpects(200))
	if err != nil {
		return nil, err
	}
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, errDecoding
	}
	projects := &Projects{}
	for _, p := range *data {
		*projects = append(*projects, &Project{
			URL:         p.URL,
			Name:        p.Name,
			Description: p.Description,
		})
	}
	return projects, nil
}

// CreateProject makes a project
func (c *Client) CreateProject(name string, properties map[string]string) (*Project, error) {
	req := &requests.ProjectCreationRequest{
		Name:   name,
		Config: &properties,
	}
	data, _ := json.Marshal(req)
	info := &responses.ProjectInfoResponse{}
	res, postErr := c.httpPost("projects", requestJSON(), withBody(bytes.NewReader(data)), requestExpects(201))
	if postErr != nil {
		return nil, postErr
	}
	if jsonErr := json.Unmarshal(res, &info); jsonErr != nil {
		return nil, errDecoding
	}
	project := &Project{
		URL:         info.URL,
		Name:        info.Name,
		Description: info.Description,
		Properties:  *info.Config,
	}
	return project, nil
}

// DeleteProject deletes a project
func (c *Client) DeleteProject(p string) error {
	url := fmt.Sprintf("project/%s", p)
	return c.httpDelete(url, requestJSON(), requestExpects(204))
}

// TODO: ProjectConfiguration
// http://rundeck.org/docs/api/index.html#project-configuration
