package saved_object

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mitchellh/mapstructure"
)

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model ksoModelV0

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := model.UpdateModelWithObject()
	if err != nil {
		resp.Diagnostics.AddError("failed to update model from object", err.Error())
		return
	}

	kibanaClient, err := r.client.GetKibanaClient()
	if err != nil {
		resp.Diagnostics.AddError("unable to get kibana client", err.Error())
		return
	}

	result, err := kibanaClient.KibanaSavedObject.Import([]byte(model.Imported.ValueString()), false, model.SpaceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to import saved objects", err.Error())
		return
	}

	var importResponse ksoResponse
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

type ksoResponse struct {
	Success        bool               `json:"success"`
	SuccessCount   int                `json:"successCount"`
	Errors         []ksoImportError   `json:"errors"`
	SuccessResults []ksoImportSuccess `json:"successResults"`
}

type ksoImportSuccess struct {
	ID            string        `tfsdk:"id" json:"id"`
	Type          string        `tfsdk:"type" json:"type"`
	DestinationID string        `tfsdk:"destination_id" json:"destinationId"`
	Meta          ksoImportMeta `tfsdk:"meta" json:"meta"`
}

type ksoImportError struct {
	ID    string             `tfsdk:"id" json:"id"`
	Type  string             `tfsdk:"type" json:"type"`
	Title string             `tfsdk:"title" json:"title"`
	Error ksoImportErrorType `tfsdk:"error" json:"error"`
	Meta  ksoImportMeta      `tfsdk:"meta" json:"meta"`
}

func (ie ksoImportError) String() string {
	title := ie.Title
	if title == "" {
		title = ie.Meta.Title
	}

	return fmt.Sprintf("[%s] error on [%s] with ID [%s] and title [%s]", ie.Error.Type, ie.Type, ie.ID, title)
}

type ksoImportErrorType struct {
	Type string `tfsdk:"type" json:"type"`
}

type ksoImportMeta struct {
	Icon  string `tfsdk:"icon" json:"icon"`
	Title string `tfsdk:"title" json:"title"`
}
