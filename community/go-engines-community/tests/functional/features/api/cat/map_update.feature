Feature: Update a map
  I need to be able to update a map
  Only admin should be able to update a map

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/cat/maps/test-map-to-update
    Then the response code should be 401

  Scenario: given update request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/maps/test-map-to-update
    Then the response code should be 403

  Scenario: given update mermaid map request should return ok
    When I am admin
    When I do PUT /api/v4/cat/maps/test-map-to-update-1:
    """json
    {
      "name": "test-map-to-update-1-name",
      "type": "mermaid",
      "parameters": {
        "code": "test-map-to-update-1-code",
        "theme": "test-map-to-update-1-theme",
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
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-map-to-update-1",
      "name": "test-map-to-update-1-name",
      "type": "mermaid",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1605263992,
      "parameters": {
        "code": "test-map-to-update-1-code",
        "theme": "test-map-to-update-1-theme",
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

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/cat/maps/test-map-not-exist:
    """json
    {
      "name": "test-map-to-not-exist-name",
      "type": "mermaid",
      "parameters": {
        "code": "test-map-to-not-exist-code",
        "theme": "test-map-to-not-exist-theme",
        "points": [
          {
            "x": 0,
            "y": 0,
            "entity": "test-resource-to-map-edit-1/test-component-default"
          }
        ]
      }
    }
    """
    Then the response code should be 404

  Scenario: given update request with missing fields should return bad request error
    When I am admin
    When I do PUT /api/v4/cat/maps/test-map-to-update:
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
    When I do PUT /api/v4/cat/maps/test-map-to-update:
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
        "type": "Type must be one of [mermaid geo treeofdeps flowchart]."
      }
    }
    """
