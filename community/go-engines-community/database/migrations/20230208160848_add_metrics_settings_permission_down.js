db.default_rights.deleteMany({
    _id: {
        $in: [
            "api_metrics_settings",
            "models_kpiCollectionSettings",
        ]
    }
});
db.default_rights.updateMany({crecord_type: "role"}, {
    $unset: {
        "rights.api_metrics_settings": "",
        "rights.models_kpiCollectionSettings": "",
    }
});
