db.map.drop();

db.default_rights.deleteMany({
    _id: {
        $in: [
            "api_map",
            "models_map",
            "map_alarmsList",
            "map_listFilters",
            "map_editFilter",
            "map_addFilter",
            "map_userFilter",
            "map_category",
        ]
    }
});
db.default_rights.updateMany({crecord_type: "role"}, {
    $unset: {
        "rights.api_map": "",
        "rights.models_map": "",
        "rights.map_alarmsList": "",
        "rights.map_listFilters": "",
        "rights.map_editFilter": "",
        "rights.map_addFilter": "",
        "rights.map_userFilter": "",
        "rights.map_category": "",
    }
});
