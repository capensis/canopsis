if (!db.default_rights.findOne({_id: "crudcontext_massEnable"})) {
    db.default_rights.insertOne({
        _id: "crudcontext_massEnable",
        crecord_name: "crudcontext_massEnable",
        crecord_type: "action",
        desc: "Rights on context: Mass enable"
    });
}

if (!db.default_rights.findOne({_id: "crudcontext_massDisable"})) {
    db.default_rights.insertOne({
        _id: "crudcontext_massDisable",
        crecord_name: "crudcontext_massDisable",
        crecord_type: "action",
        desc: "Rights on context: Mass disable"
    });
}

db.default_rights.updateOne(
    {
        _id: "admin"
    },
    {
        $set: {
            "rights.crudcontext_massEnable": {
                checksum: 1
            },
            "rights.crudcontext_massDisable": {
                checksum: 1
            }
        }
    }
);
