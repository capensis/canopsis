(function () {
    db.views.find({
        $or: [
            {"tabs.widgets.type": "AlarmsList"},
            {"tabs.widgets.type": "StatsCalendar"},
            {"tabs.widgets.type": "Counter"},
        ]
    }).forEach(function (doc) {
        if (doc.tabs) {
            doc.tabs.forEach(function (tab) {
                if (tab.widgets) {
                    tab.widgets.forEach(function (widget) {
                        if ((widget.type === "AlarmsList" || widget.type === "StatsCalendar" || widget.type === "Counter") &&
                            widget.parameters.alarmsStateFilter) {
                            var opened = null;
                            var alarmsStateFilter = widget.parameters.alarmsStateFilter;
                            if (alarmsStateFilter.opened && !alarmsStateFilter.resolved) {
                                opened = true;
                            } else if (!alarmsStateFilter.opened && alarmsStateFilter.resolved) {
                                opened = false;
                            }

                            delete widget.parameters.alarmsStateFilter;
                            widget.parameters.opened = opened;
                        }
                    });
                }
            });

            db.views.updateOne({_id: doc._id}, {
                "$set": doc,
            });
        }
    });
})();
