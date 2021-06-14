db.default_entities.updateMany(
    {
        type: "watcher"
    },
    {
        $set: {
            type: "service"
        },
        $rename: {
            "entities": "entity_patterns"
        }
    }
);