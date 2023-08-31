db.widgets.find({
    "type": "AlarmsList"
}).forEach(function (widget) {
    if (widget.parameters) {
        var updated = false,
            wp = widget.parameters;
        if (wp.moreInfoTemplate) {
            var newTemplate = wp.moreInfoTemplate.replaceAll(/alarm\.v\.ticket\.data\./g, "alarm.v.ticket_data.")
            updated = newTemplate !== wp.moreInfoTemplate;
            if (updated) {
                wp.moreInfoTemplate = newTemplate;
            }
        }
        for (var columnSet of ["widgetColumns", "widgetGroupColumns", "widgetExportColumns"]) {
            if (wp[columnSet]) {
                wp[columnSet].forEach(function (column, n) {
                    if (column.value.indexOf("v.ticket.data.") === 0) {
                        column.value = column.value.replace("v.ticket.data.", "v.ticket_data.");
                        wp[columnSet][n] = column;
                        updated = true;
                    }
                });
            }
        }    
        if (updated) {
            db.widgets.updateOne({_id: widget._id}, {$set: {"parameters": wp}});
        }
    }
})