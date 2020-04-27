// +build integration

package integration_test

import (
	"context"

	emptypb "github.com/golang/protobuf/ptypes/empty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IntegrationTests", func() {

	const (
		knownUsersKey   = "knownUsers"
		newFallbackName = "Mini-Me"
	)

	AfterEach(func() {
		// Reset the expectations of the mock servers after each test
	})

	Describe("GetStatus", func() {
		Context("always ", func() {
			It("should return OK and no error", func() {
				status, err := statusClient.GetStatus(context.Background(), &emptypb.Empty{})
				Expect(err).To(Not(HaveOccurred()))
				Expect(status).To(Not(BeNil()))
			})
		})
	})
})
