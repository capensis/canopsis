db.default_rights.deleteOne({
    _id: "api_metrics_settings"
});
db.default_rights.updateMany({crecord_type: "role"}, {
    $unset: {
        "rights.api_metrics_settings": "",
    }
});
