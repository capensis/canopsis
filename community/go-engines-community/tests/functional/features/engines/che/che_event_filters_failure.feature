Feature: modify event on event filter
  I need to be able to modify event on event filter

  @concurrent
  Scenario: given check event and invalid set_field enrich action should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-failure-1-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-failure-1"
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
            "type": "set_field",
            "name": "Resource",
            "value": 123
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
        "output": "test-output-che-event-filters-failure-1",
        "connector": "test-connector-che-event-filters-failure-1",
        "connector_name": "test-connector-name-che-event-filters-failure-1",
        "component": "test-component-che-event-filters-failure-1",
        "resource": "test-resource-che-event-filters-failure-1-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-1",
        "connector": "test-connector-che-event-filters-failure-1",
        "connector_name": "test-connector-name-che-event-filters-failure-1",
        "component": "test-component-che-event-filters-failure-1",
        "resource": "test-resource-che-event-filters-failure-1-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 4,
          "message": "action 0 cannot set \"Resource\" field: float64 value cannot be assigned to a string: 123",
          "event": null,
          "unread": true
        },
        {
          "type": 4,
          "message": "action 0 cannot set \"Resource\" field: float64 value cannot be assigned to a string: 123",
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
  Scenario: given check event and copy enrich action with invalid field path should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-failure-2-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-failure-2"
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
            "type": "copy",
            "name": "Resource",
            "value": "Event.ExtraInfos.notfound"
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
        "output": "test-output-che-event-filters-failure-2",
        "connector": "test-connector-che-event-filters-failure-2",
        "connector_name": "test-connector-name-che-event-filters-failure-2",
        "component": "test-component-che-event-filters-failure-2",
        "resource": "test-resource-che-event-filters-failure-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-2",
        "connector": "test-connector-che-event-filters-failure-2",
        "connector_name": "test-connector-name-che-event-filters-failure-2",
        "component": "test-component-che-event-filters-failure-2",
        "resource": "test-resource-che-event-filters-failure-2-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 4,
          "message": "action 0 cannot copy from \"Event.ExtraInfos.notfound\" to \"Resource\": field does not exist",
          "event": {
            "event_type": "check",
            "connector": "test-connector-che-event-filters-failure-2",
            "connector_name": "test-connector-name-che-event-filters-failure-2",
            "component": "test-component-che-event-filters-failure-2",
            "source_type": "resource"
          },
          "unread": true
        },
        {
          "type": 4,
          "message": "action 0 cannot copy from \"Event.ExtraInfos.notfound\" to \"Resource\": field does not exist",
          "event": {
            "event_type": "check",
            "connector": "test-connector-che-event-filters-failure-2",
            "connector_name": "test-connector-name-che-event-filters-failure-2",
            "component": "test-component-che-event-filters-failure-2",
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

  @concurrent
  Scenario: given check event and copy enrich action with invalid field value should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-failure-3-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-failure-3"
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
            "type": "copy",
            "name": "Resource",
            "value": "Event.State"
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
        "output": "test-output-che-event-filters-failure-3",
        "connector": "test-connector-che-event-filters-failure-3",
        "connector_name": "test-connector-name-che-event-filters-failure-3",
        "component": "test-component-che-event-filters-failure-3",
        "resource": "test-resource-che-event-filters-failure-3-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-3",
        "connector": "test-connector-che-event-filters-failure-3",
        "connector_name": "test-connector-name-che-event-filters-failure-3",
        "component": "test-component-che-event-filters-failure-3",
        "resource": "test-resource-che-event-filters-failure-3-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 4,
          "message": "action 0 cannot copy from \"Event.State\" to \"Resource\": types.CpsNumber value cannot be assigned to a string: 2",
          "event": {
            "event_type": "check",
            "connector": "test-connector-che-event-filters-failure-3",
            "connector_name": "test-connector-name-che-event-filters-failure-3",
            "component": "test-component-che-event-filters-failure-3",
            "source_type": "resource"
          },
          "unread": true
        },
        {
          "type": 4,
          "message": "action 0 cannot copy from \"Event.State\" to \"Resource\": types.CpsNumber value cannot be assigned to a string: 2",
          "event": {
            "event_type": "check",
            "connector": "test-connector-che-event-filters-failure-3",
            "connector_name": "test-connector-name-che-event-filters-failure-3",
            "component": "test-component-che-event-filters-failure-3",
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

  @concurrent
  Scenario: given check event and set_field_from_template enrich action with invalid template should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-failure-4-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-failure-4"
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
            "type": "set_field_from_template",
            "name": "Resource",
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
        "output": "test-output-che-event-filters-failure-4",
        "connector": "test-connector-che-event-filters-failure-4",
        "connector_name": "test-connector-name-che-event-filters-failure-4",
        "component": "test-component-che-event-filters-failure-4",
        "resource": "test-resource-che-event-filters-failure-4-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-4",
        "connector": "test-connector-che-event-filters-failure-4",
        "connector_name": "test-connector-name-che-event-filters-failure-4",
        "component": "test-component-che-event-filters-failure-4",
        "resource": "test-resource-che-event-filters-failure-4-2",
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
  Scenario: given check event and set_field_from_template enrich action with invalid template var should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-failure-5-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-failure-5"
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
            "type": "set_field_from_template",
            "name": "Resource",
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
        "output": "test-output-che-event-filters-failure-5",
        "connector": "test-connector-che-event-filters-failure-5",
        "connector_name": "test-connector-name-che-event-filters-failure-5",
        "component": "test-component-che-event-filters-failure-5",
        "resource": "test-resource-che-event-filters-failure-5-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-5",
        "connector": "test-connector-che-event-filters-failure-5",
        "connector_name": "test-connector-name-che-event-filters-failure-5",
        "component": "test-component-che-event-filters-failure-5",
        "resource": "test-resource-che-event-filters-failure-5-2",
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
            "connector": "test-connector-che-event-filters-failure-5",
            "connector_name": "test-connector-name-che-event-filters-failure-5",
            "component": "test-component-che-event-filters-failure-5",
            "source_type": "resource"
          },
          "unread": true
        },
        {
          "type": 1,
          "message": "cannot execute template \"Actions.0.Value\" for event: template: tpl:1:9: executing \"tpl\" at <.Event.ExtraInfos.notfound>: map has no entry for key \"notfound\"",
          "event": {
            "event_type": "check",
            "connector": "test-connector-che-event-filters-failure-5",
            "connector_name": "test-connector-name-che-event-filters-failure-5",
            "component": "test-component-che-event-filters-failure-5",
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

  @concurrent
  Scenario: given check event and set_field_from_template enrich action with invalid field var should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-failure-6-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-failure-6"
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
            "type": "set_field_from_template",
            "name": "State",
            "value": "{{ `{{ .Event.ExtraInfos.statestr }}` }}"
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
        "statestr": "2",
        "output": "test-output-che-event-filters-failure-6",
        "connector": "test-connector-che-event-filters-failure-6",
        "connector_name": "test-connector-name-che-event-filters-failure-6",
        "component": "test-component-che-event-filters-failure-6",
        "resource": "test-resource-che-event-filters-failure-6-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "statestr": "2",
        "output": "test-output-che-event-filters-failure-6",
        "connector": "test-connector-che-event-filters-failure-6",
        "connector_name": "test-connector-name-che-event-filters-failure-6",
        "component": "test-component-che-event-filters-failure-6",
        "resource": "test-resource-che-event-filters-failure-6-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 4,
          "message": "action 0 cannot set \"State\" field: string value cannot be converted to an integer: 2",
          "event": null,
          "unread": true
        },
        {
          "type": 4,
          "message": "action 0 cannot set \"State\" field: string value cannot be converted to an integer: 2",
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
  Scenario: given check event and invalid set_entity_info enrich action should create failure
    Given I am admin
    Then I save response ruleId=test-event-filter-che-event-filters-failure-7
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-7",
        "connector": "test-connector-che-event-filters-failure-7",
        "connector_name": "test-connector-name-che-event-filters-failure-7",
        "component": "test-component-che-event-filters-failure-7",
        "resource": "test-resource-che-event-filters-failure-7-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-7",
        "connector": "test-connector-che-event-filters-failure-7",
        "connector_name": "test-connector-name-che-event-filters-failure-7",
        "component": "test-component-che-event-filters-failure-7",
        "resource": "test-resource-che-event-filters-failure-7-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 4,
          "message": "action 0 cannot set \"info1\" entity info: invalid type of [1 2]",
          "event": null,
          "unread": true
        },
        {
          "type": 4,
          "message": "action 0 cannot set \"info1\" entity info: invalid type of [1 2]",
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
  Scenario: given check event and copy_to_entity_info enrich action with invalid field path should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-failure-8-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-failure-8"
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
            "type": "copy_to_entity_info",
            "name": "info1",
            "value": "Event.ExtraInfos.notfound"
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
        "output": "test-output-che-event-filters-failure-8",
        "connector": "test-connector-che-event-filters-failure-8",
        "connector_name": "test-connector-name-che-event-filters-failure-8",
        "component": "test-component-che-event-filters-failure-8",
        "resource": "test-resource-che-event-filters-failure-8-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-8",
        "connector": "test-connector-che-event-filters-failure-8",
        "connector_name": "test-connector-name-che-event-filters-failure-8",
        "component": "test-component-che-event-filters-failure-8",
        "resource": "test-resource-che-event-filters-failure-8-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 4,
          "message": "action 0 cannot copy from \"Event.ExtraInfos.notfound\" to \"info1\" entity info: field does not exist",
          "event": {
            "event_type": "check",
            "connector": "test-connector-che-event-filters-failure-8",
            "connector_name": "test-connector-name-che-event-filters-failure-8",
            "component": "test-component-che-event-filters-failure-8",
            "source_type": "resource"
          },
          "unread": true
        },
        {
          "type": 4,
          "message": "action 0 cannot copy from \"Event.ExtraInfos.notfound\" to \"info1\" entity info: field does not exist",
          "event": {
            "event_type": "check",
            "connector": "test-connector-che-event-filters-failure-8",
            "connector_name": "test-connector-name-che-event-filters-failure-8",
            "component": "test-component-che-event-filters-failure-8",
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

  @concurrent
  Scenario: given check event and copy_to_entity_info enrich action with invalid field value should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-failure-9-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-failure-9"
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
            "type": "copy_to_entity_info",
            "name": "info1",
            "value": "Event.ExtraInfos.valints"
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
        "output": "test-output-che-event-filters-failure-9",
        "valints": [1, 2],
        "connector": "test-connector-che-event-filters-failure-9",
        "connector_name": "test-connector-name-che-event-filters-failure-9",
        "component": "test-component-che-event-filters-failure-9",
        "resource": "test-resource-che-event-filters-failure-9-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-9",
        "valints": [1, 2],
        "connector": "test-connector-che-event-filters-failure-9",
        "connector_name": "test-connector-name-che-event-filters-failure-9",
        "component": "test-component-che-event-filters-failure-9",
        "resource": "test-resource-che-event-filters-failure-9-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 4,
          "message": "action 0 cannot copy from \"Event.ExtraInfos.valints\" to \"info1\" entity info: invalid type of [1 2]",
          "event": {
            "event_type": "check",
            "connector": "test-connector-che-event-filters-failure-9",
            "connector_name": "test-connector-name-che-event-filters-failure-9",
            "component": "test-component-che-event-filters-failure-9",
            "source_type": "resource"
          },
          "unread": true
        },
        {
          "type": 4,
          "message": "action 0 cannot copy from \"Event.ExtraInfos.valints\" to \"info1\" entity info: invalid type of [1 2]",
          "event": {
            "event_type": "check",
            "connector": "test-connector-che-event-filters-failure-9",
            "connector_name": "test-connector-name-che-event-filters-failure-9",
            "component": "test-component-che-event-filters-failure-9",
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
