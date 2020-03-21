package internal

import (
	"encoding/base64"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
	var utils *Utils

	BeforeEach(func() {
		utils = NewUtils()
	})

	Describe("When RandomBytes is called", func() {
		It("should return random bytes of the length requested", func() {
			bytes := utils.RandomBytes(32)

			Expect(len(bytes)).To(Equal(32))
		})
	})

	Describe("When Encode is called", func() {
		It("should return the base64 encoded representation of a string", func() {
			original, _ := base64.StdEncoding.DecodeString("hJAq49pUrsb8rOQO+bw+56LVG2cx32+SymJS4YJjNsM=")
			result := utils.Encode(original)

			Expect(result).To(Equal("hJAq49pUrsb8rOQO-bw-56LVG2cx32-SymJS4YJjNsM"))
		})
	})

	Describe("When Sha256Hash is called", func() {
		It("should return the encoded sha256 hash representation of a string", func() {
			result := utils.Sha256Hash("https://blog.golang.org/wire")

			Expect(result).To(Equal("hJAq49pUrsb8rOQO-bw-56LVG2cx32-SymJS4YJjNsM"))
		})
	})
})
