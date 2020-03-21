package internal

import (
	"errors"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Code challenger", func() {
	var verifier *CodeVerifier
	var utils *MockCodeChallengeUtils
	var mockCtrl *gomock.Controller

	AfterEach(func() {
		mockCtrl.Finish()
	})

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		utils = NewMockCodeChallengeUtils(mockCtrl)
		verifier = NewCodeVerifier(utils)
	})

	Describe("When using method = S256", func() {
		It("should return the code verifier and challenge created with sha256", func() {
			randomBytes := []byte("randomBytes")
			utils.EXPECT().
				RandomBytes(32).
				Return(randomBytes).Times(1)
			utils.EXPECT().
				Encode(randomBytes).
				Return("someVerifier").Times(1)
			utils.EXPECT().
				Sha256Hash("someVerifier").
				Return("someSha").Times(1)

			result, _ := verifier.CreateCodeChallenge("S256")

			expectedResult := &CodeChallenge{
				Verifier:  "someVerifier",
				Challenge: "someSha",
				Method:    "S256",
			}
			Expect(result).To(Equal(expectedResult))
		})
	})
	Describe("When using method = plain", func() {
		It("should return the code verifier and a challenge that is the same", func() {
			randomBytes := []byte("randomBytes")
			utils.EXPECT().
				RandomBytes(32).
				Return(randomBytes).Times(1)
			utils.EXPECT().
				Encode(randomBytes).
				Return("someVerifier").Times(1)

			result, _ := verifier.CreateCodeChallenge("plain")

			expectedResult := &CodeChallenge{
				Verifier:  "someVerifier",
				Challenge: "someVerifier",
				Method:    "plain",
			}
			Expect(result).To(Equal(expectedResult))
		})
	})
	Describe("When using unknown method", func() {
		It("should return an error if the method is not known", func() {
			_, err := verifier.CreateCodeChallenge("notKnown")

			expected := errors.New("invalid method: notKnown")
			Expect(err).To(Equal(expected))
		})
	})
})
