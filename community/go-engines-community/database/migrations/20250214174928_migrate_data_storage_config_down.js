db.configuration.find({_id: "data_storage"}).forEach(function (doc) {
    let set = {};
    if (doc.history) {
        const keys = [
            "junit",
            "remediation",
            "pbehavior",
            "health_check",
            "webhook",
            "event_filter_failure",
            "event_records",
        ];
        for (const key of keys) {
            const v = doc.history[key];
            if (v) {
                set["history." + key] = v.time;
            }
        }
    }

    let update = {
        $unset: {
            "config.event_records.enabled": "",
        }
    };
    if (Object.keys(set).length > 0) {
        update["$set"] = update;
    }

    db.configuration.updateOne({_id: "data_storage"}, update);
});
