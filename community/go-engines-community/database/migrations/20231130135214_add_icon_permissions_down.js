db.icon.drop()

db.permission.deleteMany({
    _id: {
        $in: [
            "api_icon",
            "models_stateSetting",
            "models_icon",
            "models_storageSettings",
        ]
    }
});
db.role.updateMany({}, {
    $unset: {
        "permissions.api_icon": "",
        "permissions.models_stateSetting": "",
        "permissions.models_icon": "",
        "permissions.models_storageSettings": "",
    }
});
