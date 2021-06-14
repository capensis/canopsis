Feature: update watcher weather on event
  I need to be able to see new watcher weather state on event

  Scenario: given watcher for entity should update watcher weather on entity event
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-watcher-weather-1",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-watcher-weather-1"}
      ],
      "output_template": "Test watcher weather 1"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I send an event:
    """
    {
      "connector": "test-connector-watcher-weather-1",
      "connector_name": "test-connector-name-watcher-weather-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-watcher-weather-1",
      "resource": "test-resource-watcher-weather-1",
      "state": 2,
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers?filter={"name":"test-watcher-weather-1"}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-watcher-weather-1",
          "connector": "watcher",
          "connector_name": "watcher",
          "resource": "",
          "state": {"val": 2},
          "status": {"val": 1},
          "ack": null,
          "snooze": null,
          "infos": {},
          "icon": "major",
          "secondary_icon": "",
          "color": "major",
          "is_action_required": true,
          "alarm_counters": [],
          "pbehaviors": []
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

  Scenario: given watcher for multiple entities should set watcher weather state as worst entity state
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-watcher-weather-2",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-watcher-weather-2-1"},
        {"name": "test-resource-watcher-weather-2-2"},
        {"name": "test-resource-watcher-weather-2-3"}
      ],
      "output_template": "Test watcher weather 2"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I send an event:
    """
    {
      "connector": "test-connector-watcher-weather-2",
      "connector_name": "test-connector-name-watcher-weather-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-watcher-weather-2",
      "resource": "test-resource-watcher-weather-2-1",
      "state": 1,
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-connector-watcher-weather-2",
      "connector_name": "test-connector-name-watcher-weather-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-watcher-weather-2",
      "resource": "test-resource-watcher-weather-2-2",
      "state": 2,
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-connector-watcher-weather-2",
      "connector_name": "test-connector-name-watcher-weather-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-watcher-weather-2",
      "resource": "test-resource-watcher-weather-2-3",
      "state": 3,
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers?filter={"name":"test-watcher-weather-2"}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-watcher-weather-2",
          "connector": "watcher",
          "connector_name": "watcher",
          "resource": "",
          "state": {"val": 3},
          "status": {"val": 1},
          "ack": null,
          "snooze": null,
          "infos": {},
          "icon": "critical",
          "secondary_icon": "",
          "color": "critical",
          "is_action_required": true,
          "alarm_counters": [],
          "pbehaviors": []
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

  Scenario: given watcher for entity and no open alarm should get ok watcher weather state
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-watcher-weather-3",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-watcher-weather-3"}
      ],
      "output_template": "Test watcher weather 3"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I send an event:
    """
    {
      "connector": "test-connector-watcher-weather-3",
      "connector_name": "test-connector-name-watcher-weather-3",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-watcher-weather-3",
      "resource": "test-resource-watcher-weather-3",
      "state": 0,
      "output": "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/weather-watchers?filter={"name":"test-watcher-weather-3"}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-watcher-weather-3",
          "connector": "",
          "connector_name": "",
          "component": "",
          "resource": "",
          "state": {"val": 0},
          "status": {"val": 0},
          "ack": null,
          "snooze": null,
          "infos": {},
          "icon": "ok",
          "secondary_icon": "",
          "color": "ok",
          "is_action_required": false,
          "alarm_counters": [],
          "pbehaviors": []
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

  Scenario: given watcher for multiple entities and acked alarms
    should return false in is_action_required
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-watcher-weather-4",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-watcher-weather-4-1"},
        {"name": "test-resource-watcher-weather-4-2"}
      ],
      "output_template": "Test watcher weather 4"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I send an event:
    """
    {
      "connector": "test-connector-watcher-weather-4",
      "connector_name": "test-connector-name-watcher-weather-4",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-watcher-weather-4",
      "resource": "test-resource-watcher-weather-4-1",
      "state": 2,
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-connector-watcher-weather-4",
      "connector_name": "test-connector-name-watcher-weather-4",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-watcher-weather-4",
      "resource": "test-resource-watcher-weather-4-2",
      "state": 3,
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-connector-watcher-weather-4",
      "connector_name": "test-connector-name-watcher-weather-4",
      "source_type": "resource",
      "event_type": "ack",
      "component":  "test-component-watcher-weather-4",
      "resource": "test-resource-watcher-weather-4-1",
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-connector-watcher-weather-4",
      "connector_name": "test-connector-name-watcher-weather-4",
      "source_type": "resource",
      "event_type": "ack",
      "component":  "test-component-watcher-weather-4",
      "resource": "test-resource-watcher-weather-4-2",
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers?filter={"name":"test-watcher-weather-4"}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-watcher-weather-4",
          "connector": "watcher",
          "connector_name": "watcher",
          "resource": "",
          "state": {"val": 3},
          "status": {"val": 1},
          "ack": null,
          "snooze": null,
          "infos": {},
          "icon": "critical",
          "secondary_icon": "",
          "color": "critical",
          "is_action_required": false,
          "alarm_counters": [],
          "pbehaviors": []
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