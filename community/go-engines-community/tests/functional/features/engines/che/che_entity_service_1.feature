Feature: create service entity
  I need to be able to create service entity

  @concurrent
  Scenario: given resource entity and new service entity should add resource to service on service creation
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-service-1",
      "connector_name": "test-connector-name-che-service-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-1",
      "resource": "test-resource-che-service-1",
      "state": 2,
      "output": "test-output-che-service-1"
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-1-name",
      "output_template": "test-entityservice-che-service-1-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-che-service-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I do GET /api/v4/entities?search=che-service-1&sort_by=name&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-1",
          "depends_count": 0,
          "impacts_count": 0
        },
        {
          "_id": "test-connector-che-service-1/test-connector-name-che-service-1",
          "depends_count": 0,
          "impacts_count": 0
        },
        {
          "_id": "{{ .serviceID }}",
          "category": null,
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-entityservice-che-service-1-name",
          "type": "service",
          "depends_count": 1,
          "impacts_count": 0
        },
        {
          "_id": "test-resource-che-service-1/test-component-che-service-1",
          "category": null,
          "component": "test-component-che-service-1",
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-resource-che-service-1",
          "type": "resource",
          "depends_count": 0,
          "impacts_count": 1
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 4
      }
    }
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-1/test-component-che-service-1"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-1/test-component-che-service-1
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-1/test-connector-name-che-service-1"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-1",
      "{{ .serviceID }}"
    ]
    """

  @concurrent
  Scenario: given service entity and new resource entity should add resource to service on resource creation
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-2-name",
      "output_template": "test-entityservice-che-service-2-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-che-service-2"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-2",
      "connector_name": "test-connector-name-che-service-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-2",
      "resource": "test-resource-che-service-2",
      "state": 2,
      "output": "test-output-che-service-2"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-che-service-2",
        "connector_name": "test-connector-name-che-service-2",
        "component": "test-component-che-service-2",
        "resource": "test-resource-che-service-2",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I do GET /api/v4/entities?search=che-service-2&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-2"
        },
        {
          "_id": "test-connector-che-service-2/test-connector-name-che-service-2"
        },
        {
          "_id": "{{ .serviceID }}",
          "category": null,
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-entityservice-che-service-2-name",
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-2/test-component-che-service-2",
          "category": null,
          "component": "test-component-che-service-2",
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-resource-che-service-2",
          "type": "resource"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 4
      }
    }
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-2/test-component-che-service-2"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-2/test-component-che-service-2
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-2/test-connector-name-che-service-2"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-2",
      "{{ .serviceID }}"
    ]
    """

  @concurrent
  Scenario: given service entity with updated pattern should remove now unmatched and add now matched entities on service update
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-service-3",
      "connector_name": "test-connector-name-che-service-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-3",
      "resource": "test-resource-che-service-3-1",
      "state": 2,
      "output": "test-output-che-service-3"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-service-3",
      "connector_name": "test-connector-name-che-service-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-3",
      "resource": "test-resource-che-service-3-2",
      "state": 2,
      "output": "test-output-che-service-3"
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-3-name",
      "output_template": "test-entityservice-che-service-3-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-che-service-3-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-3-1/test-component-che-service-3"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-3-1/test-component-che-service-3
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-3/test-connector-name-che-service-3"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-3",
      "{{ .serviceID }}"
    ]
    """
    When I do PUT /api/v4/entityservices/{{ .serviceID }}:
    """json
    {
      "name": "test-entityservice-che-service-3-name",
      "output_template": "test-entityservice-che-service-3-output",
      "category": null,
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-che-service-3-2"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      }
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-3-2/test-component-che-service-3"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-3-1/test-component-che-service-3
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-3/test-connector-name-che-service-3"
      ],
      "impact": [
        "test-component-che-service-3"
      ]
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-3-2/test-component-che-service-3
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-3/test-connector-name-che-service-3"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-3",
      "{{ .serviceID }}"
    ]
    """

  @concurrent
  Scenario: given service entity and resource entity and enrichment event filter should add resource to service on infos update
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-4-name",
      "output_template": "test-entityservice-che-service-4-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "infos.manager",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-manager-che-service-4"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      }
    ]
    """
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-resource-che-service-4"
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
            "field": "infos.manager",
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
            "type": "set_entity_info_from_template",
            "name": "manager",
            "description": "Manager",
            "value": "test-manager-che-service-4"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-eventfilter-che-service-4-description",
      "enabled": true,
      "priority": 2
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-4",
      "connector_name": "test-connector-name-che-service-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-4",
      "resource": "test-resource-che-service-4",
      "state": 2,
      "output": "test-output-che-service-4"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-che-service-4",
        "connector_name": "test-connector-name-che-service-4",
        "component": "test-component-che-service-4",
        "resource": "test-resource-che-service-4",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-4/test-component-che-service-4"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-4/test-component-che-service-4
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-4/test-connector-name-che-service-4"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-4",
      "{{ .serviceID }}"
    ]
    """

  @concurrent
  Scenario: given service entity and resource entity and updated enrichment event filter should remove resource from service on infos update
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-5-name",
      "output_template": "test-entityservice-che-service-5-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "infos.manager",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-manager-che-service-5"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      }
    ]
    """
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-resource-che-service-5"
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
            "field": "infos.manager",
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
            "type": "set_entity_info_from_template",
            "name": "manager",
            "description": "Manager",
            "value": "test-manager-che-service-5"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-eventfilter-che-service-5-description",
      "enabled": true,
      "priority": 2
    }
    """
    Then the response code should be 201
    When I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-5",
      "connector_name": "test-connector-name-che-service-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-5",
      "resource": "test-resource-che-service-5",
      "state": 2,
      "output": "test-output-che-service-5"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-che-service-5",
        "connector_name": "test-connector-name-che-service-5",
        "component": "test-component-che-service-5",
        "resource": "test-resource-che-service-5",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-5/test-component-che-service-5"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-5/test-component-che-service-5
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-5/test-connector-name-che-service-5"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-5",
      "{{ .serviceID }}"
    ]
    """
    When I do PUT /api/v4/eventfilter/rules/{{ .ruleID }}:
    """json
    {
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-resource-che-service-5"
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
            "name": "manager",
            "description": "Manager",
            "value": "test-another-manager-che-service-5"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-eventfilter-che-service-5-description",
      "enabled": true,
      "priority": 2
    }
    """
    Then the response code should be 200
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-5",
      "connector_name": "test-connector-name-che-service-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-5",
      "resource": "test-resource-che-service-5",
      "state": 1,
      "output": "test-output-che-service-5"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "check",
        "connector": "test-connector-che-service-5",
        "connector_name": "test-connector-name-che-service-5",
        "component": "test-component-che-service-5",
        "resource": "test-resource-che-service-5",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      }
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-5/test-component-che-service-5
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-5/test-connector-name-che-service-5"
      ],
      "impact": [
        "test-component-che-service-5"
      ]
    }
    """

  @concurrent
  Scenario: given service entity and new component entity on resource event should add component to service
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-6-name",
      "output_template": "test-entityservice-che-service-6-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-component-che-service-6"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-service-6",
      "connector_name": "test-connector-name-che-service-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-6",
      "resource": "test-resource-che-service-6",
      "state": 2,
      "output": "test-output-che-service-6"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-component-che-service-6"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-component-che-service-6
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-resource-che-service-6/test-component-che-service-6"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-connector-che-service-6/test-connector-name-che-service-6",
      "{{ .serviceID }}"
    ]
    """

  @concurrent
  Scenario: given service entity and new connector entity on resource event should add connector to service
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-7-name",
      "output_template": "test-entityservice-che-service-7-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-connector-name-che-service-7"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-service-7",
      "connector_name": "test-connector-name-che-service-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-7",
      "resource": "test-resource-che-service-7",
      "state": 2,
      "output": "test-output-che-service-7"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-7/test-connector-name-che-service-7"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-connector-che-service-7/test-connector-name-che-service-7
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-component-che-service-7"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-resource-che-service-7/test-component-che-service-7",
      "{{ .serviceID }}"
    ]
    """

  @concurrent
  Scenario: given service entity and updated resource entity by api should add resource to service on infos update
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-8-name",
      "output_template": "test-entityservice-che-service-8-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "infos.manager",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-manager-che-service-8"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-service-8",
      "connector_name": "test-connector-name-che-service-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-8",
      "resource": "test-resource-che-service-8",
      "state": 2,
      "output": "test-output-che-service-8"
    }
    """
    When I do PUT /api/v4/entitybasics?_id=test-resource-che-service-8/test-component-che-service-8:
    """json
    {
      "enabled": true,
      "impact_level": 1,
      "infos": [
        {
          "description": "Manager",
          "name": "manager",
          "value": "test-manager-che-service-8"
        }
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entityupdated",
        "connector": "test-connector-che-service-8",
        "connector_name": "test-connector-name-che-service-8",
        "component": "test-component-che-service-8",
        "resource": "test-resource-che-service-8",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-8/test-component-che-service-8"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-8/test-component-che-service-8
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-8/test-connector-name-che-service-8"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-8",
      "{{ .serviceID }}"
    ]
    """

  @concurrent
  Scenario: given deleted service entity should remove service from all links
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-9-name-1",
      "impact_level": 1,
      "output_template": "test-entityservice-che-service-9-output-1",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-che-service-9"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      }
    ]
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-9-name-2",
      "impact_level": 1,
      "output_template": "test-entityservice-che-service-9-output-2",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-che-service-9-name-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response impactServiceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .impactServiceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .impactServiceID }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-9",
      "connector_name": "test-connector-name-che-service-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-9",
      "resource": "test-resource-che-service-9",
      "state": 2,
      "output": "test-output-che-service-9"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-che-service-9",
        "connector_name": "test-connector-name-che-service-9",
        "component": "test-component-che-service-9",
        "resource": "test-resource-che-service-9",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .impactServiceID }}"
      }
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-9/test-component-che-service-9"
      ],
      "impact": [
        "{{ .impactServiceID }}"
      ]
    }
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .impactServiceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "{{ .serviceID }}"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-9/test-component-che-service-9
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-9/test-connector-name-che-service-9"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-9",
      "{{ .serviceID }}"
    ]
    """
    When I do DELETE /api/v4/entityservices/{{ .serviceID }}
    Then the response code should be 204
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .impactServiceID }}"
      }
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .impactServiceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-9/test-component-che-service-9
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-9/test-connector-name-che-service-9"
      ],
      "impact": [
        "test-component-che-service-9"
      ]
    }
    """
