if (!db.permission.findOne({_id: "models_templateData"})) {
    db.permission.insertOne({
        _id: "models_templateData",
        name: "models_templateData",
        type: "CRUD",
        description: "Template data management",
        groups: ["technical", "technical_admin", "technical_admin_general"],
        api_permissions: {api_template_data: 0}
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.models_templateData": 15
        }
    });
}
