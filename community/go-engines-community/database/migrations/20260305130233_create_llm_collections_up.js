if (db.getCollectionNames().includes("llm_config")) {
    db.runCommand({collMod: "llm_config", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("llm_config", {changeStreamPreAndPostImages: {enabled: true}})
}

db.llm_config.createIndex({name: 1}, {name: "name_1", unique: true});
db.llm_chat_history.createIndex({config: 1, user: 1}, {name: "config_1_user_1"});
db.llm_chat_history.createIndex({thread: 1}, {name: "thread_1"});
db.llm_message_history.createIndex({chat: 1, timestamp: 1}, {name: "chat_1_timestamp_1"});
db.llm_message_history.createIndex({turn: 1}, {name: "turn_1"});

if (!db.permission.findOne({_id: "api_llm_config"})) {
    db.permission.insertOne({
        _id: "api_llm_config",
        name: "api_llm_config",
        type: "CRUD",
        description: "LLM config",
        groups: ["api", "api_rules"]
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_llm_config": 15
        }
    });
}

if (!db.permission.findOne({_id: "api_llm_chat"})) {
    db.permission.insertOne({
        _id: "api_llm_chat",
        name: "api_llm_chat",
        description: "LLM chat",
        groups: ["api", "api_general"]
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_llm_chat": 1
        }
    });
}

if (db.configuration.findOne({_id: "data_storage"})) {
    db.configuration.updateOne(
        {
            _id: "data_storage",
            "config.llm_chat.delete_after": null
        },
        {
            $set: {
                "config.llm_chat.delete_after": {
                    "value": 1,
                    "unit": "M",
                    "enabled": true,
                },
            }
        },
    );
} else {
    db.configuration.insertOne({
        _id: "data_storage",
        "config.llm_chat.delete_after": {
            "value": 1,
            "unit": "M",
            "enabled": true,
        },
    });
}
