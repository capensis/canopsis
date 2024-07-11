var collectionNames = db.getCollectionNames();

if (collectionNames.includes("user")) {
    db.runCommand({collMod: "user", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("user", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("role")) {
    db.runCommand({collMod: "role", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("role", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("eventfilter")) {
    db.runCommand({collMod: "eventfilter", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("eventfilter", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("view_playlist")) {
    db.runCommand({collMod: "view_playlist", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("view_playlist", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("action_scenario")) {
    db.runCommand({collMod: "action_scenario", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("action_scenario", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("meta_alarm_rules")) {
    db.runCommand({collMod: "meta_alarm_rules", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("meta_alarm_rules", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("dynamic_infos")) {
    db.runCommand({collMod: "dynamic_infos", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("dynamic_infos", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("entity_category")) {
    db.runCommand({collMod: "entity_category", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("entity_category", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("pbehavior_type")) {
    db.runCommand({collMod: "pbehavior_type", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("pbehavior_type", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("pbehavior_reason")) {
    db.runCommand({collMod: "pbehavior_reason", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("pbehavior_reason", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("pbehavior_exception")) {
    db.runCommand({collMod: "pbehavior_exception", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("pbehavior_exception", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("pbehavior")) {
    db.runCommand({collMod: "pbehavior", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("pbehavior", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("job_config")) {
    db.runCommand({collMod: "job_config", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("job_config", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("job")) {
    db.runCommand({collMod: "job", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("job", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("state_settings")) {
    db.runCommand({collMod: "state_settings", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("state_settings", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("broadcast_message")) {
    db.runCommand({collMod: "broadcast_message", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("broadcast_message", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("idle_rule")) {
    db.runCommand({collMod: "idle_rule", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("idle_rule", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("views")) {
    db.runCommand({collMod: "views", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("views", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("viewtabs")) {
    db.runCommand({collMod: "viewtabs", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("viewtabs", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("widgets")) {
    db.runCommand({collMod: "widgets", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("widgets", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("widget_filters")) {
    db.runCommand({collMod: "widget_filters", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("widget_filters", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("widget_templates")) {
    db.runCommand({collMod: "widget_templates", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("widget_templates", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("viewgroups")) {
    db.runCommand({collMod: "viewgroups", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("viewgroups", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("resolve_rule")) {
    db.runCommand({collMod: "resolve_rule", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("resolve_rule", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("flapping_rule")) {
    db.runCommand({collMod: "flapping_rule", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("flapping_rule", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("kpi_filter")) {
    db.runCommand({collMod: "kpi_filter", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("kpi_filter", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("pattern")) {
    db.runCommand({collMod: "pattern", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("pattern", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("map")) {
    db.runCommand({collMod: "map", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("map", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("default_snmprules")) {
    db.runCommand({collMod: "default_snmprules", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("default_snmprules", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("declare_ticket_rule")) {
    db.runCommand({collMod: "declare_ticket_rule", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("declare_ticket_rule", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("link_rule")) {
    db.runCommand({collMod: "link_rule", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("link_rule", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("alarm_tag")) {
    db.runCommand({collMod: "alarm_tag", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("alarm_tag", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("color_theme")) {
    db.runCommand({collMod: "color_theme", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("color_theme", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("icon")) {
    db.runCommand({collMod: "icon", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("icon", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("instruction")) {
    db.instruction.updateMany({}, {
        $rename: {
            last_modified: "updated",
        }
    });

    db.runCommand({collMod: "instruction", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("instruction", {changeStreamPreAndPostImages: {enabled: true}})
}

if (collectionNames.includes("default_entities")) {
    db.runCommand({collMod: "default_entities", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("default_entities", {changeStreamPreAndPostImages: {enabled: true}})
}
