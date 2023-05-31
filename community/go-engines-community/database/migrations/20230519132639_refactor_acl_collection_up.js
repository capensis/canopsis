db.default_rights.find({crecord_type: "action"}).forEach(function (doc) {
    doc.name = doc.crecord_name;
    delete doc.crecord_name;
    delete doc.crecord_type;

    db.permission.updateOne(
        {name: doc.name},
        {$setOnInsert: doc},
        {upsert: true},
    );
});

db.default_rights.find({crecord_type: "role"}).forEach(function (doc) {
    doc.name = doc.crecord_name;
    delete doc.crecord_name;
    delete doc.crecord_type;
    doc.permissions = {};
    for (var permission of Object.keys(doc.rights)) {
        doc.permissions[permission] = doc.rights[permission].checksum;
    }
    delete doc.rights;

    db.role.updateOne(
        {name: doc.name},
        {$setOnInsert: doc},
        {upsert: true},
    );
});

db.default_rights.find({crecord_type: "user"}).forEach(function (doc) {
    doc.name = doc.crecord_name;
    delete doc.crecord_name;
    delete doc.crecord_type;
    doc.password = doc.shadowpasswd;
    delete doc.shadowpasswd;
    if (doc.role) {
        doc.roles = [doc.role];
    } else {
        doc.roles = [];
    }
    delete doc.role;
    db.user.updateOne(
        {name: doc.name},
        {$setOnInsert: doc},
        {upsert: true},
    );
});

db.default_rights.drop();
