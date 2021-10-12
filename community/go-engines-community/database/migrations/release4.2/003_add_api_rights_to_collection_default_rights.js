db.default_rights.insertMany([
    {
        _id: "api_execution",
        crecord_name: "api_execution",
        crecord_type: "action",
        desc: "Runs instructions",
    },
    {
        _id: "api_job_config",
        crecord_name: "api_job_config",
        crecord_type: "action",
        desc: "Job configs",
        type: "CRUD",
    },
    {
        _id: "api_job",
        crecord_name: "api_job",
        crecord_type: "action",
        desc: "Jobs",
        type: "CRUD",
    },
    {
        _id: "api_instruction",
        crecord_name: "api_instruction",
        crecord_type: "action",
        desc: "Instructions",
        type: "CRUD",
    },
    {
        _id: "api_file",
        crecord_name: "api_file",
        crecord_type: "action",
        desc: "File",
        type: "CRUD",
    },
    {
        _id: "api_metaalarmrule",
        crecord_name: "api_metaalarmrule",
        crecord_type: "action",
        desc: "Meta-alarm rules",
        type: "CRUD",
    },
]);
db.default_rights.find({
    "crecord_name": {
        "$in": [
            "api_execution",
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
db.default_rights.find({
    "crecord_name": {
        "$in": [
            "api_job_config",
            "api_job",
            "api_instruction",
            "api_file",
            "api_metaalarmrule"
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
