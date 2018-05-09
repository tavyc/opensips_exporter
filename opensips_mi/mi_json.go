package opensips_mi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type miJsonClient struct {
	url    string
	client *http.Client
}

type MIJsonConfig struct {
	HttpClient *http.Client
}

// Create a new Client for OpenSIPS mi_json interface.
func NewMIJsonClient(miJsonUrl string, config MIJsonConfig) (Client, error) {
	_, err := url.Parse(miJsonUrl)
	if err != nil {
		return nil, err
	}

	client := config.HttpClient
	if client == nil {
		client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	return &miJsonClient{
		url:    miJsonUrl,
		client: client,
	}, nil
}

// Execute an OpenSIPS MI commnad and return the resulting tree of MI nodes.
func (mj *miJsonClient) Command(cmd string, args ... string) (*MINode, error) {
	reqUrl := mj.url + "/" + cmd
	if len(args) > 0 {
		query := url.Values{}
		query.Set("params", strings.Join(args, ","))
		reqUrl = reqUrl + "?" + query.Encode()
	}

	// HTTP GET
	resp, err := mj.client.Get(reqUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("mi_json status: %s", resp.StatusCode)
	}

	// Decode the response JSON
	body := map[string]interface{}{}
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&body); err != nil {
		return nil, err
	}

	// Handle errors
	if v, ok := body["error"]; ok {
		if v := v.(map[string]interface{}); ok {
			return nil, fmt.Errorf("mi_json error: %s", v["message"])
		}
		return nil, fmt.Errorf("mi_json error")
	}

	// Parse the MI node tree
	node := &MINode{}
	if err = node.fromJson(body); err != nil {
		return nil, err
	}

	return node, nil
}

func (mj *miJsonClient) Close() error {
	return nil
}

// Convert the OpenSIPS JSON mi_tree representation to a tree of MINodes.
func (n *MINode) fromJson(value interface{}) error {
	switch value.(type) {
	case map[string]interface{}:
		mapval := value.(map[string]interface{})

		isNode := false
		if val, exists := mapval["name"]; exists {
			if s, ok := val.(string); ok {
				n.Name = s
				isNode = true
			}
		}
		if val, exists := mapval["value"]; exists {
			if s, ok := val.(string); ok {
				n.Value = s
				isNode = true
			}
		}
		if val, exists := mapval["attributes"]; exists {
			if m, ok := val.(map[string]interface{}); ok {
				n.Attrs = map[string]string{}
				for k, elem := range m {
					if vs, ok := elem.(string); ok {
						n.Attrs[k] = vs
					}
				}
				isNode = true
			}
		}
		if val, exists := mapval["children"]; exists {
			if lst, ok := val.([]interface{}); ok {
				if err := n.fromJsonList(lst); err != nil {
					return err
				}
				isNode = true
			}
			if mp, ok := val.(map[string]interface{}); ok {
				// parse as map
				n.Children = make([]*MINode, 0, len(mp))
				n.ChildValues = make(map[string]string, len(mapval))
				for k, v := range mp {
					child := &MINode{Name: k}
					if sv, ok := v.(string); ok {
						child.Value = sv
					} else if _, ok := v.(map[string]interface{}); ok {
						if err := child.fromJson(v); err != nil {
							return err
						}
					} else if lst, ok := v.([]interface{}); ok {
						if err := child.fromJsonList(lst); err != nil {
							return err
						}
					} else {
						return fmt.Errorf("Unsupported type in JSON: %+v", reflect.TypeOf(v))
					}
					n.Children = append(n.Children, child)
					if child.Name != "" {
						n.ChildValues[child.Name] = child.Value
					}
				}
				isNode = true
			}
		}

		if !isNode {
			if len(mapval) == 1 {
				for name, val := range mapval {
					if lst, ok := val.([]interface{}); ok {
						// parse as array
						n.Name = name
						if err := n.fromJsonList(lst); err != nil {
							return err
						}
						return nil
					}
				}
			}

			// parse as map
			n.Children = make([]*MINode, 0, len(mapval))
			n.ChildValues = make(map[string]string, len(mapval))
			for k, v := range mapval {
				if vs, ok := v.(string); ok {
					n.Children = append(n.Children, &MINode{Name: k, Value: vs})
					n.ChildValues[k] = vs
				}
			}
		}

	default:
		return fmt.Errorf("Unsupported type in JSON: %+v", reflect.TypeOf(value))
	}

	return nil
}

func (n *MINode) fromJsonList(lst []interface{}) error {
	n.Children = make([]*MINode, 0, len(lst))
	for _, elem := range lst {
		child := &MINode{}
		if err := child.fromJson(elem); err != nil {
			return err
		}
		n.Children = append(n.Children, child)
		if child.Name != "" {
			n.ChildValues[child.Name] = child.Value
		}
	}
	return nil
}
