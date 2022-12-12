db.widgets.find({type: "ServiceWeather"}).forEach(function (doc) {
    var set = {};
    var unset = {};

    set["parameters.counters.pbehavior_enabled"] = false;
    set["parameters.counters.pbehavior_types"] = [];
    set["parameters.counters.state_enabled"] = false;
    set["parameters.counters.state_types"] = [];
    unset["parameters.counters.enabled"] = "";
    unset["parameters.counters.types"] = "";

    if (doc.parameters && doc.parameters.counters) {
        if (doc.parameters.counters.enabled) {
            set["parameters.counters.pbehavior_enabled"] = doc.parameters.counters.enabled;
        }
        if (doc.parameters.counters.types) {
            set["parameters.counters.pbehavior_types"] = doc.parameters.counters.types;
        }
    }

    db.widgets.updateOne({_id: doc._id}, {$set: set, $unset: unset});
});

db.widgets.find({type: {$in: ["ServiceWeather", "Counter"]}}).forEach(function (doc) {
    var set = {};
    var unset = {};

    set["parameters.columnMobile"] = 2;
    set["parameters.columnTablet"] = 3;
    set["parameters.columnDesktop"] = 4;
    unset["parameters.columnSM"] = "";
    unset["parameters.columnMD"] = "";
    unset["parameters.columnLG"] = "";

    if (doc.parameters) {
        if (doc.parameters.columnSM <= 6) {
            set["parameters.columnMobile"] = 2;
        } else {
            set["parameters.columnMobile"] = 1;
        }
        if (doc.parameters.columnMD <= 3) {
            set["parameters.columnTablet"] = 4;
        } else if (doc.parameters.columnMD === 4) {
            set["parameters.columnTablet"] = 3;
        } else if (doc.parameters.columnMD <= 6) {
            set["parameters.columnTablet"] = 2;
        } else {
            set["parameters.columnTablet"] = 1;
        }
        if (doc.parameters.columnLG <= 2) {
            set["parameters.columnDesktop"] = 6;
        } else if (doc.parameters.columnLG === 3) {
            set["parameters.columnDesktop"] = 4;
        } else if (doc.parameters.columnLG <= 4) {
            set["parameters.columnDesktop"] = 3;
        } else if (doc.parameters.columnLG <= 6) {
            set["parameters.columnDesktop"] = 2;
        } else {
            set["parameters.columnDesktop"] = 1;
        }
    }

    db.widgets.updateOne({_id: doc._id}, {$set: set, $unset: unset});
});
