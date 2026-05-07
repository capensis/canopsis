db.permission.deleteOne({
    _id: "models_templateData"
});
db.role.updateMany({}, {
    $unset: {
        "permissions.models_templateData": ""
    }
});
db.role.updateMany({type: "ui"}, {
    $pull: {
        permissions: "api_template_data"
    }
});