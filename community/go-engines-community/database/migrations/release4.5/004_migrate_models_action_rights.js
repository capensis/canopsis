(function () {
    db.default_rights.deleteOne({"_id": "models_action"});

    db.default_rights.insertOne({
        _id: "models_permission",
        loader_id: "models_permission",
        crecord_name: "models_permission",
        crecord_type: "action",
        desc: "Rights",
        type: "CRUD"
    });

    db.default_rights.find({crecord_type: "role", "rights.models_action": {$exists: true}}).forEach(function (doc) {
        db.default_rights.updateOne(
            {_id: doc._id},
            {
                $set: {
                    "rights.models_permission": doc.rights.models_action,
                },
                $unset: {
                    "rights.models_action": "",
                }
            }
        )
    })
})()
