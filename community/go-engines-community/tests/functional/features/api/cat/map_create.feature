Feature: Create a map
  I need to be able to create a map
  Only admin should be able to create a map

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/cat/maps
    Then the response code should be 401

  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/maps
    Then the response code should be 403

  Scenario: given create mermaid map request should return ok
    When I am admin
    When I do POST /api/v4/cat/maps:
    """json
    {
      "name": "test-map-to-create-1-name",
      "type": "mermaid",
      "parameters": {
        "code": "test-map-to-create-1-code",
        "theme": "test-map-to-create-1-theme",
        "points": [
          {
            "x": 0,
            "y": 0,
            "entity": "test-resource-to-map-edit-1/test-component-default"
          },
          {
            "x": 100,
            "y": 100,
            "map": "test-map-to-map-edit-1"
          },
          {
            "x": 200,
            "y": 200,
            "entity": "test-resource-to-map-edit-2/test-component-default",
            "map": "test-map-to-map-edit-2"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "name": "test-map-to-create-1-name",
      "type": "mermaid",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "parameters": {
        "code": "test-map-to-create-1-code",
        "theme": "test-map-to-create-1-theme",
        "points": [
          {
            "x": 0,
            "y": 0,
            "entity": {
              "_id": "test-resource-to-map-edit-1/test-component-default",
              "name": "test-resource-to-map-edit-1"
            },
            "map": null
          },
          {
            "x": 100,
            "y": 100,
            "entity": null,
            "map": {
              "_id": "test-map-to-map-edit-1",
              "name": "test-map-to-map-edit-1-name"
            }
          },
          {
            "x": 200,
            "y": 200,
            "entity": {
              "_id": "test-resource-to-map-edit-2/test-component-default",
              "name": "test-resource-to-map-edit-2"
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
    When I do GET /api/v4/cat/maps/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-map-to-create-1-name",
      "type": "mermaid",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "parameters": {
        "code": "test-map-to-create-1-code",
        "theme": "test-map-to-create-1-theme",
        "points": [
          {
            "x": 0,
            "y": 0,
            "entity": {
              "_id": "test-resource-to-map-edit-1/test-component-default",
              "name": "test-resource-to-map-edit-1"
            },
            "map": null
          },
          {
            "x": 100,
            "y": 100,
            "entity": null,
            "map": {
              "_id": "test-map-to-map-edit-1",
              "name": "test-map-to-map-edit-1-name"
            }
          },
          {
            "x": 200,
            "y": 200,
            "entity": {
              "_id": "test-resource-to-map-edit-2/test-component-default",
              "name": "test-resource-to-map-edit-2"
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

  Scenario: given create request with missing fields should return bad request error
    When I am admin
    When I do POST /api/v4/cat/maps:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name is missing.",
        "type": "Type is missing."
      }
    }
    """
    When I do POST /api/v4/cat/maps:
    """json
    {
      "type": "unknown"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "type": "Type must be one of [mermaid]."
      }
    }
    """
    When I do POST /api/v4/cat/maps:
    """json
    {
      "type": "mermaid"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name is missing.",
        "parameters.code": "Code is missing.",
        "parameters.points": "Points is missing."
      }
    }
    """
    When I do POST /api/v4/cat/maps:
    """json
    {
      "type": "mermaid",
      "parameters": {
        "points": []
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name is missing.",
        "parameters.code": "Code is missing.",
        "parameters.points": "Points is missing."
      }
    }
    """
    When I do POST /api/v4/cat/maps:
    """json
    {
      "type": "mermaid",
      "parameters": {
        "points": [
          {
            "x": 0,
            "y": 0
          }
        ]
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name is missing.",
        "parameters.code": "Code is missing.",
        "parameters.points.0.entity": "Entity is required when Map is not present.",
        "parameters.points.0.map": "Map is required when Entity is not present."
      }
    }
    """
