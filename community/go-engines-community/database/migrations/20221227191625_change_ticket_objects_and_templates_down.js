db.periodical_alarm.find({"v.tickets": {$exists: true}}).forEach(function (doc) {
    var ticket = doc.v.ticket
    db.periodical_alarm.updateOne(
        {
            _id: doc._id
        },
        {
            $unset: {
                "v.tickets": "",
                "v.ticket.ticket": "",
                "v.ticket.ticket_data": "",
                "v.ticket.ticket_url": "",
            },
            $set: {
                "v.ticket.val": ticket.ticket,
                "v.ticket.data": ticket.ticket_data,
            }
        }
    )
});

db.resolved_alarms.find({"v.tickets": {$exists: true}}).forEach(function (doc) {
    var ticket = doc.v.ticket
    db.resolved_alarms.updateOne(
        {
            _id: doc._id
        },
        {
            $unset: {
                "v.tickets": "",
                "v.ticket.ticket": "",
                "v.ticket.ticket_data": "",
                "v.ticket.ticket_url": "",
            },
            $set: {
                "v.ticket.val": ticket.ticket,
                "v.ticket.data": ticket.ticket_data,
            }
        }
    )
});

db.widgets.updateMany(
    {
        "$and": [
            {
                "type": "AlarmsList"
            },
            {
                "$or": [
                    {"parameters.widgetColumns.value": "v.ticket.ticket"},
                    {"parameters.widgetGroupColumns.$[column].value": "v.ticket.ticket"}
                ]
            }
        ]
    },
    {
        $set: {
            "parameters.widgetColumns.$[column].value": "v.ticket.val",
            "parameters.widgetGroupColumns.$[column].value": "v.ticket.val"
        }
    },
    {
        arrayFilters: [
            {
                "column.value": "v.ticket.ticket"
            }
        ]
    }
);

db.widgets.updateMany(
    {
        "$and": [
            {
                "type": "AlarmsList"
            },
            {
                "$or": [
                    {"parameters.widgetColumns.value": "v.ticket.ticket"},
                    {"parameters.widgetGroupColumns.value": "v.ticket.ticket"},
                    {"parameters.widgetExportColumns.value": "v.ticket.ticket"}
                ]
            }
        ]
    },
    {
        $set: {
            "parameters.widgetColumns.$[column].value": "v.ticket.val",
            "parameters.widgetGroupColumns.$[column].value": "v.ticket.val",
            "parameters.widgetExportColumns.$[column].value": "v.ticket.val"
        }
    },
    {
        arrayFilters: [
            {
                "column.value": "v.ticket.ticket"
            }
        ]
    }
);

db.widgets.find({"type": "ServiceWeather"}).forEach(function (doc) {
    var set = {}

    if (typeof doc.parameters.modalTemplate === 'string') {
        set["parameters.modalTemplate"] = doc.parameters.modalTemplate.replace(/entity.ticket.ticket/g,'entity.ticket.val')
    }

    if (typeof doc.parameters.entityTemplate === 'string') {
        set["parameters.entityTemplate"] = doc.parameters.entityTemplate.replace(/entity.ticket.ticket/g,'entity.ticket.val')
    }

    set["parameters.alarmsList.widgetColumns.$[column].value"] = "v.ticket.val"

    db.widgets.updateOne(
        {
            _id: doc._id
        },
        {
            $set: set
        },
        {
            arrayFilters: [
                {
                    "column.value": "v.ticket.ticket"
                }
            ]
        }
    )
});

db.action_scenario.find({"actions.parameters.output": {$exists: true}}).forEach(function (doc) {
    var actions = doc.actions
    actions.forEach(function (action, i) {
        if (typeof action.parameters.output !== 'string') {
            return
        }

        actions[i].parameters.output = action.parameters.output.replace(/.Alarm.Value.Ticket.Ticket/g,'.Alarm.Value.Ticket.Value');
    });

    db.action_scenario.updateOne(
        {
            _id: doc._id
        },
        {
            $set: {
                "actions": actions,
            }
        }
    )
});

db.idle_rule.find({"operation.parameters.output": {$exists: true}}).forEach(function (doc) {
    var output = doc.operation.parameters.output
    if (typeof output !== 'string') {
        return
    }

    db.idle_rule.updateOne(
        {
            _id: doc._id
        },
        {
            $set: {
                "operation.parameters.output": output.replace(/.Alarm.Value.Ticket.Ticket/g,'.Alarm.Value.Ticket.Value'),
            }
        }
    )
});

db.widgets.updateMany(
    {
        "$and": [
            {
                "type": "Context"
            },
            {
                "$or": [
                    {"parameters.activeAlarmsColumns.value": "v.ticket.ticket"},
                    {"parameters.resolvedAlarmsColumns.value": "v.ticket.ticket"}
                ]
            }
        ]
    },
    {
        $set: {
            "parameters.activeAlarmsColumns.$[column].value": "v.ticket.val",
            "parameters.resolvedAlarmsColumns.$[column].value": "v.ticket.val"
        }
    },
    {
        arrayFilters: [
            {
                "column.value": "v.ticket.ticket"
            }
        ]
    }
);

db.widgets.updateMany(
    {
        "type": "Map",
        "parameters.alarms_columns.value": "v.ticket.ticket"
    },
    {
        $set: {
            "parameters.alarms_columns.$[column].value": "v.ticket.val"
        }
    },
    {
        arrayFilters: [
            {
                "column.value": "v.ticket.ticket"
            }
        ]
    }
);
