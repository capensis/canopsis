Feature: modify event on event filter
  I need to be able to modify event on event filter

  Scenario: given check event and drop event filter should drop event
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "drop",
      "patterns": [{
        "event_type": "check",
        "component": "test-component-che-event-filters-1"
      }],
      "description": "test-event-filter-che-event-filters-1-description",
      "priority": 1,
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """
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
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=che-event-filters-1
    Then the response code should be 200
    Then the response body should be:
    """
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

  Scenario: given check event and break event filter should not process next event filters
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "break",
      "patterns": [{
        "event_type": "check",
        "component": "test-component-che-event-filters-2"
      }],
      "description": "test-event-filter-che-event-filters-2-description",
      "priority": 1,
      "enabled": true
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "drop",
      "patterns": [{
        "event_type": "check",
        "component": "test-component-che-event-filters-2"
      }],
      "description": "test-event-filter-che-event-filters-1-description",
      "priority": 2,
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """
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
    When I save response createTimestamp={{ now }}
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=che-event-filters-2
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-component-che-event-filters-2",
          "category": null,
          "component": "test-component-che-event-filters-2",
          "depends": [
            "test-resource-che-event-filters-2/test-component-che-event-filters-2"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-connector-che-event-filters-2/test-connector-name-che-event-filters-2"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-component-che-event-filters-2",
          "type": "component"
        },
        {
          "_id": "test-connector-che-event-filters-2/test-connector-name-che-event-filters-2",
          "category": null,
          "depends": [
            "test-component-che-event-filters-2"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-resource-che-event-filters-2/test-component-che-event-filters-2"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-connector-name-che-event-filters-2",
          "type": "connector"
        },
        {
          "_id": "test-resource-che-event-filters-2/test-component-che-event-filters-2",
          "category": null,
          "component": "test-component-che-event-filters-2",
          "depends": [
            "test-connector-che-event-filters-2/test-connector-name-che-event-filters-2"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-component-che-event-filters-2"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-resource-che-event-filters-2",
          "type": "resource"
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

  Scenario: given check event and enrichment event filter with set_entity_info_from_template action
    should update event and entity
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "enrichment",
      "patterns": [{
        "event_type": "check",
        "component": "test-component-che-event-filters-3",
        "current_entity": {
          "infos": {
            "customer": null,
            "manager": null
          }
        }
      }],
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
      "description": "test-event-filter-che-event-filters-3-description",
      "enabled": true,
      "priority": 2
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "enrichment",
      "patterns": [{
        "event_type": "check",
        "component": "test-component-che-event-filters-3"
      }],
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
      "description": "test-event-filter-che-event-filters-3-description",
      "enabled": true,
      "priority": 3
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "enrichment",
      "patterns": [{
        "event_type": "check",
        "component": "test-component-che-event-filters-3",
        "current_entity": {
          "infos": {"output": null}
        }
      }],
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
      "description": "test-event-filter-che-event-filters-3-description",
      "enabled": true,
      "priority": 4
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """
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
    When I save response createTimestamp={{ now }}
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=che-event-filters-3
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-component-che-event-filters-3",
          "category": null,
          "component": "test-component-che-event-filters-3",
          "depends": [
            "test-resource-che-event-filters-3/test-component-che-event-filters-3"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-connector-che-event-filters-3/test-connector-name-che-event-filters-3"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-component-che-event-filters-3",
          "type": "component"
        },
        {
          "_id": "test-connector-che-event-filters-3/test-connector-name-che-event-filters-3",
          "category": null,
          "depends": [
            "test-component-che-event-filters-3"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-resource-che-event-filters-3/test-component-che-event-filters-3"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-connector-name-che-event-filters-3",
          "type": "connector"
        },
        {
          "_id": "test-resource-che-event-filters-3/test-component-che-event-filters-3",
          "category": null,
          "component": "test-component-che-event-filters-3",
          "depends": [
            "test-connector-che-event-filters-3/test-connector-name-che-event-filters-3"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-component-che-event-filters-3"
          ],
          "impact_level": 1,
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
          },
          "measurements": null,
          "name": "test-resource-che-event-filters-3",
          "type": "resource"
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

  Scenario: given resource event should fill component infos
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "enrichment",
      "patterns": [{
        "event_type": "check",
        "source_type": "component",
        "component": "test-component-che-event-filters-4",
        "current_entity": {
          "infos": {"customer": null}
        }
      }],
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
      "enabled": true,
      "priority": 2
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """
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
    When I save response createComponentTimestamp={{ now }}
    When I wait the end of event processing
    When I send an event:
    """
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
    When I save response createResourceTimestamp={{ now }}
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=che-event-filters-4
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-component-che-event-filters-4",
          "category": null,
          "component": "test-component-che-event-filters-4",
          "depends": [
            "test-resource-che-event-filters-4/test-component-che-event-filters-4"
          ],
          "enable_history": [
            {{ .createComponentTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-connector-che-event-filters-4/test-connector-name-che-event-filters-4"
          ],
          "impact_level": 1,
          "infos": {
            "customer": {
              "description": "Client",
              "name": "customer",
              "value": "test-customer-che-event-filters-4"
            }
          },
          "measurements": null,
          "name": "test-component-che-event-filters-4",
          "type": "component"
        },
        {
          "_id": "test-connector-che-event-filters-4/test-connector-name-che-event-filters-4",
          "category": null,
          "depends": [
            "test-component-che-event-filters-4"
          ],
          "enable_history": [
            {{ .createComponentTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-resource-che-event-filters-4/test-component-che-event-filters-4"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-connector-name-che-event-filters-4",
          "type": "connector"
        },
        {
          "_id": "test-resource-che-event-filters-4/test-component-che-event-filters-4",
          "category": null,
          "component": "test-component-che-event-filters-4",
          "component_infos": {
            "customer": {
              "description": "Client",
              "name": "customer",
              "value": "test-customer-che-event-filters-4"
            }
          },
          "depends": [
            "test-connector-che-event-filters-4/test-connector-name-che-event-filters-4"
          ],
          "enable_history": [
            {{ .createResourceTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-component-che-event-filters-4"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-resource-che-event-filters-4",
          "type": "resource"
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

  Scenario: given component event should fill component infos of resource entity
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "enrichment",
      "patterns": [{
        "event_type": "check",
        "source_type": "component",
        "component": "test-component-che-event-filters-5",
        "current_entity": {
          "infos": {"customer": null}
        }
      }],
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
      "enabled": true,
      "priority": 2
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """
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
    When I save response createTimestamp={{ now }}
    When I wait the end of event processing
    When I send an event:
    """
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
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=che-event-filters-5
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-component-che-event-filters-5",
          "category": null,
          "component": "test-component-che-event-filters-5",
          "depends": [
            "test-resource-che-event-filters-5/test-component-che-event-filters-5"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-connector-che-event-filters-5/test-connector-name-che-event-filters-5"
          ],
          "impact_level": 1,
          "infos": {
            "customer": {
              "description": "Client",
              "name": "customer",
              "value": "test-customer-che-event-filters-5"
            }
          },
          "measurements": null,
          "name": "test-component-che-event-filters-5",
          "type": "component"
        },
        {
          "_id": "test-connector-che-event-filters-5/test-connector-name-che-event-filters-5",
          "category": null,
          "depends": [
            "test-component-che-event-filters-5"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-resource-che-event-filters-5/test-component-che-event-filters-5"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-connector-name-che-event-filters-5",
          "type": "connector"
        },
        {
          "_id": "test-resource-che-event-filters-5/test-component-che-event-filters-5",
          "category": null,
          "component": "test-component-che-event-filters-5",
          "component_infos": {
            "customer": {
              "description": "Client",
              "name": "customer",
              "value": "test-customer-che-event-filters-5"
            }
          },
          "depends": [
            "test-connector-che-event-filters-5/test-connector-name-che-event-filters-5"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-component-che-event-filters-5"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-resource-che-event-filters-5",
          "type": "resource"
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

  Scenario: given check event and enrichment event filter should not update disabled entity
    Given I am admin
    When I send an event:
    """
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
    When I save response createTimestamp={{ now }}
    When I wait the end of event processing
    When I do PUT /api/v4/entitybasics?_id=test-resource-che-event-filters-6/test-component-che-event-filters-6:
    """
    {
      "enabled": false,
      "impact_level": 1,
      "sli_avail_state": 0,
      "infos": [],
      "impact": [
        "test-component-che-event-filters-6"
      ],
      "depends": [
        "test-connector-che-event-filters-6/test-connector-name-che-event-filters-6"
      ]
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "enrichment",
      "patterns": [{
        "event_type": "check",
        "component": "test-component-che-event-filters-6",
        "current_entity": {
          "infos": {
            "customer": null
          }
        }
      }],
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
      "enabled": true,
      "priority": 2
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """
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
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=che-event-filters-6
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-component-che-event-filters-6",
          "category": null,
          "component": "test-component-che-event-filters-6",
          "depends": [
            "test-resource-che-event-filters-6/test-component-che-event-filters-6"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-connector-che-event-filters-6/test-connector-name-che-event-filters-6"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-component-che-event-filters-6",
          "type": "component"
        },
        {
          "_id": "test-connector-che-event-filters-6/test-connector-name-che-event-filters-6",
          "category": null,
          "depends": [
            "test-component-che-event-filters-6"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-resource-che-event-filters-6/test-component-che-event-filters-6"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-connector-name-che-event-filters-6",
          "type": "connector"
        },
        {
          "_id": "test-resource-che-event-filters-6/test-component-che-event-filters-6",
          "category": null,
          "component": "test-component-che-event-filters-6",
          "depends": [
            "test-connector-che-event-filters-6/test-connector-name-che-event-filters-6"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": false,
          "impact": [
            "test-component-che-event-filters-6"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-resource-che-event-filters-6",
          "type": "resource"
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

  Scenario: given check event and enrichment event filter with set_entity_info action
  should update event and entity
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "enrichment",
      "patterns": [{
        "event_type": "check",
        "component": "test-component-che-event-filters-7",
        "current_entity": {
          "infos": {
            "testdate": null
          }
        }
      }],
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
      "description": "test-event-filter-che-event-filters-7-2-description",
      "enabled": true,
      "priority": 2
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """
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
    When I save response createTimestamp={{ now }}
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=che-event-filters-7
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-component-che-event-filters-7",
          "category": null,
          "component": "test-component-che-event-filters-7",
          "depends": [
            "test-resource-che-event-filters-7/test-component-che-event-filters-7"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-connector-che-event-filters-7/test-connector-name-che-event-filters-7"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-component-che-event-filters-7",
          "type": "component"
        },
        {
          "_id": "test-connector-che-event-filters-7/test-connector-name-che-event-filters-7",
          "category": null,
          "depends": [
            "test-component-che-event-filters-7"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-resource-che-event-filters-7/test-component-che-event-filters-7"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-connector-name-che-event-filters-7",
          "type": "connector"
        },
        {
          "_id": "test-resource-che-event-filters-7/test-component-che-event-filters-7",
          "category": null,
          "component": "test-component-che-event-filters-7",
          "depends": [
            "test-connector-che-event-filters-7/test-connector-name-che-event-filters-7"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-component-che-event-filters-7"
          ],
          "impact_level": 1,
          "infos": {
            "testdate": {
              "description": "Date",
              "name": "testdate",
              "value": 1592215337
            }
          },
          "measurements": null,
          "name": "test-resource-che-event-filters-7",
          "type": "resource"
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

  Scenario: given check event and enrichment event filter with copy_to_entity_info action
  should update event and entity
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "enrichment",
      "patterns": [{
        "event_type": "check",
        "component": "test-component-che-event-filters-8",
        "current_entity": {
          "infos": {
            "customer": null,
            "testdate": null
          }
        }
      }],
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
      "description": "test-event-filter-che-event-filters-8-2-description",
      "enabled": true,
      "priority": 2
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "enrichment",
      "patterns": [{
        "event_type": "check",
        "component": "test-component-che-event-filters-8"
      }],
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
      "description": "test-event-filter-che-event-filters-8-3-description",
      "enabled": true,
      "priority": 3
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "enrichment",
      "patterns": [{
        "event_type": "check",
        "component": "test-component-che-event-filters-8",
        "current_entity": {
          "infos": {"output": null}
        }
      }],
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
      "priority": 4,
      "description": "test-event-filter-che-event-filters-8-4-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """
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
    When I save response createTimestamp={{ now }}
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=che-event-filters-8
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-component-che-event-filters-8",
          "category": null,
          "component": "test-component-che-event-filters-8",
          "depends": [
            "test-resource-che-event-filters-8/test-component-che-event-filters-8"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-connector-che-event-filters-8/test-connector-name-che-event-filters-8"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-component-che-event-filters-8",
          "type": "component"
        },
        {
          "_id": "test-connector-che-event-filters-8/test-connector-name-che-event-filters-8",
          "category": null,
          "depends": [
            "test-component-che-event-filters-8"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-resource-che-event-filters-8/test-component-che-event-filters-8"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-connector-name-che-event-filters-8",
          "type": "connector"
        },
        {
          "_id": "test-resource-che-event-filters-8/test-component-che-event-filters-8",
          "category": null,
          "component": "test-component-che-event-filters-8",
          "depends": [
            "test-connector-che-event-filters-8/test-connector-name-che-event-filters-8"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-component-che-event-filters-8"
          ],
          "impact_level": 1,
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
          },
          "measurements": null,
          "name": "test-resource-che-event-filters-8",
          "type": "resource"
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

  Scenario: given check event and enrichment event filter should enrich from external mongo data
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "enrichment",
      "external_data": {
        "component": {
          "type": "mongo",
          "select": {
            "component_customer": "{{ `{{.Event.Component}}` }}"
          },
          "collection": "assets"
        }
      },
      "patterns": [{
        "component": "test-eventfilter-assets-customer-1"
      }],
      "description": "test-event-filter-che-event-filters-9-description",
      "priority": 1,
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "status",
            "value": "{{ `{{.ExternalData.component.component_status}}` }}",
            "description": "status from assets"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-connector-che-event-filters-9",
      "connector_name": "test-connector-name-che-event-filters-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-eventfilter-assets-customer-1",
      "resource": "test-resource-che-event-filters-9",
      "state": 2,
      "output": "test-output-che-event-filters-9"
    }
    """
    When I save response createTimestamp={{ now }}
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-che-event-filters-9
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-resource-che-event-filters-9/test-eventfilter-assets-customer-1",
          "category": null,
          "component": "test-eventfilter-assets-customer-1",
          "depends": [
            "test-connector-che-event-filters-9/test-connector-name-che-event-filters-9"
          ],
          "enabled": true,
          "impact": [
            "test-eventfilter-assets-customer-1"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "impact_level": 1,
          "infos": {
            "status": {
              "name": "status",
              "description": "status from assets",
              "value": "test-eventfilter-assets-status-1"
            }
          },
          "measurements": null,
          "name": "test-resource-che-event-filters-9",
          "type": "resource"
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

  Scenario: given check event and enrichment event filter shouldn't drop event if enrich from external mongo data failed
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "enrichment",
      "external_data": {
        "component": {
          "type": "mongo",
          "select": {
            "component_customer": "{{ `{{.Event.Component}}` }}"
          },
          "collection": "assets"
        }
      },
      "patterns": [{
        "component": "assets_customer_not_exist"
      }],
      "description": "test-event-filter-che-event-filters-10-description",
      "priority": 1,
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "status",
            "value": "{{ `{{.ExternalData.component.component_status}}` }}",
            "description": "status from assets"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-connector-che-event-filters-10",
      "connector_name": "test-connector-name-che-event-filters-10",
      "source_type": "resource",
      "event_type": "check",
      "component": "assets_customer_not_exist",
      "resource": "test-resource-che-event-filters-10",
      "state": 2,
      "output": "test-output-che-event-filters-10"
    }
    """
    When I save response createTimestamp={{ now }}
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-che-event-filters-10
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-resource-che-event-filters-10/assets_customer_not_exist",
          "category": null,
          "component": "assets_customer_not_exist",
          "depends": [
            "test-connector-che-event-filters-10/test-connector-name-che-event-filters-10"
          ],
          "enabled": true,
          "impact": [
            "assets_customer_not_exist"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-resource-che-event-filters-10",
          "type": "resource"
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
    When I do GET /api/v4/alarms?search=che-event-filters-10
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-che-event-filters-10/assets_customer_not_exist"
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

  Scenario: given check event and enrichment event filter should drop event if enrich from external mongo data failed
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "enrichment",
      "external_data": {
        "component": {
          "type": "mongo",
          "select": {
            "component_customer": "{{ `{{.Event.Component}}` }}"
          },
          "collection": "assets"
        }
      },
      "patterns": [{
        "component": "assets_customer_not_exist_2"
      }],
      "description": "test-event-filter-che-event-filters-11-description",
      "priority": 1,
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "status",
            "value": "{{ `{{.ExternalData.component.component_status}}` }}",
            "description": "status from assets"
          }
        ],
        "on_success": "pass",
        "on_failure": "drop"
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-connector-che-event-filters-11",
      "connector_name": "test-connector-name-che-event-filters-11",
      "source_type": "resource",
      "event_type": "check",
      "component": "assets_customer_not_exist_2",
      "resource": "test-resource-che-event-filters-11",
      "state": 2,
      "output": "test-output-che-event-filters-11"
    }
    """
    When I save response createTimestamp={{ now }}
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-che-event-filters-11
    Then the response code should be 200
    Then the response body should contain:
    """
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
    When I do GET /api/v4/alarms?search=che-event-filters-11
    Then the response code should be 200
    Then the response body should contain:
    """
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

  Scenario: given check event and enrichment event filter should enrich from external api data
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "enrichment",
      "external_data": {
        "title": {
          "type": "api",
          "request": {
            "url": "http://localhost:3000/api/external_data",
            "method": "GET"
          }
        }
      },
      "patterns": [{
        "component": "test-component-che-event-filters-12"
      }],
      "description": "test-event-filter-che-event-filters-12-description",
      "priority": 1,
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "title",
            "value": "{{ `{{.ExternalData.title.title}}` }}",
            "description": "title from external api"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-connector-che-event-filters-12",
      "connector_name": "test-connector-name-che-event-filters-12",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-12",
      "resource": "test-resource-che-event-filters-12",
      "state": 2,
      "output": "test-output-che-event-filters-12"
    }
    """
    When I save response createTimestamp={{ now }}
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-che-event-filters-12
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-resource-che-event-filters-12/test-component-che-event-filters-12",
          "category": null,
          "component": "test-component-che-event-filters-12",
          "depends": [
            "test-connector-che-event-filters-12/test-connector-name-che-event-filters-12"
          ],
          "enabled": true,
          "impact": [
            "test-component-che-event-filters-12"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "impact_level": 1,
          "infos": {
            "title": {
              "name": "title",
              "description": "title from external api",
              "value": "test title"
            }
          },
          "measurements": null,
          "name": "test-resource-che-event-filters-12",
          "type": "resource"
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

  Scenario: given check event and enrichment event filter should enrich from external api data and mongo
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "enrichment",
      "external_data": {
        "title": {
          "type": "api",
          "request": {
            "url": "http://localhost:3000/api/external_data",
            "method": "GET"
          }
        },
        "component": {
          "type": "mongo",
          "select": {
            "component_customer": "{{ `{{.Event.Component}}` }}"
          },
          "collection": "assets"
        }
      },
      "patterns": [{
        "component": "test-eventfilter-assets-customer-2"
      }],
      "description": "test-event-filter-che-event-filters-13-description",
      "priority": 1,
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "title",
            "value": "{{ `{{.ExternalData.title.title}}` }}",
            "description": "title from external api"
          },
          {
            "type": "set_entity_info_from_template",
            "name": "status",
            "value": "{{ `{{.ExternalData.component.component_status}}` }}",
            "description": "status from assets"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-connector-che-event-filters-13",
      "connector_name": "test-connector-name-che-event-filters-13",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-eventfilter-assets-customer-2",
      "resource": "test-resource-che-event-filters-13",
      "state": 2,
      "output": "test-output-che-event-filters-13"
    }
    """
    When I save response createTimestamp={{ now }}
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-che-event-filters-13
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-resource-che-event-filters-13/test-eventfilter-assets-customer-2",
          "category": null,
          "component": "test-eventfilter-assets-customer-2",
          "depends": [
            "test-connector-che-event-filters-13/test-connector-name-che-event-filters-13"
          ],
          "enabled": true,
          "impact": [
            "test-eventfilter-assets-customer-2"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "impact_level": 1,
          "infos": {
            "status": {
              "name": "status",
              "description": "status from assets",
              "value": "test-eventfilter-assets-status-2"
            },
            "title": {
              "name": "title",
              "description": "title from external api",
              "value": "test title"
            }
          },
          "measurements": null,
          "name": "test-resource-che-event-filters-13",
          "type": "resource"
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
