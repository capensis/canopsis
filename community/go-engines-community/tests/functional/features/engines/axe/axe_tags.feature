Feature: create and update alarm by main event stream
  I need to be able to create and update alarm on event

  Scenario: given check event should create alarm with tags
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-tags-1",
      "connector_name": "test-connector-name-axe-tags-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-tags-1",
      "resource": "test-resource-axe-tags-1",
      "state": 2,
      "output": "test-output-axe-tags-1",
      "tags": {
        "env": "prod",
        "cloud": "",
        "": "localhost"
      }
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-1
    Then the response code should be 200
    Then the response array key "data.0.tags" should contain only:
    """json
    [
      "env: prod",
      "cloud"
    ]
    """

  Scenario: given check event with another tags should add new tags
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-tags-2",
      "connector_name": "test-connector-name-axe-tags-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-tags-2",
      "resource": "test-resource-axe-tags-2",
      "state": 2,
      "output": "test-output-axe-tags-2",
      "tags": {
        "env": "prod",
        "cloud": "",
        "": "localhost"
      }
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-tags-2",
      "connector_name": "test-connector-name-axe-tags-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-tags-2",
      "resource": "test-resource-axe-tags-2",
      "state": 2,
      "output": "test-output-axe-tags-2",
      "tags": {
        "env": "test",
        "website": "",
        "cloud": "",
        "": ""
      }
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-2
    Then the response code should be 200
    Then the response array key "data.0.tags" should contain only:
    """json
    [
      "env: prod",
      "cloud",
      "env: test",
      "website"
    ]
    """

  Scenario: given check event without tags should keep old tags
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-tags-3",
      "connector_name": "test-connector-name-axe-tags-3",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-tags-3",
      "resource": "test-resource-axe-tags-3",
      "state": 2,
      "output": "test-output-axe-tags-3",
      "tags": {
        " env ": " prod ",
        "cloud": "  ",
        " ": "localhost"
      }
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-tags-3",
      "connector_name": "test-connector-name-axe-tags-3",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-tags-3",
      "resource": "test-resource-axe-tags-3",
      "state": 2,
      "output": "test-output-axe-tags-3"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-3
    Then the response code should be 200
    Then the response array key "data.0.tags" should contain only:
    """json
    [
      "env: prod",
      "cloud"
    ]
    """
