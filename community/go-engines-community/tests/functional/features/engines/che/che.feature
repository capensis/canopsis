Feature: create entities on event
  I need to be able to create entities on event

  @concurrent
  Scenario: given resource check event should create entities
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-1",
      "connector_name": "test-connector-name-che-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-1",
      "resource": "test-resource-che-1",
      "state": 2,
      "output": "test-output-che-1"
    }
    """
    When I do GET /api/v4/entities?search=che-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-1",
          "category": null,
          "connector": "test-connector-che-1/test-connector-name-che-1",
          "component": "test-component-che-1",
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-component-che-1",
          "type": "component"
        },
        {
          "_id": "test-connector-che-1/test-connector-name-che-1",
          "category": null,
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-connector-name-che-1",
          "connector_type": "test-connector-che-1",
          "type": "connector"
        },
        {
          "_id": "test-resource-che-1/test-component-che-1",
          "category": null,
          "connector": "test-connector-che-1/test-connector-name-che-1",
          "component": "test-component-che-1",
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-resource-che-1",
          "type": "resource"
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
    When I do GET /api/v4/entities/context-graph?_id=test-component-che-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "depends": [
        "test-resource-che-1/test-component-che-1"
      ],
      "impact": [
        "test-connector-che-1/test-connector-name-che-1"
      ]
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-connector-che-1/test-connector-name-che-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "depends": [
        "test-component-che-1"
      ],
      "impact": [
        "test-resource-che-1/test-component-che-1"
      ]
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-1/test-component-che-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "depends": [
        "test-connector-che-1/test-connector-name-che-1"
      ],
      "impact": [
        "test-component-che-1"
      ]
    }
    """

  @concurrent
  Scenario: given component event should create entities
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-2",
      "connector_name": "test-connector-name-che-2",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-che-2",
      "state": 2,
      "output": "test-output-che-2"
    }
    """
    When I do GET /api/v4/entities?search=che-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-2",
          "category": null,
          "connector": "test-connector-che-2/test-connector-name-che-2",
          "component": "test-component-che-2",
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-component-che-2",
          "type": "component"
        },
        {
          "_id": "test-connector-che-2/test-connector-name-che-2",
          "category": null,
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-connector-name-che-2",
          "type": "connector"
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
    When I do GET /api/v4/entities/context-graph?_id=test-component-che-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "depends": [],
      "impact": [
        "test-connector-che-2/test-connector-name-che-2"
      ]
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-connector-che-2/test-connector-name-che-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "depends": [
        "test-component-che-2"
      ],
      "impact": []
    }
    """

  @concurrent
  Scenario: given event should update entities
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-3",
      "connector_name": "test-connector-name-che-3",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-che-3",
      "state": 2,
      "output": "test-output-che-3"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-3",
      "connector_name": "test-connector-name-che-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-3",
      "resource": "test-resource-che-3",
      "state": 2,
      "output": "test-output-che-3"
    }
    """
    When I do GET /api/v4/entities?search=che-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-3",
          "category": null,
          "connector": "test-connector-che-3/test-connector-name-che-3",
          "component": "test-component-che-3",
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-component-che-3",
          "type": "component"
        },
        {
          "_id": "test-connector-che-3/test-connector-name-che-3",
          "category": null,
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-connector-name-che-3",
          "type": "connector"
        },
        {
          "_id": "test-resource-che-3/test-component-che-3",
          "category": null,
          "connector": "test-connector-che-3/test-connector-name-che-3",
          "component": "test-component-che-3",
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-resource-che-3",
          "type": "resource"
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
    When I do GET /api/v4/entities/context-graph?_id=test-component-che-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "depends": [
        "test-resource-che-3/test-component-che-3"
      ],
      "impact": [
        "test-connector-che-3/test-connector-name-che-3"
      ]
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-connector-che-3/test-connector-name-che-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "depends": [
        "test-component-che-3"
      ],
      "impact": [
        "test-resource-che-3/test-component-che-3"
      ]
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-3/test-component-che-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "depends": [
        "test-connector-che-3/test-connector-name-che-3"
      ],
      "impact": [
        "test-component-che-3"
      ]
    }
    """

  @concurrent
  Scenario: given updated component by api should update resource component infos
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-4",
      "connector_name": "test-connector-name-che-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-4",
      "resource": "test-resource-che-4",
      "state": 2,
      "output": "test-output-che-4"
    }
    """
    When I do PUT /api/v4/entitybasics?_id=test-component-che-4:
    """json
    {
      "enabled": true,
      "impact_level": 1,
      "sli_avail_state": 1,
      "infos": [
        {
          "description": "test-component-che-4-info-1-description",
          "name": "test-component-che-4-info-1-name",
          "value": "test-component-che-4-info-1-value"
        }
      ]
    }
    """
    Then the response code should be 200
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entityupdated",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-che-4",
        "source_type": "component"
      }
    ]
    """
    When I do GET /api/v4/entities?search=che-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-4",
          "infos": {
            "test-component-che-4-info-1-name": {
              "description": "test-component-che-4-info-1-description",
              "name": "test-component-che-4-info-1-name",
              "value": "test-component-che-4-info-1-value"
            }
          }
        },
        {
          "_id": "test-connector-che-4/test-connector-name-che-4"
        },
        {
          "_id": "test-resource-che-4/test-component-che-4",
          "component_infos": {
            "test-component-che-4-info-1-name": {
              "description": "test-component-che-4-info-1-description",
              "name": "test-component-che-4-info-1-name",
              "value": "test-component-che-4-info-1-value"
            }
          }
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

  @concurrent
  Scenario: given resource event with a new connector should change connector field in the entity
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-5",
      "connector_name": "test-connector-name-che-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-5",
      "resource": "test-resource-che-5",
      "state": 2,
      "output": "test-output-che-5"
    }
    """
    When I do GET /api/v4/entities?search=test-resource-che-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-che-5/test-component-che-5",
          "category": null,
          "connector": "test-connector-che-5/test-connector-name-che-5",
          "component": "test-component-che-5",
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-resource-che-5",
          "type": "resource"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-5",
      "connector_name": "test-connector-name-che-5-new",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-5",
      "resource": "test-resource-che-5",
      "state": 2,
      "output": "test-output-che-5"
    }
    """
    When I do GET /api/v4/entities?search=test-resource-che-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-che-5/test-component-che-5",
          "category": null,
          "connector": "test-connector-che-5/test-connector-name-che-5-new",
          "component": "test-component-che-5",
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-resource-che-5",
          "type": "resource"
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

  @concurrent
  Scenario: given component event with a new connector should change connector field in the entity
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-6",
      "connector_name": "test-connector-name-che-6",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-che-6",
      "state": 2,
      "output": "test-output-che-6"
    }
    """
    When I do GET /api/v4/entities?search=test-component-che-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-6",
          "category": null,
          "connector": "test-connector-che-6/test-connector-name-che-6",
          "component": "test-component-che-6",
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-component-che-6",
          "type": "component"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-6",
      "connector_name": "test-connector-name-che-6-new",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-che-6",
      "state": 2,
      "output": "test-output-che-6"
    }
    """
    When I do GET /api/v4/entities?search=test-component-che-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-6",
          "category": null,
          "connector": "test-connector-che-6/test-connector-name-che-6-new",
          "component": "test-component-che-6",
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-component-che-6",
          "type": "component"
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

  @concurrent
  Scenario: given resources check events should create entities, context graph should be valid
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-context-graph-build-che-7",
      "connector_name": "test-connector-name-context-graph-build-che-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-context-graph-build-che-7-1",
      "resource": "test-resource-context-graph-build-che-7",
      "state": 2,
      "output": "test-context-graph-build-output-che-7"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-context-graph-build-che-7",
      "connector_name": "test-connector-name-context-graph-build-che-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-context-graph-build-che-7-2",
      "resource": "test-resource-context-graph-build-che-7",
      "state": 2,
      "output": "test-context-graph-build-output-che-7"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-context-graph-build-che-7",
      "connector_name": "test-connector-name-context-graph-build-che-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-context-graph-build-che-7-3",
      "resource": "test-resource-context-graph-build-che-7",
      "state": 2,
      "output": "test-context-graph-build-output-che-7"
    }
    """
    When I do GET /api/v4/entities?search=context-graph-build-che-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-context-graph-build-che-7-1",
          "category": null,
          "connector": "test-connector-context-graph-build-che-7/test-connector-name-context-graph-build-che-7",
          "component": "test-component-context-graph-build-che-7-1",
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-component-context-graph-build-che-7-1",
          "type": "component"
        },
        {
          "_id": "test-component-context-graph-build-che-7-2",
          "category": null,
          "connector": "test-connector-context-graph-build-che-7/test-connector-name-context-graph-build-che-7",
          "component": "test-component-context-graph-build-che-7-2",
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-component-context-graph-build-che-7-2",
          "type": "component"
        },
        {
          "_id": "test-component-context-graph-build-che-7-3",
          "category": null,
          "connector": "test-connector-context-graph-build-che-7/test-connector-name-context-graph-build-che-7",
          "component": "test-component-context-graph-build-che-7-3",
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-component-context-graph-build-che-7-3",
          "type": "component"
        },
        {
          "_id": "test-connector-context-graph-build-che-7/test-connector-name-context-graph-build-che-7",
          "category": null,
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-connector-name-context-graph-build-che-7",
          "connector_type": "test-connector-context-graph-build-che-7",
          "type": "connector"
        },
        {
          "_id": "test-resource-context-graph-build-che-7/test-component-context-graph-build-che-7-1",
          "category": null,
          "connector": "test-connector-context-graph-build-che-7/test-connector-name-context-graph-build-che-7",
          "component": "test-component-context-graph-build-che-7-1",
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-resource-context-graph-build-che-7",
          "type": "resource"
        },
        {
          "_id": "test-resource-context-graph-build-che-7/test-component-context-graph-build-che-7-2",
          "category": null,
          "connector": "test-connector-context-graph-build-che-7/test-connector-name-context-graph-build-che-7",
          "component": "test-component-context-graph-build-che-7-2",
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-resource-context-graph-build-che-7",
          "type": "resource"
        },
        {
          "_id": "test-resource-context-graph-build-che-7/test-component-context-graph-build-che-7-3",
          "category": null,
          "connector": "test-connector-context-graph-build-che-7/test-connector-name-context-graph-build-che-7",
          "component": "test-component-context-graph-build-che-7-3",
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-resource-context-graph-build-che-7",
          "type": "resource"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 7
      }
    }
    """
