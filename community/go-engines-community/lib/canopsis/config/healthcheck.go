package config

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"

var defaultConfig = []EngineOrder{
	{
		From: canopsis.FIFOEngineName,
		To: canopsis.CheEngineName,
	},
	{
		From: canopsis.CheEngineName,
		To: canopsis.PBehaviorEngineName,
	},
	{
		From: canopsis.PBehaviorEngineName,
		To: canopsis.AxeEngineName,
	},
	{
		From: canopsis.AxeEngineName,
		To: canopsis.CorrelationEngineName,
	},
	{
		From: canopsis.AxeEngineName,
		To: canopsis.RemediationEngineName,
	},
	{
		From: canopsis.CorrelationEngineName,
		To: canopsis.ServiceEngineName,
	},
	{
		From: canopsis.ServiceEngineName,
		To: canopsis.DynamicInfosEngineName,
	},
	{
		From: canopsis.DynamicInfosEngineName,
		To: canopsis.ActionEngineName,
	},
	{
		From: canopsis.ActionEngineName,
		To: canopsis.WebhookEngineName,
	},
}

type HealthCheckConf struct {
	EngineOrder []EngineOrder `toml:"engine_order" bson:"engine_order"`
}

type EngineOrder struct {
	From string `toml:"from" bson:"from"`
	To   string `toml:"to" bson:"to"`
}
