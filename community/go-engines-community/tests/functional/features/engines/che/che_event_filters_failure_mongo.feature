Feature: modify event on event filter
  I need to be able to modify event on event filter

  @concurrent
  Scenario: given check event and enrich from external select mongo data with invalid query should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-failure-mongo-1-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-failure-mongo-1"
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
      "external_data": {
        "component": {
          "type": "mongo",
          "sort_by": "status",
          "sort": "asc",
          "select": {
            "$cus$tomer": "{{ `{{.Event.Component}}` }}"
          },
          "collection": "eventfilter_mongo_data"
        }
      },
      "config": {
        "actions": [
          {
            "type": "set_field_from_template",
            "name": "status",
            "value": "{{ `{{.ExternalData.component.status}}` }}"
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
        "output": "test-output-che-event-filters-failure-mongo-1",
        "connector": "test-connector-che-event-filters-failure-mongo-1",
        "connector_name": "test-connector-name-che-event-filters-failure-mongo-1",
        "component": "test-component-che-event-filters-failure-mongo-1",
        "resource": "test-resource-che-event-filters-failure-mongo-1-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-mongo-1",
        "connector": "test-connector-che-event-filters-failure-mongo-1",
        "connector_name": "test-connector-name-che-event-filters-failure-mongo-1",
        "component": "test-component-che-event-filters-failure-mongo-1",
        "resource": "test-resource-che-event-filters-failure-mongo-1-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 2,
          "message": "external data \"component\" has invalid query: (BadValue) unknown top level operator: $cus$tomer. If you have a field name that starts with a '$' symbol, consider using $getField or $setField.",
          "event": null,
          "unread": true
        },
        {
          "type": 2,
          "message": "external data \"component\" has invalid query: (BadValue) unknown top level operator: $cus$tomer. If you have a field name that starts with a '$' symbol, consider using $getField or $setField.",
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
  Scenario: given check event and enrich from external select mongo data with not found row should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-failure-mongo-2-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-failure-mongo-2"
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
      "external_data": {
        "component": {
          "type": "mongo",
          "sort_by": "status",
          "sort": "asc",
          "select": {
            "customer": "{{ `{{.Event.Component}}` }}"
          },
          "collection": "eventfilter_mongo_data"
        }
      },
      "config": {
        "actions": [
          {
            "type": "set_field_from_template",
            "name": "status",
            "value": "{{ `{{.ExternalData.component.status}}` }}"
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
        "output": "test-output-che-event-filters-failure-mongo-2",
        "connector": "test-connector-che-event-filters-failure-mongo-2",
        "connector_name": "test-connector-name-che-event-filters-failure-mongo-2",
        "component": "test-component-che-event-filters-failure-mongo-2",
        "resource": "test-resource-che-event-filters-failure-mongo-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-mongo-2",
        "connector": "test-connector-che-event-filters-failure-mongo-2",
        "connector_name": "test-connector-name-che-event-filters-failure-mongo-2",
        "component": "test-component-che-event-filters-failure-mongo-2",
        "resource": "test-resource-che-event-filters-failure-mongo-2-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 2,
          "message": "external data \"component\" cannot be find by event",
          "event": {
            "event_type": "check",
            "connector": "test-connector-che-event-filters-failure-mongo-2",
            "connector_name": "test-connector-name-che-event-filters-failure-mongo-2",
            "component": "test-component-che-event-filters-failure-mongo-2",
            "source_type": "resource"
          },
          "unread": true
        },
        {
          "type": 2,
          "message": "external data \"component\" cannot be find by event",
          "event": {
            "event_type": "check",
            "connector": "test-connector-che-event-filters-failure-mongo-2",
            "connector_name": "test-connector-name-che-event-filters-failure-mongo-2",
            "component": "test-component-che-event-filters-failure-mongo-2",
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
  Scenario: given check event and enrich from external regexp mongo data with invalid query should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-failure-mongo-3-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-failure-mongo-3"
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
      "external_data": {
        "component": {
          "type": "mongo",
          "sort_by": "status",
          "sort": "asc",
          "select": {
            "$cus$tomer": "{{ `{{.Event.Component}}` }}"
          },
          "regexp": {
            "message": "{{ `{{.Event.Output}}` }}"
          },
          "collection": "eventfilter_mongo_data_regexp"
        }
      },
      "config": {
        "actions": [
          {
            "type": "set_field_from_template",
            "name": "status",
            "value": "{{ `{{.ExternalData.component.status}}` }}"
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
        "output": "test-output-che-event-filters-failure-mongo-3",
        "connector": "test-connector-che-event-filters-failure-mongo-3",
        "connector_name": "test-connector-name-che-event-filters-failure-mongo-3",
        "component": "test-component-che-event-filters-failure-mongo-3",
        "resource": "test-resource-che-event-filters-failure-mongo-3-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-mongo-3",
        "connector": "test-connector-che-event-filters-failure-mongo-3",
        "connector_name": "test-connector-name-che-event-filters-failure-mongo-3",
        "component": "test-component-che-event-filters-failure-mongo-3",
        "resource": "test-resource-che-event-filters-failure-mongo-3-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 2,
          "message": "external data \"component\" has invalid query: (BadValue) unknown top level operator: $cus$tomer. If you have a field name that starts with a '$' symbol, consider using $getField or $setField.",
          "event": null,
          "unread": true
        },
        {
          "type": 2,
          "message": "external data \"component\" has invalid query: (BadValue) unknown top level operator: $cus$tomer. If you have a field name that starts with a '$' symbol, consider using $getField or $setField.",
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
  Scenario: given check event and enrich from external regexp mongo data with not found row should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-failure-mongo-4-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-connector-che-event-filters-failure-mongo-4"
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
      "external_data": {
        "component": {
          "type": "mongo",
          "sort_by": "status",
          "sort": "asc",
          "select": {
            "customer": "{{ `{{.Event.Component}}` }}"
          },
          "regexp": {
            "message": "{{ `{{.Event.Output}}` }}"
          },
          "collection": "eventfilter_mongo_data_regexp"
        }
      },
      "config": {
        "actions": [
          {
            "type": "set_field_from_template",
            "name": "status",
            "value": "{{ `{{.ExternalData.component.status}}` }}"
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
        "output": "test-output-che-event-filters-failure-mongo-4",
        "connector": "test-connector-che-event-filters-failure-mongo-4",
        "connector_name": "test-connector-name-che-event-filters-failure-mongo-4",
        "component": "test-eventfilter-mongo-data-regexp-1-customer",
        "resource": "test-resource-che-event-filters-failure-mongo-4-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-mongo-4",
        "connector": "test-connector-che-event-filters-failure-mongo-4",
        "connector_name": "test-connector-name-che-event-filters-failure-mongo-4",
        "component": "test-eventfilter-mongo-data-regexp-1-customer",
        "resource": "test-resource-che-event-filters-failure-mongo-4-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 2,
          "message": "external data \"component\" cannot be find by event",
          "event": {
            "event_type": "check",
            "connector": "test-connector-che-event-filters-failure-mongo-4",
            "connector_name": "test-connector-name-che-event-filters-failure-mongo-4",
            "component": "test-eventfilter-mongo-data-regexp-1-customer",
            "source_type": "resource"
          },
          "unread": true
        },
        {
          "type": 2,
          "message": "external data \"component\" cannot be find by event",
          "event": {
            "event_type": "check",
            "connector": "test-connector-che-event-filters-failure-mongo-4",
            "connector_name": "test-connector-name-che-event-filters-failure-mongo-4",
            "component": "test-eventfilter-mongo-data-regexp-1-customer",
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
  Scenario: given check event and enrich from external regexp mongo data with invalid regexp should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-failure-mongo-5-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-connector-che-event-filters-failure-mongo-5"
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
      "external_data": {
        "component": {
          "type": "mongo",
          "sort_by": "status",
          "sort": "asc",
          "select": {
            "customer": "{{ `{{.Event.Component}}` }}"
          },
          "regexp": {
            "message": "{{ `{{.Event.Output}}` }}"
          },
          "collection": "eventfilter_mongo_data_regexp"
        }
      },
      "config": {
        "actions": [
          {
            "type": "set_field_from_template",
            "name": "status",
            "value": "{{ `{{.ExternalData.component.status}}` }}"
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
        "output": "test-output-che-event-filters-failure-mongo-5",
        "connector": "test-connector-che-event-filters-failure-mongo-5",
        "connector_name": "test-connector-name-che-event-filters-failure-mongo-5",
        "component": "test-eventfilter-mongo-data-regexp-4-customer",
        "resource": "test-resource-che-event-filters-failure-mongo-5-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-mongo-5",
        "connector": "test-connector-che-event-filters-failure-mongo-5",
        "connector_name": "test-connector-name-che-event-filters-failure-mongo-5",
        "component": "test-eventfilter-mongo-data-regexp-4-customer",
        "resource": "test-resource-che-event-filters-failure-mongo-5-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 2,
          "message": "external data \"component\" has invalid regexp in collection field \"message\": error parsing regexp: missing closing ): `test-eventfilter-mongo-data-regexp-1-(.*`",
          "event": null,
          "unread": true
        },
        {
          "type": 2,
          "message": "external data \"component\" has invalid regexp in collection field \"message\": error parsing regexp: missing closing ): `test-eventfilter-mongo-data-regexp-1-(.*`",
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
  Scenario: given check event and enrich from external select mongo data with invalid template should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-failure-mongo-6-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-failure-mongo-6"
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
      "external_data": {
        "component": {
          "type": "mongo",
          "sort_by": "status",
          "sort": "asc",
          "select": {
            "customer": "{{ `{{.Event.Component` }}"
          },
          "collection": "eventfilter_mongo_data"
        }
      },
      "config": {
        "actions": [
          {
            "type": "set_field_from_template",
            "name": "status",
            "value": "{{ `{{.ExternalData.component.status}}` }}"
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
        "output": "test-output-che-event-filters-failure-mongo-6",
        "connector": "test-connector-che-event-filters-failure-mongo-6",
        "connector_name": "test-connector-name-che-event-filters-failure-mongo-6",
        "component": "test-component-che-event-filters-failure-mongo-6",
        "resource": "test-resource-che-event-filters-failure-mongo-6-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-mongo-6",
        "connector": "test-connector-che-event-filters-failure-mongo-6",
        "connector_name": "test-connector-name-che-event-filters-failure-mongo-6",
        "component": "test-component-che-event-filters-failure-mongo-6",
        "resource": "test-resource-che-event-filters-failure-mongo-6-2",
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
          "message": "invalid template \"ExternalData.component.Select.customer\": template: tpl:1: unclosed action",
          "event": null,
          "unread": true
        },
        {
          "type": 1,
          "message": "invalid template \"ExternalData.component.Select.customer\": template: tpl:1: unclosed action",
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
  Scenario: given check event and enrich from external regexp mongo data with invalid template should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-failure-mongo-7-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-failure-mongo-7"
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
      "external_data": {
        "component": {
          "type": "mongo",
          "sort_by": "status",
          "sort": "asc",
          "select": {
            "customer": "{{ `{{.Event.Component}}` }}"
          },
          "regexp": {
            "message": "{{ `{{.Event.Output` }}"
          },
          "collection": "eventfilter_mongo_data_regexp"
        }
      },
      "config": {
        "actions": [
          {
            "type": "set_field_from_template",
            "name": "status",
            "value": "{{ `{{.ExternalData.component.status}}` }}"
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
        "output": "test-output-che-event-filters-failure-mongo-7",
        "connector": "test-connector-che-event-filters-failure-mongo-7",
        "connector_name": "test-connector-name-che-event-filters-failure-mongo-7",
        "component": "test-component-che-event-filters-failure-mongo-7",
        "resource": "test-resource-che-event-filters-failure-mongo-7-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-mongo-7",
        "connector": "test-connector-che-event-filters-failure-mongo-7",
        "connector_name": "test-connector-name-che-event-filters-failure-mongo-7",
        "component": "test-component-che-event-filters-failure-mongo-7",
        "resource": "test-resource-che-event-filters-failure-mongo-7-2",
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
          "message": "invalid template \"ExternalData.component.Regexp.message\": template: tpl:1: unclosed action",
          "event": null,
          "unread": true
        },
        {
          "type": 1,
          "message": "invalid template \"ExternalData.component.Regexp.message\": template: tpl:1: unclosed action",
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
  Scenario: given check event and enrich from external select mongo data with not found template var should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-failure-mongo-8-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-failure-mongo-8"
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
      "external_data": {
        "component": {
          "type": "mongo",
          "sort_by": "status",
          "sort": "asc",
          "select": {
            "customer": "{{ `{{.Event.ExtraInfos.notfound}}` }}"
          },
          "collection": "eventfilter_mongo_data"
        }
      },
      "config": {
        "actions": [
          {
            "type": "set_field_from_template",
            "name": "status",
            "value": "{{ `{{.ExternalData.component.status}}` }}"
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
        "output": "test-output-che-event-filters-failure-mongo-8",
        "connector": "test-connector-che-event-filters-failure-mongo-8",
        "connector_name": "test-connector-name-che-event-filters-failure-mongo-8",
        "component": "test-component-che-event-filters-failure-mongo-8",
        "resource": "test-resource-che-event-filters-failure-mongo-8-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-mongo-8",
        "connector": "test-connector-che-event-filters-failure-mongo-8",
        "connector_name": "test-connector-name-che-event-filters-failure-mongo-8",
        "component": "test-component-che-event-filters-failure-mongo-8",
        "resource": "test-resource-che-event-filters-failure-mongo-8-2",
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
          "message": "cannot execute template \"ExternalData.component.Select.customer\" for event: template: tpl:1:8: executing \"tpl\" at <.Event.ExtraInfos.notfound>: map has no entry for key \"notfound\"",
          "event": {
            "event_type": "check",
            "connector": "test-connector-che-event-filters-failure-mongo-8",
            "connector_name": "test-connector-name-che-event-filters-failure-mongo-8",
            "component": "test-component-che-event-filters-failure-mongo-8",
            "source_type": "resource"
          },
          "unread": true
        },
        {
          "type": 1,
          "message": "cannot execute template \"ExternalData.component.Select.customer\" for event: template: tpl:1:8: executing \"tpl\" at <.Event.ExtraInfos.notfound>: map has no entry for key \"notfound\"",
          "event": {
            "event_type": "check",
            "connector": "test-connector-che-event-filters-failure-mongo-8",
            "connector_name": "test-connector-name-che-event-filters-failure-mongo-8",
            "component": "test-component-che-event-filters-failure-mongo-8",
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
  Scenario: given check event and enrich from external regexp mongo data with not found template var should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-failure-mongo-9-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-failure-mongo-9"
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
      "external_data": {
        "component": {
          "type": "mongo",
          "sort_by": "status",
          "sort": "asc",
          "select": {
            "customer": "{{ `{{.Event.Component}}` }}"
          },
          "regexp": {
            "message": "{{ `{{.Event.ExtraInfos.notfound}}` }}"
          },
          "collection": "eventfilter_mongo_data_regexp"
        }
      },
      "config": {
        "actions": [
          {
            "type": "set_field_from_template",
            "name": "status",
            "value": "{{ `{{.ExternalData.component.status}}` }}"
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
        "output": "test-output-che-event-filters-failure-mongo-9",
        "connector": "test-connector-che-event-filters-failure-mongo-9",
        "connector_name": "test-connector-name-che-event-filters-failure-mongo-9",
        "component": "test-component-che-event-filters-failure-mongo-9",
        "resource": "test-resource-che-event-filters-failure-mongo-9-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-che-event-filters-failure-mongo-9",
        "connector": "test-connector-che-event-filters-failure-mongo-9",
        "connector_name": "test-connector-name-che-event-filters-failure-mongo-9",
        "component": "test-component-che-event-filters-failure-mongo-9",
        "resource": "test-resource-che-event-filters-failure-mongo-9-2",
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
          "message": "cannot execute template \"ExternalData.component.Regexp.message\" for event: template: tpl:1:8: executing \"tpl\" at <.Event.ExtraInfos.notfound>: map has no entry for key \"notfound\"",
          "event": {
            "event_type": "check",
            "connector": "test-connector-che-event-filters-failure-mongo-9",
            "connector_name": "test-connector-name-che-event-filters-failure-mongo-9",
            "component": "test-component-che-event-filters-failure-mongo-9",
            "source_type": "resource"
          },
          "unread": true
        },
        {
          "type": 1,
          "message": "cannot execute template \"ExternalData.component.Regexp.message\" for event: template: tpl:1:8: executing \"tpl\" at <.Event.ExtraInfos.notfound>: map has no entry for key \"notfound\"",
          "event": {
            "event_type": "check",
            "connector": "test-connector-che-event-filters-failure-mongo-9",
            "connector_name": "test-connector-name-che-event-filters-failure-mongo-9",
            "component": "test-component-che-event-filters-failure-mongo-9",
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
