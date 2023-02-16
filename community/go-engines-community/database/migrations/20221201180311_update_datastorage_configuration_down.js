db.configuration.updateOne({"_id": "data_storage"}, {
    $rename: {
        "config.remediation.delete_after": "config.remediation.accumulate_after",
    }
});
db.configuration.updateOne({"_id": "data_storage"}, {
    $rename: {
        "config.remediation.delete_stats_after": "config.remediation.delete_after",
    }
});
