Feature: create and update alarm by main event stream
  I need to be able to create and update alarm on event

  Scenario: given a new alarm with a new connector should update an alarm
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-change-connector-1",
      "connector_name": "test-connector-name-axe-change-connector-1-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-change-connector-1",
      "resource": "test-resource-axe-change-connector-1",
      "state": 2,
      "output": "test-output-axe-change-connector-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-change-connector-1&opened=true&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-axe-change-connector-1/test-component-axe-change-connector-1"
          },
          "v": {
            "component": "test-component-axe-change-connector-1",
            "connector": "test-connector-axe-change-connector-1",
            "connector_name": "test-connector-name-axe-change-connector-1-1",
            "resource": "test-resource-axe-change-connector-1",
            "steps": [
              {
                "_t": "stateinc",
                "a": "test-connector-axe-change-connector-1.test-connector-name-axe-change-connector-1-1",
                "m": "test-output-axe-change-connector-1",
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "test-connector-axe-change-connector-1.test-connector-name-axe-change-connector-1-1",
                "m": "test-output-axe-change-connector-1",
                "val": 1
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
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-change-connector-1",
      "connector_name": "test-connector-name-axe-change-connector-1-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-change-connector-1",
      "resource": "test-resource-axe-change-connector-1",
      "state": 2,
      "output": "test-output-axe-change-connector-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-change-connector-1&opened=true&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-axe-change-connector-1/test-component-axe-change-connector-1"
          },
          "v": {
            "component": "test-component-axe-change-connector-1",
            "connector": "test-connector-axe-change-connector-1",
            "connector_name": "test-connector-name-axe-change-connector-1-1",
            "resource": "test-resource-axe-change-connector-1",
            "steps": [
              {
                "_t": "stateinc",
                "a": "test-connector-axe-change-connector-1.test-connector-name-axe-change-connector-1-1",
                "m": "test-output-axe-change-connector-1",
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "test-connector-axe-change-connector-1.test-connector-name-axe-change-connector-1-1",
                "m": "test-output-axe-change-connector-1",
                "val": 1
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
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-change-connector-1",
      "connector_name": "test-connector-name-axe-change-connector-1-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-change-connector-1",
      "resource": "test-resource-axe-change-connector-1",
      "state": 3,
      "output": "test-output-axe-change-connector-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-change-connector-1&opened=true&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-axe-change-connector-1/test-component-axe-change-connector-1"
          },
          "v": {
            "component": "test-component-axe-change-connector-1",
            "connector": "test-connector-axe-change-connector-1",
            "connector_name": "test-connector-name-axe-change-connector-1-1",
            "resource": "test-resource-axe-change-connector-1",
            "steps": [
              {
                "_t": "stateinc",
                "a": "test-connector-axe-change-connector-1.test-connector-name-axe-change-connector-1-1",
                "m": "test-output-axe-change-connector-1",
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "test-connector-axe-change-connector-1.test-connector-name-axe-change-connector-1-1",
                "m": "test-output-axe-change-connector-1",
                "val": 1
              },
              {
                "_t": "stateinc",
                "a": "test-connector-axe-change-connector-1.test-connector-name-axe-change-connector-1-1",
                "m": "test-output-axe-change-connector-1",
                "val": 3
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
