package config

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"

var defaultConfig = []EnginePair{
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
	EnginePairs []EnginePair `toml:"engine_pairs" bson:"engine_pairs"`
}

type EnginePair struct {
	From string `toml:"from" bson:"from"`
	To   string `toml:"to" bson:"to"`
}
