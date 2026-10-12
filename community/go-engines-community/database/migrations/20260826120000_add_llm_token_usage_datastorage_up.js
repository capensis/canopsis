if (db.configuration.findOne({_id: "data_storage"})) {
    db.configuration.updateOne(
        {
            _id: "data_storage",
            "config.llm_chat.delete_token_usage_after": null
        },
        {
            $set: {
                "config.llm_chat.delete_token_usage_after": {
                    "value": 90,
                    "unit": "d",
                    "enabled": true,
                },
            }
        },
    );
} else {
    db.configuration.insertOne({
        _id: "data_storage",
        "config.llm_chat.delete_token_usage_after": {
            "value": 90,
            "unit": "d",
            "enabled": true,
        },
    });
}
