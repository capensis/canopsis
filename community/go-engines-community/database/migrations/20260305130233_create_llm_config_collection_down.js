db.api_llm_config.drop();

db.permission.deleteMany({
    _id: {
        $in: [
            "api_api_llm_config",
        ]
    }
});
db.role.updateMany({}, {
    $unset: {
        "permissions.api_api_llm_config": "",
    }
});
