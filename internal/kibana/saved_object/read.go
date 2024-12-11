package saved_object

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model modelV0

	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	kibanaClient, err := r.client.GetKibanaClient()
	if err != nil {
		resp.Diagnostics.AddError("unable to get kibana client", err.Error())
		return
	}

	kibanaType, kibanaId, err := model.GetTypeAndObjectID()
	if err != nil {
		resp.Diagnostics.AddError("unable to get kibana type and object id", err.Error())
		return
	}

	spaceId := model.SpaceID
	result, err := kibanaClient.KibanaSavedObject.Get(kibanaType, kibanaId, spaceId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to get saved object", err.Error())
		return
	}

	// remove fields carrying state
	delete(result, "created_at")
	delete(result, "updated_at")
	delete(result, "version")
	delete(result, "migrationVersion")
	delete(result, "namespaces")

	object, err := json.Marshal(result)
	if err != nil {
		resp.Diagnostics.AddError("failed to marshal saved object", err.Error())
		return
	}

	resourceId := result["type"].(string) + "/" + result["id"].(string)
	if resourceId != kibanaType+"/"+kibanaId {
		resp.Diagnostics.AddError("ID changed for the resource", err.Error())
		return
	}

	model.Imported = types.StringValue(string(object))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
