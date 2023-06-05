Feature: unsnooze alarm
  I need to be able to unsnooze alarm

  @concurrent
  Scenario: given unsnooze event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
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
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
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
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response snoozeEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type" : "unsnooze",
      "connector" : "test-connector-axe-unsnooze-1",
      "connector_name" : "test-connector-name-axe-unsnooze-1",
      "source_type" : "resource",
      "component" :  "test-component-axe-unsnooze-1",
      "resource" : "test-resource-axe-unsnooze-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-unsnooze-1
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
            }
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
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
                "val": {{ .snoozeEventTimestamp | sumTime 3 }}
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      }
    ]
    """
