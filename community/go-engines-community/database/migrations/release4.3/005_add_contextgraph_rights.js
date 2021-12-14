db.default_rights.insertMany([
    {
        _id: "api_contextgraph",
        loader_id: "api_contextgraph",
        crecord_name: "api_contextgraph",
        crecord_type: "action",
        desc: "Context graph import",
        type: "CRUD"
    },
]);
db.default_rights.find({
    "crecord_name": {
        "$in": [
            "api_contextgraph",
        ],
    }
}).forEach(function (doc) {
    db.default_rights.update(
        {
            crecord_name: "admin",
            crecord_type: "role",
        },
        {
            $set: {
                ['rights.' + doc._id]: {
                    checksum: 15,
                    crecord_type: "right",
                },
            },
        }
    )
});
