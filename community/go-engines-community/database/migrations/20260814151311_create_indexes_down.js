db.default_entities.getIndexes().forEach(function (index) {
    if (index.name === "connector_1_services_1" || index.name === "type_1_enabled_1_idle_since_1") {
        db.default_entities.dropIndex(index.name);
    }
});
