if (db.getCollectionNames().includes("entity_counters")) {
    db.entity_counters.renameCollection("entity_service_counters")
}
