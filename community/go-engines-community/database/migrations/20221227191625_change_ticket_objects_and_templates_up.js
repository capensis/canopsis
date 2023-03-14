db.periodical_alarm.find({"v.ticket": {$exists: true}}).forEach(function (doc) {
    var ticket = doc.v.ticket
    var ticketNumber = ticket.val
    var ticketData = ticket.data

    var set = {
        "v.ticket.val": 0,
        "v.ticket.ticket": ticketNumber,
        "v.ticket.ticket_data": ticketData,
    }

    ticket.val = 0
    ticket.ticket = ticketNumber
    ticket.ticket_data = ticketData
    if (ticketData !== null && ticketData !== undefined && ticketData.url !== null && ticketData.url !== undefined) {
        set["v.ticket.ticket_url"] = ticketData.url
        ticket.ticket_url = ticketData.url
    }

    delete ticket.data

    db.periodical_alarm.updateOne(
        {
            _id: doc._id
        },
        {
            $push: {
                "v.tickets": ticket
            },
            $set: set,
            $unset: {
                "v.ticket.data": ""
            }
        }
    )
});

db.resolved_alarms.find({"v.ticket": {$exists: true}}).forEach(function (doc) {
    var ticket = doc.v.ticket
    var ticketNumber = ticket.val
    var ticketData = ticket.data

    var set = {
        "v.ticket.val": 0,
        "v.ticket.ticket": ticketNumber,
        "v.ticket.ticket_data": ticketData,
    }

    ticket.val = 0
    ticket.ticket = ticketNumber
    ticket.ticket_data = ticketData
    if (ticketData !== null && ticketData !== undefined && ticketData.url !== null && ticketData.url !== undefined) {
        set["v.ticket.ticket_url"] = ticketData.url
        ticket.ticket_url = ticketData.url
    }

    delete ticket.data

    db.resolved_alarms.updateOne(
        {
            _id: doc._id
        },
        {
            $push: {
                "v.tickets": ticket
            },
            $set: set,
            $unset: {
                "v.ticket.data": ""
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
                    {"parameters.widgetColumns.value": "v.ticket.val"},
                    {"parameters.widgetGroupColumns.value": "v.ticket.val"},
                    {"parameters.widgetExportColumns.value": "v.ticket.val"}
                ]
            }
        ]
    },
    {
        $set: {
            "parameters.widgetColumns.$[column].value": "v.ticket.ticket",
            "parameters.widgetGroupColumns.$[column].value": "v.ticket.ticket",
            "parameters.widgetExportColumns.$[column].value": "v.ticket.ticket"
        }
    },
    {
        arrayFilters: [
            {
                "column.value": "v.ticket.val"
            }
        ]
    }
);

db.widgets.find({"type": "ServiceWeather"}).forEach(function (doc) {
    var set = {}

    if (typeof doc.parameters.modalTemplate === 'string') {
        set["parameters.modalTemplate"] = doc.parameters.modalTemplate.replace(/entity.ticket.val/g,'entity.ticket.ticket')
    }

    if (typeof doc.parameters.entityTemplate === 'string') {
        set["parameters.entityTemplate"] = doc.parameters.entityTemplate.replace(/entity.ticket.val/g,'entity.ticket.ticket')
    }

    set["parameters.alarmsList.widgetColumns.$[column].value"] = "v.ticket.ticket"

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
                    "column.value": "v.ticket.val"
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

        actions[i].parameters.output = action.parameters.output.replace(/.Alarm.Value.Ticket.Value/g,'.Alarm.Value.Ticket.Ticket');
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
                "operation.parameters.output": output.replace(/.Alarm.Value.Ticket.Value/g,'.Alarm.Value.Ticket.Ticket'),
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
                    {"parameters.activeAlarmsColumns.value": "v.ticket.val"},
                    {"parameters.resolvedAlarmsColumns.value": "v.ticket.val"}
                ]
            }
        ]
    },
    {
        $set: {
            "parameters.activeAlarmsColumns.$[column].value": "v.ticket.ticket",
            "parameters.resolvedAlarmsColumns.$[column].value": "v.ticket.ticket"
        }
    },
    {
        arrayFilters: [
            {
                "column.value": "v.ticket.val"
            }
        ]
    }
);

db.widgets.updateMany(
    {
        "type": "Map",
        "parameters.alarms_columns.value": "v.ticket.val"
    },
    {
        $set: {
            "parameters.alarms_columns.$[column].value": "v.ticket.ticket"
        }
    },
    {
        arrayFilters: [
            {
                "column.value": "v.ticket.val"
            }
        ]
    }
);
