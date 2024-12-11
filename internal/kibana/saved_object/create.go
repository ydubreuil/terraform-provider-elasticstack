package saved_object

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mitchellh/mapstructure"
)

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model kibanaSavedObjectModelV0

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
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
	model.ID = types.StringValue(resourceId)

	// remove fields carrying state
	delete(object, "created_at")
	delete(object, "created_by")
	delete(object, "updated_at")
	delete(object, "updated_by")
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

	model.Imported = types.StringValue(string(imported))

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
		resp.Diagnostics.AddError("import failed", fmt.Sprintf("%#v\n", importResponse.Errors))
	}
}

type kibanaSavedObjectResponse struct {
	Success        bool            `json:"success"`
	SuccessCount   int             `json:"successCount"`
	Errors         []importError   `json:"errors"`
	SuccessResults []importSuccess `json:"successResults"`
}

type importSuccess struct {
	ID            string     `tfsdk:"id" json:"id"`
	Type          string     `tfsdk:"type" json:"type"`
	DestinationID string     `tfsdk:"destination_id" json:"destinationId"`
	Meta          importMeta `tfsdk:"meta" json:"meta"`
}

type importError struct {
	ID    string          `tfsdk:"id" json:"id"`
	Type  string          `tfsdk:"type" json:"type"`
	Title string          `tfsdk:"title" json:"title"`
	Error importErrorType `tfsdk:"error" json:"error"`
	Meta  importMeta      `tfsdk:"meta" json:"meta"`
}

func (ie importError) String() string {
	title := ie.Title
	if title == "" {
		title = ie.Meta.Title
	}

	return fmt.Sprintf("[%s] error on [%s] with ID [%s] and title [%s]", ie.Error.Type, ie.Type, ie.ID, title)
}

type importErrorType struct {
	Type string `tfsdk:"type" json:"type"`
}

type importMeta struct {
	Icon  string `tfsdk:"icon" json:"icon"`
	Title string `tfsdk:"title" json:"title"`
}
