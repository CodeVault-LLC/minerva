package routes

import "net/http"

func serveReDoc(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
