if (db.configuration.findOne({_id: "data_storage"})) {
    db.configuration.updateOne(
        {
            _id: "data_storage",
            "config.entity_infos_log.delete_after": null
        },
        {
            $set: {
                "config.entity_infos_log.delete_after": {
                    "value": 7,
                    "unit": "d",
                    "enabled": true,
                },
            }
        },
    );
} else {
    db.configuration.insertOne({
        _id: "data_storage",
        "config.entity_infos_log.delete_after": {
            "value": 7,
            "unit": "d",
            "enabled": true,
        },
    });
}
