db.external_data_table.drop();

db.permission.deleteMany({
    _id: {
        $in: [
            "api_external_data_table",
        ]
    }
});
db.role.updateMany({}, {
    $unset: {
        "permissions.api_external_data_table": "",
    }
});
