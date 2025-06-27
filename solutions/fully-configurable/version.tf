terraform {
  required_version = ">=1.9.0"

  # Lock DA into an exact provider version - renovate automation will keep it updated
  required_providers {
    ibm = {
      source  = "IBM-Cloud/ibm"
      version = "1.79.2"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.16.1, < 3.0.0"
    }
  }
}
