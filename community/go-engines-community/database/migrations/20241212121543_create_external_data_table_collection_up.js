if (!db.permission.findOne({_id: "api_external_data_table"})) {
    db.permission.insertOne({
        _id: "api_external_data_table",
        name: "api_external_data_table",
        type: "CRUD",
        description: "External data",
        groups: ["api", "api_rules"]
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_external_data_table": 15
        }
    });
}

if (!db.permission.findOne({_id: "models_exploitation_externalData"})) {
    db.permission.insertOne({
        _id: "models_exploitation_externalData",
        name: "models_exploitation_externalData",
        type: "CRUD",
        description: "External data",
        groups: ["technical", "technical_exploitation"],
        api_permissions: {
            api_external_data_table: 0
        }
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.models_exploitation_externalData": 15
        }
    });
}

if (!db.permission.findOne({_id: "ExternalData"})) {
    db.permission.insertOne({
        _id: "ExternalData",
        name: "ExternalData",
        description: "API permissions for ExternalData widget",
        hidden: true,
        api_permissions_bitmask: {
            4: {
                api_view: 4,
                api_associative_table: 4,
                api_external_data_table: 6
            },
            2: {
                api_view: 2,
                api_widgettemplate: 4,
                api_external_data_table: 6
            },
            1: {
                api_view: 1
            }
        }
    });
    db.permission.updateOne({_id: "models_privateView"}, {
        $set: {
            "api_permissions.api_external_data_table": 6
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
                    let hasCol = {};
                    let columns = [];
                    let columnTypes = [];
                    db.getCollection(d.collection).find({}, {}, {limit: 10}).forEach(function (item) {
                        for (const c of Object.keys(item)) {
                            if (c !== "_id" && !hasCol[c]) {
                                hasCol[c] = true;
                                columns.push(c);
                                columnTypes.push(0);
                            }
                        }
                    });

                    d.table_type = 0;
                    d.table_name = d.collection;
                    d.table_columns = columns;

                    tables.push({
                        _id: tableID,
                        type: 0,
                        name: d.collection,
                        from_config: true,
                        columns: columns,
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
