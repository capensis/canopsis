const deleteAfter = "config.metrics.delete_after";
db.configuration.updateOne({ _id: "data_storage", [deleteAfter]: null },
    {
        $set: {
            [deleteAfter]: {
                "value": 30,
                "unit": "d",
                "enabled": true
            }
        }
    }
);
