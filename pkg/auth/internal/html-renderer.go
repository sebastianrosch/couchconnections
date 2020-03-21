package internal

import (
	"html/template"
	"io"
)

type HTMLRenderer struct {
}

func (r *HTMLRenderer) RenderSuccessPage(writer io.Writer, emailAddress, logoutURL string) error {
	template, err := r.successTemplate()
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	data["EmailAddress"] = emailAddress
	data["LogoutURL"] = logoutURL

	return r.renderTemplate(writer, template, &data)
}

func (r *HTMLRenderer) RenderErrorPage(writer io.Writer, errorMessage string) error {
	template, err := r.errorTemplete()
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	data["ErrorMessage"] = errorMessage

	return r.renderTemplate(writer, template, &data)
}

func (r *HTMLRenderer) renderTemplate(writer io.Writer, template *template.Template, data *map[string]interface{}) error {
	return template.Execute(writer, data)
}

func (r *HTMLRenderer) successTemplate() (*template.Template, error) {
	const successTemplate = `<!DOCTYPE html>

	<html>
	  <head>
		<meta charset="utf-8" />
		<meta http-equiv="X-UA-Compatible" content="IE=edge" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />

		<meta name="robots" content="noindex, nofollow" />
		<title>Logged In</title>

		<link
		  rel="stylesheet"
		  href="https://auth0-ulp.herokuapp.com/static/css/main.local.min.css"
		/>
		<style id="custom-styles-container">
		  body {
			background: #000;
			font-family: ulp-font, -apple-system, BlinkMacSystemFont, Roboto,
			  Helvetica, sans-serif;
		  }
		  .main-wrapper {
			background: #000;
		  }
		  .ulp-alert.danger {
			background: #e55e3f;
		  }
		  .ulp-alert.success {
			background: #1bc462;
		  }
		  .ulp-button-default {
			background-color: #0059d6;
			color: #fff;
		  }
		  .ulp-button-success {
			background-color: #1bc462;
		  }
		  .ulp-button-error {
			background-color: #e55e3f;
		  }
		  @supports (
			mask-image:
			  url('https://auth0-ulp.herokuapp.com/static/img/branding-generic/copy-icon.svg')
		  ) {
			@supports not (-ms-ime-align: auto) {
			  .input-container.error::before {
				background-color: #e55e3f;
			  }
			}
		  }
		  .input.ulp-input-error {
			border-color: #e55e3f;
		  }
		  .error-cloud {
			background-color: #e55e3f;
		  }
		  .error-fatal {
			background-color: #e55e3f;
		  }
		  .error-local {
			background-color: #e55e3f;
		  }
		  .ulp-popover.error {
			background-color: #e55e3f;
			border-color: #e55e3f;
		  }
		  .ulp-popover.error::before {
			border-bottom-color: #e55e3f;
		  }
		  .ulp-popover.error::after {
			border-bottom-color: #e55e3f;
		  }
		  #alert-trigger {
			background-color: #e55e3f;
		  }
		</style>
		<style>
		  /* By default, hide features for javascript-disabled browsing */
		  /* We use !important to override any css with higher specificity */
		  /* It is also overriden by the styles in <noscript> in the header file */
		  .no-js {
			display: none !important;
		  }
		</style>
		<noscript>
		  <style>
			/* We use !important to override the default for js enabled */
			/* If the display should be other than block, it should be defined specifically here */
			.js-required {
			  display: none !important;
			}
			.no-js {
			  display: block !important;
			}
			div.input-container input {
			  padding-right: 16px !important;
			}
			div.input-container::after {
			  display: none !important;
			}
		  </style>
		</noscript>
	  </head>

	  <body>
		<div class="main-wrapper">
		  <main class="ulp-outer">
			<section class="ulp-box email-verification-result ulp-event-outer ">
			  <div class="ulp-box-inner">
				<div class="ulp-container ulp-event-screen" data-event-id="">
				  <div class="event-img-container">
					<span class="event-img success-lock"></span>
				  </div>

				  <section class="event-container">
					<h3 class="event-title">Logged In</h3>

					<div class="event-text">
					  Successfully logged in as <b>{{.EmailAddress}}</b>.<br />
					  You can close this page and return to your CLI.<br />
					  <a href="{{.LogoutURL}}">Logout</a>
					</div>
				  </section>
				</div>
			  </div>

			  <a
				href="https://auth0.com/?utm_source=lock&amp;utm_campaign=badge&amp;utm_medium=widget"
				target="_blank"
				rel="noopener noreferrer"
				class="ulp-auth0-badge"
				aria-label="Link to the Auth0 website"
			  >
				<img
				  src="https://auth0-ulp.herokuapp.com/static/img/theme-generic/auth0-logo.svg"
				  alt="Link to the Auth0 website"
				  class="ulp-auth0-badge-image"
				/>
			  </a>
			</section>
		  </main>
		</div>
	  </body>
	</html>
`

	return template.New("success").Parse(successTemplate)
}

