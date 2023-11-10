Feature: update service on event
  I need to be able to see new service state on event

  Scenario: given entity service and event filters and new resource entity should update service alarm on resource event
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
            "value": "test-component-entityservice-eventfilters-1"
          }
        }
      ]],
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "customer",
            "description": "Client",
            "value": "{{ `{{ .Event.Component }}` }}"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-event-filter-entityservice-eventfilters-1-description",
      "enabled": true,
      "priority": 2
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-entityservice-eventfilters-1-name",
      "output_template": "Depends: {{ `{{ .Depends }}` }}; All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-entityservice-eventfilters-1"
            }
          },
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
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
    When I send an event:
    """json
    {
      "connector": "test-connector-entityservice-eventfilters-1",
      "connector_name": "test-connector-name-entityservice-eventfilters-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-entityservice-eventfilters-1",
      "resource": "test-resource-entityservice-eventfilters-1-1",
      "state": 1,
      "output": "test-output-entityservice-eventfilters-1"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-entityservice-eventfilters-1",
      "connector_name": "test-connector-name-entityservice-eventfilters-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-entityservice-eventfilters-1",
      "resource": "test-resource-entityservice-eventfilters-1-2",
      "state": 3,
      "output": "test-output-entityservice-eventfilters-1"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-entityservice-eventfilters-1",
      "connector_name": "test-connector-name-entityservice-eventfilters-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-entityservice-eventfilters-1",
      "resource": "test-resource-entityservice-eventfilters-1-3",
      "state": 2,
      "output": "test-output-entityservice-eventfilters-1"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search={{ .serviceID }}&with_dependencies=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "{{ .serviceID }}",
            "depends_count": 3,
            "impacts_count": 0
          },
          "v": {
            "component": "{{ .serviceID }}",
            "connector": "service",
            "connector_name": "service",
            "output": "Depends: 3; All: 3; Alarms: 3; Acknowledged: 0; NotAcknowledged: 3; StateCritical: 1; StateMajor: 1; StateMinor: 1; StateInfo: 0; Pbehaviors: map[];",
            "state": {
              "val": 3
            },
            "status": {
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        },
        "with_dependencies": true
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "Depends: 1; All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "Depends: 1; All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              },
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "Depends: 2; All: 2; Alarms: 2; Acknowledged: 0; NotAcknowledged: 2; StateCritical: 1; StateMajor: 0; StateMinor: 1; StateInfo: 0; Pbehaviors: map[];",
                "val": 3
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          },
          "entity": {
            "depends_count": 3,
            "impacts_count": 0
          }
        }
      }
    ]
    """
