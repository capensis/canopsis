Feature: modify event on event filter
  I need to be able to modify event on event filter

  @concurrent
  Scenario: given check event and drop event filter should drop event
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "drop",
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
            "value": "test-component-che-event-filters-1"
          }
        }
      ]],
      "description": "test-event-filter-che-event-filters-1-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-1",
      "connector_name": "test-connector-name-che-event-filters-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-1",
      "resource": "test-resource-che-event-filters-1",
      "state": 2,
      "output": "test-output-che-event-filters-1"
    }
    """
    When I do GET /api/v4/alarms?search=che-event-filters-1
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

  @concurrent
  Scenario: given check event and break event filter should not process next event filters
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "break",
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
            "value": "test-component-che-event-filters-2"
          }
        }
      ]],
      "description": "test-event-filter-che-event-filters-2-1-description",
      "priority": 1001,
      "enabled": true
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "drop",
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
            "value": "test-component-che-event-filters-2"
          }
        }
      ]],
      "description": "test-event-filter-che-event-filters-2-2-description",
      "priority": 1002,
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-2",
      "connector_name": "test-connector-name-che-event-filters-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-2",
      "resource": "test-resource-che-event-filters-2",
      "state": 2,
      "output": "test-output-che-event-filters-2"
    }
    """
    When I do GET /api/v4/entities?search=che-event-filters-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-event-filters-2"
        },
        {
          "_id": "test-connector-che-event-filters-2/test-connector-name-che-event-filters-2"
        },
        {
          "_id": "test-resource-che-event-filters-2/test-component-che-event-filters-2"
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
  Scenario: given check event and enrichment event filter with set_entity_info_from_template action
    should update event and entity
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
            "value": "test-component-che-event-filters-3"
          }
        }
      ]],
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "customer",
            "description": "Client",
            "value": "{{ `{{ .Event.ExtraInfos.customer }}` }}"
          },
          {
            "type": "set_entity_info_from_template",
            "name": "manager",
            "description": "Manager",
            "value": "{{ `{{ .Event.ExtraInfos.manager }}` }}"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-event-filter-che-event-filters-3-1-description",
      "enabled": true,
      "priority": 1003
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "event_pattern":[[
        {
          "field": "component",
          "cond": {
            "type": "eq",
            "value": "test-component-che-event-filters-3"
          }
        }
      ]],
      "config": {
        "actions": [
          {
            "type": "set_field_from_template",
            "name": "Output",
            "value": "{{ `{{ .Event.Output }}` }} (client: {{ `{{ .Event.Entity.Infos.customer.Value }}` }})"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-event-filter-che-event-filters-3-2-description",
      "enabled": true,
      "priority": 1004
    }
    """
    Then the response code should be 201
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
            "value": "test-component-che-event-filters-3"
          }
        }
      ]],
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "output",
            "description": "Output",
            "value": "{{ `{{ .Event.Output }}` }}"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-event-filter-che-event-filters-3-3-description",
      "enabled": true,
      "priority": 1005
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-3",
      "connector_name": "test-connector-name-che-event-filters-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-3",
      "resource": "test-resource-che-event-filters-3",
      "state": 2,
      "output": "test-output-che-event-filters-3",
      "customer": "test-customer-che-event-filters-3",
      "manager": "test-manager-che-event-filters-3"
    }
    """
    When I do GET /api/v4/entities?search=che-event-filters-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-event-filters-3"
        },
        {
          "_id": "test-connector-che-event-filters-3/test-connector-name-che-event-filters-3"
        },
        {
          "_id": "test-resource-che-event-filters-3/test-component-che-event-filters-3",
          "infos": {
            "customer": {
              "description": "Client",
              "name": "customer",
              "value": "test-customer-che-event-filters-3"
            },
            "manager": {
              "description": "Manager",
              "name": "manager",
              "value": "test-manager-che-event-filters-3"
            },
            "output": {
              "description": "Output",
              "name": "output",
              "value": "test-output-che-event-filters-3 (client: test-customer-che-event-filters-3)"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-3",
      "connector_name": "test-connector-name-che-event-filters-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-3",
      "resource": "test-resource-che-event-filters-3",
      "state": 2,
      "output": "test-output-che-event-filters-3",
      "customer": "test-customer-che-event-filters-3",
      "manager": "test-manager-che-event-filters-3-updated"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-che-event-filters-3/test-component-che-event-filters-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-resource-che-event-filters-3/test-component-che-event-filters-3",
      "infos": {
        "customer": {
          "description": "Client",
          "name": "customer",
          "value": "test-customer-che-event-filters-3"
        },
        "manager": {
          "description": "Manager",
          "name": "manager",
          "value": "test-manager-che-event-filters-3-updated"
        },
        "output": {
          "description": "Output",
          "name": "output",
          "value": "test-output-che-event-filters-3 (client: test-customer-che-event-filters-3)"
        }
      }
    }
    """

  @concurrent
  Scenario: given resource event should fill component infos
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-4"
            }
          },
          {
            "field": "source_type",
            "cond": {
              "type": "eq",
              "value": "component"
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
      "entity_pattern": [
        [
          {
            "field": "infos.customer",
            "cond": {
              "type": "exist",
              "value": false
            }
          }
        ]
      ],
      "config": {
        "actions": [
          {
            "type": "set_entity_info",
            "name": "customer",
            "description": "Client",
            "value": "test-customer-che-event-filters-4"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-event-filter-che-event-filters-4-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-4",
      "connector_name": "test-connector-name-che-event-filters-4",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-che-event-filters-4",
      "state": 2,
      "output": "test-output-che-event-filters-4"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-4",
      "connector_name": "test-connector-name-che-event-filters-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-4",
      "resource": "test-resource-che-event-filters-4",
      "state": 2,
      "output": "test-output-che-event-filters-4"
    }
    """
    When I do GET /api/v4/entities?search=che-event-filters-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-event-filters-4",
          "infos": {
            "customer": {
              "description": "Client",
              "name": "customer",
              "value": "test-customer-che-event-filters-4"
            }
          }
        },
        {
          "_id": "test-connector-che-event-filters-4/test-connector-name-che-event-filters-4"
        },
        {
          "_id": "test-resource-che-event-filters-4/test-component-che-event-filters-4",
          "component_infos": {
            "customer": {
              "description": "Client",
              "name": "customer",
              "value": "test-customer-che-event-filters-4"
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
  Scenario: given component event should fill component infos of resource entity
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-5"
            }
          },
          {
            "field": "source_type",
            "cond": {
              "type": "eq",
              "value": "component"
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
      "entity_pattern": [
        [
          {
            "field": "infos.customer",
            "cond": {
              "type": "exist",
              "value": false
            }
          }
        ]
      ],
      "config": {
        "actions": [
          {
            "type": "set_entity_info",
            "name": "customer",
            "description": "Client",
            "value": "test-customer-che-event-filters-5"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-event-filter-che-event-filters-5-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-5",
      "connector_name": "test-connector-name-che-event-filters-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-5",
      "resource": "test-resource-che-event-filters-5",
      "state": 2,
      "output": "test-output-che-event-filters-5"
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-che-event-filters-5",
      "connector_name": "test-connector-name-che-event-filters-5",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-che-event-filters-5",
      "state": 2,
      "output": "test-output-che-event-filters-5"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-che-event-filters-5",
        "connector_name": "test-connector-name-che-event-filters-5",
        "component": "test-component-che-event-filters-5",
        "source_type": "component"
      }
    ]
    """
    When I do GET /api/v4/entities?search=che-event-filters-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-event-filters-5",
          "infos": {
            "customer": {
              "description": "Client",
              "name": "customer",
              "value": "test-customer-che-event-filters-5"
            }
          }
        },
        {
          "_id": "test-connector-che-event-filters-5/test-connector-name-che-event-filters-5"
        },
        {
          "_id": "test-resource-che-event-filters-5/test-component-che-event-filters-5",
          "component_infos": {
            "customer": {
              "description": "Client",
              "name": "customer",
              "value": "test-customer-che-event-filters-5"
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
  Scenario: given check event and enrichment event filter should not update disabled entity
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-6",
      "connector_name": "test-connector-name-che-event-filters-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-6",
      "resource": "test-resource-che-event-filters-6",
      "state": 2,
      "output": "test-output-che-event-filters-6"
    }
    """
    When I do PUT /api/v4/entitybasics?_id=test-resource-che-event-filters-6/test-component-che-event-filters-6:
    """json
    {
      "enabled": false,
      "impact_level": 1,
      "sli_avail_state": 0,
      "infos": []
    }
    """
    Then the response code should be 200
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "entitytoggled",
      "connector": "test-connector-che-event-filters-6",
      "connector_name": "test-connector-name-che-event-filters-6",
      "component": "test-component-che-event-filters-6",
      "resource": "test-resource-che-event-filters-6",
      "source_type": "resource"
    }
    """
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-6"
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
      "entity_pattern": [
        [
          {
            "field": "infos.customer",
            "cond": {
              "type": "exist",
              "value": false
            }
          }
        ]
      ],
      "config": {
        "actions": [
          {
            "type": "set_entity_info",
            "name": "customer",
            "description": "Client",
            "value": "test-customer-che-event-filters-6"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-event-filter-che-event-filters-6-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-6",
      "connector_name": "test-connector-name-che-event-filters-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-6",
      "resource": "test-resource-che-event-filters-6",
      "state": 2,
      "output": "test-output-che-event-filters-6"
    }
    """
    When I do GET /api/v4/entities?search=test-resource-che-event-filters-6
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-resource-che-event-filters-6/test-component-che-event-filters-6",
          "infos": {},
          "name": "test-resource-che-event-filters-6",
          "type": "resource",
          "category": null,
          "component": "test-component-che-event-filters-6",
          "connector": "test-connector-che-event-filters-6/test-connector-name-che-event-filters-6",
          "enabled": false,
          "old_entity_patterns": null,
          "impact_level": 1,
          "impact_state": 0,
          "last_event_date": {{ (index .lastResponse.data 0).last_event_date }},
          "ko_events": 1,
          "ok_events": 0,
          "state": 0,
          "status": 0
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
  Scenario: given check event and enrichment event filter with set_entity_info action
  should update event and entity
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-7"
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
      "entity_pattern": [
        [
          {
            "field": "infos.testdate",
            "cond": {
              "type": "exist",
              "value": false
            }
          }
        ]
      ],
      "config": {
        "actions": [
          {
            "type": "set_entity_info",
            "name": "testdate",
            "description": "Date",
            "value": 1592215337
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-event-filter-che-event-filters-7-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-7",
      "connector_name": "test-connector-name-che-event-filters-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-7",
      "resource": "test-resource-che-event-filters-7",
      "state": 2,
      "output": "test-output-che-event-filters-7"
    }
    """
    When I do GET /api/v4/entities?search=che-event-filters-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-event-filters-7"
        },
        {
          "_id": "test-connector-che-event-filters-7/test-connector-name-che-event-filters-7"
        },
        {
          "_id": "test-resource-che-event-filters-7/test-component-che-event-filters-7",
          "infos": {
            "testdate": {
              "description": "Date",
              "name": "testdate",
              "value": 1592215337
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
  Scenario: given check event and enrichment event filter with copy_to_entity_info action
  should update event and entity
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-8"
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
      "entity_pattern": [
        [
          {
            "field": "infos.customer",
            "cond": {
              "type": "exist",
              "value": false
            }
          },
          {
            "field": "infos.testdate",
            "cond": {
              "type": "exist",
              "value": false
            }
          }
        ]
      ],
      "config": {
        "actions": [
          {
            "type": "copy_to_entity_info",
            "name": "customer",
            "description": "Client",
            "value": "Event.ExtraInfos.customer"
          },
          {
            "type": "copy_to_entity_info",
            "name": "testdate",
            "description": "Date",
            "value": "Event.ExtraInfos.testdate"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-event-filter-che-event-filters-8-1-description",
      "enabled": true,
      "priority": 1006
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-8"
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
            "name": "Output",
            "value": "{{ `{{ .Event.Output }}` }} (client: {{ `{{ .Event.Entity.Infos.customer.Value }}` }})"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-event-filter-che-event-filters-8-2-description",
      "enabled": true,
      "priority": 1007
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-8"
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
      "entity_pattern": [
        [
          {
            "field": "infos.output",
            "cond": {
              "type": "exist",
              "value": false
            }
          }
        ]
      ],
      "config": {
        "actions": [
          {
            "type": "copy_to_entity_info",
            "name": "output",
            "description": "Output",
            "value": "Event.Output"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "priority": 1008,
      "description": "test-event-filter-che-event-filters-8-3-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-8",
      "connector_name": "test-connector-name-che-event-filters-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-8",
      "resource": "test-resource-che-event-filters-8",
      "state": 2,
      "output": "test-output-che-event-filters-8",
      "customer": "test-customer-che-event-filters-8",
      "testdate": 1592215337
    }
    """
    When I do GET /api/v4/entities?search=che-event-filters-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-event-filters-8"
        },
        {
          "_id": "test-connector-che-event-filters-8/test-connector-name-che-event-filters-8"
        },
        {
          "_id": "test-resource-che-event-filters-8/test-component-che-event-filters-8",
          "infos": {
            "customer": {
              "description": "Client",
              "name": "customer",
              "value": "test-customer-che-event-filters-8"
            },
            "testdate": {
              "description": "Date",
              "name": "testdate",
              "value": 1592215337
            },
            "output": {
              "description": "Output",
              "name": "output",
              "value": "test-output-che-event-filters-8 (client: test-customer-che-event-filters-8)"
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
  Scenario: given check event and drop event filter with several event patterns should drop event
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "drop",
      "event_pattern":[
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-9"
            }
          }
        ],
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-resource-che-event-filters-9"
            }
          }
        ]
      ],
      "description": "test-event-filter-che-event-filters-9-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-9",
      "connector_name": "test-connector-name-che-event-filters-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-9-another",
      "resource": "test-resource-che-event-filters-9",
      "state": 2,
      "output": "test-output-che-event-filters-9"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-9",
      "connector_name": "test-connector-name-che-event-filters-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-9",
      "resource": "test-resource-che-event-filters-9-another",
      "state": 2,
      "output": "test-output-che-event-filters-9"
    }
    """
    When I do GET /api/v4/alarms?search=che-event-filters-9
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
