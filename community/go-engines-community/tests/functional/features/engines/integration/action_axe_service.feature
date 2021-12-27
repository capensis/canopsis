Feature: update service when alarm is updated by action
  I need to be able to update service when action is applied to alarm.

  Scenario: given meta alarm and scenario should update meta alarm and update children
    Given I am admin
    When I do POST /api/v4/entityservices:
    """
    {
      "name": "test-entityservice-action-axe-service-1-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"name": "test-resource-action-axe-service-1"}],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/scenarios:
    """
    {
      "name": "test-scenario-action-axe-service-1-name",
      "enabled": true,
      "priority": 93,
      "triggers": ["cancel"],
      "actions": [
        {
          "entity_patterns": [
            {
              "name": "test-resource-action-axe-service-1"
            }
          ],
          "type": "changestate",
          "parameters": {
            "output": "test-output-action-axe-service-1",
            "state": 3
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-connector-action-axe-service-1",
      "connector_name": "test-connector-name-action-axe-service-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-action-axe-service-1",
      "resource": "test-resource-action-axe-service-1",
      "state": 1,
      "output": "test-output-action-axe-service-1"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-connector-action-axe-service-1",
      "connector_name": "test-connector-name-action-axe-service-1",
      "source_type": "resource",
      "event_type": "cancel",
      "component":  "test-component-action-axe-service-1",
      "resource": "test-resource-action-axe-service-1",
      "output": "test-output-action-axe-service-1"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"entity._id":"{{ .serviceID }}"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "component": "{{ .serviceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            },
            "steps": [
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              },
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 3
              }
            ]
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
