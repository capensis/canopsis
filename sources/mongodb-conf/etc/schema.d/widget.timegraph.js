{
    "type": "object",
    "categories": [{
        "title": "General",
        "keys": ["title", "human_readable", "legend", "tooltip"]
    },{
        "title": "Time Window",
        "keys": ["time_window", "time_window_offset"]
    },{
        "title": "History Navigation",
        "keys": ["timenav", "timenav_window"]
    },{
        "title": "Choose metrics",
        "keys": ["metrics"]
    },{
        "title": "Series",
        "keys": ["stacked", "points", "lines", "areas", "bars", "line_width", "bar_width"]
    }],
    "properties": {
        "title":  {
            "type": "string"
        },
        "human_readable": {
            "type": "boolean",
            "defaultValue": false,
            "label": "Human readable values"
        },
        "time_window": {
            "type": "number",
            "role": "interval",
            "defaultValue": 86400,
            "label": "Time Window to show on graph"
        },
        "time_window_offset": {
            "type": "number",
            "role": "interval",
            "defaultValue": 0,
            "label": "Time window offset from now"
        },
        "timenav": {
            "type": "boolean",
            "defaultValue": false,
            "label": "Enable history navigation"
        },
        "timenav_window": {
            "type": "number",
            "role": "interval",
            "defaultValue": 172800,
            "label": "History time window"
        },
        "stacked": {
            "type": "boolean",
            "defaultValue": false,
            "label": "Stacked graph"
        },
        "lines": {
            "type": "boolean",
            "defaultValue": true,
            "label": "Show lines"
        },
        "areas": {
            "type": "boolean",
            "defaultValue": false,
            "label": "Fill area under lines"
        },
        "points": {
            "type": "boolean",
            "defaultValue": true,
            "label": "Show points"
        },
        "bars": {
            "type": "boolean",
            "defaultValue": false,
            "label": "Show bars"
        },
        "line_width": {
            "type": "number",
            "defaultValue": 1,
            "label": "Lines width"
        },
        "bar_width": {
            "type": "number",
            "defaultValue": 10,
            "label": "Bars width"
        },
        "legend": {
            "type": "boolean",
            "defaultValue": true,
            "label": "Show legend"
        },
        "tooltip": {
            "type": "boolean",
            "defaultValue": true,
            "label": "Show tooltips"
        },
        "metrics": {
            "type": "array",
            "items": {
                "type": "string"
            },
            "role": "cmetric",
            "label": "Data source for series"
        }
    }
}