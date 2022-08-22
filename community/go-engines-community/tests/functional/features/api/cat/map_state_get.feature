Feature: Get a map's state
  I need to be able to get a map's state
  Only admin should be able to get a map's state

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/map-state/test-map-to-state-get
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/map-state/test-map-to-state-get
    Then the response code should be 403

  Scenario: given get mermaid map request should return map
    When I am admin
    When I do GET /api/v4/cat/map-state/test-map-to-state-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-map-to-state-get-1",
      "name": "test-map-to-state-get-1-name",
      "type": "mermaid",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1605263992,
      "updated": 1605263992,
      "parameters": {
        "code": "test-map-to-state-get-1-code",
        "theme": "test-map-to-state-get-1-theme",
        "points": [
          {
            "_id": "test-map-to-state-get-1-point-1",
            "x": 0,
            "y": 0,
            "entity": {
              "_id": "test-resource-to-map-state-get-1/test-component-default",
              "name": "test-resource-to-map-state-get-1",
              "type": "resource",
              "enabled": true,
              "category": {
                "_id": "test-category-to-map-state-get-1",
                "name": "test-category-to-map-state-get-1-name",
                "author": "root",
                "created": 1592215337,
                "updated": 1592215337
              },
              "connector": "test-connector-default/test-connector-default-name",
              "component": "test-component-default",
              "infos": {
                "test-resource-to-map-state-get-1-info-1-name": {
                  "name": "test-resource-to-map-state-get-1-info-1-name",
                  "description": "test-resource-to-map-state-get-1-info-1-description",
                  "value": "test-resource-to-map-state-get-1-info-1-value"
                }
              },
              "last_event_date": 1605263992,
              "impact_level": 1,
              "alarm_last_update_date": 1605263992,
              "impact_state": 3,
              "state": 3,
              "status": 1,
              "ack": {
                "_t": "ack",
                "a": "root",
                "user_id": "root",
                "m": "test-alarm-to-map-state-get-1-ack-message",
                "t": 1605263992,
                "initiator": "user",
                "val": 0
              },
              "ko_events": 0,
              "ok_events": 0
            },
            "map": null
          },
          {
            "_id": "test-map-to-state-get-1-point-2",
            "x": 100,
            "y": 100,
            "entity": null,
            "map": {
              "_id": "test-map-to-map-edit-1",
              "name": "test-map-to-map-edit-1-name"
            }
          },
          {
            "_id": "test-map-to-state-get-1-point-3",
            "x": 200,
            "y": 200,
            "entity": {
              "_id": "test-resource-to-map-state-get-2/test-component-default",
              "name": "test-resource-to-map-state-get-2",
              "type": "resource",
              "enabled": true,
              "category": null,
              "connector": "test-connector-default/test-connector-default-name",
              "component": "test-component-default",
              "infos": {},
              "impact_level": 1,
              "impact_state": 0,
              "state": 0,
              "status": 0,
              "ko_events": 0,
              "ok_events": 0
            },
            "map": {
              "_id": "test-map-to-map-edit-2",
              "name": "test-map-to-map-edit-2-name"
            }
          }
        ]
      }
    }
    """

  Scenario: given get filtered mermaid map request should return map
    When I am admin
    When I do GET /api/v4/cat/map-state/test-map-to-state-get-1?filter=test-widgetfilter-to-map-state-get-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-map-to-state-get-1",
      "parameters": {
        "points": [
          {
            "_id": "test-map-to-state-get-1-point-1",
            "entity": {
              "_id": "test-resource-to-map-state-get-1/test-component-default"
            },
            "map": null
          },
          {
            "_id": "test-map-to-state-get-1-point-2",
            "entity": null,
            "map": {
              "_id": "test-map-to-map-edit-1"
            }
          }
        ]
      }
    }
    """
    When I do GET /api/v4/cat/map-state/test-map-to-state-get-2?filter=test-widgetfilter-to-map-state-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-map-to-state-get-2",
      "name": "test-map-to-state-get-2-name",
      "type": "mermaid",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1605263992,
      "updated": 1605263992,
      "parameters": {
        "code": "test-map-to-state-get-2-code",
        "theme": "test-map-to-state-get-2-theme"
      }
    }
    """

  Scenario: given get filtered by category mermaid map request should return map
    When I am admin
    When I do GET /api/v4/cat/map-state/test-map-to-state-get-1?category=test-category-to-map-state-get-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-map-to-state-get-1",
      "parameters": {
        "points": [
          {
            "_id": "test-map-to-state-get-1-point-1",
            "entity": {
              "_id": "test-resource-to-map-state-get-1/test-component-default"
            },
            "map": null
          },
          {
            "_id": "test-map-to-state-get-1-point-2",
            "entity": null,
            "map": {
              "_id": "test-map-to-map-edit-1"
            }
          }
        ]
      }
    }
    """
    When I do GET /api/v4/cat/map-state/test-map-to-state-get-2?category=test-category-to-map-state-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-map-to-state-get-2",
      "name": "test-map-to-state-get-2-name",
      "type": "mermaid",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1605263992,
      "updated": 1605263992,
      "parameters": {
        "code": "test-map-to-state-get-2-code",
        "theme": "test-map-to-state-get-2-theme"
      }
    }
    """

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/cat/map-state/test-map-not-exist
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """
