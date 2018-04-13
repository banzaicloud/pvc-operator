package servermanagement

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"net/http"
)

// GatewayClient is the REST API for Azure Server Management Service.
type GatewayClient struct {
	BaseClient
}

// NewGatewayClient creates an instance of the GatewayClient client.
func NewGatewayClient(subscriptionID string) GatewayClient {
	return NewGatewayClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewGatewayClientWithBaseURI creates an instance of the GatewayClient client.
func NewGatewayClientWithBaseURI(baseURI string, subscriptionID string) GatewayClient {
	return GatewayClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Create creates or updates a ManagementService gateway.
//
// resourceGroupName is the resource group name uniquely identifies the resource group within the user
// subscriptionId. gatewayName is the gateway name (256 characters maximum). gatewayParameters is parameters
// supplied to the CreateOrUpdate operation.
func (client GatewayClient) Create(ctx context.Context, resourceGroupName string, gatewayName string, gatewayParameters GatewayParameters) (result GatewayCreateFuture, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `[a-zA-Z0-9]+`, Chain: nil}}},
		{TargetValue: gatewayName,
			Constraints: []validation.Constraint{{Target: "gatewayName", Name: validation.MaxLength, Rule: 256, Chain: nil},
				{Target: "gatewayName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "gatewayName", Name: validation.Pattern, Rule: `^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("servermanagement.GatewayClient", "Create", err.Error())
	}

	req, err := client.CreatePreparer(ctx, resourceGroupName, gatewayName, gatewayParameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "Create", nil, "Failure preparing request")
		return
	}

	result, err = client.CreateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "Create", result.Response(), "Failure sending request")
		return
	}

	return
}

