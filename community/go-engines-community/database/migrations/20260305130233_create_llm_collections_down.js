db.api_llm_config.drop();

db.permission.deleteMany({
    _id: {
        $in: [
            "api_llm_config",
            "api_llm_chat",
        ]
    }
});
db.role.updateMany({}, {
    $unset: {
        "permissions.api_llm_config": "",
        "permissions.api_llm_chat": "",
    }
});

db.configuration.updateOne({_id: "data_storage"}, {
    $unset: {
        "config.llm_chat": "",
    }
});
