package taikun

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/s3_credentials"
)

const testAccResourceTaikunBackupCredentialConfig = `
resource "taikun_backup_credential" "foo" {
  name            = "%s"
  lock       = %t

  s3_endpoint = "%s"
  s3_region   = "%s"
}
`

func TestAccResourceTaikunBackupCredential(t *testing.T) {
	backupCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckS3(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunBackupCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunBackupCredentialConfig,
					backupCredentialName,
					false,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBackupCredentialExists,
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "name", backupCredentialName),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_access_key_id", os.Getenv("AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_secret_access_key", os.Getenv("AWS_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_endpoint", os.Getenv("S3_ENDPOINT")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_region", os.Getenv("S3_REGION")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "is_default"),
				),
			},
		},
	})
}

func TestAccResourceTaikunBackupCredentialLock(t *testing.T) {
	backupCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckS3(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunBackupCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunBackupCredentialConfig,
					backupCredentialName,
					false,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBackupCredentialExists,
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "name", backupCredentialName),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_access_key_id", os.Getenv("AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_secret_access_key", os.Getenv("AWS_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_endpoint", os.Getenv("S3_ENDPOINT")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_region", os.Getenv("S3_REGION")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "is_default"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunBackupCredentialConfig,
					backupCredentialName,
					true,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBackupCredentialExists,
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "name", backupCredentialName),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_access_key_id", os.Getenv("AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_secret_access_key", os.Getenv("AWS_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_endpoint", os.Getenv("S3_ENDPOINT")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_region", os.Getenv("S3_REGION")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "lock", "true"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "is_default"),
				),
			},
		},
	})
}

func TestAccResourceTaikunBackupCredentialRename(t *testing.T) {
	backupCredentialName := randomTestName()
	newBackupCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckS3(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunBackupCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunBackupCredentialConfig,
					backupCredentialName,
					false,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBackupCredentialExists,
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "name", backupCredentialName),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_access_key_id", os.Getenv("AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_secret_access_key", os.Getenv("AWS_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_endpoint", os.Getenv("S3_ENDPOINT")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_region", os.Getenv("S3_REGION")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "is_default"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunBackupCredentialConfig,
					newBackupCredentialName,
					false,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBackupCredentialExists,
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "name", newBackupCredentialName),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_access_key_id", os.Getenv("AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_secret_access_key", os.Getenv("AWS_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_endpoint", os.Getenv("S3_ENDPOINT")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_region", os.Getenv("S3_REGION")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "is_default"),
				),
			},
		},
	})
}

func testAccCheckTaikunBackupCredentialExists(state *terraform.State) error {
	client := testAccProvider.Meta().(*taikungoclient.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_backup_credential" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := s3_credentials.NewS3CredentialsListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.Client.S3Credentials.S3CredentialsList(params, client)
		if err != nil || response.Payload.TotalCount != 1 {
			return fmt.Errorf("backup credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunBackupCredentialDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*taikungoclient.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_backup_credential" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)
			params := s3_credentials.NewS3CredentialsListParams().WithV(ApiVersion).WithID(&id)

			response, err := client.Client.S3Credentials.S3CredentialsList(params, client)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.Payload.TotalCount != 0 {
				return resource.RetryableError(errors.New("backup credential still exists ()"))
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("backup credential still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
