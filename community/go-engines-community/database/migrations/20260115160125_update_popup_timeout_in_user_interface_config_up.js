function toSeconds(d) {
    let value = 0;
    switch (d.unit) {
        case 'h':
            value = d.value * 60 * 60;
            break;
        case 'm':
            value = d.value * 60;
            break;
    }

    if (value === 0) {
        // replace invalid duration with default value
        value = 3;
    }

    return {
        value: value,
        unit: 's',
    };
}

const conf = db.configuration.findOne({"_id": "user_interface"});
if (conf && conf.popup_timeout) {
    let set = {};
    if (conf.popup_timeout.info && conf.popup_timeout.info.unit !== 's') {
        set['popup_timeout.info'] = toSeconds(conf.popup_timeout.info);
    }

    if (conf.popup_timeout.error && conf.popup_timeout.error.unit !== 's') {
        set['popup_timeout.error'] = toSeconds(conf.popup_timeout.error);
    }

    if (Object.keys(set).length > 0) {
        db.configuration.updateOne({"_id": "user_interface"}, {$set: set});
    }
}
