---
page_title: "taikun_cloud_credential_azure Resource - terraform-provider-taikun"
subcategory: ""
description: |-   Taikun Azure Cloud Credential
---

# taikun_cloud_credential_azure (Resource)

Taikun Azure Cloud Credential

~> **Role Requirement** To use the `taikun_cloud_credential_azure` resource, you need a Manager or Partner account.

-> **Organization ID** `organization_id` can be specified for the Partner role, it otherwise defaults to the user's organization.

## Example Usage

```terraform
resource "taikun_cloud_credential_azure" "foo" {
  name = "foo"

  client_id         = "client_id"
  client_secret     = "client_secret"
  subscription_id   = "subscription_id"
  tenant_id         = "tenant_id"
  location          = "location"
  availability_zone = "availability_zone"

  organization_id = "42"
  lock            = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `availability_zone` (String) The Azure availability zone for the location.
- `client_id` (String, Sensitive) The Azure client ID.
- `client_secret` (String, Sensitive) The Azure client secret.
- `location` (String) The Azure location.
- `name` (String) The name of the Azure cloud credential.
- `subscription_id` (String) The Azure subscription ID.
- `tenant_id` (String) The Azure tenant ID.

### Optional

- `lock` (Boolean) Indicates whether to lock the Azure cloud credential. Defaults to `false`.
- `organization_id` (String) The ID of the organization which owns the Azure cloud credential.

### Read-Only

- `created_by` (String) The creator of the Azure cloud credential.
- `id` (String) The ID of the Azure cloud credential.
- `is_default` (Boolean) Indicates whether the Azure cloud credential is the default one.
- `last_modified` (String) Time and date of last modification.
- `last_modified_by` (String) The last user to have modified the Azure cloud credential.
- `organization_name` (String) The name of the organization which owns the Azure cloud credential.
