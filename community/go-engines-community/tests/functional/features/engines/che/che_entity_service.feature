Feature: create service entity
  I need to be able to create service entity

  Scenario: given resource entity and new service entity should add resource to service on service creation
    Given I am admin
    When I send an event:
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
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-1-name",
      "output_template": "test-entityservice-che-service-1-output",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"name": "test-resource-che-service-1"}],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities?search=che-service-1&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-1",
          "category": null,
          "component": "test-component-che-service-1",
          "depends": [
            "test-resource-che-service-1/test-component-che-service-1"
          ],
          "enabled": true,
          "impact": [
            "test-connector-che-service-1/test-connector-name-che-service-1"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-component-che-service-1",
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-1/test-connector-name-che-service-1",
          "category": null,
          "depends": [
            "test-component-che-service-1"
          ],
          "enabled": true,
          "impact": [
            "test-resource-che-service-1/test-component-che-service-1"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-connector-name-che-service-1",
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "category": null,
          "depends": [
            "test-resource-che-service-1/test-component-che-service-1"
          ],
          "enabled": true,
          "impact": [],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-entityservice-che-service-1-name",
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-1/test-component-che-service-1",
          "category": null,
          "component": "test-component-che-service-1",
          "depends": [
            "test-connector-che-service-1/test-connector-name-che-service-1"
          ],
          "enabled": true,
          "impact": [
            "test-component-che-service-1",
            "{{ .serviceID }}"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-resource-che-service-1",
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

  Scenario: given service entity and new resource entity should add resource to service on resource creation
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-2-name",
      "output_template": "test-entityservice-che-service-2-output",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"name": "test-resource-che-service-2"}],
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
          "_id": "test-component-che-service-2",
          "category": null,
          "component": "test-component-che-service-2",
          "depends": [
            "test-resource-che-service-2/test-component-che-service-2"
          ],
          "enabled": true,
          "impact": [
            "test-connector-che-service-2/test-connector-name-che-service-2"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-component-che-service-2",
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-2/test-connector-name-che-service-2",
          "category": null,
          "depends": [
            "test-component-che-service-2"
          ],
          "enabled": true,
          "impact": [
            "test-resource-che-service-2/test-component-che-service-2"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-connector-name-che-service-2",
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "category": null,
          "depends": [
            "test-resource-che-service-2/test-component-che-service-2"
          ],
          "enabled": true,
          "impact": [],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-entityservice-che-service-2-name",
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-2/test-component-che-service-2",
          "category": null,
          "component": "test-component-che-service-2",
          "depends": [
            "test-connector-che-service-2/test-connector-name-che-service-2"
          ],
          "enabled": true,
          "impact": [
            "test-component-che-service-2",
            "{{ .serviceID }}"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
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
      "entity_patterns": [{"name": "test-resource-che-service-3-1"}],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities?search=che-service-3&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-3",
          "depends": [
            "test-resource-che-service-3-1/test-component-che-service-3",
            "test-resource-che-service-3-2/test-component-che-service-3"
          ],
          "impact": [
            "test-connector-che-service-3/test-connector-name-che-service-3"
          ],
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-3/test-connector-name-che-service-3",
          "depends": [
            "test-component-che-service-3"
          ],
          "impact": [
            "test-resource-che-service-3-1/test-component-che-service-3",
            "test-resource-che-service-3-2/test-component-che-service-3"
          ],
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [
            "test-resource-che-service-3-1/test-component-che-service-3"
          ],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-3-1/test-component-che-service-3",
          "depends": [
            "test-connector-che-service-3/test-connector-name-che-service-3"
          ],
          "impact": [
            "test-component-che-service-3",
            "{{ .serviceID }}"
          ],
          "type": "resource"
        },
        {
          "_id": "test-resource-che-service-3-2/test-component-che-service-3",
          "depends": [
            "test-connector-che-service-3/test-connector-name-che-service-3"
          ],
          "impact": [
            "test-component-che-service-3"
          ],
          "type": "resource"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 5
      }
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
      "entity_patterns": [{"name": "test-resource-che-service-3-2"}],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities?search=che-service-3&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-3",
          "depends": [
            "test-resource-che-service-3-1/test-component-che-service-3",
            "test-resource-che-service-3-2/test-component-che-service-3"
          ],
          "impact": [
            "test-connector-che-service-3/test-connector-name-che-service-3"
          ],
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-3/test-connector-name-che-service-3",
          "depends": [
            "test-component-che-service-3"
          ],
          "impact": [
            "test-resource-che-service-3-1/test-component-che-service-3",
            "test-resource-che-service-3-2/test-component-che-service-3"
          ],
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [
            "test-resource-che-service-3-2/test-component-che-service-3"
          ],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-3-1/test-component-che-service-3",
          "depends": [
            "test-connector-che-service-3/test-connector-name-che-service-3"
          ],
          "impact": [
            "test-component-che-service-3"
          ],
          "type": "resource"
        },
        {
          "_id": "test-resource-che-service-3-2/test-component-che-service-3",
          "depends": [
            "test-connector-che-service-3/test-connector-name-che-service-3"
          ],
          "impact": [
            "test-component-che-service-3",
            "{{ .serviceID }}"
          ],
          "type": "resource"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 5
      }
    }
    """

  Scenario: given service entity and resource entity with extra infos should remove resource from service on infos update
    Given I am admin
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
      "output": "test-output-che-service-4",
      "client": "test-client-che-service-4"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-4-name",
      "output_template": "test-entityservice-che-service-4-output",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"infos": {
        "client": {"value": "test-client-che-service-4"}
      }}],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities?search=che-service-4&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-4",
          "depends": [
            "test-resource-che-service-4/test-component-che-service-4"
          ],
          "impact": [
            "test-connector-che-service-4/test-connector-name-che-service-4"
          ],
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-4/test-connector-name-che-service-4",
          "depends": [
            "test-component-che-service-4"
          ],
          "impact": [
            "test-resource-che-service-4/test-component-che-service-4"
          ],
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [
            "test-resource-che-service-4/test-component-che-service-4"
          ],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-4/test-component-che-service-4",
          "depends": [
            "test-connector-che-service-4/test-connector-name-che-service-4"
          ],
          "impact": [
            "test-component-che-service-4",
            "{{ .serviceID }}"
          ],
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
      "output": "test-output-che-service-4",
      "client": "test-another-client-che-service-4"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities?search=che-service-4&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-4",
          "depends": [
            "test-resource-che-service-4/test-component-che-service-4"
          ],
          "impact": [
            "test-connector-che-service-4/test-connector-name-che-service-4"
          ],
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-4/test-connector-name-che-service-4",
          "depends": [
            "test-component-che-service-4"
          ],
          "impact": [
            "test-resource-che-service-4/test-component-che-service-4"
          ],
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-4/test-component-che-service-4",
          "depends": [
            "test-connector-che-service-4/test-connector-name-che-service-4"
          ],
          "impact": [
            "test-component-che-service-4"
          ],
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

  Scenario: given service entity and updated resource entity with extra infos should add resource to service on infos update
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-5-name",
      "output_template": "test-entityservice-che-service-5-output",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"infos": {
        "client": {"value": "test-client-che-service-5"},
        "company": {"value": "test-company-che-service-5"}
      }}],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
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
      "output": "test-output-che-service-5",
      "client": "test-client-che-service-5"
    }
    """
    When I wait the end of event processing
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
      "output": "test-output-che-service-5",
      "company": "test-company-che-service-5"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities?search=che-service-5&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-5",
          "depends": [
            "test-resource-che-service-5/test-component-che-service-5"
          ],
          "impact": [
            "test-connector-che-service-5/test-connector-name-che-service-5"
          ],
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-5/test-connector-name-che-service-5",
          "depends": [
            "test-component-che-service-5"
          ],
          "impact": [
            "test-resource-che-service-5/test-component-che-service-5"
          ],
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [
            "test-resource-che-service-5/test-component-che-service-5"
          ],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-5/test-component-che-service-5",
          "depends": [
            "test-connector-che-service-5/test-connector-name-che-service-5"
          ],
          "impact": [
            "test-component-che-service-5",
            "{{ .serviceID }}"
          ],
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

  Scenario: given service entity and resource entity and enrichment event filter should add resource to service on infos update
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-6-name",
      "output_template": "test-entityservice-che-service-6-output",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"infos": {
        "manager": {"value": "test-manager-che-service-6"}
      }}],
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
      "patterns": [{
        "event_type": "check",
        "resource": "test-resource-che-service-6"
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
      "on_success": "pass",
      "on_failure": "pass",
      "description": "test-eventfilter-che-service-6-description",
      "enabled": true,
      "priority": 1
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "patterns": [{
        "event_type": "check",
        "resource": "test-resource-che-service-6",
        "current_entity": {
          "infos": {
            "manager": null
          }
        }
      }],
      "actions": [
        {
          "type": "set_entity_info_from_template",
          "name": "manager",
          "description": "Manager",
          "value": "test-manager-che-service-6"
        }
      ],
      "description": "test-eventfilter-che-service-6-description",
      "enabled": true,
      "priority": 2,
      "on_success": "pass",
      "on_failure": "pass"
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
    When I do GET /api/v4/entities?search=che-service-6&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-6",
          "depends": [
            "test-resource-che-service-6/test-component-che-service-6"
          ],
          "impact": [
            "test-connector-che-service-6/test-connector-name-che-service-6"
          ],
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-6/test-connector-name-che-service-6",
          "depends": [
            "test-component-che-service-6"
          ],
          "impact": [
            "test-resource-che-service-6/test-component-che-service-6"
          ],
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [
            "test-resource-che-service-6/test-component-che-service-6"
          ],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-6/test-component-che-service-6",
          "depends": [
            "test-connector-che-service-6/test-connector-name-che-service-6"
          ],
          "impact": [
            "test-component-che-service-6",
            "{{ .serviceID }}"
          ],
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

  Scenario: given service entity and resource entity and updated enrichment event filter should remove resource from service on infos update
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-7-name",
      "output_template": "test-entityservice-che-service-7-output",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"infos": {
        "manager": {"value": "test-manager-che-service-7"}
      }}],
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
      "patterns": [{
        "event_type": "check",
        "resource": "test-resource-che-service-7"
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
      "on_success": "pass",
      "on_failure": "pass",
      "description": "test-eventfilter-che-service-7-description",
      "enabled": true,
      "priority": 1
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "patterns": [{
        "event_type": "check",
        "resource": "test-resource-che-service-7",
        "current_entity": {
          "infos": {
            "manager": null
          }
        }
      }],
      "actions": [
        {
          "type": "set_entity_info_from_template",
          "name": "manager",
          "description": "Manager",
          "value": "test-manager-che-service-7"
        }
      ],
      "description": "test-eventfilter-che-service-7-description",
      "enabled": true,
      "priority": 2,
      "on_success": "pass",
      "on_failure": "pass"
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
    When I do GET /api/v4/entities?search=che-service-7&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-7",
          "depends": [
            "test-resource-che-service-7/test-component-che-service-7"
          ],
          "impact": [
            "test-connector-che-service-7/test-connector-name-che-service-7"
          ],
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-7/test-connector-name-che-service-7",
          "depends": [
            "test-component-che-service-7"
          ],
          "impact": [
            "test-resource-che-service-7/test-component-che-service-7"
          ],
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [
            "test-resource-che-service-7/test-component-che-service-7"
          ],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-7/test-component-che-service-7",
          "depends": [
            "test-connector-che-service-7/test-connector-name-che-service-7"
          ],
          "impact": [
            "test-component-che-service-7",
            "{{ .serviceID }}"
          ],
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
    When I do PUT /api/v4/eventfilter/rules/{{ .ruleID }}:
    """json
    {
      "type": "enrichment",
      "patterns": [{
        "event_type": "check",
        "resource": "test-resource-che-service-7"
      }],
      "actions": [
        {
          "type": "set_entity_info_from_template",
          "name": "manager",
          "description": "Manager",
          "value": "test-another-manager-che-service-7"
        }
      ],
      "description": "test-eventfilter-che-service-7-description",
      "enabled": true,
      "priority": 2,
      "on_success": "pass",
      "on_failure": "pass"
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
    When I do GET /api/v4/entities?search=che-service-7&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-7",
          "depends": [
            "test-resource-che-service-7/test-component-che-service-7"
          ],
          "impact": [
            "test-connector-che-service-7/test-connector-name-che-service-7"
          ],
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-7/test-connector-name-che-service-7",
          "depends": [
            "test-component-che-service-7"
          ],
          "impact": [
            "test-resource-che-service-7/test-component-che-service-7"
          ],
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-7/test-component-che-service-7",
          "depends": [
            "test-connector-che-service-7/test-connector-name-che-service-7"
          ],
          "impact": [
            "test-component-che-service-7"
          ],
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

  Scenario: given service entity and new component entity on resource event should add component to service
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-8-name",
      "output_template": "test-entityservice-che-service-8-output",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"name": "test-component-che-service-8"}],
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
    When I do GET /api/v4/entities?search=che-service-8&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-8",
          "depends": [
            "test-resource-che-service-8/test-component-che-service-8"
          ],
          "impact": [
            "test-connector-che-service-8/test-connector-name-che-service-8",
            "{{ .serviceID }}"
          ],
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-8/test-connector-name-che-service-8",
          "depends": [
            "test-component-che-service-8"
          ],
          "impact": [
            "test-resource-che-service-8/test-component-che-service-8"
          ],
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [
            "test-component-che-service-8"
          ],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-8/test-component-che-service-8",
          "depends": [
            "test-connector-che-service-8/test-connector-name-che-service-8"
          ],
          "impact": [
            "test-component-che-service-8"
          ],
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

  Scenario: given service entity and new connector entity on resource event should add connector to service
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-9-name",
      "output_template": "test-entityservice-che-service-9-output",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"name": "test-connector-name-che-service-9"}],
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
    When I do GET /api/v4/entities?search=che-service-9&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-9",
          "depends": [
            "test-resource-che-service-9/test-component-che-service-9"
          ],
          "impact": [
            "test-connector-che-service-9/test-connector-name-che-service-9"
          ],
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-9/test-connector-name-che-service-9",
          "depends": [
            "test-component-che-service-9"
          ],
          "impact": [
            "test-resource-che-service-9/test-component-che-service-9",
            "{{ .serviceID }}"
          ],
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [
            "test-connector-che-service-9/test-connector-name-che-service-9"
          ],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-9/test-component-che-service-9",
          "depends": [
            "test-connector-che-service-9/test-connector-name-che-service-9"
          ],
          "impact": [
            "test-component-che-service-9"
          ],
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

  Scenario: given service entity and enrichment event filter and new component entity on resource event should add component to service
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-10-name",
      "output_template": "test-entityservice-che-service-10-output",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"name": "test-component-che-service-10"}],
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
      "patterns": [{
        "event_type": "check",
        "resource": "test-resource-che-service-10"
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
      "on_success": "pass",
      "on_failure": "pass",
      "description": "test-eventfilter-che-service-10-description",
      "enabled": true,
      "priority": 1
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "patterns": [{
        "event_type": "check",
        "resource": "test-resource-che-service-10",
        "current_entity": {
          "infos": {
            "manager": null
          }
        }
      }],
      "actions": [
        {
          "type": "set_entity_info_from_template",
          "name": "manager",
          "description": "Manager",
          "value": "test-manager-che-service-10"
        }
      ],
      "priority": 2,
      "on_success": "pass",
      "on_failure": "pass",
      "description": "test-eventfilter-che-service-10-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-10",
      "connector_name": "test-connector-name-che-service-10",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-10",
      "resource": "test-resource-che-service-10",
      "state": 2,
      "output": "test-output-che-service-10"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=che-service-10&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-10",
          "depends": [
            "test-resource-che-service-10/test-component-che-service-10"
          ],
          "impact": [
            "test-connector-che-service-10/test-connector-name-che-service-10",
            "{{ .serviceID }}"
          ],
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-10/test-connector-name-che-service-10",
          "depends": [
            "test-component-che-service-10"
          ],
          "impact": [
            "test-resource-che-service-10/test-component-che-service-10"
          ],
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [
            "test-component-che-service-10"
          ],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-10/test-component-che-service-10",
          "depends": [
            "test-connector-che-service-10/test-connector-name-che-service-10"
          ],
          "impact": [
            "test-component-che-service-10"
          ],
          "infos": {
            "manager": {
              "description": "Manager",
              "name": "manager",
              "value": "test-manager-che-service-10"
            }
          },
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

  Scenario: given service entity and resource entity with component infos by extra infos should add resource to service on component event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-11-name",
      "impact_level": 1,
      "output_template": "test-entityservice-che-service-11-output",
      "enabled": true,
      "entity_patterns": [{"component_infos": {
        "client": {"value": "test-client-che-service-11"}
      }}],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-11",
      "connector_name": "test-connector-name-che-service-11",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-11",
      "resource": "test-resource-che-service-11",
      "state": 2,
      "output": "test-output-che-service-11"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-11",
      "connector_name": "test-connector-name-che-service-11",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-che-service-11",
      "state": 2,
      "output": "test-output-che-service-11",
      "client": "test-client-che-service-11"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/entities?search=che-service-11&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-11",
          "depends": [
            "test-resource-che-service-11/test-component-che-service-11"
          ],
          "impact": [
            "test-connector-che-service-11/test-connector-name-che-service-11"
          ],
          "infos": {
            "client": {
              "name": "client",
              "value": "test-client-che-service-11"
            }
          },
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-11/test-connector-name-che-service-11",
          "depends": [
            "test-component-che-service-11"
          ],
          "impact": [
            "test-resource-che-service-11/test-component-che-service-11"
          ],
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [
            "test-resource-che-service-11/test-component-che-service-11"
          ],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-11/test-component-che-service-11",
          "component_infos": {
            "client": {
              "name": "client",
              "value": "test-client-che-service-11"
            }
          },
          "depends": [
            "test-connector-che-service-11/test-connector-name-che-service-11"
          ],
          "impact": [
            "test-component-che-service-11",
            "{{ .serviceID }}"
          ],
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

  Scenario: given service entity and resource entity with component infos by extra infos should add resource to service on resource event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-12-name",
      "impact_level": 1,
      "output_template": "test-entityservice-che-service-12-output",
      "enabled": true,
      "entity_patterns": [{"component_infos": {
        "client": {"value": "test-client-che-service-12"}
      }}],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-12",
      "connector_name": "test-connector-name-che-service-12",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-che-service-12",
      "state": 2,
      "output": "test-output-che-service-12",
      "client": "test-client-che-service-12"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-12",
      "connector_name": "test-connector-name-che-service-12",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-12",
      "resource": "test-resource-che-service-12",
      "state": 2,
      "output": "test-output-che-service-12"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities?search=che-service-12&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-12",
          "depends": [
            "test-resource-che-service-12/test-component-che-service-12"
          ],
          "impact": [
            "test-connector-che-service-12/test-connector-name-che-service-12"
          ],
          "infos": {
            "client": {
              "name": "client",
              "value": "test-client-che-service-12"
            }
          },
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-12/test-connector-name-che-service-12",
          "depends": [
            "test-component-che-service-12"
          ],
          "impact": [
            "test-resource-che-service-12/test-component-che-service-12"
          ],
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [
            "test-resource-che-service-12/test-component-che-service-12"
          ],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-12/test-component-che-service-12",
          "component_infos": {
            "client": {
              "name": "client",
              "value": "test-client-che-service-12"
            }
          },
          "depends": [
            "test-connector-che-service-12/test-connector-name-che-service-12"
          ],
          "impact": [
            "test-component-che-service-12",
            "{{ .serviceID }}"
          ],
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

  Scenario: given service entity and resource entity with component infos by event filter should add resource to service on component event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-13-name",
      "impact_level": 1,
      "output_template": "test-entityservice-che-service-13-output",
      "enabled": true,
      "entity_patterns": [{"component_infos": {
        "manager": {"value": "test-manager-che-service-13"}
      }}],
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
      "patterns": [{
        "event_type": "check",
        "component": "test-component-che-service-13"
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
      "on_success": "pass",
      "on_failure": "pass",
      "description": "test-eventfilter-che-service-13-description",
      "enabled": true,
      "priority": 1
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "patterns": [{
        "event_type": "check",
        "source_type": "component",
        "component": "test-component-che-service-13",
        "current_entity": {
          "infos": {
            "manager": null
          }
        }
      }],
      "actions": [
        {
          "type": "set_entity_info_from_template",
          "name": "manager",
          "description": "Manager",
          "value": "test-manager-che-service-13"
        }
      ],
      "priority": 2,
      "on_success": "pass",
      "on_failure": "pass",
      "description": "test-eventfilter-che-service-13-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-13",
      "connector_name": "test-connector-name-che-service-13",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-13",
      "resource": "test-resource-che-service-13",
      "state": 2,
      "output": "test-output-che-service-13"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-13",
      "connector_name": "test-connector-name-che-service-13",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-che-service-13",
      "state": 2,
      "output": "test-output-che-service-13"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/entities?search=che-service-13&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-13",
          "depends": [
            "test-resource-che-service-13/test-component-che-service-13"
          ],
          "impact": [
            "test-connector-che-service-13/test-connector-name-che-service-13"
          ],
          "infos": {
            "manager": {
              "name": "manager",
              "value": "test-manager-che-service-13"
            }
          },
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-13/test-connector-name-che-service-13",
          "depends": [
            "test-component-che-service-13"
          ],
          "impact": [
            "test-resource-che-service-13/test-component-che-service-13"
          ],
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [
            "test-resource-che-service-13/test-component-che-service-13"
          ],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-13/test-component-che-service-13",
          "component_infos": {
            "manager": {
              "name": "manager",
              "value": "test-manager-che-service-13"
            }
          },
          "depends": [
            "test-connector-che-service-13/test-connector-name-che-service-13"
          ],
          "impact": [
            "test-component-che-service-13",
            "{{ .serviceID }}"
          ],
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

  Scenario: given service entity and resource entity with component infos by event filter should add resource to service on resource event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-14-name",
      "impact_level": 1,
      "output_template": "test-entityservice-che-service-14-output",
      "enabled": true,
      "entity_patterns": [{"component_infos": {
        "manager": {"value": "test-manager-che-service-14"}
      }}],
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
      "patterns": [{
        "event_type": "check",
        "component": "test-component-che-service-14"
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
      "on_success": "pass",
      "on_failure": "pass",
      "description": "test-eventfilter-che-service-14-description",
      "enabled": true,
      "priority": 1
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "patterns": [{
        "event_type": "check",
        "source_type": "component",
        "component": "test-component-che-service-14",
        "current_entity": {
          "infos": {
            "manager": null
          }
        }
      }],
      "actions": [
        {
          "type": "set_entity_info_from_template",
          "name": "manager",
          "description": "Manager",
          "value": "test-manager-che-service-14"
        }
      ],
      "description": "test-eventfilter-che-service-14-description",
      "enabled": true,
      "priority": 2,
      "on_success": "pass",
      "on_failure": "pass"
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-14",
      "connector_name": "test-connector-name-che-service-14",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-che-service-14",
      "state": 2,
      "output": "test-output-che-service-14"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-14",
      "connector_name": "test-connector-name-che-service-14",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-14",
      "resource": "test-resource-che-service-14",
      "state": 2,
      "output": "test-output-che-service-14"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities?search=che-service-14&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-14",
          "depends": [
            "test-resource-che-service-14/test-component-che-service-14"
          ],
          "impact": [
            "test-connector-che-service-14/test-connector-name-che-service-14"
          ],
          "infos": {
            "manager": {
              "name": "manager",
              "value": "test-manager-che-service-14"
            }
          },
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-14/test-connector-name-che-service-14",
          "depends": [
            "test-component-che-service-14"
          ],
          "impact": [
            "test-resource-che-service-14/test-component-che-service-14"
          ],
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [
            "test-resource-che-service-14/test-component-che-service-14"
          ],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-14/test-component-che-service-14",
          "component_infos": {
            "manager": {
              "name": "manager",
              "value": "test-manager-che-service-14"
            }
          },
          "depends": [
            "test-connector-che-service-14/test-connector-name-che-service-14"
          ],
          "impact": [
            "test-component-che-service-14",
            "{{ .serviceID }}"
          ],
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

  Scenario: given service entity and resource entity with component infos by extra infos should remove resource from service on component event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-15-name",
      "impact_level": 1,
      "output_template": "test-entityservice-che-service-15-output",
      "enabled": true,
      "entity_patterns": [{"component_infos": {
        "client": {"value": "test-client-che-service-15"}
      }}],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-15",
      "connector_name": "test-connector-name-che-service-15",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-15",
      "resource": "test-resource-che-service-15",
      "state": 2,
      "output": "test-output-che-service-15"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-15",
      "connector_name": "test-connector-name-che-service-15",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-che-service-15",
      "state": 2,
      "output": "test-output-che-service-15",
      "client": "test-client-che-service-15"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/entities?search=che-service-15&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-15",
          "depends": [
            "test-resource-che-service-15/test-component-che-service-15"
          ],
          "impact": [
            "test-connector-che-service-15/test-connector-name-che-service-15"
          ],
          "infos": {
            "client": {
              "name": "client",
              "value": "test-client-che-service-15"
            }
          },
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-15/test-connector-name-che-service-15",
          "depends": [
            "test-component-che-service-15"
          ],
          "impact": [
            "test-resource-che-service-15/test-component-che-service-15"
          ],
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [
            "test-resource-che-service-15/test-component-che-service-15"
          ],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-15/test-component-che-service-15",
          "component_infos": {
            "client": {
              "name": "client",
              "value": "test-client-che-service-15"
            }
          },
          "depends": [
            "test-connector-che-service-15/test-connector-name-che-service-15"
          ],
          "impact": [
            "test-component-che-service-15",
            "{{ .serviceID }}"
          ],
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
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-15",
      "connector_name": "test-connector-name-che-service-15",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-che-service-15",
      "state": 2,
      "output": "test-output-che-service-15",
      "client": "test-another-client-che-service-15"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/entities?search=che-service-15&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-15",
          "depends": [
            "test-resource-che-service-15/test-component-che-service-15"
          ],
          "impact": [
            "test-connector-che-service-15/test-connector-name-che-service-15"
          ],
          "infos": {
            "client": {
              "name": "client",
              "value": "test-another-client-che-service-15"
            }
          },
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-15/test-connector-name-che-service-15",
          "depends": [
            "test-component-che-service-15"
          ],
          "impact": [
            "test-resource-che-service-15/test-component-che-service-15"
          ],
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-15/test-component-che-service-15",
          "component_infos": {
            "client": {
              "name": "client",
              "value": "test-another-client-che-service-15"
            }
          },
          "depends": [
            "test-connector-che-service-15/test-connector-name-che-service-15"
          ],
          "impact": [
            "test-component-che-service-15"
          ],
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

  Scenario: given service entity and updated resource entity by api should add resource to service on infos update
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-16-name",
      "output_template": "test-entityservice-che-service-16-output",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"infos": {"manager": {"value": "test-manager-che-service-16"}}}],
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
    When I do GET /api/v4/entities?search=che-service-16&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-16",
          "depends": [
            "test-resource-che-service-16/test-component-che-service-16"
          ],
          "impact": [
            "test-connector-che-service-16/test-connector-name-che-service-16"
          ],
          "type": "component"
        },
        {
          "_id": "test-connector-che-service-16/test-connector-name-che-service-16",
          "depends": [
            "test-component-che-service-16"
          ],
          "impact": [
            "test-resource-che-service-16/test-component-che-service-16"
          ],
          "type": "connector"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [
            "test-resource-che-service-16/test-component-che-service-16"
          ],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-16/test-component-che-service-16",
          "depends": [
            "test-connector-che-service-16/test-connector-name-che-service-16"
          ],
          "impact": [
            "test-component-che-service-16",
            "{{ .serviceID }}"
          ],
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

  Scenario: given deleted service entity should remove service from all links
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-17-name-1",
      "impact_level": 1,
      "output_template": "test-entityservice-che-service-17-output-1",
      "enabled": true,
      "entity_patterns": [{"name": "test-resource-che-service-17"}],
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
      "entity_patterns": [{"name": "test-entityservice-che-service-17-name-1"}],
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
    When I do GET /api/v4/entities?search=che-service-17&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-17"
        },
        {
          "_id": "test-connector-che-service-17/test-connector-name-che-service-17"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [
            "test-resource-che-service-17/test-component-che-service-17"
          ],
          "impact": [
            "{{ .impactServiceID }}"
          ],
          "type": "service"
        },
        {
          "_id": "{{ .impactServiceID }}",
          "depends": [
            "{{ .serviceID }}"
          ],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-17/test-component-che-service-17",
          "depends": [
            "test-connector-che-service-17/test-connector-name-che-service-17"
          ],
          "impact": [
            "test-component-che-service-17",
            "{{ .serviceID }}"
          ],
          "type": "resource"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 5
      }
    }
    """
    When I do DELETE /api/v4/entityservices/{{ .serviceID }}
    Then the response code should be 204
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities?search=che-service-17&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-17"
        },
        {
          "_id": "test-connector-che-service-17/test-connector-name-che-service-17"
        },
        {
          "_id": "{{ .impactServiceID }}",
          "depends": [],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-17/test-component-che-service-17",
          "depends": [
            "test-connector-che-service-17/test-connector-name-che-service-17"
          ],
          "impact": [
            "test-component-che-service-17"
          ],
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

  Scenario: given disabled service entity should not update service context graph
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-18-name",
      "impact_level": 1,
      "output_template": "test-entityservice-che-service-18-output",
      "enabled": true,
      "entity_patterns": [
        {"name": "test-resource-che-service-18-1"},
        {"name": "test-resource-che-service-18-2"}
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
    When I do GET /api/v4/entities?search=che-service-18&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-18"
        },
        {
          "_id": "test-connector-che-service-18/test-connector-name-che-service-18"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [
            "test-resource-che-service-18-1/test-component-che-service-18"
          ],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-18-1/test-component-che-service-18",
          "depends": [
            "test-connector-che-service-18/test-connector-name-che-service-18"
          ],
          "impact": [
            "test-component-che-service-18",
            "{{ .serviceID }}"
          ],
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
    When I do PUT /api/v4/entityservices/{{ .serviceID }}:
    """json
    {
      "name": "test-entityservice-che-service-18-name",
      "impact_level": 1,
      "output_template": "test-entityservice-che-service-18-output",
      "enabled": false,
      "entity_patterns": [{"name": "test-resource-che-service-18"}],
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
    When I do GET /api/v4/entities?search=che-service-18&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-che-service-18"
        },
        {
          "_id": "test-connector-che-service-18/test-connector-name-che-service-18"
        },
        {
          "_id": "{{ .serviceID }}",
          "depends": [
            "test-resource-che-service-18-1/test-component-che-service-18"
          ],
          "impact": [],
          "type": "service"
        },
        {
          "_id": "test-resource-che-service-18-1/test-component-che-service-18",
          "depends": [
            "test-connector-che-service-18/test-connector-name-che-service-18"
          ],
          "impact": [
            "test-component-che-service-18",
            "{{ .serviceID }}"
          ],
          "type": "resource"
        },
        {
          "_id": "test-resource-che-service-18-2/test-component-che-service-18",
          "depends": [
            "test-connector-che-service-18/test-connector-name-che-service-18"
          ],
          "impact": [
            "test-component-che-service-18"
          ],
          "type": "resource"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 5
      }
    }
    """
