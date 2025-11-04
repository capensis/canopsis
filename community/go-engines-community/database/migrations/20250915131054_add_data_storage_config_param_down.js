db.configuration.updateOne({_id: "data_storage"}, {
    $unset: {
        "config.entity_infos_log": "",
    }
});
