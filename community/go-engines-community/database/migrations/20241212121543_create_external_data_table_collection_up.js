if (!db.permission.findOne({_id: "api_external_data_table"})) {
    db.permission.insertOne({
        _id: "api_external_data_table",
        name: "api_external_data_table",
        description: "External data",
        groups: ["api", "api_rules"]
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_external_data_table": 15
        }
    });
}

db.external_data_table.createIndex({type: 1, name: 1}, {name: "type_1_name_1", unique: true});
const now = Math.ceil((new Date()).getTime() / 1000);
let tableByName = {};

for (const ruleCollName of ["eventfilter", "link_rule"]) {
    const ruleCollection = db.getCollection(ruleCollName);
    ruleCollection.find({
        external_data: {
            $nin: [null, {}],
            $not: {$type: "array"}
        },
    }).forEach(function (doc) {
        let externalData = [];
        let tables = [];
        for (const ref of Object.keys(doc.external_data)) {
            let d = doc.external_data[ref];
            d.reference = ref;
            if (d.type === "mongo") {
                d.type = "table";
                let tableID = tableByName[d.collection];
                if (tableID === undefined) {
                    tableID = genID();
                    tableByName[d.collection] = tableID;
                    let columnTypes = {};
                    db.getCollection(d.collection).find({}, {}, {limit: 10}).forEach(function (item) {
                        for (const c of Object.keys(item)) {
                            if (c !== "_id") {
                                columnTypes[c] = 0
                            }
                        }
                    });

                    tables.push({
                        _id: tableID,
                        type: 0,
                        name: d.collection,
                        from_config: true,
                        column_types: columnTypes,
                        created: now,
                        updated: now,
                    });
                }

                d.table = tableID;
                delete d.collection;
            }

            externalData.push(d);
        }

        if (tables.length > 0) {
            db.external_data_table.insertMany(tables);
        }

        ruleCollection.updateOne({_id: doc._id}, {$set: {external_data: externalData}});
    });

    ruleCollection.updateMany({external_data: {}}, {$set: {external_data: []}});
}
