package saved_object

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ksoModelV0 struct {
	ID       types.String `tfsdk:"id"`
	SpaceID  types.String `tfsdk:"space_id"`
	Object   types.String `tfsdk:"object"`
	Imported types.String `tfsdk:"imported"`
	Type     types.String `tfsdk:"type"`
}
