package saved_object

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

	if err := kibanaClient.KibanaSavedObject.Delete(kibanaType, kibanaId, model.SpaceID.ValueString()); err != nil {
		resp.Diagnostics.AddError("failed to delete saved object", err.Error())
		return
	}
}
