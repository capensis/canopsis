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

let collectionByTable = {};
for (const ruleCollName of ["eventfilter", "link_rule"]) {
    const ruleCollection = db.getCollection(ruleCollName);
    ruleCollection.find({
        external_data: {
            $nin: [null, []],
            $type: "array",
        },
    }).forEach(function (doc) {
        let externalData = {};
        for (const d of doc.external_data) {
            if (d.type === "table") {
                d.type = "mongo";
                let collectionName = collectionByTable[d.table];
                if (collectionName === undefined) {
                    collectionName = db.external_data_table.findOne({_id: d.table}).name;
                    collectionName[d.table] = collectionName;
                }

                d.collection = collectionName;
                delete d.table;
            }

            const ref = d.reference
            delete d.reference;
            externalData[ref] = d;
        }

        ruleCollection.updateOne({_id: doc._id}, {$set: {external_data: externalData}});
    });

    ruleCollection.updateMany({external_data: []}, {$set: {external_data: {}}});
}

db.external_data_table.drop();
