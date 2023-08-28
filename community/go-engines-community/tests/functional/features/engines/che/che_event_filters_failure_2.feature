Feature: modify event on event filter
  I need to be able to modify event on event filter

  @concurrent
  Scenario: given check event and set_entity_info_from_template enrich action with invalid template should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-failure-second-1-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-failure-second-1"
            }
          },
          {
            "field": "event_type",
            "cond": {
              "type": "eq",
              "value": "check"
            }
          }
        ]
      ],
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "info1",
            "value": "{{ `{{ .Event.Output` }}"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      }
    }
    """
    Then the response code should be 201
    Then I save response ruleId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-second-1",
        "connector": "test-connector-che-event-filters-failure-second-1",
        "connector_name": "test-connector-name-che-event-filters-failure-second-1",
        "component": "test-component-che-event-filters-failure-second-1",
        "resource": "test-resource-che-event-filters-failure-second-1-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-second-1",
        "connector": "test-connector-che-event-filters-failure-second-1",
        "connector_name": "test-connector-name-che-event-filters-failure-second-1",
        "component": "test-component-che-event-filters-failure-second-1",
        "resource": "test-resource-che-event-filters-failure-second-1-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 1,
          "message": "invalid template \"Actions.0.Value\": template: tpl:1: unclosed action",
          "event": null,
          "unread": true
        },
        {
          "type": 1,
          "message": "invalid template \"Actions.0.Value\": template: tpl:1: unclosed action",
          "event": null,
          "unread": true
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

  @concurrent
  Scenario: given check event and set_entity_info_from_template enrich action with invalid template var should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-failure-second-2-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-failure-second-2"
            }
          },
          {
            "field": "event_type",
            "cond": {
              "type": "eq",
              "value": "check"
            }
          }
        ]
      ],
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "info1",
            "value": "{{ `{{ .Event.ExtraInfos.notfound }}` }}"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      }
    }
    """
    Then the response code should be 201
    Then I save response ruleId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-second-2",
        "connector": "test-connector-che-event-filters-failure-second-2",
        "connector_name": "test-connector-name-che-event-filters-failure-second-2",
        "component": "test-component-che-event-filters-failure-second-2",
        "resource": "test-resource-che-event-filters-failure-second-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-second-2",
        "connector": "test-connector-che-event-filters-failure-second-2",
        "connector_name": "test-connector-name-che-event-filters-failure-second-2",
        "component": "test-component-che-event-filters-failure-second-2",
        "resource": "test-resource-che-event-filters-failure-second-2-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 1,
          "message": "cannot execute template \"Actions.0.Value\" for event: template: tpl:1:9: executing \"tpl\" at <.Event.ExtraInfos.notfound>: map has no entry for key \"notfound\"",
          "event": {
            "event_type": "check",
            "connector": "test-connector-che-event-filters-failure-second-2",
            "connector_name": "test-connector-name-che-event-filters-failure-second-2",
            "component": "test-component-che-event-filters-failure-second-2",
            "source_type": "resource"
          },
          "unread": true
        },
        {
          "type": 1,
          "message": "cannot execute template \"Actions.0.Value\" for event: template: tpl:1:9: executing \"tpl\" at <.Event.ExtraInfos.notfound>: map has no entry for key \"notfound\"",
          "event": {
            "event_type": "check",
            "connector": "test-connector-che-event-filters-failure-second-2",
            "connector_name": "test-connector-name-che-event-filters-failure-second-2",
            "component": "test-component-che-event-filters-failure-second-2",
            "source_type": "resource"
          },
          "unread": true
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

  @standalone
  Scenario: given check event and rules with invalid patterns should create failure
    Given I am admin
    Then I save response ruleId1=test-event-filter-che-event-filters-failure-second-3-1
    Then I save response ruleId2=test-event-filter-che-event-filters-failure-second-3-2
    Then I save response ruleId3=test-event-filter-che-event-filters-failure-second-3-3
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-che-event-filters-failure-second-3",
      "connector": "test-connector-che-event-filters-failure-second-3",
      "connector_name": "test-connector-name-che-event-filters-failure-second-3",
      "component": "test-component-che-event-filters-failure-second-3",
      "resource": "test-resource-che-event-filters-failure-second-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId1 }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 0,
          "message": "invalid event pattern: invalid condition for \"event_type\" field: wrong condition value",
          "event": null,
          "unread": true
        },
        {
          "type": 0,
          "message": "invalid event pattern: invalid condition for \"event_type\" field: wrong condition value",
          "event": null,
          "unread": true
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
    When I do GET /api/v4/eventfilter/{{ .ruleId2 }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 0,
          "message": "invalid entity pattern: invalid condition for \"type\" field: wrong condition value",
          "event": null,
          "unread": true
        },
        {
          "type": 0,
          "message": "invalid entity pattern: invalid condition for \"type\" field: wrong condition value",
          "event": null,
          "unread": true
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
    When I do GET /api/v4/eventfilter/{{ .ruleId3 }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 0,
          "message": "invalid old pattern",
          "event": null,
          "unread": true
        },
        {
          "type": 0,
          "message": "invalid old pattern",
          "event": null,
          "unread": true
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
