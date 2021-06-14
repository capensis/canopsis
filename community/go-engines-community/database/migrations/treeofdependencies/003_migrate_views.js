(function () {
    // Fix AlarmsList widget's parameters
    db.views.find({"tabs.widgets.type": "AlarmsList"}).forEach(function (doc) {
        if (doc.tabs) {
            doc.tabs.forEach(function (tab) {
                if (tab.widgets) {
                    tab.widgets.forEach(function (widget) {
                        if (widget.type === "AlarmsList") {
                            widget.parameters.widgetColumns.forEach(function (column) {
                                if (column.isState) {
                                    column.colorIndicator = 'state';
                                    delete column.isState;
                                }
                            });
                        }
                    });
                }
            });

            db.views.updateOne({_id: doc._id}, {
                "$set": doc,
            });
        }
    });

    // Migrate tab rows
    db.views.find().forEach(function (doc) {
        var tabs = [];
        var tabsChanged = false;
        if (doc.tabs) {
            doc.tabs.forEach(function (tab) {
                if (!tab.widgets && tab.rows) {
                    tabsChanged = true;
                    var widgets = [];
                    var rowIndex = 0
                    tab.rows.forEach(function (row) {
                        const prevEnd = {
                            ["sm"]: 0,
                            ["md"]: 0,
                            ["lg"]: 0,
                        };

                        row.widgets.forEach(function (widget) {
                            widget.gridParameters = {
                                "mobile": {
                                    "x": prevEnd["sm"],
                                    "y": rowIndex,
                                    "w": widget.size["sm"],
                                    "h": 1,
                                    "autoHeight": true
                                },
                                "tablet": {
                                    "x": prevEnd["md"],
                                    "y": rowIndex,
                                    "w": widget.size["md"],
                                    "h": 1,
                                    "autoHeight": true
                                },
                                "desktop": {
                                    "x": prevEnd["lg"],
                                    "y": rowIndex,
                                    "w": widget.size["lg"],
                                    "h": 1,
                                    "autoHeight": true
                                }
                            }

                            prevEnd["sm"] += widget.size["sm"]
                            prevEnd["md"] += widget.size["md"]
                            prevEnd["lg"] += widget.size["lg"]

                            delete widget.size;
                            widgets.push(widget);
                        });

                        rowIndex++
                    });

                    delete tab.rows;
                    tab.widgets = widgets;
                }

                tabs.push(tab);
            });
        }
        if (tabsChanged) {
            db.views.updateOne({_id: doc._id}, {$set: {tabs: tabs}});
        }
    });
})();
