Feature: modify event on event filter
  I need to be able to modify event on event filter

  Scenario: given check event and enrichment event filter should enrich from external mongo data
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
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
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-mongo-data-1-customer"
            }
          }
        ]
      ],
      "description": "test-event-filter-che-event-filters-mongo-1-description",
      "priority": 1,
      "enabled": true,
      "config": {
        "actions": [
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
      "connector": "test-connector-che-event-filters-mongo-1",
      "connector_name": "test-connector-name-che-event-filters-mongo-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-eventfilter-mongo-data-1-customer",
      "resource": "test-resource-che-event-filters-mongo-1",
      "state": 2,
      "output": "test-output-che-event-filters-mongo-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-che-event-filters-mongo-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-che-event-filters-mongo-1/test-eventfilter-mongo-data-1-customer",
          "infos": {
            "status": {
              "name": "status",
              "description": "status from external collection",
              "value": "test-eventfilter-mongo-data-1-status"
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

  Scenario: given check event and enrichment event filter shouldn't drop event if enrich from external mongo data failed
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "external_data": {
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
              "value": "test-eventfilter-mongo-data-not-exist-customer"
            }
          }
        ]
      ],
      "description": "test-event-filter-che-event-filters-mongo-2-description",
      "priority": 1,
      "enabled": true,
      "config": {
        "actions": [
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
      "connector": "test-connector-che-event-filters-mongo-2",
      "connector_name": "test-connector-name-che-event-filters-mongo-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-eventfilter-mongo-data-not-exist-customer",
      "resource": "test-resource-che-event-filters-mongo-2",
      "state": 2,
      "output": "test-output-che-event-filters-mongo-2"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-che-event-filters-mongo-2
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-resource-che-event-filters-mongo-2/test-eventfilter-mongo-data-not-exist-customer",
          "infos": {},
          "name": "test-resource-che-event-filters-mongo-2",
          "type": "resource",
          "category": null,
          "component": "test-eventfilter-mongo-data-not-exist-customer",
          "connector": "test-connector-che-event-filters-mongo-2/test-connector-name-che-event-filters-mongo-2",
          "enabled": true,
          "old_entity_patterns": null,
          "impact_level": 1,
          "impact_state": 2,
          "last_event_date": {{ (index .lastResponse.data 0).last_event_date }},
          "alarm_last_update_date": {{ (index .lastResponse.data 0).alarm_last_update_date }},
          "ko_events": 1,
          "ok_events": 0,
          "state": 2,
          "status": 1
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
    When I do GET /api/v4/alarms?search=che-event-filters-mongo-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-che-event-filters-mongo-2/test-eventfilter-mongo-data-not-exist-customer"
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
    """json
    {
      "type": "enrichment",
      "external_data": {
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
              "value": "test-eventfilter-mongo-data-not-exist-2-customer"
            }
          }
        ]
      ],
      "description": "test-event-filter-che-event-filters-mongo-3-description",
      "priority": 1,
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "status",
            "value": "{{ `{{.ExternalData.component.status}}` }}",
            "description": "status from external collection"
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
    """json
    {
      "connector": "test-connector-che-event-filters-mongo-3",
      "connector_name": "test-connector-name-che-event-filters-mongo-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-eventfilter-mongo-data-not-exist-2-customer",
      "resource": "test-resource-che-event-filters-mongo-3",
      "state": 2,
      "output": "test-output-che-event-filters-mongo-3"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-che-event-filters-mongo-3
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-resource-che-event-filters-mongo-3/test-eventfilter-mongo-data-not-exist-2-customer",
          "infos": {},
          "name": "test-resource-che-event-filters-mongo-3",
          "type": "resource",
          "category": null,
          "component": "test-eventfilter-mongo-data-not-exist-2-customer",
          "connector": "test-connector-che-event-filters-mongo-3/test-connector-name-che-event-filters-mongo-3",
          "enabled": true,
          "old_entity_patterns": null,
          "impact_level": 1,
          "impact_state": 0,
          "last_event_date": {{ (index .lastResponse.data 0).last_event_date }},
          "ko_events": 0,
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
    When I do GET /api/v4/alarms?search=che-event-filters-mongo-3
    Then the response code should be 200
    Then the response body should contain:
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

  Scenario: given check event and enrichment event filter should enrich from external mongo data by regexp
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
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
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-connector-che-event-filters-mongo-4"
            }
          }
        ]
      ],
      "description": "test-event-filter-che-event-filters-mongo-4-description",
      "priority": 1,
      "enabled": true,
      "config": {
        "actions": [
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
    [
      {
        "connector": "test-connector-che-event-filters-mongo-4",
        "connector_name": "test-connector-name-che-event-filters-mongo-4",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-eventfilter-mongo-data-regexp-1-customer",
        "resource": "test-resource-che-event-filters-mongo-4-1",
        "state": 2,
        "output": "test-eventfilter-mongo-data-regexp-1-message"
      },
      {
        "connector": "test-connector-che-event-filters-mongo-4",
        "connector_name": "test-connector-name-che-event-filters-mongo-4",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-eventfilter-mongo-data-regexp-1-customer",
        "resource": "test-resource-che-event-filters-mongo-4-2",
        "state": 2,
        "output": "test-eventfilter-mongo-data-regexp-2-message"
      }
    ]
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities?search=test-resource-che-event-filters-mongo-4&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-che-event-filters-mongo-4-1/test-eventfilter-mongo-data-regexp-1-customer",
          "infos": {
            "status": {
              "name": "status",
              "description": "status from external collection",
              "value": "test-eventfilter-mongo-data-regexp-1-status"
            }
          }
        },
        {
          "_id": "test-resource-che-event-filters-mongo-4-2/test-eventfilter-mongo-data-regexp-1-customer",
          "infos": {
            "status": {
              "name": "status",
              "description": "status from external collection",
              "value": "test-eventfilter-mongo-data-regexp-2-status"
            }
          }
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

  Scenario: given int field in regexp external data should not update entity
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "external_data": {
        "component": {
          "type": "mongo",
          "sort_by": "status",
          "sort": "asc",
          "select": {
            "customer": "{{ `{{.Event.Component}}` }}"
          },
          "regexp": {
            "message": "{{ `{{.Event.Output}}` }}",
            "state": "{{ `{{.Event.State}}` }}"
          },
          "collection": "eventfilter_mongo_data_regexp"
        }
      },
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-connector-che-event-filters-mongo-5"
            }
          }
        ]
      ],
      "description": "test-event-filter-che-event-filters-mongo-5-description",
      "priority": 1,
      "enabled": true,
      "config": {
        "actions": [
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
      "connector": "test-connector-che-event-filters-mongo-5",
      "connector_name": "test-connector-name-che-event-filters-mongo-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-eventfilter-mongo-data-regexp-1-customer",
      "resource": "test-resource-che-event-filters-mongo-5",
      "state": 1,
      "output": "test-eventfilter-mongo-data-regexp-1-message"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-che-event-filters-mongo-5
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-resource-che-event-filters-mongo-5/test-eventfilter-mongo-data-regexp-1-customer",
          "infos": {},
          "name": "test-resource-che-event-filters-mongo-5",
          "type": "resource",
          "category": null,
          "component": "test-eventfilter-mongo-data-regexp-1-customer",
          "connector": "test-connector-che-event-filters-mongo-5/test-connector-name-che-event-filters-mongo-5",
          "enabled": true,
          "old_entity_patterns": null,
          "impact_level": 1,
          "impact_state": 1,
          "last_event_date": {{ (index .lastResponse.data 0).last_event_date }},
          "alarm_last_update_date": {{ (index .lastResponse.data 0).alarm_last_update_date }},
          "ko_events": 1,
          "ok_events": 0,
          "state": 1,
          "status": 1
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
