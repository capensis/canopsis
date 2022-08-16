db.map.drop();

db.default_rights.deleteMany({
    _id: {
        $in: ["api_map", "models_map"]
    }
});
db.default_rights.updateMany({crecord_type: "role"}, {
    $unset: {
        "rights.api_map": "",
        "rights.models_map": "",
    }
});
