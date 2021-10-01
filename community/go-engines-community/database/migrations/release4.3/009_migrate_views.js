(function () {
    var now = Math.ceil((new Date()).getTime() / 1000);

    db.views.find().forEach(function (doc) {
        var set = {};
        var unset = {};

        // Migrate periodical refresh
        if (doc.periodicRefresh) {
            var enabled = doc.periodicRefresh.enabled;
            var unit = doc.periodicRefresh.unit;
            var seconds = parseInt(doc.periodicRefresh.interval, 10);
            if (isNaN(seconds)) {
                seconds = 0;
            }
            if (enabled === undefined) {
                enabled = false;
            }
            if (unit === "m") {
                seconds *= 60;
            } else if (unit === "h") {
                seconds *= 60 * 60;
            }

            set["periodic_refresh"] = {
                enabled: enabled,
                seconds: seconds,
                unit: unit,
            };
            unset["periodicRefresh"] = "";
        }

        if (doc.tabs) {
            doc.tabs.forEach(function (tab, tabIndex) {
                if (tab.widgets) {
                    tab.widgets.forEach(function (widget, widgetIndex) {
                        if (widget.parameters.periodicRefresh) {
                            var key = 'tabs.' + tabIndex + '.widgets.' + widgetIndex + '.parameters.';
                            var enabled = widget.parameters.periodicRefresh.enabled;
                            var unit = widget.parameters.periodicRefresh.unit;
                            var seconds = parseInt(widget.parameters.periodicRefresh.interval, 10);
                            if (isNaN(seconds)) {
                                seconds = 0;
                            }
                            if (enabled === undefined) {
                                enabled = false;
                            }
                            if (unit === "m") {
                                seconds *= 60;
                            } else if (unit === "h") {
                                seconds *= 60 * 60;
                            }

                            set[key + 'periodic_refresh'] = {
                                enabled: enabled,
                                seconds: seconds,
                                unit: unit,
                            };
                            unset[key + 'periodicRefresh'] = "";
                        }

                        if (widget.gridParameters) {
                            set['tabs.' + tabIndex + '.widgets.' + widgetIndex + ".grid_parameters"] = widget.gridParameters;
                            unset['tabs.' + tabIndex + '.widgets.' + widgetIndex + ".gridParameters"] = "";
                        }
                    });
                }
            });
        }

        if (!doc.title || doc.title === "") {
            set["title"] = doc.name;
            unset["name"] = "";
        }

        if (!doc.created) {
            set["created"] = now;
            set["updated"] = now;
        }
        if (!doc.author) {
            set["author"] = "root";
        }

        var update = {};
        if (Object.keys(set).length > 0) {
            update["$set"] = set;
        }
        if (Object.keys(unset).length > 0) {
            update["$unset"] = unset;
        }

        if (Object.keys(update).length > 0) {
            db.views.updateOne({_id: doc._id}, update);
        }
    });

    db.viewgroups.find().forEach(function (doc) {
        var set = {};
        var unset = {};

        if (!doc.title || doc.title === "") {
            set["title"] = doc.name;
            unset["name"] = "";
        }

        if (!doc.created) {
            set["created"] = now;
            set["updated"] = now;
        }
        if (!doc.author) {
            set["author"] = "root";
        }

        var update = {};
        if (Object.keys(set).length > 0) {
            update["$set"] = set;
        }
        if (Object.keys(unset).length > 0) {
            update["$unset"] = unset;
        }

        if (Object.keys(update).length > 0) {
            db.viewgroups.updateOne({_id: doc._id}, update);
        }
    });
})();
