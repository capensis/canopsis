Feature: Get entities
  I need to be able to get a entities

  Scenario: given get search request should return entities only
  with string in name or type fields
    When I am admin
    When I do GET /api/v4/entities?search=test-entity-to-get
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "data": [
        {
          "_id": "test-entity-to-get-component",
          "category": {
            "_id": "test-category-to-entity-get-2",
            "name": "test-category-to-entity-get-2-name",
            "author": "test-category-to-entity-get-2-author",
            "created": 1592215337,
            "updated": 1592215337
          },
          "depends": [
            "test-entity-to-get-resource/test-entity-to-get-component"
          ],
          "enable_history": [
            1597030220
          ],
          "enabled": true,
          "impact": [
            "test-entity-to-get-component"
          ],
          "impact_level": 2,
          "infos": {},
          "measurements": null,
          "name": "test-entity-to-get-component",
          "type": "component",
          "ok_events": 0,
          "ko_events": 0,
          "state": 0
        },
        {
          "_id": "test-entity-to-get-connector/test-entity-to-get-connector-name",
          "category": {
            "_id": "test-category-to-entity-get-1",
            "name": "test-category-to-entity-get-1-name",
            "author": "test-category-to-entity-get-1-author",
            "created": 1592215337,
            "updated": 1592215337
          },
          "depends": [
            "test-entity-to-get-component"
          ],
          "enable_history": [
            1597030220
          ],
          "enabled": true,
          "impact": [
            "test-entity-to-get-resource/test-entity-to-get-component"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-entity-to-get-connector-name",
          "type": "connector",
          "ok_events": 0,
          "ko_events": 0,
          "state": 0
        },
        {
          "_id": "test-entity-to-get-resource/test-entity-to-get-component",
          "category": {
            "_id": "test-category-to-entity-get-2",
            "name": "test-category-to-entity-get-2-name",
            "author": "test-category-to-entity-get-2-author",
            "created": 1592215337,
            "updated": 1592215337
          },
          "depends": [
            "test-entity-to-get-connector/test-entity-to-get-connector-name"
          ],
          "enable_history": [
            1597030220
          ],
          "enabled": true,
          "impact": [
            "test-entity-to-get-component"
          ],
          "impact_level": 3,
          "infos": {
            "test-entity-to-get-info-1": {
              "name": "test-entity-to-get-info-1-name",
              "description": "test-entity-to-get-info-1-description",
              "value": "test-entity-to-get-info-1-value"
            },
            "test-entity-to-get-info-2": {
              "name": "test-entity-to-get-info-2-name",
              "description": "test-entity-to-get-info-2-description",
              "value": false
            },
            "test-entity-to-get-info-3": {
              "name": "test-entity-to-get-info-3-name",
              "description": "test-entity-to-get-info-3-description",
              "value": 1022
            },
            "test-entity-to-get-info-4": {
              "name": "test-entity-to-get-info-4-name",
              "description": "test-entity-to-get-info-4-description",
              "value": 10.45
            },
            "test-entity-to-get-info-5": {
              "name": "test-entity-to-get-info-5-name",
              "description": "test-entity-to-get-info-5-description",
              "value": null
            },
            "test-entity-to-get-info-6": {
              "name": "test-entity-to-get-info-6-name",
              "description": "test-entity-to-get-info-6-description",
              "value": ["test-entity-to-get-info-6-value", false, 1022, 10.45, null]
            },
            "test-entity-to-get-info-7": {
              "name": "test-entity-to-get-info-7",
              "description": "test-entity-to-get-info-7-description",
              "value": "test-entity-to-get-info-7-value"
            }
          },
          "measurements": null,
          "name": "test-entity-to-get-resource",
          "type": "resource",
          "ok_events": 0,
          "ko_events": 0,
          "state": 0
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: given get sort request should return sorted entities
    When I am admin
    When I do GET /api/v4/entities?search=test-entity-to-get&sort_by=impact_level&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-entity-to-get-resource/test-entity-to-get-component"
        },
        {
          "_id": "test-entity-to-get-component"
        },
        {
          "_id": "test-entity-to-get-connector/test-entity-to-get-connector-name"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: given get all request should return entities with deletable field
    When I am admin
    When I do GET /api/v4/entities?search=test-entity-to-check-deletable&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-entity-to-check-deletable-component",
          "deletable": true
        },
        {
          "_id": "test-entity-to-check-deletable-connector/test-entity-to-check-deletable-connector-name",
          "deletable": true
        },
        {
          "_id": "test-entity-to-check-deletable-resource/test-entity-to-check-deletable-component",
          "deletable": false
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: given get category request should return entities which are matched to the category
    When I am admin
    When I do GET /api/v4/entities?category=test-category-to-entity-get-2
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-entity-to-get-component"
        },
        {
          "_id": "test-entity-to-get-resource/test-entity-to-get-component"
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
    When I do GET /api/v4/entities?filter={"depends":{"$in":["test-entity-to-get-component"]}}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-entity-to-get-connector/test-entity-to-get-connector-name"
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

  Scenario: given get search expression request should return entities which are matched
  to expression filter
    When I am admin
    When I do GET /api/v4/entities?search=_id="test-entity-to-get-resource/test-entity-to-get-component"
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-entity-to-get-resource/test-entity-to-get-component"
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

  Scenario: given get search expression request should return entities which are matched
  to expression filter
    When I am admin
    When I do GET /api/v4/entities?search=_id%20LIKE%20"test-entity-to-get-component"
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-entity-to-get-component"
        },
        {
          "_id": "test-entity-to-get-resource/test-entity-to-get-component"
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

  Scenario: given get search expression request should return entities which are matched
  to expression filter
    When I am admin
    When I do GET /api/v4/entities?search=depends%20CONTAINS%20"test-entity-to-get-resource/test-entity-to-get-component"
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-entity-to-get-component"
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

  Scenario: given get search expression request should return entities some of entities with idle since
    When I am admin
    When I do GET /api/v4/entities?search=test-idle-since-resource
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-idle-since-resource-1/test-idle-since-component",
          "idle_since": 123
        },
        {
          "_id": "test-idle-since-resource-2/test-idle-since-component",
          "idle_since": 321
        },
        {
          "_id": "test-idle-since-resource-without/test-idle-since-component"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: given get search expression request should return entities only if idle_since is present
    When I am admin
    When I do GET /api/v4/entities?search=test-idle-since-resource&no_events=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-idle-since-resource-1/test-idle-since-component",
          "idle_since": 123
        },
        {
          "_id": "test-idle-since-resource-2/test-idle-since-component",
          "idle_since": 321
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

  Scenario: given get search expression request should return entities only if idle_since is present with sort by idle_since
    When I am admin
    When I do GET /api/v4/entities?search=test-idle-since-resource&no_events=true&sort_by=idle_since&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-idle-since-resource-1/test-idle-since-component",
          "idle_since": 123
        },
        {
          "_id": "test-idle-since-resource-2/test-idle-since-component",
          "idle_since": 321
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
    When I do GET /api/v4/entities?search=test-idle-since-resource&no_events=true&sort_by=idle_since&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-idle-since-resource-2/test-idle-since-component",
          "idle_since": 321
        },
        {
          "_id": "test-idle-since-resource-1/test-idle-since-component",
          "idle_since": 123
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
