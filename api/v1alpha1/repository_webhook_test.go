/*
Copyright 2024.

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

package v1alpha1

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pacv1alpha1 "github.com/openshift-pipelines/pipelines-as-code/pkg/apis/pipelinesascode/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Repository Webhook", func() {

	var defaultNS string = "default"
	var defaultRepoName string = "default-repo"

	Context("Creating Repository under Validating Webhook", func() {
		repository := &pacv1alpha1.Repository{
			ObjectMeta: metav1.ObjectMeta{
				Name:      defaultRepoName,
				Namespace: defaultNS,
			},
			Spec: pacv1alpha1.RepositorySpec{
				URL: "https://github.com/org/repo",
			},
		}

		It("Should successfully create a repository", func() {
			By("Having a prefix in the url allow list the matches the repository url", func() {
				urlValidator.URLPrefixAllowList = []string{
					"https://github.com/org",
					"https//gitlab.com/org/group",
				}
			})
			err := k8sClient.Create(ctx, repository)
			Expect(err).ToNot(HaveOccurred())
		})

		It("Should successfully delete a repository", func() {
			err := k8sClient.Delete(ctx, repository)
			Expect(err).ToNot(HaveOccurred())
		})

		It("Should fail to create a repository", func() {
			By("Not having a prefix in the url allow list the matches the repository url", func() {
				urlValidator.URLPrefixAllowList = []string{
					"https://github.com/other-org",
					"https//gitlab.com/org/group",
				}
			})
			err := k8sClient.Create(ctx, repository)
			Expect(err).To(HaveOccurred())
		})

	})

})
