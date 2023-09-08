db.color_theme.createIndex({name: 1}, {name: "name_1", unique: true});

if (!db.permission.findOne({_id: "api_color_theme"})) {
    db.permission.insertOne({
        _id: "api_color_theme",
        name: "api_color_theme",
        type: "CRUD",
        description: "Api color themes"
    });
    db.role.updateMany({_id: "admin"}, {$set: {"permissions.api_color_theme": 15}});
}

if (!db.permission.findOne({_id: "models_profile_color_theme"})) {
    db.permission.insertOne({
        _id: "models_profile_color_theme",
        name: "models_profile_color_theme",
        type: "CRUD",
        description: "Models color themes"
    });
    db.role.updateMany({_id: "admin"}, {$set: {"permissions.models_profile_color_theme": 15}});
}

var now = Math.ceil((new Date()).getTime() / 1000);

if (!db.color_theme.findOne({_id: "canopsis"})) {
    db.color_theme.insertOne({
        _id: "canopsis",
        name: "Canopsis",
        colors: {
            main: {
                primary: '#2fab63',
                secondary: '#2b3e4f',
                accent: '#82b1ff',
                error: '#ff5252',
                info: '#2196f3',
                success: '#4caf50',
                warning: '#fb8c00',
                background: '#ffffff',
                active_color: '#000',
                font_size: 2
            },
            state: {
                ok: '#00a65a',
                minor: '#fcdc00',
                major: '#ff9900',
                critical: '#f56954',
            },
            table: {
                background: '#fff',
                row_color: '#fff',
                hover_row_color: '#eee',
            }
        },
        updated: now,
        deletable: false
    });
}

if (!db.color_theme.findOne({_id: "canopsis_dark"})) {
    db.color_theme.insertOne({
        _id: "canopsis_dark",
        name: "Canopsis dark",
        colors: {
            main: {
                primary: '#2fab63',
                secondary: '#2b3e4f',
                accent: '#82b1ff',
                error: '#ff8b8b',
                info: '#2196f3',
                success: '#4caf50',
                warning: '#fb8c00',
                background: '#303030',
                active_color: '#fff',
                font_size: 2
            },
            state: {
                ok: '#00a65a',
                minor: '#fcdc00',
                major: '#ff9900',
                critical: '#f56954',
            },
            table: {
                background: '#424242',
                row_color: '#424242',
                hover_row_color: '#616161',
            }
        },
        updated: now,
        deletable: false
    });
}

if (!db.color_theme.findOne({_id: "color_blind"})) {
    db.color_theme.insertOne({
        _id: "color_blind",
        name: "Color blind",
        colors: {
            main: {
                primary: '#2196f3',
                secondary: '#2b3e4f',
                accent: '#82b1ff',
                error: '#ff5252',
                info: '#2196f3',
                success: '#4caf50',
                warning: '#fb8c00',
                background: '#ffffff',
                active_color: '#000',
                font_size: 2
            },
            state: {
                ok: '#00a65a',
                minor: '#fcdc00',
                major: '#ff9900',
                critical: '#f56954',
            },
            table: {
                background: '#fff',
                row_color: '#fff',
                hover_row_color: '#eee',
            }
        },
        updated: now,
        deletable: false,
    });
}

if (!db.color_theme.findOne({_id: "color_blind_dark"})) {
    db.color_theme.insertOne({
        _id: "color_blind_dark",
        name: "Color blind dark",
        colors: {
            main: {
                primary: '#2196f3',
                secondary: '#2b3e4f',
                accent: '#82b1ff',
                error: '#ff8b8b',
                info: '#2196f3',
                success: '#4caf50',
                warning: '#fb8c00',
                background: '#303030',
                active_color: '#fff',
                font_size: 2
            },
            state: {
                ok: '#00a65a',
                minor: '#fcdc00',
                major: '#ff9900',
                critical: '#f56954',
            },
            table: {
                background: '#424242',
                row_color: '#424242',
                hover_row_color: '#616161',
            }
        },
        updated: now,
        deletable: false,
    });
}

db.user.updateMany({"ui_theme": {$in: ["", null]}}, {$set:{"ui_theme": "canopsis"}})
db.user.updateMany({"ui_theme": "canopsisDark"}, {$set:{"ui_theme": "canopsis_dark"}})
db.user.updateMany({"ui_theme": "colorBlind"}, {$set:{"ui_theme": "color_blind"}})
db.user.updateMany({"ui_theme": "colorBlindDark"}, {$set:{"ui_theme": "color_blind_dark"}})
