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
            "url": "http://localhost:3000/api/external_data",
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
          "category": null,
          "component": "test-component-che-event-filters-api-1",
          "depends": [
            "test-connector-che-event-filters-api-1/test-connector-name-che-event-filters-api-1"
          ],
          "enabled": true,
          "impact": [
            "test-component-che-event-filters-api-1"
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
          "name": "test-resource-che-event-filters-api-1",
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
    """json
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
          "category": null,
          "component": "test-eventfilter-mongo-data-2-customer",
          "depends": [
            "test-connector-che-event-filters-api-2/test-connector-name-che-event-filters-api-2"
          ],
          "enabled": true,
          "impact": [
            "test-eventfilter-mongo-data-2-customer"
          ],
          "impact_level": 1,
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
          },
          "measurements": null,
          "name": "test-resource-che-event-filters-api-2",
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
