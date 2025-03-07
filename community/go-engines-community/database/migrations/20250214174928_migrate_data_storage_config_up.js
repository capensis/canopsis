db.configuration.find({_id: "data_storage"}).forEach(function (doc) {
    let update = {};
    if (doc.history) {
        const keys = Object.keys(doc.history)
        for (const key of keys) {
            const v = doc.history[key];
            if (isInt(v)) {
                update["history." + key] = {
                    time: v,
                };
            }
        }
    }

    if (doc.config && doc.config.event_records && doc.config.event_records.enabled === undefined) {
        update["config.event_records.enabled"] = true;
    }

    if (Object.keys(update).length > 0) {
        db.configuration.updateOne({_id: "data_storage"}, {$set: update});
    }
});
