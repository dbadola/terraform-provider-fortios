package forticlient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type ZoneMultiValue struct {
	InterfaceName string `json:"interface-name"`
}

type ZoneMultiValues []ZoneMultiValue

// ExpandZone extracts Zone value from result and put them into a string array,
// and return the string array
func ExpandZone(members []ZoneMultiValue) []string {
	vs := make([]string, 0, len(members))
	for _, v := range members {
		c := v.InterfaceName
		vs = append(vs, c)
	}
	return vs
}

// JSONSystemZone contains the parameters for Create and Update API function
type JSONSystemZone struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Intrazone   string          `json:"intrazone"`
	Interface   ZoneMultiValues `json:"interface"`
}

// JSONCreateSystemZoneOutput contains the output results for Create API function
type JSONCreateSystemZoneOutput struct {
	Vdom       string  `json:"vdom"`
	Mkey       string  `json:"mkey"`
	Status     string  `json:"status"`
	HTTPStatus float64 `json:"http_status"`
}

// JSONUpdateSystemZoneOutput contains the output results for Update API function
// Attention: Considering scalability, the previous structure and the current structure may change differently
type JSONUpdateSystemZoneOutput struct {
	Vdom       string  `json:"vdom"`
	Mkey       string  `json:"mkey"`
	Status     string  `json:"status"`
	HTTPStatus float64 `json:"http_status"`
}

// CreateSystemZone API operation for FortiOS
func (c *FortiSDKClient) CreateSystemZone(params *JSONSystemZone) (output *JSONCreateSystemZoneOutput, err error) {
	return
}

// UpdateSystemZone API operation for FortiOS set the dns server.
// Returns the execution result when the request executes successfully.
// Returns error for service API and SDK errors.
// See the system - dns chapter in the FortiOS Handbook - CLI Reference.
func (c *FortiSDKClient) UpdateSystemZone(params *JSONSystemZone, mkey string) (output *JSONUpdateSystemZoneOutput, err error) {
	HTTPMethod := "POST"
	path := "/api/v2/cmdb/system/zone"

	output = &JSONUpdateSystemZoneOutput{}
	locJSON, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
		return
	}
	bytes := bytes.NewBuffer(locJSON)
	req := c.NewRequest(HTTPMethod, path, nil, bytes)
	err = req.Send()
	if err != nil || req.HTTPResponse == nil {
		err = fmt.Errorf("cannot send request %s", err)
		return
	}

	body, err := ioutil.ReadAll(req.HTTPResponse.Body)
	if err != nil || body == nil {
		err = fmt.Errorf("cannot get response body %s", err)
		return
	}
	log.Printf("FOS-fortios response: %s", string(body))

	var result map[string]interface{}
	json.Unmarshal([]byte(string(body)), &result)

	req.HTTPResponse.Body.Close()

	if result != nil {
		if result["vdom"] != nil {
			output.Vdom = result["vdom"].(string)
		}
		if result["mkey"] != nil {
			output.Mkey = result["mkey"].(string)
		}
		if result["status"] != nil {
			if result["status"] != "success" {
				if result["error"] != nil {
					err = fmt.Errorf("status is %s and error no is %.0f", result["status"], result["error"])
				} else {
					err = fmt.Errorf("status is %s and error no is not found", result["status"])
				}

				if result["http_status"] != nil {
					err = fmt.Errorf("%s and http_status no is %.0f", err, result["http_status"])
				} else {
					err = fmt.Errorf("%s and and http_status no is not found", err)
				}

				return
			}
			output.Status = result["status"].(string)
		} else {
			err = fmt.Errorf("cannot get status from the response")
			return
		}
		if result["http_status"] != nil {
			output.HTTPStatus = result["http_status"].(float64)
		}
	} else {
		err = fmt.Errorf("cannot get the right response")
		return
	}

	return

}

// DeleteSystemZone API operation for FortiOS
func (c *FortiSDKClient) DeleteSystemZone(mkey string) (err error) {
	return
}

// ReadSystemZone API operation for FortiOS gets the dns server setting.
// Returns the requested dns server value when the request executes successfully.
// Returns error for service API and SDK errors.
// See the system - dns chapter in the FortiOS Handbook - CLI Reference.
func (c *FortiSDKClient) ReadSystemZone(mkey string) (output *JSONSystemZone, err error) {
	HTTPMethod := "GET"
	path := "/api/v2/cmdb/system/zone"

	output = &JSONSystemZone{}
	req := c.NewRequest(HTTPMethod, path, nil, nil)
	err = req.Send()
	if err != nil || req.HTTPResponse == nil {
		err = fmt.Errorf("cannot send request %s", err)
		return
	}

	body, err := ioutil.ReadAll(req.HTTPResponse.Body)
	if err != nil || body == nil {
		err = fmt.Errorf("cannot get response body %s", err)
		return
	}
	log.Printf("FOS-fortios reading response: %s", string(body))

	var result map[string]interface{}
	json.Unmarshal([]byte(string(body)), &result)

	req.HTTPResponse.Body.Close()

	if result != nil {
		if result["http_status"] == nil {
			err = fmt.Errorf("cannot get http_status from the response")
			return
		}

		if result["http_status"] == 404.0 {
			output = nil
			return
		}

		if result["status"] == nil {
			err = fmt.Errorf("cannot get status from the response")
			return
		}

		if result["status"] != "success" {
			if result["error"] != nil {
				err = fmt.Errorf("status is %s and error no is %.0f", result["status"], result["error"])
			} else {
				err = fmt.Errorf("status is %s and error no is not found", result["status"])
			}

			if result["http_status"] != nil {
				err = fmt.Errorf("%s and http_status no is %.0f", err, result["http_status"])
			} else {
				err = fmt.Errorf("%s and and http_status no is not found", err)
			}

			return
		}
		mapTmp := (result["results"].([]interface{}))[0].(map[string]interface{})

		if mapTmp == nil {
			err = fmt.Errorf("cannot get the results from the response")
			return
		}

		if mapTmp["name"] != nil {
			output.Name = mapTmp["name"].(string)
		}
		if mapTmp["description"] != nil {
			output.Description = mapTmp["description"].(string)
		}
		if mapTmp["intrazone"] != nil {
			output.Intrazone = mapTmp["intrazone"].(string)
		}
		if mapTmp["interface"] != nil {
			member := mapTmp["interface"].([]interface{})

			var members []ZoneMultiValue
			for _, v := range member {
				c := v.(map[string]interface{})

				members = append(members,
					ZoneMultiValue{
						InterfaceName: c["interface-name"].(string),
					})
				output.Interface = members
			}
		}
	} else {
		err = fmt.Errorf("cannot get the right response")
		return
	}

	return

}
