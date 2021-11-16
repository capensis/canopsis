Feature: unsnooze alarm
  I need to be able to unsnooze alarm
  
  Scenario: given unsnooze event should update alarm
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-unsnooze-1",
      "connector_name" : "test-connector-name-axe-unsnooze-1",
      "source_type" : "resource",
      "component" :  "test-component-axe-unsnooze-1",
      "resource" : "test-resource-axe-unsnooze-1",
      "state" : 2,
      "output" : "test-output-axe-unsnooze-1",
      "long_output" : "test-long-output-axe-unsnooze-1",
      "author" : "test-author-axe-unsnooze-1",
      "timestamp": {{ (now.Add (parseDuration "-10s")).UTC.Unix }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "snooze",
      "duration": 3,
      "connector" : "test-connector-axe-unsnooze-1",
      "connector_name" : "test-connector-name-axe-unsnooze-1",
      "source_type" : "resource",
      "component" :  "test-component-axe-unsnooze-1",
      "resource" : "test-resource-axe-unsnooze-1",
      "output" : "test-output-axe-unsnooze-1",
      "long_output" : "test-long-output-axe-unsnooze-1",
      "author" : "test-author-axe-unsnooze-1",
      "timestamp": {{ (now.Add (parseDuration "-5s")).UTC.Unix }}
    }
    """
    When I save response snoozeEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-axe-unsnooze-1"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-unsnooze-1",
            "connector": "test-connector-axe-unsnooze-1",
            "connector_name": "test-connector-name-axe-unsnooze-1",
            "last_update_date": {{ .checkEventTimestamp }},
            "resource": "test-resource-axe-unsnooze-1",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "snooze",
                "a": "test-author-axe-unsnooze-1",
                "m": "test-output-axe-unsnooze-1",
                "t": {{ .snoozeEventTimestamp }},
                "val": {{ .snoozeEventTimestamp | sum 3 }}
              }
            ]
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    Then the response key "data.0.v.snooze" should not exist

  Scenario: given unsnooze event should update alarm snooze duration
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-unsnooze-2",
      "connector_name" : "test-connector-name-axe-unsnooze-2",
      "source_type" : "resource",
      "component" :  "test-component-axe-unsnooze-2",
      "resource" : "test-resource-axe-unsnooze-2",
      "state" : 2,
      "output" : "test-output-axe-unsnooze-2"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "snooze",
      "duration": 3,
      "connector" : "test-connector-axe-unsnooze-2",
      "connector_name" : "test-connector-name-axe-unsnooze-2",
      "source_type" : "resource",
      "component" :  "test-component-axe-unsnooze-2",
      "resource" : "test-resource-axe-unsnooze-2",
      "output" : "test-output-axe-unsnooze-2"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-axe-unsnooze-2"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "snooze_duration": 0
          }
        }
      ]
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-axe-unsnooze-2"}]}&with_steps=true
    Then the response code should be 200
    Then the response key "data.0.v.snooze_duration" should not be "0"
