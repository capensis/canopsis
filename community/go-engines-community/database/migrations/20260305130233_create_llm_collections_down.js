db.api_llm_config.drop();

db.permission.deleteMany({
    _id: {
        $in: [
            "api_llm_config",
        ]
    }
});
db.role.updateMany({}, {
    $unset: {
        "permissions.api_llm_config": "",
    }
});

db.configuration.updateOne({_id: "data_storage"}, {
    $unset: {
        "config.llm_chat": "",
    }
});
