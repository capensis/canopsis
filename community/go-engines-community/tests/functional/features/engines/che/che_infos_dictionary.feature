Feature: get infos dictionary
  I need to be able to get infos dictionary

  Scenario: given requests should return entity infos dictionary
    Given I am admin
    When I wait the next periodical process
    When I do GET /api/v4/entity-infos-dictionary/keys?search=test-entity-infos-dictionary-key
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-entity-infos-dictionary-key-1"
        },
        {
          "value": "test-entity-infos-dictionary-key-2"
        },
        {
          "value": "test-entity-infos-dictionary-key-3"
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
    When I do GET /api/v4/entity-infos-dictionary/keys?search=test-entity-infos-dictionary-key-2
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-entity-infos-dictionary-key-2"
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
    When I do GET /api/v4/entity-infos-dictionary/values?key=test-entity-infos-dictionary-key-2
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-entity-infos-dictionary-value-11"
        },
        {
          "value": "test-entity-infos-dictionary-value-2"
        },
        {
          "value": "test-entity-infos-dictionary-value-5"
        },
        {
          "value": "test-entity-infos-dictionary-value-8"
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
    When I do GET /api/v4/entity-infos-dictionary/values?key=test-entity-infos-dictionary-key-2&search=test-entity-infos-dictionary-value-5
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-entity-infos-dictionary-value-5"
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
    When I do GET /api/v4/entity-infos-dictionary/keys?search=test-entity-infos-dictionary-should-be-ignored
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-entity-infos-dictionary-should-be-ignored-key-1"
        },
        {
          "value": "test-entity-infos-dictionary-should-be-ignored-key-2"
        },
        {
          "value": "test-entity-infos-dictionary-should-be-ignored-key-3"
        },
        {
          "value": "test-entity-infos-dictionary-should-be-ignored-key-4"
        },
        {
          "value": "test-entity-infos-dictionary-should-be-ignored-key-5"
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
    When I do GET /api/v4/entity-infos-dictionary/values?key=test-entity-infos-dictionary-should-be-ignored-key-1
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
    When I do GET /api/v4/entity-infos-dictionary/values?key=test-entity-infos-dictionary-should-be-ignored-key-2
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
    When I do GET /api/v4/entity-infos-dictionary/values?key=test-entity-infos-dictionary-should-be-ignored-key-3
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
    When I do GET /api/v4/entity-infos-dictionary/values?key=test-entity-infos-dictionary-should-be-ignored-key-4
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
    When I do GET /api/v4/entity-infos-dictionary/values?key=test-entity-infos-dictionary-should-be-ignored-key-5
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

  Scenario: given requests should update entity infos dictionary
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "event_pattern":[[
        {
          "field": "event_type",
          "cond": {
            "type": "eq",
            "value": "check"
          }
        },
        {
          "field": "component",
          "cond": {
            "type": "eq",
            "value": "test-component-entity-infos-dictionary-1"
          }
        }
      ]],
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "test-resource-entity-infos-dictionary-1-key-1",
            "description": "Client",
            "value": "{{ `{{ .Event.ExtraInfos.customer }}` }}"
          },
          {
            "type": "set_entity_info_from_template",
            "name": "test-resource-entity-infos-dictionary-1-key-2",
            "description": "Manager",
            "value": "{{ `{{ .Event.ExtraInfos.manager }}` }}"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-event-filter-che-event-filters-3-description",
      "enabled": true,
      "priority": 2
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-infos-dictionary-1",
      "connector_name": "test-connector-name-entity-infos-dictionary-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-entity-infos-dictionary-1",
      "resource": "test-resource-entity-infos-dictionary-1",
      "state": 0,
      "output": "test-output-entity-infos-dictionary-1",
      "customer": "test-resource-entity-infos-dictionary-1-val-1",
      "manager": "test-resource-entity-infos-dictionary-1-val-2"
    }
    """
    When I wait the end of event processing
    When I wait the next periodical process
    When I do GET /api/v4/entity-infos-dictionary/keys?search=test-resource-entity-infos-dictionary
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-resource-entity-infos-dictionary-1-key-1"
        },
        {
          "value": "test-resource-entity-infos-dictionary-1-key-2"
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
    When I do GET /api/v4/entity-infos-dictionary/values?key=test-resource-entity-infos-dictionary-1-key-1
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-resource-entity-infos-dictionary-1-val-1"
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
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-infos-dictionary-1",
      "connector_name": "test-connector-name-entity-infos-dictionary-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-entity-infos-dictionary-1",
      "resource": "test-resource-entity-infos-dictionary-1",
      "state": 0,
      "output": "test-output-entity-infos-dictionary-1",
      "customer": "test-resource-entity-infos-dictionary-1-val-1-updated",
      "manager": "test-resource-entity-infos-dictionary-1-val-2"
    }
    """
    When I wait the end of event processing
    When I wait the next periodical process
    When I do GET /api/v4/entity-infos-dictionary/keys?search=test-resource-entity-infos-dictionary
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-resource-entity-infos-dictionary-1-key-1"
        },
        {
          "value": "test-resource-entity-infos-dictionary-1-key-2"
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
    When I do GET /api/v4/entity-infos-dictionary/values?key=test-resource-entity-infos-dictionary-1-key-1
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-resource-entity-infos-dictionary-1-val-1-updated"
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
    When I do DELETE /api/v4/entitybasics?_id=test-resource-entity-infos-dictionary-1/test-component-entity-infos-dictionary-1
    Then the response code should be 204
    When I wait the next periodical process
    When I do GET /api/v4/entity-infos-dictionary/keys?search=test-resource-entity-infos-dictionary
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
