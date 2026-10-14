db.configuration.updateOne({_id: "data_storage"}, {
    $unset: {
        "config.llm_chat.delete_token_usage_after": "",
    }
});
