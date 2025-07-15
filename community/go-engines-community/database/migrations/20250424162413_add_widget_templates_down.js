db.widget_templates.deleteMany({
    type: {
        $in: [
            "alarm_quick_actions",
            "alarm_mass_quick_actions"
        ]
    }
});