func (r *HTMLRenderer) errorTemplete() (*template.Template, error) {
	const errorTemplate = `<!DOCTYPE html>

	<html>
	  <head>
		<meta charset="utf-8" />
		<meta http-equiv="X-UA-Compatible" content="IE=edge" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />

		<meta name="robots" content="noindex, nofollow" />
		<title>Error</title>

		<link
		  rel="stylesheet"
		  href="https://auth0-ulp.herokuapp.com/static/css/main.local.min.css"
		/>
		<style id="custom-styles-container">
		  body {
			background: #000;
			font-family: ulp-font, -apple-system, BlinkMacSystemFont, Roboto,
			  Helvetica, sans-serif;
		  }
		  .main-wrapper {
			background: #000;
		  }
		  .ulp-alert.danger {
			background: #e55e3f;
		  }
		  .ulp-alert.success {
			background: #1bc462;
		  }
		  .ulp-button-default {
			background-color: #0059d6;
			color: #fff;
		  }
		  .ulp-button-success {
			background-color: #1bc462;
		  }
		  .ulp-button-error {
			background-color: #e55e3f;
		  }
		  @supports (
			mask-image:
			  url('https://auth0-ulp.herokuapp.com/static/img/branding-generic/copy-icon.svg')
		  ) {
			@supports not (-ms-ime-align: auto) {
			  .input-container.error::before {
				background-color: #e55e3f;
			  }
			}
		  }
		  .input.ulp-input-error {
			border-color: #e55e3f;
		  }
		  .error-cloud {
			background-color: #e55e3f;
		  }
		  .error-fatal {
			background-color: #e55e3f;
		  }
		  .error-local {
			background-color: #e55e3f;
		  }
		  .ulp-popover.error {
			background-color: #e55e3f;
			border-color: #e55e3f;
		  }
		  .ulp-popover.error::before {
			border-bottom-color: #e55e3f;
		  }
		  .ulp-popover.error::after {
			border-bottom-color: #e55e3f;
		  }
		  #alert-trigger {
			background-color: #e55e3f;
		  }
		</style>
		<style>
		  /* By default, hide features for javascript-disabled browsing */
		  /* We use !important to override any css with higher specificity */
		  /* It is also overriden by the styles in <noscript> in the header file */
		  .no-js {
			display: none !important;
		  }
		</style>
		<noscript>
		  <style>
			/* We use !important to override the default for js enabled */
			/* If the display should be other than block, it should be defined specifically here */
			.js-required {
			  display: none !important;
			}
			.no-js {
			  display: block !important;
			}
			div.input-container input {
			  padding-right: 16px !important;
			}
			div.input-container::after {
			  display: none !important;
			}
		  </style>
		</noscript>
	  </head>

	  <body>
		<div class="main-wrapper">
		  <main class="ulp-outer">
			<section class="ulp-box email-verification-result ulp-event-outer ">
			  <div class="ulp-box-inner">
				<div class="ulp-container ulp-event-screen" data-event-id="">
				  <div class="event-img-container">
					<span class="event-img error-cross"></span>
				  </div>

				  <section class="event-container">
					<h3 class="event-title">Error</h3>

					<div class="event-text">
					  {{.ErrorMessage}}<br />
					  You can close this page now.
					</div>
				  </section>
				</div>
			  </div>

			  <a
				href="https://auth0.com/?utm_source=lock&amp;utm_campaign=badge&amp;utm_medium=widget"
				target="_blank"
				rel="noopener noreferrer"
				class="ulp-auth0-badge"
				aria-label="Link to the Auth0 website"
			  >
				<img
				  src="https://auth0-ulp.herokuapp.com/static/img/theme-generic/auth0-logo.svg"
				  alt="Link to the Auth0 website"
				  class="ulp-auth0-badge-image"
				/>
			  </a>
			</section>
		  </main>
		</div>
	  </body>
	</html>
`
	return template.New("error").Parse(errorTemplate)
}

// NewHTMLRenderer func
func NewHTMLRenderer() *HTMLRenderer {
	return &HTMLRenderer{}
}
