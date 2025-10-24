if (db.getCollectionNames().includes("entity_infos_property")) {
    db.runCommand({collMod: "entity_infos_property", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("entity_infos_property", {changeStreamPreAndPostImages: {enabled: true}})
}

db.entity_infos_property.createIndex({name: 1}, {name: "name_1", unique: true});
db.entity_infos_property.createIndex({alias: 1}, {name: "alias_1", unique: true, partialFilterExpression: {alias: {$exists: true}}});
db.entity_infos_property.createIndex({type: 1}, {name: "type_1"});

db.alarm_tag.createIndex({aliases: 1}, {name: "aliases_1"});
db.default_entities.createIndex({aliases: 1}, {name: "aliases_1"});
db.eventfilter.createIndex({aliases: 1}, {name: "aliases_1"});
db.flapping_rule.createIndex({aliases: 1}, {name: "aliases_1"});
db.idle_rule.createIndex({aliases: 1}, {name: "aliases_1"});
db.link_rule.createIndex({aliases: 1}, {name: "aliases_1"});
db.pattern.createIndex({aliases: 1}, {name: "aliases_1"});
db.pbehavior.createIndex({aliases: 1}, {name: "aliases_1"});
db.resolve_rule.createIndex({aliases: 1}, {name: "aliases_1"});
db.action_scenario.createIndex({aliases: 1}, {name: "aliases_1"});
db.widget_filters.createIndex({aliases: 1}, {name: "aliases_1"});
db.meta_alarm_rules.createIndex({aliases: 1}, {name: "aliases_1"});
db.declare_ticket_rule.createIndex({aliases: 1}, {name: "aliases_1"});
db.dynamic_infos.createIndex({aliases: 1}, {name: "aliases_1"});
db.instruction.createIndex({aliases: 1}, {name: "aliases_1"});
db.kpi_filter.createIndex({aliases: 1}, {name: "aliases_1"});

if (!db.permission.findOne({_id: "api_entity_info_property"})) {
    db.permission.insertOne({
        _id: "api_entity_info_property",
        name: "api_entity_info_property",
        type: "CRUD",
        description: "Entity info properties",
        groups: ["api", "api_general"]
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_entity_info_property": 15
        }
    });
}

if (!db.permission.findOne({_id: "models_exploitation_entityInfoProperty"})) {
    db.permission.insertOne({
        _id: "models_exploitation_entityInfoProperty",
        name: "models_exploitation_entityInfoProperty",
        type: "CRUD",
        description: "Entity info properties",
        groups: ["technical", "technical_exploitation"],
        api_permissions: {
            api_entity_info_property: 0
        }
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.models_exploitation_entityInfoProperty": 15
        }
    });
}
