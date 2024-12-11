package saved_object

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mitchellh/mapstructure"
)

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model kibanaSavedObjectModelV0

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var object map[string]any
	err := json.Unmarshal([]byte(model.Object.ValueString()), &object)
	if err != nil {
		resp.Diagnostics.AddError("invalid JSON in object", err.Error())
		return
	}

	var objType any
	var objId any
	var ok bool
	if objType, ok = object["type"]; !ok {
		resp.Diagnostics.AddError("missing 'type' field in JSON object", err.Error())
		return
	}
	if objId, ok = object["id"]; !ok {
		resp.Diagnostics.AddError("missing 'id' field in JSON object", err.Error())
		return
	}

	resourceId := objType.(string) + "/" + objId.(string)
	if err != nil {
		resp.Diagnostics.AddError("unable to compute ID", err.Error())
		return
	}

	kibanaType, kibanaId, err := model.GetTypeAndObjectID()
	if err != nil {
		resp.Diagnostics.AddError("unable to get kibana type and object id", err.Error())
		return
	}

	if resourceId != kibanaType+"/"+kibanaId {
		resp.Diagnostics.AddError("ID changed for the resource", fmt.Sprintf("Old: '%s', New: '%s'", kibanaType+"/"+kibanaId, resourceId))
		return
	}

	// remove fields carrying state
	delete(object, "created_at")
	delete(object, "updated_at")
	delete(object, "version")
	delete(object, "migrationVersion")
	delete(object, "namespaces")

	kibanaClient, err := r.client.GetKibanaClient()
	if err != nil {
		resp.Diagnostics.AddError("unable to get kibana client", err.Error())
		return
	}

	imported, err := json.Marshal(object)
	if err != nil {
		resp.Diagnostics.AddError("unable to marshal object", err.Error())
		return
	}

	result, err := kibanaClient.KibanaSavedObject.Import(imported, true, model.SpaceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to import saved objects", err.Error())
		return
	}

	var importResponse kibanaSavedObjectResponse
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &importResponse,
		TagName: "json",
	})
	if err != nil {
		resp.Diagnostics.AddError("failed to create model decoder", err.Error())
		return
	}

	err = decoder.Decode(result)
	if err != nil {
		resp.Diagnostics.AddError("failed to decode response", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !importResponse.Success {
		resp.Diagnostics.AddError("import failed", fmt.Sprintf("%v\n", importResponse.Errors))
	}
}
