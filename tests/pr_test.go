// Tests in this file are run in the PR pipeline and the continuous testing pipeline
package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
	"github.com/terraform-ibm-modules/ibmcloud-terratest-wrapper/common"
	"github.com/terraform-ibm-modules/ibmcloud-terratest-wrapper/testhelper"
	"github.com/terraform-ibm-modules/ibmcloud-terratest-wrapper/testschematic"
)

const solutionsDir = "solutions/standard"

func validateEnvVariable(t *testing.T, varName string) string {
	val, present := os.LookupEnv(varName)
	require.True(t, present, "%s environment variable not set", varName)
	require.NotEqual(t, "", val, "%s environment variable is empty", varName)
	return val
}

func setupTerraform(t *testing.T, prefix, realTerraformDir string) (string, *terraform.Options) {
	tempTerraformDir, err := files.CopyTerraformFolderToTemp(realTerraformDir, prefix)
	require.NoError(t, err, "Failed to create temporary Terraform folder")
	apiKey := validateEnvVariable(t, "TF_VAR_ibmcloud_api_key") // pragma: allowlist secret
	region, err := testhelper.GetBestVpcRegion(apiKey, "../common-dev-assets/common-go-assets/cloudinfo-region-vpc-gen2-prefs.yaml", "eu-de")
	require.NoError(t, err, "Failed to get best VPC region")

	existingTerraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: tempTerraformDir,
		Vars: map[string]interface{}{
			"prefix": prefix,
			"region": region,
		},
		// Set Upgrade to true to ensure latest version of providers and modules are used by terratest.
		// This is the same as setting the -upgrade=true flag with terraform.
		Upgrade: true,
	})

	terraform.WorkspaceSelectOrNew(t, existingTerraformOptions, prefix)
	_, err = terraform.InitAndApplyE(t, existingTerraformOptions)
	require.NoError(t, err, "Init and Apply of temp existing resource failed")

	return region, existingTerraformOptions
}

func cleanupTerraform(t *testing.T, options *terraform.Options, prefix string) {
	if t.Failed() && strings.ToLower(os.Getenv("DO_NOT_DESTROY_ON_FAILURE")) == "true" {
		fmt.Println("Terratest failed. Debug the test and delete resources manually.")
		return
	}
	logger.Log(t, "START: Destroy (existing resources)")
	terraform.Destroy(t, options)
	terraform.WorkspaceDelete(t, options, prefix)
	logger.Log(t, "END: Destroy (existing resources)")
}

func TestDAInSchematics(t *testing.T) {
	t.Parallel()

	common.UniqueId()
	prefix := fmt.Sprintf("ocp-%s", common.UniqueId())
	region, existingTerraformOptions := setupTerraform(t, prefix, "./existing-resources")

	namespaces := []map[string]interface{}{
		{
			"name": fmt.Sprintf("namespace1-%s", common.UniqueId()),
		},
		{
			"name": fmt.Sprintf("namespace2-%s", common.UniqueId()),
		},
	}

	options := testschematic.TestSchematicOptionsDefault(&testschematic.TestSchematicOptions{
		Testing:               t,
		TarIncludePatterns:    []string{"*.tf", solutionsDir + "/*.*"},
		TemplateFolder:        solutionsDir,
		Tags:                  []string{"test-schematic"},
		DeleteWorkspaceOnFail: false,
	})

	options.TerraformVars = []testschematic.TestSchematicTerraformVar{
		{Name: "ibmcloud_api_key", Value: options.RequiredEnvironmentVars["TF_VAR_ibmcloud_api_key"], DataType: "string", Secure: true},
		{Name: "region", Value: region, DataType: "string"},
		{Name: "cluster_name", Value: terraform.Output(t, existingTerraformOptions, "cluster_name"), DataType: "string"},
		{Name: "namespaces", Value: namespaces, DataType: "list(object{})"},
	}

	require.NoError(t, options.RunSchematicTest(), "This should not have errored")
	cleanupTerraform(t, existingTerraformOptions, prefix)
}

func TestDAUpgradeInSchematics(t *testing.T) {
	t.Parallel()

	common.UniqueId()
	prefix := fmt.Sprintf("ocp-upg-%s", common.UniqueId())
	region, existingTerraformOptions := setupTerraform(t, prefix, "./existing-resources")

	namespaces := []map[string]interface{}{
		{
			"name": fmt.Sprintf("ns-upg-1-%s", common.UniqueId()),
		},
		{
			"name": fmt.Sprintf("ns-upg-2-%s", common.UniqueId()),
		},
	}

	options := testschematic.TestSchematicOptionsDefault(&testschematic.TestSchematicOptions{
		Testing:               t,
		TarIncludePatterns:    []string{"*.tf", solutionsDir + "/*.*"},
		TemplateFolder:        solutionsDir,
		Tags:                  []string{"test-schematic"},
		DeleteWorkspaceOnFail: false,
	})

	options.TerraformVars = []testschematic.TestSchematicTerraformVar{
		{Name: "ibmcloud_api_key", Value: options.RequiredEnvironmentVars["TF_VAR_ibmcloud_api_key"], DataType: "string", Secure: true},
		{Name: "region", Value: region, DataType: "string"},
		{Name: "cluster_name", Value: terraform.Output(t, existingTerraformOptions, "cluster_name"), DataType: "string"},
		{Name: "namespaces", Value: namespaces, DataType: "list(object{})"},
	}

	require.NoError(t, options.RunSchematicUpgradeTest(), "This should not have errored")
	cleanupTerraform(t, existingTerraformOptions, prefix)
}
