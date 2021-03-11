/*
© 2021 Red Hat, Inc. and others.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/submariner-io/submariner-operator/pkg/internal/cli"
)

var validateK8sVersionCmd = &cobra.Command{
	Use:   "k8s-version",
	Short: "Validate if Submariner supports the Kubernetes version.",
	Long:  "This command validates if Submariner can be deployed on the Kubernetes version used in the cluster.",
	Run:   validateK8sVersion,
}

func init() {
	validateCmd.AddCommand(validateK8sVersionCmd)
}

func validateK8sVersion(cmd *cobra.Command, args []string) {
	configs, err := getMultipleRestConfigs(kubeConfig, kubeContext)
	exitOnError("Error getting REST config for cluster", err)

	for _, item := range configs {
		message := fmt.Sprintf("Validating if Submariner supports Kubernetes version used in %q", item.clusterName)
		status.Start(message)

		failedRequirements, err := checkRequirements(item.config)
		if len(failedRequirements) > 0 {
			status.QueueFailureMessage("The target cluster fails to meet Submariner's requirements:")
			for i := range failedRequirements {
				message = fmt.Sprintf("* %s\n", (failedRequirements)[i])
				status.QueueFailureMessage(message)
			}
			status.End(cli.Failure)
			continue
		}
		if err != nil {
			status.QueueFailureMessage(err.Error())
			status.End(cli.Failure)
			continue
		}
		message = fmt.Sprintf("The Kubernetes version on %q meets Submariner's requirements", item.clusterName)
		status.QueueSuccessMessage(message)
		status.End(cli.Success)
	}
}
