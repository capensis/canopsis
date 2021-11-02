(function () {
    // Add Enabled field to existing dynamic_infos records
    db.dynamic_infos.updateMany({
        "enabled": { "$exists": false }
    }, {
        "$set": {
            "enabled": true
        }
    });
})();