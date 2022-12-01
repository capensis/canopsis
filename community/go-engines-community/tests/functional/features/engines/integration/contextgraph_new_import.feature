Feature: new-import entities
  I need to be able to new-import entities

  Scenario: given service and new entity by new-import should update service
    When I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-new-import-partial-1-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-component-new-import-partial-1"
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
    When I do PUT /api/v4/contextgraph-import-partial?source=test-new-import-partial-1-source:
    """json
    {
      "cis": [
        {
          "action": "set",
          "name": "test-component-new-import-partial-1",
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
    When I do GET /api/v4/entities/context-graph?_id=test-component-new-import-partial-1
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
        "test-component-new-import-partial-1"
      ],
      "impact": []
    }
    """

  Scenario: given service and updated entity by new-import should update service
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-2",
      "connector_name": "test-connector-name-new-import-partial-2",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-new-import-partial-2",
      "state": 1,
      "output": "test-output-new-import-partial-2"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-new-import-partial-2-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "infos.test-component-new-import-partial-2-infos-1",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-component-new-import-partial-2-infos-1-value"
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
    When I do PUT /api/v4/contextgraph-import-partial?source=test-new-import-partial-2-source:
    """json
    {
      "cis": [
        {
          "action": "set",
          "name": "test-component-new-import-partial-2",
          "type": "component",
          "infos": {
            "test-component-new-import-partial-2-infos-1": {
              "name": "test-component-new-import-partial-2-infos-1",
              "value": "test-component-new-import-partial-2-infos-1-value"
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
    When I do GET /api/v4/entities/context-graph?_id=test-component-new-import-partial-2
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [],
      "impact": [
        "test-connector-new-import-partial-2/test-connector-name-new-import-partial-2",
        "{{ .serviceID }}"
      ]
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-connector-new-import-partial-2/test-connector-name-new-import-partial-2
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-component-new-import-partial-2"
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
        "test-component-new-import-partial-2"
      ],
      "impact": []
    }
    """

  Scenario: given disabled entity by new-import should resolve alarm
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-3",
      "connector_name": "test-connector-name-new-import-partial-3",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-new-import-partial-3",
      "state": 1,
      "output": "test-output-new-import-partial-3"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-component-new-import-partial-3&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-new-import-partial-3"
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
    When I do PUT /api/v4/contextgraph-import-partial?source=test-new-import-partial-3-source:
    """json
    {
      "cis": [
        {
          "action": "disable",
          "name": "test-component-new-import-partial-3",
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
    When I do GET /api/v4/alarms?search=test-component-new-import-partial-3&opened=true
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

  Scenario: given deleted entity by new import should resolve alarm
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-4",
      "connector_name": "test-connector-name-new-import-partial-4",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-new-import-partial-4",
      "state": 1,
      "output": "test-output-new-import-partial-4"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-component-new-import-partial-4&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-new-import-partial-4"
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
    When I do PUT /api/v4/contextgraph-import-partial?source=test-new-import-partial-4-source:
    """json
    {
      "cis": [
        {
          "action": "delete",
          "name": "test-component-new-import-partial-4",
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
    When I do GET /api/v4/alarms?search=test-component-new-import-partial-4&opened=true
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

  Scenario: given deleted component by new import should resolve alarm for resources
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-5",
      "connector_name": "test-connector-name-new-import-partial-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-5",
      "resource": "test-resource-new-import-partial-5-1",
      "state": 1,
      "output": "test-output-new-import-partial-5"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-5",
      "connector_name": "test-connector-name-new-import-partial-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-5",
      "resource": "test-resource-new-import-partial-5-2",
      "state": 1,
      "output": "test-output-new-import-partial-5"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-new-import-partial-5-1&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-new-import-partial-5-1"
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
    When I do GET /api/v4/alarms?search=test-resource-new-import-partial-5-2&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-new-import-partial-5-2"
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
    When I do PUT /api/v4/contextgraph-import-partial?source=test-new-import-partial-5-source:
    """json
    {
      "cis": [
        {
          "action": "delete",
          "name": "test-component-new-import-partial-5",
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-new-import-partial-5-1&opened=true
    Then the response code should be 200
    Then the response body should contain:
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
    When I do GET /api/v4/alarms?search=test-resource-new-import-partial-5-2&opened=true
    Then the response code should be 200
    Then the response body should contain:
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
