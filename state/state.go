package state

import (
	"database/sql"

	"github.com/aes421/cliStandup/models"
)

var Updates []models.Update
var Db *sql.DB
var LLMConnector models.LLM
var Config models.Config
