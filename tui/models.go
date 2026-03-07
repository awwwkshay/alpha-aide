package main

import agentmodels "github.com/awwwkshay/alpha-aide/agent/models"

// Re-export the agent's model catalog for use within the TUI.
type ModelDef = agentmodels.Model

var KnownModels = agentmodels.All
