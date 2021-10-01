(function () {
    var deprecatedFilterKeyMap = {
        "entity_id": "_id",
        "display_name": "name",
        "tileIcon": "icon",
        "tileSecondaryIcon": "secondary_icon",
        "watcher_pbehavior": "pbehaviors",
        "isActionRequired": "has_open_alarm",
        "isAllEntitiesPaused": "is_grey",
        "isWatcherPaused": "is_grey",
    };

    db.views.find().forEach(function (doc) {
        if (!doc.tabs || doc.tabs.length === 0) {
            return;
        }

        var set = {};
        var unset = {};
        doc.tabs.forEach(function (tab, tabIndex) {
            if (!tab.widgets || tab.widgets.length === 0) {
                return;
            }

            tab.widgets.forEach(function (widget, widgetIndex) {
                if (widget.type !== "ServiceWeather") {
                    return;
                }

                if (widget.parameters.mfilter) {
                    var key = "tabs." + tabIndex + ".widgets." + widgetIndex + ".parameters";
                    if (widget.parameters.mfilter.filter) {
                        var filter = JSON.parse(widget.parameters.mfilter.filter);
                        var newFilter = replaceFilters(filter);
                        set[key + ".mainFilter"] = {
                            title: "Default",
                            filter: JSON.stringify(newFilter),
                        };
                        set[key + ".viewFilters"] = [set[key + ".mainFilter"]];
                        set[key + ".mainFilterUpdatedAt"] = 1;
                    }

                    unset[key + ".mfilter"] = ""
                }
                ["blockTemplate", "entityTemplate", "modalTemplate"].forEach(function (templateItem, idx) {
                    var newTemplate = widget.parameters[templateItem];
                    if (newTemplate) {
                        for(var oldName in deprecatedFilterKeyMap) {
                            newTemplate = newTemplate.replace(oldName, deprecatedFilterKeyMap[oldName]);
                        }
                        if (newTemplate !== widget.parameters[templateItem]) {
                            set[key + "." + templateItem] = newTemplate;
                        }
                    }
                })
            });
        });

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

    function replaceFilters(filter) {
        if (filter && Array.isArray(filter)) {
            var newFilter = [];
            filter.forEach(function (item) {
                newFilter.push(replaceFilters(item));
            })

            return newFilter;
        }

        if (filter && typeof filter === "object") {
            var newFilter = {};

            Object.keys(filter).forEach(function (key) {
                var val = filter[key];
                var newKey, newVal;

                if (key === "color" || key === "tileColor") {
                    if (val === "pause") {
                        newKey = "is_grey";
                        newVal = true;
                    } else {
                        newKey = "impact_state";
                        newVal = replaceColorValue(val);
                    }

                    newFilter[newKey] = newVal;
                    return;
                }

                newKey = deprecatedFilterKeyMap[key];
                if (newKey) {
                    newFilter[newKey] = val;
                    return;
                }

                newFilter[key] = replaceFilters(val);
            });

            return newFilter;
        }

        return filter;
    }

    function replaceColorValue(val) {
        var colorValuesMap = {
            "ok": 0,
            "minor": 1,
            "major": 2,
            "critical": 3,
        };
        var newVal;

        if (val && Array.isArray(val)) {
            newVal = [];
            val.forEach(function (item) {
                newVal.push(replaceColorValue(item));
            })

            return newVal;
        }

        if (val && typeof val === "object") {
            newVal = {};
            Object.keys(val).forEach(function (key) {
                newVal[key] = replaceColorValue(val[key]);
            })

            return newVal;
        }

        return colorValuesMap[val];
    }
})();
