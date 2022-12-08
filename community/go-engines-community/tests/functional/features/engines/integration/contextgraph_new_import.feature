Feature: new-import entities
  I need to be able to new-import entities


  Scenario: given delete import action should delete service and resolve it's alarm
    When I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entityservice-new-import-partial-9-1",
      "name": "test-entityservice-new-import-partial-9-1-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-new-import-partial-9-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entityservice-new-import-partial-9-2",
      "name": "test-entityservice-new-import-partial-9-2-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-new-import-partial-9-1-name"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-new-import-partial-9-2"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "test-entityservice-new-import-partial-9-1"
      },
      {
        "event_type": "check",
        "component": "test-entityservice-new-import-partial-9-1"
      },
      {
        "event_type": "recomputeentityservice",
        "component": "test-entityservice-new-import-partial-9-2"
      },
      {
        "event_type": "check",
        "component": "test-entityservice-new-import-partial-9-2"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-new-import-partial-9",
      "connector_name": "test-connector-name-new-import-partial-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-9",
      "resource": "test-resource-new-import-partial-9-1",
      "state": 3,
      "output": "test-output-import-partial-9"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "component": "test-entityservice-new-import-partial-9-1"
      },
      {
        "event_type": "activate",
        "component": "test-entityservice-new-import-partial-9-2"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-new-import-partial-9",
      "connector_name": "test-connector-name-new-import-partial-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-9",
      "resource": "test-resource-new-import-partial-9-2",
      "state": 1,
      "output": "test-output-import-partial-9"
    }
    """
    When I do GET /api/v4/alarms?search=test-entityservice-new-import-partial-9-2&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-entityservice-new-import-partial-9-2",
            "state": {
              "val": 3
            }
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
    When I do PUT /api/v4/contextgraph-import-partial?source=test-new-import-partial-9-source:
    """json
    [
      {
        "action": "delete",
        "name": "test-entityservice-new-import-partial-9-1-name",
        "type": "service",
        "enabled": true
      }
    ]
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
       "status": "done"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "test-entityservice-new-import-partial-9-1"
      },
      {
        "event_type": "check",
        "component": "test-entityservice-new-import-partial-9-2"
      }
    ]
    """
    When I do GET /api/v4/entityservices/test-entityservice-new-import-partial-9-1
    Then the response code should be 404
    When I do GET /api/v4/entityservice-dependencies?_id=test-entityservice-new-import-partial-9-1
    Then the response code should be 404
    When I do GET /api/v4/entityservice-impacts?_id=test-entityservice-new-import-partial-9-1
    Then the response code should be 404
    When I do GET /api/v4/alarms?search=test-entityservice-new-import-partial-9-2&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-entityservice-new-import-partial-9-2",
            "state": {
              "val": 1
            }
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
    When I wait the next periodical process
    Then an entity test-entityservice-new-import-partial-9-1 should not be in the db
