db.eventfilter_failure.createIndex({rule: 1}, {name: "rule_1"});

db.configuration.updateOne({_id: "data_storage", "config.event_filter_failure.delete_after": null}, {
    $set: {
        "config.event_filter_failure.delete_after": {
            "value": 30,
            "unit": "d",
            "enabled": true,
        },
    }
}, {upsert: true});
