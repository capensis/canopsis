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
