db.configuration.updateOne({_id: "user_interface"}, {$unset: {default_color_theme: ""}})
