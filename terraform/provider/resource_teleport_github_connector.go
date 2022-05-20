// Code generated by _gen/main.go DO NOT EDIT
/*
Copyright 2015-2022 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/gravitational/teleport-plugins/lib/backoff"
	"github.com/gravitational/teleport-plugins/terraform/tfschema"
	apitypes "github.com/gravitational/teleport/api/types"
	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"
)

// resourceTeleportGithubConnectorType is the resource metadata type
type resourceTeleportGithubConnectorType struct{}

// resourceTeleportGithubConnector is the resource
type resourceTeleportGithubConnector struct {
	p Provider
}

// GetSchema returns the resource schema
func (r resourceTeleportGithubConnectorType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfschema.GenSchemaGithubConnectorV3(ctx)
}

// NewResource creates the empty resource
func (r resourceTeleportGithubConnectorType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceTeleportGithubConnector{
		p: *(p.(*Provider)),
	}, nil
}

// Create creates the provision token
func (r resourceTeleportGithubConnector) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.IsConfigured(resp.Diagnostics) {
		return
	}

	var plan types.Object
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	githubConnector := &apitypes.GithubConnectorV3{}
	diags = tfschema.CopyGithubConnectorV3FromTerraform(ctx, plan, githubConnector)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	

	_, err := r.p.Client.GetGithubConnector(ctx, githubConnector.Metadata.Name, true)
	if !trace.IsNotFound(err) {
		if err == nil {
			n := githubConnector.Metadata.Name
			existErr := fmt.Sprintf("GithubConnector exists in Teleport. Either remove it (tctl rm github/%v)"+
				" or import it to the existing state (terraform import teleport_app.%v %v)", n, n, n)

			resp.Diagnostics.Append(diagFromErr("GithubConnector exists in Teleport", trace.Errorf(existErr)))
			return
		}

		resp.Diagnostics.Append(diagFromWrappedErr("Error reading GithubConnector", trace.Wrap(err), "github"))
		return
	}

	err = githubConnector.CheckAndSetDefaults()
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error setting GithubConnector defaults", trace.Wrap(err), "github"))
		return
	}

	err = r.p.Client.UpsertGithubConnector(ctx, githubConnector)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error creating GithubConnector", trace.Wrap(err), "github"))
		return
	}

	id := githubConnector.Metadata.Name
	var githubConnectorI apitypes.GithubConnector

	tries := 0
	backoff := backoff.NewDecorr(r.p.RetryConfig.Base, r.p.RetryConfig.Cap, clockwork.NewRealClock())
	for {
		tries = tries + 1
		githubConnectorI, err = r.p.Client.GetGithubConnector(ctx, id, true)
		if trace.IsNotFound(err) {
			if bErr := backoff.Do(ctx); bErr != nil {
				resp.Diagnostics.Append(diagFromWrappedErr("Error reading GithubConnector", trace.Wrap(err), "github"))
				return
			}
			if tries >= r.p.RetryConfig.MaxTries {
				diagMessage := fmt.Sprintf("Error reading GithubConnector (tried %d times)", tries)
				resp.Diagnostics.Append(diagFromWrappedErr(diagMessage, trace.Wrap(err), "github"))
				return
			}
			continue
		}
		break
	}

	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading GithubConnector", trace.Wrap(err), "github"))
		return
	}

	githubConnector, ok := githubConnectorI.(*apitypes.GithubConnectorV3)
	if !ok {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading GithubConnector", trace.Errorf("Can not convert %T to GithubConnectorV3", githubConnectorI), "github"))
		return
	}

	diags = tfschema.CopyGithubConnectorV3ToTerraform(ctx, *githubConnector, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Attrs["id"] = types.String{Value: githubConnector.Metadata.Name}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads teleport GithubConnector
func (r resourceTeleportGithubConnector) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state types.Object
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id types.String
	diags = req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("metadata").WithAttributeName("name"), &id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	githubConnectorI, err := r.p.Client.GetGithubConnector(ctx, id.Value, true)
	if trace.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading GithubConnector", trace.Wrap(err), "github"))
		return
	}

	githubConnector := githubConnectorI.(*apitypes.GithubConnectorV3)
	diags = tfschema.CopyGithubConnectorV3ToTerraform(ctx, *githubConnector, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates teleport GithubConnector
func (r resourceTeleportGithubConnector) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	if !r.p.IsConfigured(resp.Diagnostics) {
		return
	}

	var plan types.Object
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	githubConnector := &apitypes.GithubConnectorV3{}
	diags = tfschema.CopyGithubConnectorV3FromTerraform(ctx, plan, githubConnector)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := githubConnector.Metadata.Name

	err := githubConnector.CheckAndSetDefaults()
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error updating GithubConnector", err, "github"))
		return
	}

	githubConnectorBefore, err := r.p.Client.GetGithubConnector(ctx, name, true)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading GithubConnector", err, "github"))
		return
	}

	err = r.p.Client.UpsertGithubConnector(ctx, githubConnector)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error updating GithubConnector", err, "github"))
		return
	}

	var githubConnectorI apitypes.GithubConnector

	tries := 0
	backoff := backoff.NewDecorr(r.p.RetryConfig.Base, r.p.RetryConfig.Cap, clockwork.NewRealClock())
	for {
		tries = tries + 1
		githubConnectorI, err = r.p.Client.GetGithubConnector(ctx, name, true)
		if err != nil {
			resp.Diagnostics.Append(diagFromWrappedErr("Error reading GithubConnector", err, "github"))
			return
		}
		if githubConnectorBefore.GetMetadata().ID != githubConnectorI.GetMetadata().ID || true {
			break
		}

		if err := backoff.Do(ctx); err != nil {
			resp.Diagnostics.Append(diagFromWrappedErr("Error reading GithubConnector", trace.Wrap(err), "github"))
			return
		}
		if tries >= r.p.RetryConfig.MaxTries {
			diagMessage := fmt.Sprintf("Error reading GithubConnector (tried %d times)", tries)
			resp.Diagnostics.AddError(diagMessage, "github")
			return
		}
	}

	githubConnector = githubConnectorI.(*apitypes.GithubConnectorV3)
	diags = tfschema.CopyGithubConnectorV3ToTerraform(ctx, *githubConnector, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes Teleport GithubConnector
func (r resourceTeleportGithubConnector) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var id types.String
	diags := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("metadata").WithAttributeName("name"), &id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.p.Client.DeleteGithubConnector(ctx, id.Value)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error deleting GithubConnectorV3", trace.Wrap(err), "github"))
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports GithubConnector state
func (r resourceTeleportGithubConnector) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	githubConnectorI, err := r.p.Client.GetGithubConnector(ctx, req.ID, true)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading GithubConnector", trace.Wrap(err), "github"))
		return
	}

	githubConnector := githubConnectorI.(*apitypes.GithubConnectorV3)

	var state types.Object

	diags := resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = tfschema.CopyGithubConnectorV3ToTerraform(ctx, *githubConnector, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Attrs["id"] = types.String{Value: githubConnector.Metadata.Name}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
