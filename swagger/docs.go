package swagger

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

//go:generate ./download.sh
//go:embed dist/*.js dist/*.css dist/oauth2-redirect.html
var files embed.FS

func Register(e *echo.Echo, spec *openapi3.T) {
	g := e.Group("/docs")

	spec.Servers = openapi3.Servers{{
		URL:         "../",
		Description: "Current server",
	}}
	specJSON, err := spec.MarshalJSON()
	if err != nil {
		log.Fatal().Err(err).Msg("could not marshal swagger")
	}

	subFS, err := fs.Sub(files, "dist")
	if err != nil {
		log.Fatal().Err(err).Msg("Sub fs")
	}
	dist := http.StripPrefix("/docs", http.FileServer(http.FS(subFS)))
	index := func(c echo.Context) error {
		return c.HTML(http.StatusOK, indexTemplate)
	}
	g.GET("/", index)
	g.GET("/index.html", index)
	g.GET("/openapi.json", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/json", specJSON)
	})
	g.GET("/*", echo.WrapHandler(dist))
}

const indexTemplate = `<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="./swagger-ui.css" />
    <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
    <style>
      html
      {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }

      *,
      *:before,
      *:after
      {
        box-sizing: inherit;
      }

      body
      {
        margin:0;
        background: #fafafa;
      }
    </style>
  </head>

  <body>
    <div id="swagger-ui"></div>

    <script src="./swagger-ui-bundle.js" charset="UTF-8"> </script>
    <script src="./swagger-ui-standalone-preset.js" charset="UTF-8"> </script>
    <script>
    window.onload = function() {
      // Begin Swagger UI call region
      const ui = SwaggerUIBundle({
        url: "openapi.json",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      });

      window.ui = ui;
    };
  </script>
  </body>
</html>`
