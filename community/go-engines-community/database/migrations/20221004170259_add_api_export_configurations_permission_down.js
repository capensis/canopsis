db.default_rights.deleteOne({
    _id: "api_export_configurations"
});
db.default_rights.updateMany({crecord_type: "role"}, {
    $unset: {
        "rights.api_export_configurations": "",
    }
});
