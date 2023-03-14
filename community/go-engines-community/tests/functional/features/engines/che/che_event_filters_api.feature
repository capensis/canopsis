Feature: modify event on event filter
  I need to be able to modify event on event filter

  Scenario: given check event and enrichment event filter should enrich from external api data
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "external_data": {
        "title": {
          "type": "api",
          "request": {
            "url": "{{ .dummyApiURL }}/api/external_data",
            "method": "GET"
          }
        }
      },
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-api-1"
            }
          }
        ]
      ],
      "description": "test-event-filter-che-event-filters-api-1-description",
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
    """json
    {
      "connector": "test-connector-che-event-filters-api-1",
      "connector_name": "test-connector-name-che-event-filters-api-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-api-1",
      "resource": "test-resource-che-event-filters-api-1",
      "state": 2,
      "output": "test-output-che-event-filters-api-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-che-event-filters-api-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-che-event-filters-api-1/test-component-che-event-filters-api-1",
          "infos": {
            "title": {
              "name": "title",
              "description": "title from external api",
              "value": "test title"
            }
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

  Scenario: given check event and enrichment event filter should enrich from external api data and mongo
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "external_data": {
        "title": {
          "type": "api",
          "request": {
            "url": "{{ .dummyApiURL }}/api/external_data",
            "method": "GET"
          }
        },
        "component": {
          "type": "mongo",
          "select": {
            "customer": "{{ `{{.Event.Component}}` }}"
          },
          "collection": "eventfilter_mongo_data"
        }
      },
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-mongo-data-2-customer"
            }
          }
        ]
      ],
      "description": "test-event-filter-che-event-filters-api-2-description",
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
            "value": "{{ `{{.ExternalData.component.status}}` }}",
            "description": "status from external collection"
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
    """json
    {
      "connector": "test-connector-che-event-filters-api-2",
      "connector_name": "test-connector-name-che-event-filters-api-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-eventfilter-mongo-data-2-customer",
      "resource": "test-resource-che-event-filters-api-2",
      "state": 2,
      "output": "test-output-che-event-filters-api-2"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-che-event-filters-api-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-che-event-filters-api-2/test-eventfilter-mongo-data-2-customer",
          "infos": {
            "status": {
              "name": "status",
              "description": "status from external collection",
              "value": "test-eventfilter-mongo-data-2-status"
            },
            "title": {
              "name": "title",
              "description": "title from external api",
              "value": "test title"
            }
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

  Scenario: given check event and enrichment event filter should enrich from external api data where response is a document with an array
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "external_data": {
        "title": {
          "type": "api",
          "request": {
            "url": "{{ .dummyApiURL }}/api/external_data_document_with_array",
            "method": "GET"
          }
        }
      },
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-assets-customer-3"
            }
          }
        ]
      ],
      "description": "test-event-filter-che-event-filters-api-3-description",
      "priority": 1,
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "title",
            "value": "{{ `{{ index .ExternalData.title \"array.1.title\" }}` }}",
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
    """json
    {
      "connector": "test-connector-che-event-filters-api-3",
      "connector_name": "test-connector-name-che-event-filters-api-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-eventfilter-assets-customer-3",
      "resource": "test-resource-che-event-filters-api-3",
      "state": 2,
      "output": "test-output-che-event-filters-api-3"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-che-event-filters-api-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-che-event-filters-api-3/test-eventfilter-assets-customer-3",
          "infos": {
            "title": {
              "name": "title",
              "description": "title from external api",
              "value": "test title 2"
            }
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

  Scenario: given check event and enrichment event filter should enrich from external api data where response is an array
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "external_data": {
        "title": {
          "type": "api",
          "request": {
            "url": "{{ .dummyApiURL }}/api/external_data_response_is_array",
            "method": "GET"
          }
        }
      },
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-assets-customer-4"
            }
          }
        ]
      ],
      "description": "test-event-filter-che-event-filters-17-description",
      "priority": 1,
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "title",
            "value": "{{ `{{ index .ExternalData.title \"1.title\" }}` }}",
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
    """json
    {
      "connector": "test-connector-che-event-filters-api-4",
      "connector_name": "test-connector-name-che-event-filters-api-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-eventfilter-assets-customer-4",
      "resource": "test-resource-che-event-filters-api-4",
      "state": 2,
      "output": "test-output-che-event-filters-api-4"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-che-event-filters-api-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-che-event-filters-api-4/test-eventfilter-assets-customer-4",
          "infos": {
            "title": {
              "name": "title",
              "description": "title from external api",
              "value": "test title 2"
            }
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

  Scenario: given check event and enrichment event filter should enrich by regexp from external api data
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "external_data": {
        "name": {
          "type": "api",
          "request": {
            "url": "{{ .dummyApiURL }}/api/external_data_response_is_nested_documents",
            "method": "GET"
          }
        }
      },
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-assets-customer-5"
            }
          }
        ]
      ],
      "description": "test-resource-che-event-filters-api-5-description",
      "priority": 1,
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "name",
            "value": "{{ `{{ regex_map_key .ExternalData.name \"object.*.fields.name\" }}` }}",
            "description": "name from external api"
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
    """json
    {
      "connector": "test-connector-che-event-filters-api-5",
      "connector_name": "test-connector-name-che-event-filters-api-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-eventfilter-assets-customer-5",
      "resource": "test-resource-che-event-filters-api-5",
      "state": 2,
      "output": "test-output-che-event-filters-api-5"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-che-event-filters-api-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-che-event-filters-api-5/test-eventfilter-assets-customer-5",
          "infos": {
            "name": {
              "name": "name",
              "description": "name from external api",
              "value": "test name"
            }
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

  Scenario: given check event and enrichment event filter should execute templates in payload
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "external_data": {
        "name": {
          "type": "api",
          "request": {
            "url": "{{ .apiURL }}/api/v4/scenarios",
            "method": "POST",
            "auth": {
              "username": "root",
              "password": "test"
            },
            "headers": {"Content-Type": "application/json"},
            "payload": "{\"priority\": 10039,\"name\":\"{{ `{{ .Event.Component }}` }}\",\"enabled\":true,\"triggers\":[\"create\"],\"actions\":[{\"entity_pattern\":[[{\"field\":\"name\",\"cond\":{\"type\": \"eq\", \"value\": \"test-eventfilter-assets-customer-6-pattern\"}}]],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
          }
        }
      },
      "event_pattern": [
        [
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
              "value": "test-eventfilter-assets-customer-6"
            }
          }
        ]
      ],
      "description": "test-resource-che-event-filters-api-6-description",
      "priority": 1,
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "name",
            "value": "{{ `{{ index .ExternalData.name \"name\" }}` }}",
            "description": "name from external api"
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
    """json
    {
      "connector": "test-connector-che-event-filters-api-6",
      "connector_name": "test-connector-name-che-event-filters-api-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-eventfilter-assets-customer-6",
      "resource": "test-resource-che-event-filters-api-6",
      "state": 2,
      "output": "test-output-che-event-filters-api-6"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-che-event-filters-api-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-che-event-filters-api-6/test-eventfilter-assets-customer-6",
          "infos": {
            "name": {
              "name": "name",
              "description": "name from external api",
              "value": "test-eventfilter-assets-customer-6"
            }
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
