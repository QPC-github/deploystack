package deploystack

import (
	"context"
	"fmt"
	"sort"

	"google.golang.org/api/cloudfunctions/v1"
)

var cloudfunctionsService *cloudfunctions.Service

func getCloudFunctionsService(project string) (*cloudfunctions.Service, error) {
	if cloudfunctionsService != nil {
		return cloudfunctionsService, nil
	}

	if err := EnableService(project, "cloudfunctions.googleapis.com"); err != nil {
		return nil, fmt.Errorf("error activating service for polling: %s", err)
	}

	ctx := context.Background()
	svc, err := cloudfunctions.NewService(ctx, opts)
	if err != nil {
		return nil, err
	}

	cloudfunctionsService = svc
	svc.UserAgent = defaultUserAgent

	return svc, nil
}

// RegionsFunctionsList will return a list of regions for Cloud Functions
func RegionsFunctionsList(project string) ([]string, error) {
	resp := []string{}

	if err := EnableService(project, "cloudfunctions.googleapis.com"); err != nil {
		return resp, fmt.Errorf("error activating service for polling: %s", err)
	}

	svc, err := getCloudFunctionsService(project)
	if err != nil {
		return resp, err
	}

	results, err := svc.Projects.Locations.List("projects/" + project).Do()
	if err != nil {
		return resp, err
	}

	for _, v := range results.Locations {
		resp = append(resp, v.LocationId)
	}

	sort.Strings(resp)

	return resp, nil
}

// DeployFunction deploys a Cloud Function.
func DeployFunction(project, region string, f cloudfunctions.CloudFunction) error {
	svc, err := getCloudFunctionsService(project)
	if err != nil {
		return err
	}

	location := fmt.Sprintf("projects/%s/locations/%s", project, region)
	if _, err := svc.Projects.Locations.Functions.Create(location, &f).Do(); err != nil {
		return fmt.Errorf("could not create function: %s", err)
	}

	return nil
}

// DeleteFunction deletes a Cloud Function.
func DeleteFunction(project, region, name string) error {
	svc, err := getCloudFunctionsService(project)
	if err != nil {
		return err
	}
	fname := fmt.Sprintf("projects/%s/locations/%s/functions/%s", project, region, name)
	if _, err := svc.Projects.Locations.Functions.Delete(fname).Do(); err != nil {
		return fmt.Errorf("could not create function: %s", err)
	}

	return nil
}

// GetFunction gets the details of a Cloud Function.
func GetFunction(project, region, name string) (*cloudfunctions.CloudFunction, error) {
	svc, err := getCloudFunctionsService(project)
	if err != nil {
		return nil, err
	}

	fname := fmt.Sprintf("projects/%s/locations/%s/functions/%s", project, region, name)
	result, err := svc.Projects.Locations.Functions.Get(fname).Do()
	if err != nil {
		return nil, fmt.Errorf("could not get function: %s", err)
	}

	return result, nil
}

func GenerateFunctionSignedURL(project, region string) (string, error) {
	location := fmt.Sprintf("projects/%s/locations/%s", project, region)
	svc, err := getCloudFunctionsService(project)
	if err != nil {
		return "", err
	}

	req := &cloudfunctions.GenerateUploadUrlRequest{}

	result, err := svc.Projects.Locations.Functions.GenerateUploadUrl(location, req).Do()
	if err != nil {
		return "", err
	}

	return result.UploadUrl, nil
}
