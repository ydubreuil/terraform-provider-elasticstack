package saved_object

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mitchellh/mapstructure"
)

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model ksoModelV0

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	kibanaClient, err := r.client.GetKibanaClient()
	if err != nil {
		resp.Diagnostics.AddError("unable to get kibana client", err.Error())
		return
	}

	result, err := kibanaClient.KibanaSavedObject.Import([]byte(model.Imported.ValueString()), true, model.SpaceID.ValueString())
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
		resp.Diagnostics.AddError("import failed", fmt.Sprintf("%v\n", importResponse.Errors))
	}
}
