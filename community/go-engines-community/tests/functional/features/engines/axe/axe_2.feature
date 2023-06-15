Feature: create and update alarm by main event stream
  I need to be able to create and update alarm on event

  @concurrent
  Scenario: given events with different connector for the same resource should use right connectors in alarm steps
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-second-1",
      "connector_name": "test-connector-name-axe-second-1",
      "source_type": "resource",
      "component":  "test-component-axe-second-1",
      "resource": "test-resource-axe-second-1",
      "state": 1,
      "output": "test-output-axe-second-1",
      "long_output": "test-long-output-axe-second-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-second-1",
      "connector_name": "test-connector-name-axe-second-1-new",
      "source_type": "resource",
      "component":  "test-component-axe-second-1",
      "resource": "test-resource-axe-second-1",
      "state": 2,
      "output": "test-output-axe-second-1",
      "long_output": "test-long-output-axe-second-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-second-1-new",
      "connector_name": "test-connector-name-axe-second-1",
      "source_type": "resource",
      "component":  "test-component-axe-second-1",
      "resource": "test-resource-axe-second-1",
      "state": 3,
      "output": "test-output-axe-second-1",
      "long_output": "test-long-output-axe-second-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-second-1-new",
      "connector_name": "test-connector-name-axe-second-1-new",
      "source_type": "resource",
      "component":  "test-component-axe-second-1",
      "resource": "test-resource-axe-second-1",
      "state": 0,
      "output": "test-output-axe-second-1",
      "long_output": "test-long-output-axe-second-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-second-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-second-1",
            "connector": "test-connector-axe-second-1",
            "connector_name": "test-connector-name-axe-second-1",
            "resource": "test-resource-axe-second-1"
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
                "val": 1,
                "a": "test-connector-axe-second-1.test-connector-name-axe-second-1"
              },
              {
                "_t": "statusinc",
                "val": 1,
                "a": "test-connector-axe-second-1.test-connector-name-axe-second-1"
              },
              {
                "_t": "stateinc",
                "val": 2,
                "a": "test-connector-axe-second-1.test-connector-name-axe-second-1-new"
              },
              {
                "_t": "stateinc",
                "val": 3,
                "a": "test-connector-axe-second-1-new.test-connector-name-axe-second-1"
              },
              {
                "_t": "statedec",
                "val": 0,
                "a": "test-connector-axe-second-1-new.test-connector-name-axe-second-1-new"
              },
              {
                "_t": "statusdec",
                "val": 0,
                "a": "test-connector-axe-second-1-new.test-connector-name-axe-second-1-new"
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 6
            }
          }
        }
      }
    ]
    """
