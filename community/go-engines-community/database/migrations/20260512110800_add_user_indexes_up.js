db.user.createIndex({authkey: 1}, {name: "authkey_1"});
db.user.createIndex(
    {external_id: 1, source: 1},
    {
        name: "external_id_1_source_1",
        partialFilterExpression: {external_id: {$exists: true}},
    }
);
