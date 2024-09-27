function genID() {
    var hex = UUID().toString('hex');

    return hex.replace(/^(.{8})(.{4})(.{4})(.{4})(.{12})$/, "$1-$2-$3-$4-$5")
}

function isInt(value) {
    return typeof value === "number" || value instanceof Long || value instanceof Int32;
}

function toInt(value) {
    return Long(value);
}
