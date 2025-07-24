db.template_data.dropIndex("name_1");

db.permission.deleteMany({
    _id: {
        $in: [
            "api_template_data",
        ]
    }
});

db.role.updateMany({}, {
    $unset: {
        "permissions.api_template_data": "",
    }
});
