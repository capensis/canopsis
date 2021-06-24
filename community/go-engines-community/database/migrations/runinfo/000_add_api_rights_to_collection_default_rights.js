db.default_rights.insertMany([
    {
        _id: "api_engine",
        loader_id: "api_engine",
        crecord_name: "api_engine",
        crecord_type: "action",
        desc: "Engine Info",
    },
]);
db.default_rights.find({
    "crecord_name": {
        "$in": [
            "api_engine",
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
                    checksum: 1,
                    crecord_type: "right",
                },
            },
        }
    )
});
