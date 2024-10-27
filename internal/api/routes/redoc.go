package routes

import "github.com/gofiber/fiber/v2"

func serveReDoc(c *fiber.Ctx) error {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>ReDoc</title>
		<!-- ReDoc CSS -->
		<link href="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.css" rel="stylesheet">
	</head>
	<body>
		<redoc spec-url='/swagger/swagger.json'></redoc>
		<script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"></script>
	</body>
	</html>
	`

	c.Set("Content-Type", "text/html")
	c.SendString(html)
	return nil
}
