package nodes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ResourceAPI/Interface/resource"
	"github.com/ResourceAPI/Interface/schema"
	"github.com/Vilsol/GoLib"
	"github.com/gorilla/mux"
)

func RegisterResourceRoutes(router GoLib.RegisterRoute) {
	router("GET", "/resource/{resource}", getResource)
	router("POST", "/resource/{resource}", storeResource)
	router("DELETE", "/resource/{resource}", deleteResource)
}

func getResource(r *http.Request) (interface{}, *GoLib.ErrorResponse) {
	resourceName := mux.Vars(r)["resource"]

	if !schema.GetProcessor().ResourceExists(resourceName) {
		return nil, &ErrorResourceDoesNotExist
	}

	// TODO Filters
	resources, _ := (*resource.GetProcessor().GetStore(resourceName)).GetResources(resourceName, nil)
	return resources, nil
}

func storeResource(r *http.Request) (interface{}, *GoLib.ErrorResponse) {
	resourceName := mux.Vars(r)["resource"]

	if !schema.GetProcessor().ResourceExists(resourceName) {
		return nil, &ErrorResourceDoesNotExist
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, &ErrorCouldNotReadBody
	}

	valid, err := schema.GetProcessor().ResourceValid(resourceName, string(body))

	if !valid {
		resp := ErrorResourceInvalid
		resp.Message += err.Error()
		return nil, &resp
	}

	var data map[string]interface{}
	json.Unmarshal(body, &data)

	// TODO Filters
	(*resource.GetProcessor().GetStore(resourceName)).CreateResources(resourceName, []map[string]interface{}{data})

	return nil, nil
}

func deleteResource(_ *http.Request) (interface{}, *GoLib.ErrorResponse) {
	return nil, nil
}
