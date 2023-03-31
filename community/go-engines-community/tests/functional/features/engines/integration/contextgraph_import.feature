Feature: Import entities
  I need to be able to import entities

  Scenario: given service and new entity by import should update service
    When I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-import-partial-1-name",
      "output_template": "test-entityservice-import-partial-1-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-component-import-partial-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do PUT /api/v4/contextgraph/import-partial?source=test-import-partial-1-source:
    """json
    {
      "cis": [
        {
          "action": "create",
          "_id": "test-component-import-partial-1",
          "name": "test-component-import-partial-1",
          "type": "component",
          "enabled": true
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
       "status": "done"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/entities/context-graph?_id=test-component-import-partial-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [],
      "impact": [
        "{{ .serviceID }}"
      ]
    }
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-component-import-partial-1"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/weather-services/{{ .serviceID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-import-partial-1",
          "import_source": "test-import-partial-1-source"
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
    Then the response key "data.0.imported" should be greater than 0

  Scenario: given service and updated entity by import should update service
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-import-partial-2",
      "connector_name": "test-connector-name-import-partial-2",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-import-partial-2",
      "state": 1,
      "output": "test-output-import-partial-2"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-import-partial-2-name",
      "output_template": "test-entityservice-import-partial-2-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "infos.test-component-import-partial-2-infos-1",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-component-import-partial-2-infos-1-value"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do PUT /api/v4/contextgraph/import-partial?source=test-import-partial-2-source:
    """json
    {
      "cis": [
        {
          "action": "set",
          "_id": "test-component-import-partial-2",
          "name": "test-component-import-partial-2",
          "type": "component",
          "infos": {
            "test-component-import-partial-2-infos-1": {
              "name": "test-component-import-partial-2-infos-1",
              "value": "test-component-import-partial-2-infos-1-value"
            }
          },
          "enabled": true
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
       "status": "done"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities/context-graph?_id=test-component-import-partial-2
    Then the response array key "impact" should contain:
    """json
    [
      "test-connector-import-partial-2/test-connector-name-import-partial-2",
      "{{ .serviceID }}"
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id=test-connector-import-partial-2/test-connector-name-import-partial-2
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-component-import-partial-2"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-component-import-partial-2"
      ],
      "impact": []
    }
    """

  Scenario: given disabled entity by import should resolve alarm
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-import-partial-3",
      "connector_name": "test-connector-name-import-partial-3",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-import-partial-3",
      "state": 1,
      "output": "test-output-import-partial-3"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-component-import-partial-3&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-import-partial-3"
          }
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
    When I do PUT /api/v4/contextgraph/import-partial?source=test-import-partial-3-source:
    """json
    {
      "cis": [
        {
          "action": "disable",
          "_id": "test-component-import-partial-3",
          "name": "test-component-import-partial-3",
          "type": "component",
          "enabled": true
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
       "status": "done"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-component-import-partial-3&opened=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
