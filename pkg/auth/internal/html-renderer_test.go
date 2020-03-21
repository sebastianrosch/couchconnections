package internal

import (
	"net/http/httptest"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HTML renderer", func() {
	var renderer *HTMLRenderer
	var mockCtrl *gomock.Controller
	var w *httptest.ResponseRecorder
	var email string
	var logoutURL string

	AfterEach(func() {
		mockCtrl.Finish()
	})

	BeforeEach(func() {
		email = "someEmail@auth0.com"
		logoutURL = "http://logout"
		w = httptest.NewRecorder()
	})

	Describe("when RenderSuccessPage is called", func() {
		whenRenderSuccessPageIsCalled := func() error {
			renderer = NewHTMLRenderer()

			return renderer.RenderSuccessPage(w, email, logoutURL)
		}

		It("should not fail and render email and logoutURL when rendering success template", func() {
			err := whenRenderSuccessPageIsCalled()

			Expect(err).ToNot(HaveOccurred())
			Expect(w.Body.String()).To(ContainSubstring(email))
			Expect(w.Body.String()).To(ContainSubstring(logoutURL))
		})
	})
	Describe("when RenderErrorPage is called", func() {
		whenRenderErrorPageIsCalled := func(errorMsg string) error {
			renderer = NewHTMLRenderer()

			return renderer.RenderErrorPage(w, errorMsg)
		}

		It("should not fail and render error message when rendering error template", func() {
			errorMsg := "some error message"

			err := whenRenderErrorPageIsCalled(errorMsg)

			Expect(err).ToNot(HaveOccurred())
			Expect(w.Body.String()).To(ContainSubstring(errorMsg))
		})
	})
})
