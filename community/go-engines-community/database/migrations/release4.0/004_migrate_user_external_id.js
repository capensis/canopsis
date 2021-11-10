(function () {
    db.default_rights.find(
        {
            "crecord_type": "user", "external": true 
        }).forEach(doc => {
            db.default_rights.updateMany({ _id: doc._id }, { "$set": { "crecord_name": doc._id, "source": "ldap", "external_id": doc._id } })
            db.default_rights.updateMany({ _id: doc._id }, { "$unset": { "rights": 1 } })
        });
})()
