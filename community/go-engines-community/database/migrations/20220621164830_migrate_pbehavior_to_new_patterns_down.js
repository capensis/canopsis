// down script can be applied only if no rule has been updated by API
db.pbehavior.updateMany({}, {
    $rename: {
        old_mongo_query: "filter",
    },
    $unset: {
        entity_pattern: "",
    },
});
