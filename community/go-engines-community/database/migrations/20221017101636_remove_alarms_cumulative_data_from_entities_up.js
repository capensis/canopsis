db.default_entities.updateMany({}, {$unset: {alarms_cumulative_data: ""}});
