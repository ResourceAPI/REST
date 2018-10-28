package nodes

import (
	"encoding/json"
	"fmt"
	"github.com/StratoAPI/Interface/filter"
	"github.com/StratoAPI/Interface/resource"
	"github.com/StratoAPI/Interface/schema"
	"github.com/labstack/echo"
)

func RegisterResourceRoutes(router *echo.Group) {
	router.GET("/resource/:resource", getResource)
	router.PUT("/resource/:resource", updateResource)
	router.POST("/resource/:resource", storeResource)
	router.DELETE("/resource/:resource", deleteResource)
}

func getResource(c echo.Context) error {
	resourceName := c.Param("resource")

	if !schema.GetProcessor().ResourceExists(resourceName) {
		return PrepError(c, ErrorResourceDoesNotExist)
	}

	resultFilters, err := processFilters(c.QueryParams()["filters"])

	if err != nil {
		return PrepError(c, *err)
	}

	resources, errG := resource.GetProcessor().GetResources(resourceName, resultFilters)

	if errG != nil {
		resp := ErrorFetchingResource
		resp.Message += errG.Error()
		return PrepError(c, resp)
	}

	return c.JSON(200, ResponseResource{
		Success: true,
		Data:    &resources,
	})
}

func updateResource(c echo.Context) error {
	resourceName := c.Param("resource")

	if !schema.GetProcessor().ResourceExists(resourceName) {
		return PrepError(c, ErrorResourceDoesNotExist)
	}

	resultFilters, err := processFilters(c.QueryParams()["filters"])

	if err != nil {
		return PrepError(c, *err)
	}

	var rawData map[string]interface{}
	errG := c.Bind(&rawData)

	if errG != nil {
		resp := ErrorCouldNotReadBody
		resp.Message += fmt.Sprintf("%s", errG.(*echo.HTTPError).Message)
		return PrepError(c, resp)
	}

	valid, errG := schema.GetProcessor().ResourceValidGo(resourceName, rawData, false)

	if !valid {
		resp := ErrorResourceInvalid
		resp.Message += errG.Error()
		return PrepError(c, resp)
	}

	errG = resource.GetProcessor().UpdateResources(resourceName, rawData, resultFilters)

	if errG != nil {
		resp := ErrorUpdatingResource
		resp.Message += errG.Error()
		return PrepError(c, resp)
	}

	return c.JSON(200, Response{
		Success: true,
	})
}

func storeResource(c echo.Context) error {
	resourceName := c.Param("resource")

	if !schema.GetProcessor().ResourceExists(resourceName) {
		return PrepError(c, ErrorResourceDoesNotExist)
	}

	var rawData interface{}
	err := c.Bind(&rawData)

	if err != nil {
		resp := ErrorCouldNotReadBody
		resp.Message += fmt.Sprintf("%s", err.(*echo.HTTPError).Message)
		return PrepError(c, resp)
	}

	finalResources := make([]map[string]interface{}, 0)

	if arrData, ok := rawData.([]map[string]interface{}); ok {
		for _, d := range arrData {
			valid, err := schema.GetProcessor().ResourceValidGo(resourceName, d, true)

			if !valid {
				resp := ErrorResourceInvalid
				resp.Message += err.Error()
				return PrepError(c, resp)
			}

			finalResources = append(finalResources, d)
		}
	} else if singleData, ok := rawData.(map[string]interface{}); ok {
		valid, err := schema.GetProcessor().ResourceValidGo(resourceName, singleData, true)

		if !valid {
			resp := ErrorResourceInvalid
			resp.Message += err.Error()
			return PrepError(c, resp)
		}

		finalResources = append(finalResources, singleData)
	} else {
		resp := ErrorCouldNotReadBody
		resp.Message += "not of type object or array of object"
		return PrepError(c, resp)
	}

	err = (*resource.GetProcessor().GetStore(resourceName)).CreateResources(resourceName, finalResources)

	if err != nil {
		resp := ErrorCreatingResource
		resp.Message += err.Error()
		return PrepError(c, resp)
	}

	return c.JSON(200, Response{
		Success: true,
	})
}

func deleteResource(c echo.Context) error {
	resourceName := c.Param("resource")

	if !schema.GetProcessor().ResourceExists(resourceName) {
		return PrepError(c, ErrorResourceDoesNotExist)
	}

	resultFilters, err := processFilters(c.QueryParams()["filters"])

	if err != nil {
		return PrepError(c, *err)
	}

	errG := resource.GetProcessor().DeleteResources(resourceName, resultFilters)

	if errG != nil {
		resp := ErrorDeletingResource
		resp.Message += errG.Error()
		return PrepError(c, resp)
	}

	return c.JSON(200, Response{
		Success: true,
	})
}

func processFilters(filters []string) ([]filter.ProcessedFilter, *Error) {
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
