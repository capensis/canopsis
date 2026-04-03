const updatePerms = {
    models_maintenance: {
        groups: ["technical", "technical_admin", "technical_admin_maintenance"]
    },
    models_planningType: {
        description: "Planning - types",
        groups: ["technical", "technical_admin", "technical_admin_maintenance"]
    },
    models_planningReason: {
        description: "Planning - reasons",
        groups: ["technical", "technical_admin", "technical_admin_maintenance"]
    },
    models_planningExceptions: {
        description: "Planning - exception dates",
        groups: ["technical", "technical_admin", "technical_admin_maintenance"]
    },
    models_eventsRecord: {
        description: "Event records"
    },
    models_templateTesting: {
        description: "Go templates"
    },
    models_externalAuthTokens: {
        description: "External auth tokens",
        groups: ["technical", "technical_admin", "technical_admin_customobjects"]
    },
    models_icon: {
        description: "Icons",
        groups: ["technical", "technical_admin", "technical_admin_customobjects"]
    },
    models_map: {
        groups: ["technical", "technical_admin", "technical_admin_customobjects"]
    },
    models_tag: {
        description: "Tags",
        groups: ["technical", "technical_admin", "technical_admin_customobjects"]
    },
    models_view_import_export: {
        description: "Import / export",
        groups: ["technical", "technical_admin", "technical_admin_settings"]
    },
    models_notification_common: {
        description: "Notification",
        groups: ["technical", "technical_admin", "technical_admin_settings"]
    },
    models_stateSetting: {
        description: "State",
        groups: ["technical", "technical_admin", "technical_admin_settings"]
    },
    models_storageSettings: {
        description: "Storage",
        groups: ["technical", "technical_admin", "technical_admin_settings"]
    },
    models_parameters: {
        description: "User interface",
        groups: ["technical", "technical_admin", "technical_admin_settings"]
    },
    models_widgetTemplate: {
        description: "Widget templates",
        groups: ["technical", "technical_admin", "technical_admin_settings"]
    },
};

const replacePerms = {
    models_exploitation_entityInfoProperty: {
        _id: "models_entityInfoProperty",
        description: "Entity infos",
        groups: ["technical", "technical_admin", "technical_admin_customobjects"]
    },
    models_exploitation_externalData: {
        _id: "models_externalData",
        description: "External data",
        groups: ["technical", "technical_admin", "technical_admin_customobjects"]
    },
    models_remediationInstruction: {
        _id: "models_exploitation_remediationInstruction",
        description: "Instructions - instructions",
        groups: ["technical", "technical_exploitation"]
    },
    models_remediationInstructionApprove: {
        _id: "models_exploitation_remediationInstructionApprove",
        description: "Instructions - instruction approve",
        groups: ["technical", "technical_exploitation"]
    },
    models_remediationJob: {
        _id: "models_exploitation_remediationJob",
        description: "Instructions - jobs",
        groups: ["technical", "technical_exploitation"]
    },
    models_remediationConfiguration: {
        _id: "models_exploitation_remediationConfiguration",
        description: "Instructions - configurations",
        groups: ["technical", "technical_exploitation"]
    },
    models_remediationStatistic: {
        _id: "models_exploitation_remediationStatistic",
        description: "Instructions - remediation statistics",
        groups: ["technical", "technical_exploitation"]
    },
    models_instructionStats: {
        _id: "models_exploitation_instructionStats",
        description: "Instructions - instruction rating",
        groups: ["technical", "technical_exploitation"]
    }
};

for (const id in updatePerms) {
    const n = updatePerms[id]
    let set = {};
    if (n.description) {
        set["description"] = n.description;
    }

    if (n.groups) {
        set["groups"] = n.groups;
    }

    db.permission.updateOne({_id: id}, {$set: set});
}

for (const id in replacePerms) {
    let p = db.permission.findOneAndDelete({_id: id});
    if (!p) {
        continue;
    }

    const {_id: newID, description, groups} = replacePerms[id];
    if (db.permission.findOne({_id: newID})) {
        db.role.updateMany({["permissions." + id]: {$ne: null}}, {
            $unset: {["permissions." + id]: ""}
        });

        continue;
    }

    p._id = newID;
    p.name = newID;
    p.description = description;
    p.groups = groups;
    db.permission.insertOne(p);
    db.role.updateMany({["permissions." + id]: {$ne: null}}, [
        {
            $set: {
                ["permissions." + newID]: "$permissions." + id
            }
        },
        {
            $unset: ["permissions." + id]
        }
    ]);
}

function insertGroupAfter(newGroup, prevGroup) {
    if (db.permission_group.findOne({_id: newGroup})) {
        return;
    }

    const prev = db.permission_group.findOne({_id: prevGroup});
    const position = prev.position + 1;
    db.permission_group.updateMany({position: {$gte: position}}, [
        {
            $set: {
                position: {"$sum": ["$position", 1]}
            }
        }
    ]);
    db.permission_group.insertOne({
        _id: newGroup,
        name: newGroup,
        position: position
    });
}

insertGroupAfter("technical_admin_maintenance", "technical_admin_access");
insertGroupAfter("technical_admin_customobjects", "technical_admin_general");
insertGroupAfter("technical_admin_settings", "technical_admin_customobjects");
