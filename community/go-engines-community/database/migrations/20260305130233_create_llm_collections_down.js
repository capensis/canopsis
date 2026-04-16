db.api_llm_config.drop();
db.llm_chat_history.drop();
db.llm_message_history.drop();

db.permission.deleteMany({
    _id: {
        $in: [
            "api_llm_config",
            "models_llm",
        ]
    }
});
db.role.updateMany({}, {
    $unset: {
        "permissions.api_llm_config": "",
        "permissions.models_llm": "",
    }
});

db.configuration.updateOne({_id: "data_storage"}, {
    $unset: {
        "config.llm_chat": "",
    }
});
