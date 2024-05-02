db.permission.deleteMany({
    _id: {
        $in: [
            "availability_interval",
            "availability_listFilters",
            "availability_editFilter",
            "availability_addFilter",
            "availability_userFilter",
            "availability_exportAsCsv",
        ]
    }
});
db.role.updateMany({}, {
    $unset: {
        "permissions.availability_interval": "",
        "permissions.availability_listFilters": "",
        "permissions.availability_editFilter": "",
        "permissions.availability_addFilter": "",
        "permissions.availability_userFilter": "",
        "permissions.availability_exportAsCsv": "",
    }
});
