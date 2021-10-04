(function () {
    db.default_rights.insertMany([
        {
            _id: "api_healthcheck",
            loader_id: "api_healthcheck",
            crecord_name: "api_healthcheck",
            crecord_type: "action",
            desc: "Healthcheck",
        },
    ]);
    db.default_rights.deleteOne({_id: "api_engine"})
    db.default_rights.updateMany(
        {
            crecord_name: "admin",
            crecord_type: "role",
            "rights.api_engine.checksum": 1,
        },
        {
            $unset: {"rights.api_engine": ""},
            $set: {
                "rights.api_healthcheck": {
                    checksum: 1,
                    crecord_type: "right",
                },
            },
        }
    );
    db.default_rights.updateMany(
        {
            crecord_name: "admin",
            crecord_type: "role",
            "rights.api_engine": {$exists: true},
        },
        {
            $unset: {"rights.api_engine": ""},
        }
    );
})();
