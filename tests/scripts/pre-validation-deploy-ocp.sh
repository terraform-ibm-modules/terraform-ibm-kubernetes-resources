#!/bin/bash

############################################################################################################
## This script is used by the catalog pipeline to deploy the VPC instance
## which are the prerequisites for the fully-configurable OCP VPC Cluster.
############################################################################################################

set -e

DA_DIR="solutions/standard"
TERRAFORM_SOURCE_DIR="tests/existing-resources"
JSON_FILE="${DA_DIR}/catalogValidationValues.json"
TF_VARS_FILE="terraform.tfvars"

(
  cwd=$(pwd)
  cd ${TERRAFORM_SOURCE_DIR}
  echo "Provisioning pre-requisite OCP instance .."
  terraform init || exit 1

  # $VALIDATION_APIKEY is available in the catalog runtime
  {
    echo "ibmcloud_api_key=\"${VALIDATION_APIKEY}\""
    echo "prefix=\"ocp-$(openssl rand -hex 2)\""
  } >> ${TF_VARS_FILE}
  terraform apply -input=false -auto-approve -var-file=${TF_VARS_FILE} || exit 1

  cluster_name_key="cluster_name"
  cluster_name_value=$(terraform output -state=terraform.tfstate -raw cluster_name)

  cluster_resource_group_id_key="cluster_resource_group_id"
  cluster_resource_group_id_value=$(terraform output -state=terraform.tfstate -raw resource_group_id)

  echo "Appending '${cluster_name_key}' '${cluster_resource_group_id_key}' input variable values to ${JSON_FILE}.."

  cd "${cwd}"
  jq -r \
  --arg cluster_name_key "${cluster_name_key}" \
  --arg cluster_name_value "${cluster_name_value}" \
  --arg cluster_resource_group_id_key "${cluster_resource_group_id_key}" \
  --arg cluster_resource_group_id_value "${cluster_resource_group_id_value}" \
  '. + {($cluster_name_key): $cluster_name_value, ($cluster_resource_group_id_key): $cluster_resource_group_id_value}' \
  "${JSON_FILE}" > tmpfile && mv tmpfile "${JSON_FILE}" || exit 1

  echo "Pre-validation completed successfully."
)
