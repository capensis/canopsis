var existingIndexes = {}; // these indexes can be created independently of the migration
db.default_entities.getIndexes().forEach(function (index) {
    if (index.name === "connector_1_services_1" || index.name === "type_1_enabled_1_idle_since_1") {
        existingIndexes[index.name] = true;
    }
});

if (!existingIndexes["connector_1_services_1"]) {
  db.default_entities.createIndex(
      {connector: 1, services: 1},
      {
          name: "connector_1_services_1",
          partialFilterExpression: {
              enabled: true,
              connector: {$gt: ""},
              "services.0": {$exists: true}
          }
      }
  );
}

if (!existingIndexes["type_1_enabled_1_idle_since_1"]) {
  db.default_entities.createIndex(
    { "type": 1, "enabled": 1, "idle_since": 1 },
    {
      name: "type_1_enabled_1_idle_since_1",
      partialFilterExpression: {
        "type": { "$in": ["resource", "component", "connector"] },
        "enabled": true,
        "idle_since": { "$gt": 0 }
      }
    }
  );
}
