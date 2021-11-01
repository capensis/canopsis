(function () {
    db.configuration.find({_id: "data_storage"}).forEach(function (doc) {
        db.configuration.updateOne({_id: doc._id}, {
            $set: {
                "config.junit.delete_after": migrateDurWithEnabled(doc.config.junit.delete_after),
                "config.remediation.accumulate_after": migrateDurWithEnabled(doc.config.junit.delete_after),
                "config.remediation.delete_after": migrateDurWithEnabled(doc.config.junit.delete_after),
                "config.alarm.archive_after": migrateDurWithEnabled(doc.config.junit.delete_after),
                "config.alarm.delete_after": migrateDurWithEnabled(doc.config.junit.delete_after),
                "config.pbehavior.delete_after": migrateDurWithEnabled(doc.config.junit.delete_after),
                "config.health_check.delete_after": migrateDurWithEnabled(doc.config.junit.delete_after),
            }
        });
    });
    db.configuration.find({_id: "user_interface"}).forEach(function (doc) {
        var popupTimeout = doc.popup_timeout;

        if (!popupTimeout) {
            return;
        }

        Object.keys(popupTimeout).forEach(function (key) {
            var value = popupTimeout[key].interval;
            delete popupTimeout[key].interval;
            popupTimeout[key].value = value;
        });

        db.configuration.updateOne({_id: doc._id}, {
            $set: {
                "popup_timeout": popupTimeout,
            }
        });
    });

    db.action_scenario.find({
        $or: [
            {"delay.seconds": {"$gt": 0}},
            {"actions.type": {$in: ["webhook", "pbehavior"]}},
        ]
    }).forEach(function (doc) {
        var set = {};

        if (doc.delay) {
            set["delay"] = migrateDurWithUnit(doc.delay);
        }

        doc.actions.forEach(function (action, actionIndex) {
            var paramPrefix = "actions." + actionIndex + ".parameters.";

            switch (action.type) {
                case "webhook":
                    if (action.parameters.retry_delay) {
                        set[paramPrefix + "retry_delay"] = migrateDurWithUnit(action.parameters.retry_delay);
                    }
                    break;
                case "snooze":
                    set[paramPrefix + "duration"] = migrateDurWithUnit(action.parameters.duration);
                    break;
                case "pbehavior":
                    if (action.parameters.duration) {
                        set[paramPrefix + "duration"] = migrateDurWithUnit(action.parameters.duration);
                    }
                    break;
            }
        });

        if (Object.keys(set).length > 0) {
            db.action_scenario.updateOne({_id: doc._id}, {
                $set: set,
            });
        }
    });

    db.flapping_rule.find().forEach(function (doc) {
        db.flapping_rule.updateOne({_id: doc._id}, {
            $set: {
                duration: migrateDurWithUnit(doc.duration),
            }
        });
    });

    db.resolve_rule.find().forEach(function (doc) {
        db.resolve_rule.updateOne({_id: doc._id}, {
            $set: {
                duration: migrateDurWithUnit(doc.duration),
            }
        });
    });

    db.idle_rule.find().forEach(function (doc) {
        var set = {
            duration: migrateDurWithUnit(doc.duration),
        };

        if (doc.operation) {
            switch (doc.operation.type) {
                case "snooze":
                    set["operation.parameters.duration"] = migrateDurWithUnit(doc.operation.parameters.duration);
                    break;
                case "pbehavior":
                    if (doc.operation.parameters.duration) {
                        set["operation.parameters.duration"] = migrateDurWithUnit(doc.operation.parameters.duration);
                    }
                    break;
            }
        }

        db.idle_rule.updateOne({_id: doc._id}, {
            $set: set,
        });
    });

    db.instruction.find().forEach(function (doc) {
        var set = {
            timeout_after_execution: migrateDurWithUnit(doc.timeout_after_execution),
        };

        if (doc.steps) {
            doc.steps.forEach(function (step, stepIndex) {
                step.operations.forEach(function (operation, operationsIndex) {
                    set["steps." + stepIndex + ".operations." + operationsIndex + ".time_to_complete"] = migrateDurWithUnit(operation.time_to_complete);
                });
            });
        }

        db.instruction.updateOne({_id: doc._id}, {$set: set});
    });

    db.notification.find().forEach(function (doc) {
        if (doc.instruction && doc.instruction.rate_frequency) {
            db.notification.updateOne({_id: doc._id}, {
                $set: {
                    "instruction.rate_frequency": migrateDurWithUnit(doc.instruction.rate_frequency),
                }
            });
        }
    });

    db.view_playlist.find().forEach(function (doc) {
        db.view_playlist.updateOne({_id: doc._id}, {
            $set: {
                interval: migrateDurWithUnit(doc.interval),
            }
        });
    });

    db.views.find().forEach(function (doc) {
        var set = {};

        if (doc.periodic_refresh) {
            set.periodic_refresh = migrateDurWithEnabled(doc.periodic_refresh);
        }

        if (doc.tabs) {
            doc.tabs.forEach(function (tab, tabIndex) {
                if (tab.widgets) {
                    tab.widgets.forEach(function (widget, widgetIndex) {
                        if (widget.parameters.periodic_refresh) {
                            var key = 'tabs.' + tabIndex + '.widgets.' + widgetIndex + '.parameters.periodic_refresh';
                            set[key] = migrateDurWithEnabled(widget.parameters.periodic_refresh);
                        }
                    });
                }
            });
        }

        if (Object.keys(set).length > 0) {
            db.views.updateOne({_id: doc._id}, {$set: set});
        }
    });

    function migrateDurWithUnit(d) {
        var value = d.seconds;
        var unit = d.unit;

        switch (unit) {
            case "s":
                break;
            case "m":
                value = Math.ceil(value / 60);
                break;
            case "h":
                value = Math.ceil(value / (60 * 60));
                break;
            case "d":
                value = Math.ceil(value / (60 * 60 * 24));
                break;
            case "w":
                value = Math.ceil(value / (60 * 60 * 24 * 7));
                break;
            case "M":
                value = Math.ceil(value / (60 * 60 * 24 * 30));
                break;
            case "y":
                value = Math.ceil(value / (60 * 60 * 24 * 365));
                break;
        }

        return {
            value: value,
            unit: unit,
        };
    }

    function migrateDurWithEnabled(d) {
        var v = migrateDurWithUnit(d);

        return {
            value: v.value,
            unit: v.unit,
            enabled: d.enabled,
        };
    }
})();
