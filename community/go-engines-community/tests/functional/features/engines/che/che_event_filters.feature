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
    When I do GET /api/v4/entities?search=che-event-filters-1
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
        "component": "test-component-che-event-filters-3"
      }],
      "external_data": {
        "entity": {
          "type": "entity"
        }
      },
      "actions": [
        {
          "type": "copy",
          "from": "ExternalData.entity",
          "to": "Entity"
        }
      ],
      "description": "test-event-filter-che-event-filters-3-description",
      "enabled": true,
      "on_success": "pass",
      "on_failure": "pass",
      "priority": 1
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
          "infos": {
            "customer": null,
            "manager": null
          }
        }
      }],
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
      "description": "test-event-filter-che-event-filters-3-description",
      "enabled": true,
      "priority": 2,
      "on_success": "pass",
      "on_failure": "pass"
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
      "actions": [
        {
          "type": "set_field_from_template",
          "name": "Output",
          "value": "{{ `{{ .Event.Output }}` }} (client: {{ `{{ .Event.Entity.Infos.customer.Value }}` }})"
        }
      ],
      "description": "test-event-filter-che-event-filters-3-description",
      "enabled": true,
      "priority": 3,
      "on_success": "pass",
      "on_failure": "pass"
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
      "actions": [
        {
          "type": "set_entity_info_from_template",
          "name": "output",
          "description": "Output",
          "value": "{{ `{{ .Event.Output }}` }}"
        }
      ],
      "description": "test-event-filter-che-event-filters-3-description",
      "enabled": true,
      "priority": 4,
      "on_success": "pass",
      "on_failure": "pass"
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
        "component": "test-component-che-event-filters-4"
      }],
      "external_data": {
        "entity": {
          "type": "entity"
        }
      },
      "actions": [
        {
          "type": "copy",
          "from": "ExternalData.entity",
          "to": "Entity"
        }
      ],
      "description": "test-event-filter-che-event-filters-4-description",
      "enabled": true,
      "on_success": "pass",
      "on_failure": "pass",
      "priority": 1
    }
    """
    Then the response code should be 201
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
      "actions": [
        {
          "type": "set_entity_info",
          "name": "customer",
          "description": "Client",
          "value": "test-customer-che-event-filters-4"
        }
      ],
      "description": "test-event-filter-che-event-filters-4-description",
      "enabled": true,
      "priority": 2,
      "on_success": "pass",
      "on_failure": "pass"
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
        "component": "test-component-che-event-filters-5"
      }],
      "external_data": {
        "entity": {
          "type": "entity"
        }
      },
      "actions": [
        {
          "type": "copy",
          "from": "ExternalData.entity",
          "to": "Entity"
        }
      ],
      "description": "test-event-filter-che-event-filters-5-description",
      "enabled": true,
      "on_success": "pass",
      "on_failure": "pass",
      "priority": 1
    }
    """
    Then the response code should be 201
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
      "actions": [
        {
          "type": "set_entity_info",
          "name": "customer",
          "description": "Client",
          "value": "test-customer-che-event-filters-5"
        }
      ],
      "description": "test-event-filter-che-event-filters-5-description",
      "enabled": true,
      "priority": 2,
      "on_success": "pass",
      "on_failure": "pass"
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
    When I wait the end of 2 events processing
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
        "component": "test-component-che-event-filters-6"
      }],
      "external_data": {
        "entity": {
          "type": "entity"
        }
      },
      "actions": [
        {
          "type": "copy",
          "from": "ExternalData.entity",
          "to": "Entity"
        }
      ],
      "description": "test-event-filter-che-event-filters-6-description",
      "enabled": true,
      "on_success": "pass",
      "on_failure": "pass",
      "priority": 1
    }
    """
    Then the response code should be 201
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
      "actions": [
        {
          "type": "set_entity_info",
          "name": "customer",
          "description": "Client",
          "value": "test-customer-che-event-filters-6"
        }
      ],
      "description": "test-event-filter-che-event-filters-6-description",
      "enabled": true,
      "priority": 2,
      "on_success": "pass",
      "on_failure": "pass"
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
    Then the response body should be:
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
        "component": "test-component-che-event-filters-7"
      }],
      "external_data": {
        "entity": {
          "type": "entity"
        }
      },
      "actions": [
        {
          "type": "copy",
          "from": "ExternalData.entity",
          "to": "Entity"
        }
      ],
      "description": "test-event-filter-che-event-filters-7-1-description",
      "enabled": true,
      "on_success": "pass",
      "on_failure": "pass",
      "priority": 1
    }
    """
    Then the response code should be 201
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
      "actions": [
        {
          "type": "set_entity_info",
          "name": "testdate",
          "description": "Date",
          "value": 1592215337
        }
      ],
      "description": "test-event-filter-che-event-filters-7-2-description",
      "enabled": true,
      "priority": 2,
      "on_success": "pass",
      "on_failure": "pass"
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
        "component": "test-component-che-event-filters-8"
      }],
      "external_data": {
        "entity": {
          "type": "entity"
        }
      },
      "actions": [
        {
          "type": "copy",
          "from": "ExternalData.entity",
          "to": "Entity"
        }
      ],
      "description": "test-event-filter-che-event-filters-8-1-description",
      "enabled": true,
      "on_success": "pass",
      "on_failure": "pass",
      "priority": 1
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
          "infos": {
            "customer": null,
            "testdate": null
          }
        }
      }],
      "actions": [
        {
          "type": "copy_to_entity_info",
          "name": "customer",
          "description": "Client",
          "from": "Event.ExtraInfos.customer"
        },
        {
          "type": "copy_to_entity_info",
          "name": "testdate",
          "description": "Date",
          "from": "Event.ExtraInfos.testdate"
        }
      ],
      "description": "test-event-filter-che-event-filters-8-2-description",
      "enabled": true,
      "priority": 2,
      "on_success": "pass",
      "on_failure": "pass"
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
      "actions": [
        {
          "type": "set_field_from_template",
          "name": "Output",
          "value": "{{ `{{ .Event.Output }}` }} (client: {{ `{{ .Event.Entity.Infos.customer.Value }}` }})"
        }
      ],
      "description": "test-event-filter-che-event-filters-8-3-description",
      "enabled": true,
      "priority": 3,
      "on_success": "pass",
      "on_failure": "pass"
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
      "actions": [
        {
          "type": "copy_to_entity_info",
          "name": "output",
          "description": "Output",
          "from": "Event.Output"
        }
      ],
      "priority": 4,
      "description": "test-event-filter-che-event-filters-8-4-description",
      "enabled": true,
      "on_success": "pass",
      "on_failure": "pass"
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