// CreatePreparer prepares the Create request.
func (client GatewayClient) CreatePreparer(ctx context.Context, resourceGroupName string, gatewayName string, gatewayParameters GatewayParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"gatewayName":       autorest.Encode("path", gatewayName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-07-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServerManagement/gateways/{gatewayName}", pathParameters),
		autorest.WithJSON(gatewayParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateSender sends the Create request. The method will close the
// http.Response Body if it receives an error.
func (client GatewayClient) CreateSender(req *http.Request) (future GatewayCreateFuture, err error) {
	sender := autorest.DecorateSender(client, azure.DoRetryWithRegistration(client.Client))
	future.Future = azure.NewFuture(req)
	future.req = req
	_, err = future.Done(sender)
	if err != nil {
		return
	}
	err = autorest.Respond(future.Response(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated))
	return
}

// CreateResponder handles the response to the Create request. The method always
// closes the http.Response Body.
func (client GatewayClient) CreateResponder(resp *http.Response) (result GatewayResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete deletes a gateway from a resource group.
//
// resourceGroupName is the resource group name uniquely identifies the resource group within the user
// subscriptionId. gatewayName is the gateway name (256 characters maximum).
func (client GatewayClient) Delete(ctx context.Context, resourceGroupName string, gatewayName string) (result autorest.Response, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `[a-zA-Z0-9]+`, Chain: nil}}},
		{TargetValue: gatewayName,
			Constraints: []validation.Constraint{{Target: "gatewayName", Name: validation.MaxLength, Rule: 256, Chain: nil},
				{Target: "gatewayName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "gatewayName", Name: validation.Pattern, Rule: `^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("servermanagement.GatewayClient", "Delete", err.Error())
	}

	req, err := client.DeletePreparer(ctx, resourceGroupName, gatewayName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client GatewayClient) DeletePreparer(ctx context.Context, resourceGroupName string, gatewayName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"gatewayName":       autorest.Encode("path", gatewayName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-07-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServerManagement/gateways/{gatewayName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client GatewayClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client GatewayClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets a gateway.
//
// resourceGroupName is the resource group name uniquely identifies the resource group within the user
// subscriptionId. gatewayName is the gateway name (256 characters maximum) expand is gets subscription credentials
// which uniquely identify Microsoft Azure subscription. The subscription ID forms part of the URI for every
// service call.
func (client GatewayClient) Get(ctx context.Context, resourceGroupName string, gatewayName string, expand GatewayExpandOption) (result GatewayResource, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `[a-zA-Z0-9]+`, Chain: nil}}},
		{TargetValue: gatewayName,
			Constraints: []validation.Constraint{{Target: "gatewayName", Name: validation.MaxLength, Rule: 256, Chain: nil},
				{Target: "gatewayName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "gatewayName", Name: validation.Pattern, Rule: `^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("servermanagement.GatewayClient", "Get", err.Error())
	}

	req, err := client.GetPreparer(ctx, resourceGroupName, gatewayName, expand)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client GatewayClient) GetPreparer(ctx context.Context, resourceGroupName string, gatewayName string, expand GatewayExpandOption) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"gatewayName":       autorest.Encode("path", gatewayName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-07-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(string(expand)) > 0 {
		queryParameters["$expand"] = autorest.Encode("query", expand)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServerManagement/gateways/{gatewayName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client GatewayClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client GatewayClient) GetResponder(resp *http.Response) (result GatewayResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetProfile gets a gateway profile.
//
// resourceGroupName is the resource group name uniquely identifies the resource group within the user
// subscriptionId. gatewayName is the gateway name (256 characters maximum).
func (client GatewayClient) GetProfile(ctx context.Context, resourceGroupName string, gatewayName string) (result GatewayGetProfileFuture, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `[a-zA-Z0-9]+`, Chain: nil}}},
		{TargetValue: gatewayName,
			Constraints: []validation.Constraint{{Target: "gatewayName", Name: validation.MaxLength, Rule: 256, Chain: nil},
				{Target: "gatewayName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "gatewayName", Name: validation.Pattern, Rule: `^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("servermanagement.GatewayClient", "GetProfile", err.Error())
	}

	req, err := client.GetProfilePreparer(ctx, resourceGroupName, gatewayName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "GetProfile", nil, "Failure preparing request")
		return
	}

	result, err = client.GetProfileSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "GetProfile", result.Response(), "Failure sending request")
		return
	}

	return
}

// GetProfilePreparer prepares the GetProfile request.
func (client GatewayClient) GetProfilePreparer(ctx context.Context, resourceGroupName string, gatewayName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"gatewayName":       autorest.Encode("path", gatewayName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-07-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServerManagement/gateways/{gatewayName}/profile", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetProfileSender sends the GetProfile request. The method will close the
// http.Response Body if it receives an error.
func (client GatewayClient) GetProfileSender(req *http.Request) (future GatewayGetProfileFuture, err error) {
	sender := autorest.DecorateSender(client, azure.DoRetryWithRegistration(client.Client))
	future.Future = azure.NewFuture(req)
	future.req = req
	_, err = future.Done(sender)
	if err != nil {
		return
	}
	err = autorest.Respond(future.Response(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted))
	return
}

// GetProfileResponder handles the response to the GetProfile request. The method always
// closes the http.Response Body.
func (client GatewayClient) GetProfileResponder(resp *http.Response) (result GatewayProfile, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List returns gateways in a subscription.
func (client GatewayClient) List(ctx context.Context) (result GatewayResourcesPage, err error) {
	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.gr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "List", resp, "Failure sending request")
		return
	}

	result.gr, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client GatewayClient) ListPreparer(ctx context.Context) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-07-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.ServerManagement/gateways", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client GatewayClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client GatewayClient) ListResponder(resp *http.Response) (result GatewayResources, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client GatewayClient) listNextResults(lastResults GatewayResources) (result GatewayResources, err error) {
	req, err := lastResults.gatewayResourcesPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client GatewayClient) ListComplete(ctx context.Context) (result GatewayResourcesIterator, err error) {
	result.page, err = client.List(ctx)
	return
}

// ListForResourceGroup returns gateways in a resource group.
//
// resourceGroupName is the resource group name uniquely identifies the resource group within the user
// subscriptionId.
func (client GatewayClient) ListForResourceGroup(ctx context.Context, resourceGroupName string) (result GatewayResourcesPage, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `[a-zA-Z0-9]+`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("servermanagement.GatewayClient", "ListForResourceGroup", err.Error())
	}

	result.fn = client.listForResourceGroupNextResults
	req, err := client.ListForResourceGroupPreparer(ctx, resourceGroupName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "ListForResourceGroup", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListForResourceGroupSender(req)
	if err != nil {
		result.gr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "ListForResourceGroup", resp, "Failure sending request")
		return
	}

	result.gr, err = client.ListForResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "ListForResourceGroup", resp, "Failure responding to request")
	}

	return
}

// ListForResourceGroupPreparer prepares the ListForResourceGroup request.
func (client GatewayClient) ListForResourceGroupPreparer(ctx context.Context, resourceGroupName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-07-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServerManagement/gateways", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListForResourceGroupSender sends the ListForResourceGroup request. The method will close the
// http.Response Body if it receives an error.
func (client GatewayClient) ListForResourceGroupSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ListForResourceGroupResponder handles the response to the ListForResourceGroup request. The method always
// closes the http.Response Body.
func (client GatewayClient) ListForResourceGroupResponder(resp *http.Response) (result GatewayResources, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listForResourceGroupNextResults retrieves the next set of results, if any.
func (client GatewayClient) listForResourceGroupNextResults(lastResults GatewayResources) (result GatewayResources, err error) {
	req, err := lastResults.gatewayResourcesPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "listForResourceGroupNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListForResourceGroupSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "listForResourceGroupNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListForResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "listForResourceGroupNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListForResourceGroupComplete enumerates all values, automatically crossing page boundaries as required.
func (client GatewayClient) ListForResourceGroupComplete(ctx context.Context, resourceGroupName string) (result GatewayResourcesIterator, err error) {
	result.page, err = client.ListForResourceGroup(ctx, resourceGroupName)
	return
}

// RegenerateProfile regenerate a gateway's profile
//
// resourceGroupName is the resource group name uniquely identifies the resource group within the user
// subscriptionId. gatewayName is the gateway name (256 characters maximum).
func (client GatewayClient) RegenerateProfile(ctx context.Context, resourceGroupName string, gatewayName string) (result GatewayRegenerateProfileFuture, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `[a-zA-Z0-9]+`, Chain: nil}}},
		{TargetValue: gatewayName,
			Constraints: []validation.Constraint{{Target: "gatewayName", Name: validation.MaxLength, Rule: 256, Chain: nil},
				{Target: "gatewayName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "gatewayName", Name: validation.Pattern, Rule: `^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("servermanagement.GatewayClient", "RegenerateProfile", err.Error())
	}

	req, err := client.RegenerateProfilePreparer(ctx, resourceGroupName, gatewayName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "RegenerateProfile", nil, "Failure preparing request")
		return
	}

	result, err = client.RegenerateProfileSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "RegenerateProfile", result.Response(), "Failure sending request")
		return
	}

	return
}

// RegenerateProfilePreparer prepares the RegenerateProfile request.
func (client GatewayClient) RegenerateProfilePreparer(ctx context.Context, resourceGroupName string, gatewayName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"gatewayName":       autorest.Encode("path", gatewayName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-07-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServerManagement/gateways/{gatewayName}/regenerateprofile", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// RegenerateProfileSender sends the RegenerateProfile request. The method will close the
// http.Response Body if it receives an error.
func (client GatewayClient) RegenerateProfileSender(req *http.Request) (future GatewayRegenerateProfileFuture, err error) {
	sender := autorest.DecorateSender(client, azure.DoRetryWithRegistration(client.Client))
	future.Future = azure.NewFuture(req)
	future.req = req
	_, err = future.Done(sender)
	if err != nil {
		return
	}
	err = autorest.Respond(future.Response(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted))
	return
}

// RegenerateProfileResponder handles the response to the RegenerateProfile request. The method always
// closes the http.Response Body.
func (client GatewayClient) RegenerateProfileResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Update updates a gateway belonging to a resource group.
//
// resourceGroupName is the resource group name uniquely identifies the resource group within the user
// subscriptionId. gatewayName is the gateway name (256 characters maximum). gatewayParameters is parameters
// supplied to the Update operation.
func (client GatewayClient) Update(ctx context.Context, resourceGroupName string, gatewayName string, gatewayParameters GatewayParameters) (result GatewayUpdateFuture, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `[a-zA-Z0-9]+`, Chain: nil}}},
		{TargetValue: gatewayName,
			Constraints: []validation.Constraint{{Target: "gatewayName", Name: validation.MaxLength, Rule: 256, Chain: nil},
				{Target: "gatewayName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "gatewayName", Name: validation.Pattern, Rule: `^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("servermanagement.GatewayClient", "Update", err.Error())
	}

	req, err := client.UpdatePreparer(ctx, resourceGroupName, gatewayName, gatewayParameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "Update", nil, "Failure preparing request")
		return
	}

	result, err = client.UpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "Update", result.Response(), "Failure sending request")
		return
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client GatewayClient) UpdatePreparer(ctx context.Context, resourceGroupName string, gatewayName string, gatewayParameters GatewayParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"gatewayName":       autorest.Encode("path", gatewayName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-07-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServerManagement/gateways/{gatewayName}", pathParameters),
		autorest.WithJSON(gatewayParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdateSender sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (client GatewayClient) UpdateSender(req *http.Request) (future GatewayUpdateFuture, err error) {
	sender := autorest.DecorateSender(client, azure.DoRetryWithRegistration(client.Client))
	future.Future = azure.NewFuture(req)
	future.req = req
	_, err = future.Done(sender)
	if err != nil {
		return
	}
	err = autorest.Respond(future.Response(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted))
	return
}

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client GatewayClient) UpdateResponder(resp *http.Response) (result GatewayResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Upgrade upgrades a gateway.
//
// resourceGroupName is the resource group name uniquely identifies the resource group within the user
// subscriptionId. gatewayName is the gateway name (256 characters maximum).
func (client GatewayClient) Upgrade(ctx context.Context, resourceGroupName string, gatewayName string) (result GatewayUpgradeFuture, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `[a-zA-Z0-9]+`, Chain: nil}}},
		{TargetValue: gatewayName,
			Constraints: []validation.Constraint{{Target: "gatewayName", Name: validation.MaxLength, Rule: 256, Chain: nil},
				{Target: "gatewayName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "gatewayName", Name: validation.Pattern, Rule: `^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("servermanagement.GatewayClient", "Upgrade", err.Error())
	}

	req, err := client.UpgradePreparer(ctx, resourceGroupName, gatewayName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "Upgrade", nil, "Failure preparing request")
		return
	}

	result, err = client.UpgradeSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servermanagement.GatewayClient", "Upgrade", result.Response(), "Failure sending request")
		return
	}

	return
}

// UpgradePreparer prepares the Upgrade request.
func (client GatewayClient) UpgradePreparer(ctx context.Context, resourceGroupName string, gatewayName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"gatewayName":       autorest.Encode("path", gatewayName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-07-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServerManagement/gateways/{gatewayName}/upgradetolatest", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpgradeSender sends the Upgrade request. The method will close the
// http.Response Body if it receives an error.
func (client GatewayClient) UpgradeSender(req *http.Request) (future GatewayUpgradeFuture, err error) {
	sender := autorest.DecorateSender(client, azure.DoRetryWithRegistration(client.Client))
	future.Future = azure.NewFuture(req)
	future.req = req
	_, err = future.Done(sender)
	if err != nil {
		return
	}
	err = autorest.Respond(future.Response(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted))
	return
}

// UpgradeResponder handles the response to the Upgrade request. The method always
// closes the http.Response Body.
func (client GatewayClient) UpgradeResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}
