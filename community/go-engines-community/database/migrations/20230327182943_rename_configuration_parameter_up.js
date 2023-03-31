db.configuration.updateOne({_id: "global_config"}, {
    $rename: {
        "metrics.enabledmanualinstructions": "metrics.enabledinstructions"
    }
});
