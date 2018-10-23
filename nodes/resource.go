package nodes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/StratoAPI/Interface/filter"
	"github.com/StratoAPI/Interface/resource"
	"github.com/StratoAPI/Interface/schema"

	"github.com/Vilsol/GoLib"
	"github.com/gorilla/mux"
)

func RegisterResourceRoutes(router GoLib.RegisterRoute) {
	router("GET", "/resource/{resource}", getResource)
	router("PUT", "/resource/{resource}", updateResource)
	router("POST", "/resource/{resource}", storeResource)
	router("DELETE", "/resource/{resource}", deleteResource)
}

func getResource(r *http.Request) (interface{}, *GoLib.ErrorResponse) {
	resourceName := mux.Vars(r)["resource"]

	if !schema.GetProcessor().ResourceExists(resourceName) {
		return nil, &ErrorResourceDoesNotExist
	}

	resultFilters, err := processFilters(r.URL.Query()["filters"])

	if err != nil {
		return nil, err
	}

	resources, errG := resource.GetProcessor().GetResources(resourceName, resultFilters)

	if errG != nil {
		resp := ErrorFetchingResource
		resp.Message += errG.Error()
		return nil, &resp
	}

	return resources, nil
}

func updateResource(r *http.Request) (interface{}, *GoLib.ErrorResponse) {
	resourceName := mux.Vars(r)["resource"]

	if !schema.GetProcessor().ResourceExists(resourceName) {
		return nil, &ErrorResourceDoesNotExist
	}

	resultFilters, err := processFilters(r.URL.Query()["filters"])

	if err != nil {
		return nil, err
	}

	body, errG := ioutil.ReadAll(r.Body)
	if errG != nil {
		return nil, &ErrorCouldNotReadBody
	}

	valid, errG := schema.GetProcessor().ResourceValid(resourceName, string(body), false)

	if !valid {
		resp := ErrorResourceInvalid
		resp.Message += errG.Error()
		return nil, &resp
	}

	var data map[string]interface{}
	errG = json.Unmarshal(body, &data)

	if errG != nil {
		resp := ErrorResourceInvalid
		resp.Message += errG.Error()
		return nil, &resp
	}

	errG = resource.GetProcessor().UpdateResources(resourceName, data, resultFilters)

	if errG != nil {
		resp := ErrorUpdatingResource
		resp.Message += errG.Error()
		return nil, &resp
	}

	return nil, nil
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

	finalResources := make([]map[string]interface{}, 0)

	if body[0] == '[' {
		var data []map[string]interface{}
		err := json.Unmarshal(body, &data)

		if err != nil {
			resp := ErrorResourceInvalid
			resp.Message += err.Error()
			return nil, &resp
		}

		for _, d := range data {
			valid, err := schema.GetProcessor().ResourceValidGo(resourceName, d, true)

			if !valid {
				resp := ErrorResourceInvalid
				resp.Message += err.Error()
				return nil, &resp
			}

			finalResources = append(finalResources, d)
		}
	} else {
		valid, err := schema.GetProcessor().ResourceValid(resourceName, string(body), true)

		if !valid {
			resp := ErrorResourceInvalid
			resp.Message += err.Error()
			return nil, &resp
		}

		var data map[string]interface{}

		err = json.Unmarshal(body, &data)

		if err != nil {
			resp := ErrorResourceInvalid
			resp.Message += err.Error()
			return nil, &resp
		}

		finalResources = append(finalResources, data)
	}

	errG := (*resource.GetProcessor().GetStore(resourceName)).CreateResources(resourceName, finalResources)

	if errG != nil {
		resp := ErrorCreatingResource
		resp.Message += errG.Error()
		return nil, &resp
	}

	return nil, nil
}

func deleteResource(r *http.Request) (interface{}, *GoLib.ErrorResponse) {
	resourceName := mux.Vars(r)["resource"]

	if !schema.GetProcessor().ResourceExists(resourceName) {
		return nil, &ErrorResourceDoesNotExist
	}

	resultFilters, err := processFilters(r.URL.Query()["filters"])

	if err != nil {
		return nil, err
	}

	errG := resource.GetProcessor().DeleteResources(resourceName, resultFilters)

	if errG != nil {
		resp := ErrorDeletingResource
		resp.Message += errG.Error()
		return nil, &resp
	}

	return nil, nil
}

func processFilters(filters []string) ([]filter.ProcessedFilter, *GoLib.ErrorResponse) {
	resultFilters := make([]filter.ProcessedFilter, 0)
	for _, f := range filters {
		var objFilter filter.EncodedFilter
		err := json.Unmarshal([]byte(f), &objFilter)
		if err != nil {
			resp := ErrorFilterInvalid
			resp.Message += err.Error()
			return nil, &resp
		}

		if !filter.GetProcessor().FilterExists(objFilter.Type) {
			return nil, &ErrorFilterDoesntExist
		}

		filterData := filter.GetProcessor().CreateFilter(objFilter.Type)
		err = json.Unmarshal(objFilter.Data, &filterData)

		if err != nil {
			resp := ErrorFilterInvalid
			resp.Message += err.Error()
			return nil, &resp
		}

		processedFilter := filter.ProcessedFilter{
			Type: objFilter.Type,
			Data: filterData,
		}

		valid, err := filter.GetProcessor().ValidateFilter(processedFilter)

		if !valid {
			resp := ErrorFilterInvalid
			if err != nil {
				resp.Message += err.Error()
			}
			return nil, &resp
		}

		resultFilters = append(resultFilters, processedFilter)
	}

	return resultFilters, nil
}
