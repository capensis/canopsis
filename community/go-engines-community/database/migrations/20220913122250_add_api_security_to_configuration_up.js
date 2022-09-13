if (!db.configuration.findOne({_id: "api_security"})) {
    db.configuration.insertOne({
        _id: "api_security",
        basic: {
            expiration_interval: {
                value: 1,
                unit: "M"
            }
        }
    });
}
