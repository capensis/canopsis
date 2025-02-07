db.configuration.updateOne({_id: "global_config"}, {$unset: {"metrics.enabledslimetrics": ""}})
