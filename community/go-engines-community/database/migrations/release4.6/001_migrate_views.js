(function () {
    var tabIds = {};
    var widgetIds = {};
    db.views.find().forEach(function (doc) {
        if (!doc.tabs) {
            return;
        }

        var tabs = [];
        var widgets = [];

        var now = Math.ceil((new Date()).getTime() / 1000);
        var created = doc.created;
        if (!created) {
            created = now;
        }
        var updated = doc.updated;
        if (!updated) {
            updated = now;
        }
        var author = doc.author;
        if (!author) {
            author = "";
        }

        doc.tabs.forEach(function (tab, tabIndex) {
            var tabId = tab._id;
            if (tabIds[tabId] > 0) {
                tabIds[tabId]++;
                tabId += "_" + tabIds[tabId];
            } else {
                tabIds[tabId] = 1;
            }

            if (tab.widgets) {
                tab.widgets.forEach(function (widget) {
                    var widgetId = widget._id;
                    if (widgetIds[widgetId] > 0) {
                        widgetIds[widgetId]++;
                        widgetId += "_" + widgetIds[widgetId];
                    } else {
                        widgetIds[widgetId] = 1;
                    }

                    if (widget.parameters.periodicRefresh) {
                        widget.parameters.periodic_refresh = widget.parameters.periodicRefresh;
                        delete widget.parameters.periodicRefresh;
                    }

                    widget._id = widgetId;
                    widget.loader_id = widgetId;
                    widget.tab = tabId;
                    widget.author = author;
                    widget.created = created;
                    widget.updated = updated;
                    widgets.push(widget);
                });
            }

            delete tab.widgets;
            tab._id = tabId;
            tab.loader_id = tabId;
            tab.view = doc._id;
            tab.position = tabIndex;
            tab.author = author;
            tab.created = created;
            tab.updated = updated;
            tabs.push(tab);
        });

        db.views.updateOne({_id: doc._id}, {$unset: {tabs: ""}});
        db.viewtabs.insertMany(tabs);
        if (widgets.length > 0) {
            db.widgets.insertMany(widgets);
        }
    });
})();
