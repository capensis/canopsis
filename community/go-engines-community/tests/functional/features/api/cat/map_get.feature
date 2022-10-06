Feature: Get a map
  I need to be able to get a map
  Only admin should be able to get a map

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/cat/maps
    Then the response code should be 401

  Scenario: given get all request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/maps
    Then the response code should be 403

  Scenario: given search request should return maps
    When I am admin
    When I do GET /api/v4/cat/maps?search=test-map-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-map-to-get-1",
          "name": "test-map-to-get-1-name",
          "type": "mermaid",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1605263992,
          "updated": 1605263992
        },
        {
          "_id": "test-map-to-get-2",
          "name": "test-map-to-get-2-name",
          "type": "geo",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1605263992,
          "updated": 1605263992
        },
        {
          "_id": "test-map-to-get-3",
          "name": "test-map-to-get-3-name",
          "type": "flowchart",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1605263992,
          "updated": 1605263992
        },
        {
          "_id": "test-map-to-get-4",
          "name": "test-map-to-get-4-name",
          "type": "treeofdeps",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1605263992,
          "updated": 1605263992
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 4
      }
    }
    """

  Scenario: given get all request should return maps with flags
    When I am admin
    When I do GET /api/v4/cat/maps?search=test-map-to-get&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-map-to-get-1",
          "deletable": true
        },
        {
          "_id": "test-map-to-get-2",
          "deletable": true
        },
        {
          "_id": "test-map-to-get-3",
          "deletable": true
        },
        {
          "_id": "test-map-to-get-4",
          "deletable": true
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 4
      }
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/maps/test-map-to-get
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/maps/test-map-to-get
    Then the response code should be 403

  Scenario: given get request should return mermaid map
    When I am admin
    When I do GET /api/v4/cat/maps/test-map-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-map-to-get-1",
      "name": "test-map-to-get-1-name",
      "type": "mermaid",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1605263992,
      "updated": 1605263992,
      "parameters": {
        "code": "test-map-to-get-1-code",
        "theme": "test-map-to-get-1-theme",
        "points": [
          {
            "_id": "test-map-to-get-1-point-1",
            "x": 0,
            "y": 0,
            "entity": {
              "_id": "test-resource-to-map-edit-1/test-component-default",
              "name": "test-resource-to-map-edit-1",
              "type": "resource"
            },
            "map": null
          },
          {
            "_id": "test-map-to-get-1-point-2",
            "x": 100,
            "y": 100,
            "entity": null,
            "map": {
              "_id": "test-map-to-map-edit-1",
              "name": "test-map-to-map-edit-1-name"
            }
          },
          {
            "_id": "test-map-to-get-1-point-3",
            "x": 200,
            "y": 200,
            "entity": {
              "_id": "test-resource-to-map-edit-2/test-component-default",
              "name": "test-resource-to-map-edit-2",
              "type": "resource"
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

  Scenario: given get request should return geo map
    When I am admin
    When I do GET /api/v4/cat/maps/test-map-to-get-2
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-map-to-get-2",
      "name": "test-map-to-get-2-name",
      "type": "geo",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1605263992,
      "updated": 1605263992,
      "parameters": {
        "points": [
          {
            "_id": "test-map-to-get-2-point-1",
            "coordinates": {
              "lat": 62.34960927573042,
              "lng": 74.02834455685206
            },
            "entity": {
              "_id": "test-resource-to-map-edit-1/test-component-default",
              "name": "test-resource-to-map-edit-1",
              "type": "resource",
              "coordinates": {
                "lat": 62.34960927573042,
                "lng": 74.02834455685206
              }
            },
            "map": null
          },
          {
            "_id": "test-map-to-get-2-point-2",
            "coordinates": {
              "lat": 63.93737246791484,
              "lng": 34.991989666087385
            },
            "entity": null,
            "map": {
              "_id": "test-map-to-map-edit-1",
              "name": "test-map-to-map-edit-1-name"
            }
          },
          {
            "_id": "test-map-to-get-2-point-4",
            "coordinates": {
              "lat": 64.52269494598361,
              "lng": 54.037685420804365
            },
            "entity": {
              "_id": "test-resource-to-map-edit-3/test-component-default",
              "name": "test-resource-to-map-edit-3",
              "type": "resource"
            },
            "map": null
          }
        ]
      }
    }
    """

  Scenario: given get request should return flowchart map
    When I am admin
    When I do GET /api/v4/cat/maps/test-map-to-get-3
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-map-to-get-3",
      "name": "test-map-to-get-3-name",
      "type": "flowchart",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1605263992,
      "updated": 1605263992,
      "parameters": {
        "background_color": "#efa8ca",
        "shapes": [
          {
            "_id": "test-map-to-get-1-shape-1",
            "type": "rect",
            "x": 200,
            "y": 200,
            "width": 100,
            "height": 100,
            "connections": [
              {
                "shape": "test-map-to-get-1-shape-2"
              }
            ],
            "connectedTo": [],
            "text": "test-map-to-get-1-shape-1-text",
            "textProperties": {
              "color": "#b61c1c"
            },
            "properties": {
              "fill": "#b61c1c"
            },
            "aspectRatio": true
          },
          {
            "_id": "test-map-to-get-1-shape-2",
            "type": "rect",
            "x": 300,
            "y": 300,
            "width": 100,
            "height": 100,
            "connections": [],
            "connectedTo": ["test-map-to-get-1-shape-1"],
            "text": "test-map-to-get-1-shape-2-text",
            "textProperties": {
              "color": "#b61c1c"
            },
            "properties": {
              "fill": "#b61c1c"
            },
            "aspectRatio": true
          }
        ],
        "points": [
          {
            "_id": "test-map-to-get-1-point-1",
            "x": 0,
            "y": 0,
            "entity": {
              "_id": "test-resource-to-map-edit-1/test-component-default",
              "name": "test-resource-to-map-edit-1",
              "type": "resource"
            },
            "map": null
          },
          {
            "_id": "test-map-to-get-1-point-2",
            "x": 100,
            "y": 100,
            "entity": null,
            "map": {
              "_id": "test-map-to-map-edit-1",
              "name": "test-map-to-map-edit-1-name"
            }
          },
          {
            "_id": "test-map-to-get-1-point-3",
            "shape": "test-map-to-get-1-shape-1",
            "entity": {
              "_id": "test-resource-to-map-edit-2/test-component-default",
              "name": "test-resource-to-map-edit-2",
              "type": "resource"
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

  Scenario: given get request should return treeofdeps map
    When I am admin
    When I do GET /api/v4/cat/maps/test-map-to-get-4
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-map-to-get-4",
      "name": "test-map-to-get-4-name",
      "type": "treeofdeps",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1605263992,
      "updated": 1605263992,
      "parameters": {
        "type": "treeofdeps",
        "entities": [
          {
            "entity": {
              "_id": "test-resource-to-map-edit-1/test-component-default",
              "name": "test-resource-to-map-edit-1",
              "type": "resource",
              "depends_count": 0
            },
            "pinned_entities": []
          }
        ]
      }
    }
    """

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/cat/maps/test-map-not-exist
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """
