Feature: Get entities
  I need to be able to get a entities

  Scenario: given get search request should return entities only with string in name or type fields
    When I am admin
    When I do GET /api/v4/entities?search=to-entity-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-component-to-entity-get",
          "category": {
            "_id": "test-category-to-entity-get-2",
            "name": "test-category-to-entity-get-2-name"
          },
          "enabled": true,
          "old_entity_patterns": null,
          "connector": "test-connector-to-entity-get/test-connector-name-to-entity-get",
          "impact_level": 2,
          "infos": {},
          "name": "test-component-to-entity-get",
          "component": "test-component-to-entity-get",
          "type": "component",
          "ok_events": 0,
          "ko_events": 0,
          "state": 0,
          "status": 0,
          "impact_state": 0
        },
        {
          "_id": "test-connector-to-entity-get/test-connector-name-to-entity-get",
          "category": {
            "_id": "test-category-to-entity-get-1",
            "name": "test-category-to-entity-get-1-name"
          },
          "enabled": true,
          "old_entity_patterns": null,
          "impact_level": 1,
          "infos": {},
          "name": "test-connector-name-to-entity-get",
          "connector_type": "test-connector-to-entity-get",
          "type": "connector",
          "ok_events": 0,
          "ko_events": 0,
          "state": 0,
          "status": 0,
          "impact_state": 0
        },
        {
          "_id": "test-resource-to-entity-get-1/test-component-to-entity-get",
          "category": {
            "_id": "test-category-to-entity-get-2",
            "name": "test-category-to-entity-get-2-name"
          },
          "connector": "test-connector-to-entity-get/test-connector-name-to-entity-get",
          "enabled": true,
          "old_entity_patterns": null,
          "component": "test-component-to-entity-get",
          "impact_level": 3,
          "infos": {
            "test-resource-to-entity-get-1-info-1": {
              "name": "test-resource-to-entity-get-1-info-1-name",
              "description": "test-resource-to-entity-get-1-info-1-description",
              "value": "test-resource-to-entity-get-1-info-1-value"
            },
            "test-resource-to-entity-get-1-info-2": {
              "name": "test-resource-to-entity-get-1-info-2-name",
              "description": "test-resource-to-entity-get-1-info-2-description",
              "value": false
            },
            "test-resource-to-entity-get-1-info-3": {
              "name": "test-resource-to-entity-get-1-info-3-name",
              "description": "test-resource-to-entity-get-1-info-3-description",
              "value": 1022
            },
            "test-resource-to-entity-get-1-info-4": {
              "name": "test-resource-to-entity-get-1-info-4-name",
              "description": "test-resource-to-entity-get-1-info-4-description",
              "value": 10.45
            },
            "test-resource-to-entity-get-1-info-5": {
              "name": "test-resource-to-entity-get-1-info-5-name",
              "description": "test-resource-to-entity-get-1-info-5-description",
              "value": null
            },
            "test-resource-to-entity-get-1-info-6": {
              "name": "test-resource-to-entity-get-1-info-6-name",
              "description": "test-resource-to-entity-get-1-info-6-description",
              "value": ["test-resource-to-entity-get-1-info-6-value", false, 1022, 10.45, null]
            },
            "test-resource-to-entity-get-1-info-7": {
              "name": "test-resource-to-entity-get-1-info-7",
              "description": "test-resource-to-entity-get-1-info-7-description",
              "value": "test-resource-to-entity-get-1-info-7-value"
            }
          },
          "coordinates": {
            "lat": 64.52269494598361,
            "lng": 54.037685420804365
          },
          "name": "test-resource-to-entity-get-1",
          "type": "resource",
          "ok_events": 0,
          "ko_events": 0,
          "state": 3,
          "status": 1,
          "impact_state": 9,
          "alarm_last_update_date": 1597030219
        },
        {
          "_id": "test-resource-to-entity-get-2/test-component-to-entity-get",
          "category": null,
          "connector": "test-connector-to-entity-get/test-connector-name-to-entity-get",
          "enabled": true,
          "old_entity_patterns": null,
          "component": "test-component-to-entity-get",
          "impact_level": 1,
          "infos": {},
          "name": "test-resource-to-entity-get-2",
          "type": "resource",
          "ok_events": 0,
          "ko_events": 0,
          "state": 0,
          "status": 0,
          "impact_state": 0
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

  Scenario: given get sort request should return sorted entities
    When I am admin
    When I do GET /api/v4/entities?search=to-entity-get&sort_by=impact_level&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-to-entity-get-1/test-component-to-entity-get"
        },
        {
          "_id": "test-component-to-entity-get"
        },
        {
          "_id": "test-connector-to-entity-get/test-connector-name-to-entity-get"
        },
        {
          "_id": "test-resource-to-entity-get-2/test-component-to-entity-get"
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

  Scenario: given get all request should return entities with deletable flag and counts
    When I am admin
    When I do GET /api/v4/entities?search=to-entity-get&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-to-entity-get",
          "deletable": true,
          "impacts_count": 0,
          "depends_count": 0
        },
        {
          "_id": "test-connector-to-entity-get/test-connector-name-to-entity-get",
          "deletable": true,
          "impacts_count": 0,
          "depends_count": 0
        },
        {
          "_id": "test-resource-to-entity-get-1/test-component-to-entity-get",
          "deletable": false,
          "impacts_count": 0,
          "depends_count": 0
        },
        {
          "_id": "test-resource-to-entity-get-2/test-component-to-entity-get",
          "deletable": true,
          "impacts_count": 0,
          "depends_count": 0
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

  Scenario: given get category request should return entities which are matched to the category
    When I am admin
    When I do GET /api/v4/entities?category=test-category-to-entity-get-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-to-entity-get"
        },
        {
          "_id": "test-resource-to-entity-get-1/test-component-to-entity-get"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  Scenario: given get type request should return entities which are matched to the types
    When I am admin
    When I do GET /api/v4/entities?search=to-entity-get&type[]=component&type[]=connector
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-to-entity-get"
        },
        {
          "_id": "test-connector-to-entity-get/test-connector-name-to-entity-get"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  Scenario: given get filter request should return entities which are matched to the filter
    When I am admin
    When I do GET /api/v4/entities?filters[]=test-widgetfilter-to-entity-get-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-to-entity-get-1/test-component-to-entity-get"
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

  Scenario: given get filter request with old mongo query should return entities which are matched to the filter
    When I am admin
    When I do GET /api/v4/entities?filters[]=test-widgetfilter-to-entity-get-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-to-entity-get-1/test-component-to-entity-get"
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

  Scenario: given get search expression request should return entities which are matched to the expression
    When I am admin
    When I do GET /api/v4/entities?search=name%20LIKE%20"to-entity-get-1"
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-to-entity-get-1/test-component-to-entity-get"
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

  Scenario: old event statistics shouldn't be returned
    When I am admin
    When I do GET /api/v4/entities?search=test-resource-entity-api-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "ok_events": 0,
          "ko_events": 0
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """

  Scenario: given get unauth request should not allow access
    When I do GET /api/v4/entities
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/entities
    Then the response code should be 403

  Scenario: given get context graph request should return context graph
    When I am admin
    When I do GET /api/v4/entities/context-graph?_id=test-resource-to-entity-get-1/test-component-to-entity-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-to-entity-get/test-connector-name-to-entity-get"
      ],
      "impact": [
        "test-component-to-entity-get"
      ]
    }
    """

  Scenario: given get context graph request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/entities/context-graph?_id=test-entity-not-exist
    Then the response code should be 404

  Scenario: given get context graph unauth request should not allow access
    When I do GET /api/v4/entities/context-graph
    Then the response code should be 401

  Scenario: given get context graph request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/entities/context-graph
    Then the response code should be 403
