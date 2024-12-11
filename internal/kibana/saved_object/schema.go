package saved_object

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/elastic/terraform-provider-elasticstack/internal/clients"
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
	// Fill in logic.
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
					stringplanmodifier.UseStateForUnknown(),
				},
			},
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

type modelV0 struct {
	ID       types.String `tfsdk:"id"`
	SpaceID  types.String `tfsdk:"space_id"`
	Object   types.String `tfsdk:"object"`
	Imported types.String `tfsdk:"imported"`
}

func (m modelV0) GetTypeAndObjectID() (string, string, error) {
	parts := strings.SplitN(m.ID.ValueString(), "/", 2)
	if len(parts) == 2 {
		return parts[0], parts[1], nil
	} else {
		return "", "", errors.New(fmt.Sprintf("ID format is wrong: %s", m.ID.ValueString()))
	}
}
