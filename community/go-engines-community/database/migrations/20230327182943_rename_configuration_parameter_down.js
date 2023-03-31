db.configuration.updateOne({_id: "global_config"}, {
    $rename: {
        "metrics.enabledinstructions": "metrics.enabledmanualinstructions"
    }
});
