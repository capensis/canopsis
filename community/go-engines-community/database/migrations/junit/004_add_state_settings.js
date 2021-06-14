(function () {
    db.state_settings.insertMany([
        {
            "_id": "junit",
            "type": "junit",
            "method": "worst_of_share",
            "junit_thresholds": {
                "skipped": {
                    "minor": 10,
                    "major": 20,
                    "critical": 30,
                    "type": 1
                },
                "errors": {
                    "minor": 10,
                    "major": 20,
                    "critical": 30,
                    "type": 1
                },
                "failures": {
                    "minor": 10,
                    "major": 20,
                    "critical": 30,
                    "type": 1
                }
            }
        }
    ]);
})();
