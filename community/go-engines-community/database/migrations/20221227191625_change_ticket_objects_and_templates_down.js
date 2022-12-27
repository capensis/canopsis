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
                "v.ticket.ticket_data": ""
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
        "type": "AlarmsList",
        "parameters.widgetColumns.value": "v.ticket.ticket"
    },
    {
        $set: {
            "parameters.widgetColumns.$[column].value": "v.ticket.val"
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
        "type": "AlarmsList",
        "parameters.widgetColumns.value": "v.ticket.ticket"
    },
    {
        $set: {
            "parameters.widgetColumns.$[column].value": "v.ticket.val"
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

    if (typeof doc.parameters.blockTemplate === 'string') {
        set["parameters.blockTemplate"] = doc.parameters.blockTemplate.replace(/entity.ticket.ticket/g,'entity.ticket.val')
    }

    if (set === {}) {
        return
    }

    db.widgets.updateOne(
        {
            _id: doc._id
        },
        {
            $set: set
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
                "operation.parameters.output": output.replace(/{{.Alarm.Value.Ticket.Ticket}}/g,'{{.Alarm.Value.Ticket.Value}}'),
            }
        }
    )
});
