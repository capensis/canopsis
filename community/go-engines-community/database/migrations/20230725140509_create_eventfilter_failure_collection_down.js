db.eventfilter_failure.dropIndex("rule_1");

db.configuration.updateOne({_id: "data_storage"}, {
    $unset: {
        "config.event_filter_failure": "",
    }
});
