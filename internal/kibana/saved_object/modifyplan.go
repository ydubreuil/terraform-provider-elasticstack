package saved_object

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *Resource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var configData ksoModelV0

	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &configData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var object map[string]any
	err := json.Unmarshal([]byte(configData.Object.ValueString()), &object)
	if err != nil {
		resp.Diagnostics.AddError("invalid JSON in object", err.Error())
		return
	}

	var objType any
	var objId any
	var ok bool
	if objType, ok = object["type"]; !ok {
		resp.Diagnostics.AddError("missing 'type' field in JSON object", "")
		return
	}
	if objId, ok = object["id"]; !ok {
		resp.Diagnostics.AddError("missing 'id' field in JSON object", "")
		return
	}

	// remove fields carrying state
	delete(object, "created_at")
	delete(object, "created_by")
	delete(object, "updated_at")
	delete(object, "updated_by")
	delete(object, "version")
	delete(object, "migrationVersion")
	delete(object, "namespaces")

	imported, err := json.Marshal(object)
	if err != nil {
		resp.Diagnostics.AddError("unable to marshal object", err.Error())
		return
	}

	resp.Plan.SetAttribute(ctx, path.Root("id"), types.StringValue(objId.(string)))
	resp.Plan.SetAttribute(ctx, path.Root("type"), types.StringValue(objType.(string)))
	resp.Plan.SetAttribute(ctx, path.Root("imported"), types.StringValue(string(imported)))
}
