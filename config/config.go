package config

import "os"

var DVCApiServerEndpoint = os.Getenv("DVC_API_SERVER_ENDPOINT")
