db.widgets.find({type: "ServiceWeather"}).forEach(function (doc) {
    var set = {};
    var unset = {};

    unset["parameters.counters.pbehavior_enabled"] = "";
    unset["parameters.counters.pbehavior_types"] = "";
    unset["parameters.counters.state_enabled"] = "";
    unset["parameters.counters.state_types"] = "";
    set["parameters.counters.enabled"] = false;
    set["parameters.counters.types"] = [];

    if (doc.parameters && doc.parameters.counters) {
        if (doc.parameters.counters.pbehavior_enabled) {
            set["parameters.counters.enabled"] = doc.parameters.counters.pbehavior_enabled;
        }
        if (doc.parameters.counters.pbehavior_types) {
            set["parameters.counters.types"] = doc.parameters.counters.pbehavior_types;
        }
    }

    db.widgets.updateOne({_id: doc._id}, {$set: set, $unset: unset});
});

db.widgets.find({type: {$in: ["ServiceWeather", "Counter"]}}).forEach(function (doc) {
    var set = {};
    var unset = {};

    set["parameters.columnSM"] = 6;
    set["parameters.columnMD"] = 4;
    set["parameters.columnLG"] = 3;
    unset["parameters.columnMobile"] = "";
    unset["parameters.columnTablet"] = "";
    unset["parameters.columnDesktop"] = "";

    if (doc.parameters) {
        if (doc.parameters.columnMobile === 2) {
            set["parameters.columnSM"] = 6;
        } else {
            set["parameters.columnSM"] = 12;
        }
        if (doc.parameters.columnTablet === 4) {
            set["parameters.columnMD"] = 3;
        } else if (doc.parameters.columnTablet === 3) {
            set["parameters.columnMD"] = 4;
        } else if (doc.parameters.columnTablet === 2) {
            set["parameters.columnMD"] = 6;
        } else {
            set["parameters.columnMD"] = 12;
        }
        if (doc.parameters.columnDesktop === 6) {
            set["parameters.columnLG"] = 2;
        } else if (doc.parameters.columnDesktop === 4) {
            set["parameters.columnLG"] = 3;
        } else if (doc.parameters.columnDesktop === 3) {
            set["parameters.columnLG"] = 4;
        } else if (doc.parameters.columnDesktop === 2) {
            set["parameters.columnLG"] = 6;
        } else {
            set["parameters.columnLG"] = 12;
        }
    }

    db.widgets.updateOne({_id: doc._id}, {$set: set, $unset: unset});
});
