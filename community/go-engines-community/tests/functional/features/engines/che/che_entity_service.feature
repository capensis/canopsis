Feature: create service entity
  I need to be able to create service entity

  Scenario: given resource entity and new service entity should add resource to service on service creation
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-first",
      "connector_name": "test-connector-name-che-service-first",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-first",
      "resource": "test-resource-che-service-first",
      "state": 2,
      "output": "test-output-che-service-first"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-first-name",
      "output_template": "test-entityservice-che-service-first-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-che-service-first"
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
    When I do GET /api/v4/entities?search=che-service-first&sort_by=name&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-first",
          "depends_count": 0,
          "impacts_count": 0
        },
        {
          "_id": "test-connector-che-service-first/test-connector-name-che-service-first",
          "depends_count": 0,
          "impacts_count": 0
        },
        {
          "_id": "{{ .serviceID }}",
          "category": null,
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-entityservice-che-service-first-name",
          "type": "service",
          "depends_count": 1,
          "impacts_count": 0
        },
        {
          "_id": "test-resource-che-service-first/test-component-che-service-first",
          "category": null,
          "component": "test-component-che-service-first",
          "enabled": true,
          "impact_level": 1,
          "infos": {},
          "name": "test-resource-che-service-first",
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
        "test-resource-che-service-first/test-component-che-service-first"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-first/test-component-che-service-first
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-first/test-connector-name-che-service-first"
      ],
      "impact": [
        "test-component-che-service-first",
        "{{ .serviceID }}"
      ]
    }
    """

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
    When I wait the end of 2 events processing
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
    When I wait the end of 2 events processing
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
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-2/test-connector-name-che-service-2"
      ],
      "impact": [
        "test-component-che-service-2",
        "{{ .serviceID }}"
      ]
    }
    """

  Scenario: given service entity with updated pattern should remove now unmatched and add now matched entities on service update
    Given I am admin
    When I send an event:
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
    When I wait the end of event processing
    When I send an event:
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
    When I wait the end of event processing
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
    When I wait the end of 2 events processing
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
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-3/test-connector-name-che-service-3"
      ],
      "impact": [
        "test-component-che-service-3",
        "{{ .serviceID }}"
      ]
    }
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
    When I wait the end of 2 events processing
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
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-3/test-connector-name-che-service-3"
      ],
      "impact": [
        "test-component-che-service-3",
        "{{ .serviceID }}"
      ]
    }
    """

  Scenario: given service entity and resource entity and enrichment event filter should add resource to service on infos update
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
            "field": "infos.manager",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-manager-che-service-6"
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
              "value": "test-resource-che-service-6"
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
            "value": "test-manager-che-service-6"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-eventfilter-che-service-6-description",
      "enabled": true,
      "priority": 2
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-6/test-component-che-service-6"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-6/test-component-che-service-6
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-6/test-connector-name-che-service-6"
      ],
      "impact": [
        "test-component-che-service-6",
        "{{ .serviceID }}"
      ]
    }
    """

  Scenario: given service entity and resource entity and updated enrichment event filter should remove resource from service on infos update
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
            "field": "infos.manager",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-manager-che-service-7"
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
              "value": "test-resource-che-service-7"
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
            "value": "test-manager-che-service-7"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-eventfilter-che-service-7-description",
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-7/test-component-che-service-7"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-7/test-component-che-service-7
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-7/test-connector-name-che-service-7"
      ],
      "impact": [
        "test-component-che-service-7",
        "{{ .serviceID }}"
      ]
    }
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
              "value": "test-resource-che-service-7"
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
            "value": "test-another-manager-che-service-7"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-eventfilter-che-service-7-description",
      "enabled": true,
      "priority": 2
    }
    """
    Then the response code should be 200
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-7",
      "connector_name": "test-connector-name-che-service-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-7",
      "resource": "test-resource-che-service-7",
      "state": 1,
      "output": "test-output-che-service-7"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-7/test-component-che-service-7
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-7/test-connector-name-che-service-7"
      ],
      "impact": [
        "test-component-che-service-7"
      ]
    }
    """

  Scenario: given service entity and new component entity on resource event should add component to service
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
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-component-che-service-8"
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
    When I wait the end of event processing
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-component-che-service-8"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-component-che-service-8
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-8/test-component-che-service-8"
      ],
      "impact": [
        "test-connector-che-service-8/test-connector-name-che-service-8",
        "{{ .serviceID }}"
      ]
    }
    """

  Scenario: given service entity and new connector entity on resource event should add connector to service
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-9-name",
      "output_template": "test-entityservice-che-service-9-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-connector-name-che-service-9"
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
    When I wait the end of event processing
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-9/test-connector-name-che-service-9"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-connector-che-service-9/test-connector-name-che-service-9
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-component-che-service-9"
      ],
      "impact": [
        "test-resource-che-service-9/test-component-che-service-9",
        "{{ .serviceID }}"
      ]
    }
    """

  Scenario: given service entity and updated resource entity by api should add resource to service on infos update
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-16-name",
      "output_template": "test-entityservice-che-service-16-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "infos.manager",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-manager-che-service-16"
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
      "connector": "test-connector-che-service-16",
      "connector_name": "test-connector-name-che-service-16",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-16",
      "resource": "test-resource-che-service-16",
      "state": 2,
      "output": "test-output-che-service-16"
    }
    """
    When I wait the end of event processing
    When I do PUT /api/v4/entitybasics?_id=test-resource-che-service-16/test-component-che-service-16:
    """json
    {
      "enabled": true,
      "impact_level": 1,
      "infos": [
        {
          "description": "Manager",
          "name": "manager",
          "value": "test-manager-che-service-16"
        }
      ],
      "impact": [
        "test-component-che-service-16"
      ],
      "depends": [
        "test-connector-che-service-16/test-connector-name-che-service-16"
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-16/test-component-che-service-16"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-16/test-component-che-service-16
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-16/test-connector-name-che-service-16"
      ],
      "impact": [
        "test-component-che-service-16",
        "{{ .serviceID }}"
      ]
    }
    """

  Scenario: given deleted service entity should remove service from all links
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-17-name-1",
      "impact_level": 1,
      "output_template": "test-entityservice-che-service-17-output-1",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-che-service-17"
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
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-17-name-2",
      "impact_level": 1,
      "output_template": "test-entityservice-che-service-17-output-2",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-che-service-17-name-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response impactServiceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-17",
      "connector_name": "test-connector-name-che-service-17",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-17",
      "resource": "test-resource-che-service-17",
      "state": 2,
      "output": "test-output-che-service-17"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-17/test-component-che-service-17"
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
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-17/test-component-che-service-17
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-17/test-connector-name-che-service-17"
      ],
      "impact": [
        "test-component-che-service-17",
        "{{ .serviceID }}"
      ]
    }
    """
    When I do DELETE /api/v4/entityservices/{{ .serviceID }}
    Then the response code should be 204
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities/context-graph?_id={{ .impactServiceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-17/test-component-che-service-17
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-17/test-connector-name-che-service-17"
      ],
      "impact": [
        "test-component-che-service-17"
      ]
    }
    """

  Scenario: given disabled service entity should not update service context graph
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-18-name",
      "impact_level": 1,
      "output_template": "test-entityservice-che-service-18-output",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-che-service-18-1",
                "test-resource-che-service-18-2"
              ]
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
      "connector": "test-connector-che-service-18",
      "connector_name": "test-connector-name-che-service-18",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-18",
      "resource": "test-resource-che-service-18-1",
      "state": 2,
      "output": "test-output-che-service-18"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-18-1/test-component-che-service-18"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-18-1/test-component-che-service-18
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-18/test-connector-name-che-service-18"
      ],
      "impact": [
        "test-component-che-service-18",
        "{{ .serviceID }}"
      ]
    }
    """
    When I do PUT /api/v4/entityservices/{{ .serviceID }}:
    """json
    {
      "name": "test-entityservice-che-service-18-name",
      "impact_level": 1,
      "output_template": "test-entityservice-che-service-18-output",
      "enabled": false,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-che-service-18"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-18",
      "connector_name": "test-connector-name-che-service-18",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-18",
      "resource": "test-resource-che-service-18-2",
      "state": 1,
      "output": "test-output-che-service-18"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-18-1/test-component-che-service-18"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-18-1/test-component-che-service-18
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-18/test-connector-name-che-service-18"
      ],
      "impact": [
        "test-component-che-service-18",
        "{{ .serviceID }}"
      ]
    }
    """

  Scenario: given service with old pattern should update service
    Given I am admin
    When I do PUT /api/v4/entityservices/test-entityservice-che-service-19:
    """json
    {
      "name": "test-entityservice-che-service-19-name",
      "output_template": "test-entityservice-che-service-19-output",
      "impact_level": 1,
      "enabled": true,
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-19",
      "connector_name": "test-connector-name-che-service-19",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-19",
      "resource": "test-resource-che-service-19",
      "state": 2,
      "output": "test-output-che-service-19"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities/context-graph?_id=test-entityservice-che-service-19
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-19/test-component-che-service-19"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-19/test-component-che-service-19
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-19/test-connector-name-che-service-19"
      ],
      "impact": [
        "test-component-che-service-19",
        "test-entityservice-che-service-19"
      ]
    }
    """

  Scenario: given service with corporate pattern should update service on pattern update
    Given I am admin
    When I do POST /api/v4/patterns:
    """json
    {
      "title": "test-pattern-che-service-20",
      "type": "entity",
      "is_corporate": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-che-service-20-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response patternID={{ .lastResponse._id }}
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-20-name",
      "output_template": "test-entityservice-che-service-20-output",
      "impact_level": 1,
      "enabled": true,
      "sli_avail_state": 0,
      "corporate_entity_pattern": "{{ .patternID }}"
    }
    """
    Then the response code should be 201
    Then I save response serviceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-che-service-20",
        "connector_name": "test-connector-name-che-service-20",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-che-service-20",
        "resource": "test-resource-che-service-20-1",
        "state": 2,
        "output": "test-output-che-service-20"
      },
      {
        "connector": "test-connector-che-service-20",
        "connector_name": "test-connector-name-che-service-20",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-che-service-20",
        "resource": "test-resource-che-service-20-2",
        "state": 3,
        "output": "test-output-che-service-20"
      }
    ]
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-20-1/test-component-che-service-20"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-20-1/test-component-che-service-20
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-20/test-connector-name-che-service-20"
      ],
      "impact": [
        "test-component-che-service-20",
        "{{ .serviceID }}"
      ]
    }
    """
    When I do PUT /api/v4/patterns/{{ .patternID }}:
    """json
    {
      "title": "test-pattern-che-service-20",
      "type": "entity",
      "is_corporate": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-che-service-20-1",
                "test-resource-che-service-20-2"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-20-1/test-component-che-service-20",
        "test-resource-che-service-20-2/test-component-che-service-20"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-20-1/test-component-che-service-20
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-20/test-connector-name-che-service-20"
      ],
      "impact": [
        "test-component-che-service-20",
        "{{ .serviceID }}"
      ]
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-20-2/test-component-che-service-20
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-20/test-connector-name-che-service-20"
      ],
      "impact": [
        "test-component-che-service-20",
        "{{ .serviceID }}"
      ]
    }
    """

  Scenario: given resource entity and new service entity should update context graph on entity disable or enable
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-21",
      "connector_name": "test-connector-name-che-service-21",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-21",
      "resource": "test-resource-che-service-21",
      "state": 0,
      "output": "test-output-che-service-21"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entityservice-che-service-21",
      "name": "test-entityservice-che-service-21-name",
      "output_template": "test-entityservice-che-service-21-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-che-service-21"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities/context-graph?_id=test-entityservice-che-service-21
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-21/test-component-che-service-21"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-21/test-component-che-service-21
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-21/test-connector-name-che-service-21"
      ],
      "impact": [
        "test-component-che-service-21",
        "test-entityservice-che-service-21"
      ]
    }
    """
    When I do PUT /api/v4/entitybasics?_id=test-resource-che-service-21/test-component-che-service-21:
    """json
    {
      "enabled": false,
      "impact_level": 1,
      "impact": [
        "test-component-che-service-21"
      ],
      "depends": [
        "test-connector-che-service-21/test-connector-name-che-service-21"
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/entities/context-graph?_id=test-entityservice-che-service-21
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-21/test-component-che-service-21
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-21/test-connector-name-che-service-21"
      ],
      "impact": [
        "test-component-che-service-21"
      ]
    }
    """
    When I do PUT /api/v4/entitybasics?_id=test-resource-che-service-21/test-component-che-service-21:
    """json
    {
      "enabled": true,
      "impact_level": 1,
      "impact": [
        "test-component-che-service-21"
      ],
      "depends": [
        "test-connector-che-service-21/test-connector-name-che-service-21"
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/entities/context-graph?_id=test-entityservice-che-service-21
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-21/test-component-che-service-21"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-21/test-component-che-service-21
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-21/test-connector-name-che-service-21"
      ],
      "impact": [
        "test-component-che-service-21",
        "test-entityservice-che-service-21"
      ]
    }
    """

  Scenario: given resource entity and new service entity should update context graph on entity mass disable or enable
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-22",
      "connector_name": "test-connector-name-che-service-22",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-22",
      "resource": "test-resource-che-service-22-1",
      "state": 0,
      "output": "test-output-che-service-22"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-22",
      "connector_name": "test-connector-name-che-service-22",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-22",
      "resource": "test-resource-che-service-22-2",
      "state": 0,
      "output": "test-output-che-service-22"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entityservice-che-service-22",
      "name": "test-entityservice-che-service-22-name",
      "output_template": "test-entityservice-che-service-22-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-che-service-22-1",
                "test-resource-che-service-22-2"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities/context-graph?_id=test-entityservice-che-service-22
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-22-1/test-component-che-service-22",
        "test-resource-che-service-22-2/test-component-che-service-22"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-22-1/test-component-che-service-22
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-22/test-connector-name-che-service-22"
      ],
      "impact": [
        "test-component-che-service-22",
        "test-entityservice-che-service-22"
      ]
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-22-2/test-component-che-service-22
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-22/test-connector-name-che-service-22"
      ],
      "impact": [
        "test-component-che-service-22",
        "test-entityservice-che-service-22"
      ]
    }
    """
    When I do PUT /api/v4/bulk/entities/disable:
    """json
    [
      {
        "_id": "test-resource-che-service-22-1/test-component-che-service-22"
      },
      {
        "_id": "test-resource-che-service-22-2/test-component-che-service-22"
      }
    ]
    """
    Then the response code should be 207
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities/context-graph?_id=test-entityservice-che-service-22
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-22-1/test-component-che-service-22
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-22/test-connector-name-che-service-22"
      ],
      "impact": [
        "test-component-che-service-22"
      ]
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-22-2/test-component-che-service-22
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-22/test-connector-name-che-service-22"
      ],
      "impact": [
        "test-component-che-service-22"
      ]
    }
    """
    When I do PUT /api/v4/bulk/entities/enable:
    """json
    [
      {
        "_id": "test-resource-che-service-22-1/test-component-che-service-22"
      },
      {
        "_id": "test-resource-che-service-22-2/test-component-che-service-22"
      }
    ]
    """
    Then the response code should be 207
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities/context-graph?_id=test-entityservice-che-service-22
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-22-1/test-component-che-service-22",
        "test-resource-che-service-22-2/test-component-che-service-22"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-22-1/test-component-che-service-22
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-22/test-connector-name-che-service-22"
      ],
      "impact": [
        "test-component-che-service-22",
        "test-entityservice-che-service-22"
      ]
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-22-2/test-component-che-service-22
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-22/test-connector-name-che-service-22"
      ],
      "impact": [
        "test-component-che-service-22",
        "test-entityservice-che-service-22"
      ]
    }
    """
