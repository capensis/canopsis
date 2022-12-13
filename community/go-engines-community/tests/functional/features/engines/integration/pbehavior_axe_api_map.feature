Feature: Get a map's state
  I need to be able to get a map's state
  Only admin should be able to get a map's state

  Scenario: given mermaid map with snoozed alarm should return alarm
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-axe-api-map-1",
      "connector_name": "test-connector-name-pbehavior-axe-api-map-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-axe-api-map-1",
      "resource": "test-resource-pbehavior-axe-api-map-1",
      "state": 2,
      "output": "test-output-pbehavior-axe-api-map-1"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-axe-api-map-1",
      "connector_name": "test-connector-name-pbehavior-axe-api-map-1",
      "source_type": "resource",
      "event_type": "snooze",
      "component":  "test-component-pbehavior-axe-api-map-1",
      "resource": "test-resource-pbehavior-axe-api-map-1",
      "duration": 6000,
      "output": "test-output-pbehavior-axe-api-map-1-snooze",
      "initiator": "user"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/cat/maps:
    """json
    {
      "name": "test-map-pbehavior-axe-api-map-1-name",
      "type": "mermaid",
      "parameters": {
        "code": "test-map-pbehavior-axe-api-map-1-code",
        "theme": "test-map-pbehavior-axe-api-map-1-theme",
        "points": [
          {
            "x": 0,
            "y": 0,
            "entity": "test-resource-pbehavior-axe-api-map-1/test-component-pbehavior-axe-api-map-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/cat/map-state/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-map-pbehavior-axe-api-map-1-name",
      "type": "mermaid",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "parameters": {
        "code": "test-map-pbehavior-axe-api-map-1-code",
        "theme": "test-map-pbehavior-axe-api-map-1-theme",
        "points": [
          {
            "x": 0,
            "y": 0,
            "entity": {
              "_id": "test-resource-pbehavior-axe-api-map-1/test-component-pbehavior-axe-api-map-1",
              "name": "test-resource-pbehavior-axe-api-map-1",
              "type": "resource",
              "category": null,
              "connector": "test-connector-pbehavior-axe-api-map-1/test-connector-name-pbehavior-axe-api-map-1",
              "component":  "test-component-pbehavior-axe-api-map-1",
              "infos": {},
              "impact_level": 1,
              "impact_state": 2,
              "state": 2,
              "status": 1,
              "snooze": {
                "_t": "snooze",
                "a": "root",
                "user_id": "root",
                "m": "test-output-pbehavior-axe-api-map-1-snooze",
                "initiator": "user"
              },
              "ko_events": 1,
              "ok_events": 0
            },
            "map": null
          }
        ]
      }
    }
    """

  Scenario: given mermaid map with acked alarm should return alarm
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-axe-api-map-2",
      "connector_name": "test-connector-name-pbehavior-axe-api-map-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-axe-api-map-2",
      "resource": "test-resource-pbehavior-axe-api-map-2",
      "state": 2,
      "output": "test-output-pbehavior-axe-api-map-2"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-axe-api-map-2",
      "connector_name": "test-connector-name-pbehavior-axe-api-map-2",
      "source_type": "resource",
      "event_type": "ack",
      "component":  "test-component-pbehavior-axe-api-map-2",
      "resource": "test-resource-pbehavior-axe-api-map-2",
      "output": "test-output-pbehavior-axe-api-map-2-ack",
      "initiator": "user"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/cat/maps:
    """json
    {
      "name": "test-map-pbehavior-axe-api-map-2-name",
      "type": "mermaid",
      "parameters": {
        "code": "test-map-pbehavior-axe-api-map-2-code",
        "theme": "test-map-pbehavior-axe-api-map-2-theme",
        "points": [
          {
            "x": 0,
            "y": 0,
            "entity": "test-resource-pbehavior-axe-api-map-2/test-component-pbehavior-axe-api-map-2"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/cat/map-state/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-map-pbehavior-axe-api-map-2-name",
      "type": "mermaid",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "parameters": {
        "code": "test-map-pbehavior-axe-api-map-2-code",
        "theme": "test-map-pbehavior-axe-api-map-2-theme",
        "points": [
          {
            "x": 0,
            "y": 0,
            "entity": {
              "_id": "test-resource-pbehavior-axe-api-map-2/test-component-pbehavior-axe-api-map-2",
              "name": "test-resource-pbehavior-axe-api-map-2",
              "type": "resource",
              "category": null,
              "connector": "test-connector-pbehavior-axe-api-map-2/test-connector-name-pbehavior-axe-api-map-2",
              "component":  "test-component-pbehavior-axe-api-map-2",
              "infos": {},
              "impact_level": 1,
              "impact_state": 2,
              "state": 2,
              "status": 1,
              "ack": {
                "_t": "ack",
                "a": "root",
                "user_id": "root",
                "m": "test-output-pbehavior-axe-api-map-2-ack",
                "initiator": "user",
                "val": 0
              },
              "ko_events": 1,
              "ok_events": 0
            },
            "map": null
          }
        ]
      }
    }
    """

  Scenario: given mermaid map with alarm in pbehavior should return alarm
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-axe-api-map-3",
      "connector_name": "test-connector-name-pbehavior-axe-api-map-3",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-axe-api-map-3",
      "resource": "test-resource-pbehavior-axe-api-map-3",
      "state": 2,
      "output": "test-output-pbehavior-axe-api-map-3"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-pbehavior-axe-api-map-3",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-axe-api-map-3/test-component-pbehavior-axe-api-map-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do POST /api/v4/cat/maps:
    """json
    {
      "name": "test-map-pbehavior-axe-api-map-3-name",
      "type": "mermaid",
      "parameters": {
        "code": "test-map-pbehavior-axe-api-map-3-code",
        "theme": "test-map-pbehavior-axe-api-map-3-theme",
        "points": [
          {
            "x": 0,
            "y": 0,
            "entity": "test-resource-pbehavior-axe-api-map-3/test-component-pbehavior-axe-api-map-3"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/cat/map-state/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-map-pbehavior-axe-api-map-3-name",
      "type": "mermaid",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "parameters": {
        "code": "test-map-pbehavior-axe-api-map-3-code",
        "theme": "test-map-pbehavior-axe-api-map-3-theme",
        "points": [
          {
            "x": 0,
            "y": 0,
            "entity": {
              "_id": "test-resource-pbehavior-axe-api-map-3/test-component-pbehavior-axe-api-map-3",
              "name": "test-resource-pbehavior-axe-api-map-3",
              "type": "resource",
              "category": null,
              "connector": "test-connector-pbehavior-axe-api-map-3/test-connector-name-pbehavior-axe-api-map-3",
              "component":  "test-component-pbehavior-axe-api-map-3",
              "infos": {},
              "pbehavior_info": {
                "canonical_type": "maintenance",
                "icon_name": "build",
                "name": "test-pbehavior-pbehavior-axe-api-map-3",
                "reason": "test-reason-to-engine",
                "reason_name": "Test Engine",
                "type": "test-maintenance-type-to-engine",
                "type_name": "Engine maintenance"
              },
              "impact_level": 1,
              "impact_state": 2,
              "state": 2,
              "status": 1,
              "ko_events": 1,
              "ok_events": 0
            },
            "map": null
          }
        ]
      }
    }
    """
