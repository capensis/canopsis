(function () {
    db.notification.insertOne({
        "_id": "notification",
        "instruction": {
            "rate": true,
            "rate_frequency": {
                "seconds": 604800,
                "unit": "s"
            }
        }
    });
})();