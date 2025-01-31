package taikun

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/organizations"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunOrganizationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address": {
			Description: "Address.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"billing_email": {
			Description: "Billing email.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"city": {
			Description: "City.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"country": {
			Description: "Country.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"created_at": {
			Description: "Time and date of creation.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"discount_rate": {
			Description:  "Discount rate, must be between 0 and 100 (included).",
			Type:         schema.TypeFloat,
			Optional:     true,
			Default:      100,
			ValidateFunc: validation.FloatBetween(0, 100),
		},
		"email": {
			Description: "Email.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"full_name": {
			Description:  "Full name.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"id": {
			Description: "Organization's ID.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"is_read_only": {
			Description: "Whether the organization is in read-only mode.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"lock": {
			Description: "Indicates whether to lock the organization.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"managers_can_change_subscription": {
			Description: "Allow subscription to be changed by managers.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
		},
		"name": {
			Description: "Organization's name.",
			Type:        schema.TypeString,
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 30),
				validation.StringMatch(
					regexp.MustCompile("^[a-z0-9-_.]+$"),
					"expected only alpha numeric characters or non alpha numeric (_-.)",
				),
			),
		},
		"partner_id": {
			Description: "ID of the organization's partner.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"partner_name": {
			Description: "Name of the organization's partner.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"phone": {
			Description: "Phone number.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"vat_number": {
			Description: "VAT number.",
			Type:        schema.TypeString,
			Optional:    true,
		},
	}
}

func resourceTaikunOrganization() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Organization",
		CreateContext: resourceTaikunOrganizationCreate,
		ReadContext:   generateResourceTaikunOrganizationReadWithoutRetries(),
		UpdateContext: resourceTaikunOrganizationUpdate,
		DeleteContext: resourceTaikunOrganizationDelete,
		Schema:        resourceTaikunOrganizationSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)

	body := &models.OrganizationCreateCommand{
		Address:                      d.Get("address").(string),
		BillingEmail:                 d.Get("billing_email").(string),
		City:                         d.Get("city").(string),
		Country:                      d.Get("country").(string),
		DiscountRate:                 d.Get("discount_rate").(float64),
		Email:                        d.Get("email").(string),
		FullName:                     d.Get("full_name").(string),
		IsEligibleUpdateSubscription: d.Get("managers_can_change_subscription").(bool),
		Name:                         d.Get("name").(string),
		Phone:                        d.Get("phone").(string),
		VatNumber:                    d.Get("vat_number").(string),
	}

	params := organizations.NewOrganizationsCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.Client.Organizations.OrganizationsCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := atoi32(createResult.GetPayload().ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.GetPayload().ID)

	if isLocked, isLockedIsSet := d.GetOk("lock"); isLockedIsSet {
		updateLockBody := &models.UpdateOrganizationCommand{
			Address:                      body.Address,
			BillingEmail:                 body.BillingEmail,
			City:                         body.City,
			Country:                      body.Country,
			DiscountRate:                 body.DiscountRate,
			Email:                        body.Email,
			FullName:                     body.FullName,
			ID:                           id,
			IsEligibleUpdateSubscription: body.IsEligibleUpdateSubscription,
			IsLocked:                     isLocked.(bool),
			Name:                         body.Name,
			Phone:                        body.Phone,
			VatNumber:                    body.VatNumber,
		}
		updateLockParams := organizations.NewOrganizationsUpdateParams().WithV(ApiVersion).WithBody(updateLockBody)
		_, err := apiClient.Client.Organizations.OrganizationsUpdate(updateLockParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunOrganizationReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunOrganizationReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunOrganizationRead(true)
}
func generateResourceTaikunOrganizationReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunOrganizationRead(false)
}
func generateResourceTaikunOrganizationRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*taikungoclient.Client)
		id := d.Id()
		id32, _ := atoi32(d.Id())
		d.SetId("")

		response, err := apiClient.Client.Organizations.OrganizationsList(organizations.NewOrganizationsListParams().WithV(ApiVersion).WithID(&id32), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(response.Payload.Data) != 1 {
			if withRetries {
				d.SetId(id)
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawOrganization := response.GetPayload().Data[0]

		err = setResourceDataFromMap(d, flattenTaikunOrganization(rawOrganization))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(id)

		return nil
	}
}

func resourceTaikunOrganizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)

	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := &models.UpdateOrganizationCommand{
		Address:                      d.Get("address").(string),
		BillingEmail:                 d.Get("billing_email").(string),
		City:                         d.Get("city").(string),
		Country:                      d.Get("country").(string),
		DiscountRate:                 d.Get("discount_rate").(float64),
		Email:                        d.Get("email").(string),
		FullName:                     d.Get("full_name").(string),
		ID:                           id,
		IsEligibleUpdateSubscription: d.Get("managers_can_change_subscription").(bool),
		IsLocked:                     d.Get("lock").(bool),
		Name:                         d.Get("name").(string),
		Phone:                        d.Get("phone").(string),
		VatNumber:                    d.Get("vat_number").(string),
	}

	updateLockParams := organizations.NewOrganizationsUpdateParams().WithV(ApiVersion).WithBody(body)
	_, err = apiClient.Client.Organizations.OrganizationsUpdate(updateLockParams, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	return readAfterUpdateWithRetries(generateResourceTaikunOrganizationReadWithRetries(), ctx, d, meta)
}

func resourceTaikunOrganizationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	params := organizations.NewOrganizationsDeleteParams().WithV(ApiVersion).WithOrganizationID(id)
	_, _, err = apiClient.Client.Organizations.OrganizationsDelete(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func flattenTaikunOrganization(rawOrganization *models.OrganizationDetailsDto) map[string]interface{} {
	return map[string]interface{}{
		"address":                          rawOrganization.Address,
		"billing_email":                    rawOrganization.BillingEmail,
		"city":                             rawOrganization.City,
		"country":                          rawOrganization.Country,
		"created_at":                       rawOrganization.CreatedAt,
		"discount_rate":                    rawOrganization.DiscountRate,
		"email":                            rawOrganization.Email,
		"full_name":                        rawOrganization.FullName,
		"id":                               i32toa(rawOrganization.ID),
		"managers_can_change_subscription": rawOrganization.IsEligibleUpdateSubscription,
		"lock":                             rawOrganization.IsLocked,
		"is_read_only":                     rawOrganization.IsReadOnly,
		"name":                             rawOrganization.Name,
		"partner_id":                       i32toa(rawOrganization.PartnerID),
		"partner_name":                     rawOrganization.PartnerName,
		"phone":                            rawOrganization.Phone,
		"vat_number":                       rawOrganization.VatNumber,
	}
}
