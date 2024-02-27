package main

import (
	"nero/server"
)

// BasePath /api/v0.1
// @title	Nero application docs [user and orders api]
// @version 1.0
// @description backen server of nero app. https://github.com/Furqanalimir/nero
// @contact.name API Support
// @contact.url  https://furqanali.vercel.app/
// @contact.email mrifurqan89@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:5050
// @query.collection.format multi

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	server.Init()
}
