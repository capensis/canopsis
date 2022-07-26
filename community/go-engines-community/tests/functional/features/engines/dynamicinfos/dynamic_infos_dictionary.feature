Feature: get infos dictionary
  I need to be able to get infos dictionary

  Scenario: given requests should return entity infos dictionary
    Given I am admin
    When I wait the next periodical process
    When I do GET /api/v4/cat/dynamic-infos-dictionary/keys?search=test-dynamic-infos-dictionary-key
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-dynamic-infos-dictionary-key-1"
        },
        {
          "value": "test-dynamic-infos-dictionary-key-2"
        },
        {
          "value": "test-dynamic-infos-dictionary-key-3"
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 3
      }
    }
    """
    When I do GET /api/v4/cat/dynamic-infos-dictionary/keys?search=test-dynamic-infos-dictionary-key-2
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-dynamic-infos-dictionary-key-2"
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
    When I do GET /api/v4/cat/dynamic-infos-dictionary/values?key=test-dynamic-infos-dictionary-key-2
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-dynamic-infos-dictionary-value-11"
        },
        {
          "value": "test-dynamic-infos-dictionary-value-2"
        },
        {
          "value": "test-dynamic-infos-dictionary-value-5"
        },
        {
          "value": "test-dynamic-infos-dictionary-value-8"
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 4
      }
    }
    """
    When I do GET /api/v4/cat/dynamic-infos-dictionary/values?key=test-dynamic-infos-dictionary-key-2&search=test-dynamic-infos-dictionary-value-5
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-dynamic-infos-dictionary-value-5"
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
    When I do GET /api/v4/cat/dynamic-infos-dictionary/keys?search=test-dynamic-infos-dictionary-should-be-ignored
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-dynamic-infos-dictionary-should-be-ignored-key-1"
        },
        {
          "value": "test-dynamic-infos-dictionary-should-be-ignored-key-2"
        },
        {
          "value": "test-dynamic-infos-dictionary-should-be-ignored-key-3"
        },
        {
          "value": "test-dynamic-infos-dictionary-should-be-ignored-key-4"
        },
        {
          "value": "test-dynamic-infos-dictionary-should-be-ignored-key-5"
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 5
      }
    }
    """
    When I do GET /api/v4/cat/dynamic-infos-dictionary/values?key=test-dynamic-infos-dictionary-should-be-ignored-key-1
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 0
      }
    }
    """
    When I do GET /api/v4/cat/dynamic-infos-dictionary/values?key=test-dynamic-infos-dictionary-should-be-ignored-key-2
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 0
      }
    }
    """
    When I do GET /api/v4/cat/dynamic-infos-dictionary/values?key=test-dynamic-infos-dictionary-should-be-ignored-key-3
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 0
      }
    }
    """
    When I do GET /api/v4/cat/dynamic-infos-dictionary/values?key=test-dynamic-infos-dictionary-should-be-ignored-key-4
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 0
      }
    }
    """
    When I do GET /api/v4/cat/dynamic-infos-dictionary/values?key=test-dynamic-infos-dictionary-should-be-ignored-key-5
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 0
      }
    }
    """

  Scenario: given requests should update dynamic infos dictionary
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-dynamic-infos-dictionary-1",
      "connector_name": "test-connector-name-dynamic-infos-dictionary-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-dynamic-infos-dictionary-1",
      "resource": "test-resource-dynamic-infos-dictionary-1",
      "state": 1,
      "output": "test-output-dynamic-infos-dictionary-1"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "_id": "test-dynamic-infos-dictionary-1",
      "name": "test-dynamic-infos-dictionary-1-name",
      "description": "test-dynamic-infos-dictionary-1-description",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-dynamic-infos-dictionary-1"
            }
          }
        ]
      ],
      "infos": [
        {
          "name": "test-dynamic-infos-1-name",
          "value": "test-dynamic-infos-1-value-1"
        },
        {
          "name": "test-dynamic-infos-2-name",
          "value": "test-dynamic-infos-2-value-1"
        }
      ],
      "enabled": true
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "_id": "test-dynamic-infos-dictionary-2",
      "name": "test-dynamic-infos-dictionary-2-name",
      "description": "test-dynamic-infos-dictionary-2-description",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-dynamic-infos-dictionary-1"
            }
          }
        ]
      ],
      "infos": [
        {
          "name": "test-dynamic-infos-1-name",
          "value": "test-dynamic-infos-1-value-2"
        },
        {
          "name": "test-dynamic-infos-2-name",
          "value": "test-dynamic-infos-2-value-1"
        }
      ],
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I do GET /api/v4/cat/dynamic-infos-dictionary/values?key=test-dynamic-infos-1-name
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-dynamic-infos-1-value-1"
        },
        {
          "value": "test-dynamic-infos-1-value-2"
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 2
      }
    }
    """
    When I do GET /api/v4/cat/dynamic-infos-dictionary/values?key=test-dynamic-infos-2-name
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-dynamic-infos-2-value-1"
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
    When I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-dictionary-1:
    """json
    {
      "name": "test-dynamic-infos-dictionary-1-name",
      "description": "test-dynamic-infos-dictionary-1-description",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-dynamic-infos-dictionary-1"
            }
          }
        ]
      ],
      "infos": [
        {
          "name": "test-dynamic-infos-1-name",
          "value": "test-dynamic-infos-1-value-1-updated"
        },
        {
          "name": "test-dynamic-infos-2-name",
          "value": "test-dynamic-infos-2-value-1"
        }
      ],
      "enabled": true
    }
    """
    Then the response code should be 200
    When I wait the next periodical process
    When I do GET /api/v4/cat/dynamic-infos-dictionary/values?key=test-dynamic-infos-1-name
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-dynamic-infos-1-value-1-updated"
        },
        {
          "value": "test-dynamic-infos-1-value-2"
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 2
      }
    }
    """
    When I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-dictionary-1:
    """json
    {
      "name": "test-dynamic-infos-dictionary-1-name",
      "description": "test-dynamic-infos-dictionary-1-description",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-dynamic-infos-dictionary-1"
            }
          }
        ]
      ],
      "infos": [
        {
          "name": "another",
          "value": "another"
        }
      ],
      "enabled": true
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-dictionary-2:
    """json
    {
      "name": "test-dynamic-infos-dictionary-2-name",
      "description": "test-dynamic-infos-dictionary-2-description",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-dynamic-infos-dictionary-1"
            }
          }
        ]
      ],
      "infos": [
        {
          "name": "another",
          "value": "another"
        }
      ],
      "enabled": true
    }
    """
    Then the response code should be 200
    When I wait the next periodical process
    When I do GET /api/v4/cat/dynamic-infos-dictionary/values?key=test-dynamic-infos-1-name
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 0
      }
    }
    """
    When I do GET /api/v4/cat/dynamic-infos-dictionary/values?key=test-dynamic-infos-2-name
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 0
      }
    }
    """
