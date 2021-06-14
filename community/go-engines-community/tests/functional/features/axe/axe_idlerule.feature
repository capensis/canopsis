Feature: update alarm on idle rule
  I need to be able to update alarm on idle rule

  Scenario: given idle rule and no events for alarm should update alarm
    Given I am admin
    When I do POST /api/v2/idle-rule:
    """
    {
      "name": "test-idlerule-axe-idlerule-name-1",
      "description": "test-idlerule-axe-idlerule-description-1",
      "type": "last_event",
      "duration": "3s",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-component-axe-idlerule-1",
            "snooze": null
          }
        }
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-output-axe-idlerule-1",
          "duration": {
            "seconds": 600,
            "unit": "m"
          },
          "author":"test-author-axe-idlerule-1"
        }
      },
      "author":"test-author-axe-idlerule-1"
    }
    """
    Then the response code should be 200
    When I wait the next periodical process
    When I send an event:
    """
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-idlerule-1",
      "connector_name" : "test-connector-name-axe-idlerule-1",
      "source_type" : "resource",
      "component" :  "test-component-axe-idlerule-1",
      "resource" : "test-resource-axe-idlerule-1",
      "state" : 2,
      "output" : "test-output-axe-idlerule-1",
      "long_output" : "test-long-output-axe-idlerule-1",
      "author" : "test-author-axe-idlerule-1",
      "timestamp": {{ (now.Add (parseDuration "-5s")).UTC.Unix }}
    }
    """
    When I wait the end of event processing
    When I wait 3s
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-axe-idlerule-1"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-1",
            "connector": "test-connector-axe-idlerule-1",
            "connector_name": "test-connector-name-axe-idlerule-1",
            "resource": "test-resource-axe-idlerule-1",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "snooze": {
              "_t": "snooze",
              "a": "test-author-axe-idlerule-1",
              "m": "test-output-axe-idlerule-1"
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
                "a": "test-author-axe-idlerule-1",
                "m": "test-output-axe-idlerule-1"
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

  Scenario: given idle rule and no update for alarm should update alarm
    Given I am admin
    When I do POST /api/v2/idle-rule:
    """
    {
      "name": "test-idlerule-axe-idlerule-name-2",
      "description": "test-idlerule-axe-idlerule-description-2",
      "type": "last_update",
      "duration": "5s",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-component-axe-idlerule-2",
            "snooze": null
          }
        }
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-output-axe-idlerule-2",
          "duration": {
            "seconds": 600,
            "unit": "m"
          },
          "author":"test-author-axe-idlerule-2"
        }
      },
      "author":"test-author-axe-idlerule-2"
    }
    """
    Then the response code should be 200
    When I wait the next periodical process
    When I send an event:
    """
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-idlerule-2",
      "connector_name" : "test-connector-name-axe-idlerule-2",
      "source_type" : "resource",
      "component" :  "test-component-axe-idlerule-2",
      "resource" : "test-resource-axe-idlerule-2",
      "state" : 2,
      "output" : "test-output-axe-idlerule-2",
      "long_output" : "test-long-output-axe-idlerule-2",
      "author" : "test-author-axe-idlerule-2",
      "timestamp": {{ (now.Add (parseDuration "-5s")).UTC.Unix }}
    }
    """
    When I wait the end of event processing
    When I wait 3s
    When I send an event:
    """
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-idlerule-2",
      "connector_name" : "test-connector-name-axe-idlerule-2",
      "source_type" : "resource",
      "component" :  "test-component-axe-idlerule-2",
      "resource" : "test-resource-axe-idlerule-2",
      "state" : 2,
      "output" : "test-output-axe-idlerule-2",
      "long_output" : "test-long-output-axe-idlerule-2",
      "author" : "test-author-axe-idlerule-2",
      "timestamp": {{ (now.Add (parseDuration "-2s")).UTC.Unix }}
    }
    """
    When I wait the end of event processing
    When I wait 2s
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-axe-idlerule-2"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-2",
            "connector": "test-connector-axe-idlerule-2",
            "connector_name": "test-connector-name-axe-idlerule-2",
            "resource": "test-resource-axe-idlerule-2",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "snooze": {
              "_t": "snooze",
              "a": "test-author-axe-idlerule-2",
              "m": "test-output-axe-idlerule-2"
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
                "a": "test-author-axe-idlerule-2",
                "m": "test-output-axe-idlerule-2"
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
