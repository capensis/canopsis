if (db.getCollectionNames().includes("entity_service_counters")) {
    db.entity_service_counters.renameCollection("entity_counters")
}
