db.permission.deleteOne({
    _id: "api_techmetrics_settings"
});
db.role.updateMany({}, {
    $unset: {
        "permissions.api_techmetrics_settings": "",
    }
});
