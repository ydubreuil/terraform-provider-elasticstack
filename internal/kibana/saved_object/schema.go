package saved_object

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/elastic/terraform-provider-elasticstack/internal/clients"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &Resource{}
var _ resource.ResourceWithConfigure = &Resource{}
var _ resource.ResourceWithModifyPlan = &Resource{}

func (r *Resource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var configData kibanaSavedObjectModelV0

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
		resp.Diagnostics.AddError("missing 'type' field in JSON object", err.Error())
		return
	}
	if objId, ok = object["id"]; !ok {
		resp.Diagnostics.AddError("missing 'id' field in JSON object", err.Error())
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

	resourceId := objType.(string) + "/" + objId.(string)
	if err != nil {
		resp.Diagnostics.AddError("unable to compute ID", err.Error())
		return
	}
	imported, err := json.Marshal(object)
	if err != nil {
		resp.Diagnostics.AddError("unable to marshal object", err.Error())
		return
	}

	resp.Plan.SetAttribute(ctx, path.Root("id"), types.StringValue(resourceId))
	resp.Plan.SetAttribute(ctx, path.Root("imported"), types.StringValue(string(imported)))
}

// var _ resource.ResourceWithConfigValidators = &Resource{}

// func (r *Resource) ConfigValidators(context.Context) []resource.ConfigValidator {
// 	return []resource.ConfigValidator{
// 		resourcevalidator.Conflicting(
// 			path.MatchRoot("create_new_copies"),
// 			path.MatchRoot("overwrite"),
// 			path.MatchRoot("compatibility_mode"),
// 		),
// 	}
// }

func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Import a Kibana saved object",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Generated ID for the import.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// This is needed as the user provided object cannot be cleaned up
			// see https://discuss.hashicorp.com/t/using-modifyplan-with-a-custom-provider/47690
			// We thus store a copy of the object with some fields removed
			"imported": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Kibana object imported.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"space_id": schema.StringAttribute{
				Description: "An identifier for the space. If space_id is not provided, the default space is used.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"object": schema.StringAttribute{
				Description: "Kibana object to import in JSON format",
				Required:    true,
			},
		},
	}
}

type Resource struct {
	client *clients.ApiClient
}

func (r *Resource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	client, diags := clients.ConvertProviderData(request.ProviderData)
	response.Diagnostics.Append(diags...)
	r.client = client
}

func (r *Resource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_kibana_saved_object"
}

type kibanaSavedObjectModelV0 struct {
	ID       types.String `tfsdk:"id"`
	SpaceID  types.String `tfsdk:"space_id"`
	Object   types.String `tfsdk:"object"`
	Imported types.String `tfsdk:"imported"`
}

func (m kibanaSavedObjectModelV0) GetTypeAndObjectID() (string, string, error) {
	parts := strings.SplitN(m.ID.ValueString(), "/", 2)
	if len(parts) == 2 {
		return parts[0], parts[1], nil
	} else {
		return "", "", errors.New(fmt.Sprintf("ID format is wrong: %s", m.ID.ValueString()))
	}
}
