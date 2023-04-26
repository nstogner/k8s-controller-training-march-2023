package v1

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("LoadTest webhook server", func() {
	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating a valid LoadTest", func() {
		It("Should allow the request", func() {

			By("By creating a LoadTest that looks just right")

			validLT := &LoadTest{
				TypeMeta: metav1.TypeMeta{
					Kind:       "LoadTest",
					APIVersion: "tests.tbd.com/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "valid-loadtest",
					Namespace: "default",
				},
				Spec: LoadTestSpec{
					Method:  "GET",
					Address: "http://localhost:1111/not-used",
					Duration: metav1.Duration{
						Duration: 1 * time.Minute,
					},
				},
			}

			Expect(k8sClient.Create(ctx, validLT)).Should(Succeed())
		})
	})

	Context("When creating an invalid LoadTest", func() {
		It("Should block the request", func() {

			By("By creating a LoadTest that has a long duration")

			invalidLT := &LoadTest{
				TypeMeta: metav1.TypeMeta{
					Kind:       "LoadTest",
					APIVersion: "tests.tbd.com/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "invalid-loadtest-too-long",
					Namespace: "default",
				},
				Spec: LoadTestSpec{
					Method:  "GET",
					Address: "http://localhost:1111/not-used",
					Duration: metav1.Duration{
						Duration: 2 * time.Hour,
					},
				},
			}

			err := k8sClient.Create(ctx, invalidLT)
			Expect(err).ShouldNot(BeNil())
		})
	})
})
