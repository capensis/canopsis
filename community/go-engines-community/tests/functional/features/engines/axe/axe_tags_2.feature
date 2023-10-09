Feature: create and update alarm by main event stream
  I need to be able to create and update alarm on event

  @concurrent
  Scenario: given new external tag should remove internal tag with the same name
    Given I am admin
    When I do POST /api/v4/alarm-tags:
    """json
    {
      "value": "test-tag-axe-tags-second-1-1",
      "color": "#AABBCC",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-axe-tags-second-1-1",
                "test-resource-axe-tags-second-1-2",
                "test-resource-axe-tags-second-1-3"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/alarm-tags:
    """json
    {
      "value": "test-tag-axe-tags-second-1-2",
      "color": "#AABBCC",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-axe-tags-second-1-1",
                "test-resource-axe-tags-second-1-2"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-tags-second-1",
        "connector": "test-connector-axe-tags-second-1",
        "connector_name": "test-connector-name-axe-tags-second-1",
        "component": "test-component-axe-tags-second-1",
        "resource": "test-resource-axe-tags-second-1-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-tags-second-1",
        "connector": "test-connector-axe-tags-second-1",
        "connector_name": "test-connector-name-axe-tags-second-1",
        "component": "test-component-axe-tags-second-1",
        "resource": "test-resource-axe-tags-second-1-2",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-tags-second-1",
        "connector": "test-connector-axe-tags-second-1",
        "connector_name": "test-connector-name-axe-tags-second-1",
        "component": "test-component-axe-tags-second-1",
        "resource": "test-resource-axe-tags-second-1-3",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-second-1-1 until response code is 200 and response array key "data.0.tags" contains only:
    """json
    [
      "test-tag-axe-tags-second-1-1",
      "test-tag-axe-tags-second-1-2"
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-second-1-2 until response code is 200 and response array key "data.0.tags" contains only:
    """json
    [
      "test-tag-axe-tags-second-1-1",
      "test-tag-axe-tags-second-1-2"
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-second-1-3 until response code is 200 and response array key "data.0.tags" contains only:
    """json
    [
      "test-tag-axe-tags-second-1-1"
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "tags": {
        "test-tag-axe-tags-second-1-1": ""
      },
      "state": 2,
      "output": "test-output-axe-tags-second-1",
      "connector": "test-connector-axe-tags-second-1",
      "connector_name": "test-connector-name-axe-tags-second-1",
      "component": "test-component-axe-tags-second-1",
      "resource": "test-resource-axe-tags-second-1-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-second-1-3 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "tags": []
        }
      ]
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-second-1-2 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "tags": [
            "test-tag-axe-tags-second-1-2"
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-second-1-1 until response code is 200 and response array key "data.0.tags" contains only:
    """json
    [
      "test-tag-axe-tags-second-1-1",
      "test-tag-axe-tags-second-1-2"
    ]
    """
