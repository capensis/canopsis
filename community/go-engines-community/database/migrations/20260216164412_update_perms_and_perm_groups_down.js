const updatePerms = {
    models_maintenance: {
        groups: ["technical", "technical_admin", "technical_admin_general"]
    },
    models_planningType: {
        description: "Planning type (Pbehavior)",
        groups: ["technical", "technical_admin", "technical_admin_general"]
    },
    models_planningReason: {
        description: "Planning reason (Pbehavior)",
        groups: ["technical", "technical_admin", "technical_admin_general"]
    },
    models_planningExceptions: {
        description: "Planning exceptions dates (Pbehavior)",
        groups: ["technical", "technical_admin", "technical_admin_general"]
    },
    models_eventsRecord: {
        description: "Events recording"
    },
    models_templateTesting: {
        description: "Models template testing"
    },
    models_externalAuthTokens: {
        description: "Webhook token rule",
        groups: ["technical", "technical_admin", "technical_admin_general"]
    },
    models_icon: {
        description: "Parameters - icons",
        groups: ["technical", "technical_admin", "technical_admin_general"]
    },
    models_map: {
        groups: ["technical", "technical_admin", "technical_admin_general"]
    },
    models_tag: {
        description: "Tags management",
        groups: ["technical", "technical_admin", "technical_admin_general"]
    },
    models_view_import_export: {
        description: "Parameters - import / export",
        groups: ["technical", "technical_admin", "technical_admin_general"]
    },
    models_notification_common: {
        description: "Parameters - notification settings",
        groups: ["technical", "technical_admin", "technical_admin_general"]
    },
    models_stateSetting: {
        description: "State",
        groups: ["technical", "technical_admin", "technical_admin_general"]
    },
    models_storageSettings: {
        description: "Storage settings",
        groups: ["technical", "technical_admin", "technical_admin_general"]
    },
    models_parameters: {
        description: "Parameters - parameters tab",
        groups: ["technical", "technical_admin", "technical_admin_general"]
    },
    models_widgetTemplate: {
        description: "Parameters - widget templates",
        groups: ["technical", "technical_admin", "technical_admin_general"]
    },
};

const replacePerms = {
    models_entityInfoProperty: {
        _id: "models_exploitation_entityInfoProperty",
        description: "Entity info properties",
        groups: ["technical", "technical_exploitation"]
    },
    models_externalData: {
        _id: "models_exploitation_externalData",
        description: "External data",
        groups: ["technical", "technical_exploitation"]
    },
    models_exploitation_remediationInstruction: {
        _id: "models_remediationInstruction",
        description: "Instructions - instructions tab",
        groups: ["technical", "technical_admin", "technical_admin_general"]
    },
    models_exploitation_remediationInstructionApprove: {
        _id: "models_remediationInstructionApprove",
        description: "Instructions - instruction approve",
        groups: ["technical", "technical_admin", "technical_admin_general"]
    },
    models_exploitation_remediationJob: {
        _id: "models_remediationJob",
        description: "Instructions - jobs tab",
        groups: ["technical", "technical_admin", "technical_admin_general"]
    },
    models_exploitation_remediationConfiguration: {
        _id: "models_remediationConfiguration",
        description: "Instructions - configurations tab",
        groups: ["technical", "technical_admin", "technical_admin_general"]
    },
    models_exploitation_remediationStatistic: {
        _id: "models_remediationStatistic",
        description: "Instructions - remediation statistics tab",
        groups: ["technical", "technical_admin", "technical_admin_general"]
    },
    models_exploitation_instructionStats: {
        _id: "models_instructionStats",
        description: "Instructions - instructions stats tab",
        groups: ["technical", "technical_admin", "technical_admin_general"]
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

db.permission_group.deleteMany({
    _id: {
        $in: [
            "technical_admin_maintenance",
            "technical_admin_customobjects",
            "technical_admin_settings",
        ]
    }
});
