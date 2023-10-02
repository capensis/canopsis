db.default_entities.createIndex({type: 1}, {name: "type_service_1", partialFilterExpression: {type: {$eq: "service" }}})
