(function () {
    db.default_rights.find(
        {
            "crecord_type": "user", "external_id": { "$in": [null, ""] },
            "external": true
        }).forEach(doc => {
            db.default_rights.update({ _id: doc._id }, { "$set": { "external_id": doc.crecord_name } })
        });
})()