db.permission.find().forEach(function (doc) {
    doc.crecord_name = doc.name;
    delete doc.name;
    doc.crecord_type = "action";

    db.default_rights.updateOne(
        {crecord_name: doc.crecord_name},
        {$setOnInsert: doc},
        {upsert: true},
    );
});

db.role.find().forEach(function (doc) {
    doc.crecord_name = doc.name;
    delete doc.name;
    doc.crecord_type = "role";
    doc.rights = {};
    for (var permission of Object.keys(doc.permissions)) {
        doc.rights[permission] = {checksum: doc.permissions[permission]};
    }
    delete doc.permissions;

    db.default_rights.updateOne(
        {crecord_name: doc.crecord_name},
        {$setOnInsert: doc},
        {upsert: true},
    );
});

db.user.find().forEach(function (doc) {
    doc.crecord_name = doc.name;
    delete doc.name;
    doc.crecord_type = "user";
    doc.shadowpasswd = doc.password;
    delete doc.password;

    db.default_rights.updateOne(
        {crecord_name: doc.name},
        {$setOnInsert: doc},
        {upsert: true},
    );
});

db.permission.drop();
db.role.drop();
db.user.drop();
