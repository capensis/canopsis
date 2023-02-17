Feature: Get entities
  I need to be able to get a entities

  Scenario: given services should be able to retrieve dependencies and impacts
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-service-entity-get-1",
      "name": "test-service-entity-get-1",
      "output_template": "test-service-entity-get-1",
      "category": "",
      "enabled": true,
      "impact_level": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-service-entity-get-1-nested-1",
                "test-service-entity-get-1-nested-2",
                "test-service-entity-get-1-nested-3"
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
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-service-entity-get-1-nested-1",
      "name": "test-service-entity-get-1-nested-1",
      "output_template": "test-service-entity-get-1-nested-1",
      "category": "",
      "enabled": true,
      "impact_level": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-entity-get-1-1",
                "test-resource-entity-get-1-2"
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
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-service-entity-get-1-nested-2",
      "name": "test-service-entity-get-1-nested-2",
      "output_template": "test-service-entity-get-1-nested-2",
      "category": "",
      "enabled": true,
      "impact_level": 2,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-entity-get-1-3",
                "test-resource-entity-get-1-4"
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
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-service-entity-get-1-nested-3",
      "name": "test-service-entity-get-1-nested-3",
      "output_template": "test-service-entity-get-1-nested-3",
      "category": "",
      "enabled": true,
      "impact_level": 3,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-entity-get-1-5",
                "test-resource-entity-get-1-6"
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
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-entity-get-1",
      "connector_name": "test-connector-name-entity-get-1",
      "component": "test-component-entity-get-1",
      "resource": "test-resource-entity-get-1-1",
      "source_type": "resource",
      "state": 1
    }
    """
    When I wait the end of 3 events processing
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-entity-get-1",
      "connector_name": "test-connector-name-entity-get-1",
      "component": "test-component-entity-get-1",
      "resource": "test-resource-entity-get-1-2",
      "source_type": "resource",
      "state": 2
    }
    """
    When I wait the end of 3 events processing
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-entity-get-1",
      "connector_name": "test-connector-name-entity-get-1",
      "component": "test-component-entity-get-1",
      "resource": "test-resource-entity-get-1-3",
      "source_type": "resource",
      "state": 3
    }
    """
    When I wait the end of 3 events processing
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-entity-get-1",
      "connector_name": "test-connector-name-entity-get-1",
      "component": "test-component-entity-get-1",
      "resource": "test-resource-entity-get-1-4",
      "source_type": "resource",
      "state": 0
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-entity-get-1",
      "connector_name": "test-connector-name-entity-get-1",
      "component": "test-component-entity-get-1",
      "resource": "test-resource-entity-get-1-5",
      "source_type": "resource",
      "state": 0
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-entity-get-1",
      "connector_name": "test-connector-name-entity-get-1",
      "component": "test-component-entity-get-1",
      "resource": "test-resource-entity-get-1-6",
      "source_type": "resource",
      "state": 0
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/entityservice-dependencies?_id=test-service-entity-get-1&with_flags=true&sort_by=impact_state&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-service-entity-get-1-nested-2",
          "name": "test-service-entity-get-1-nested-2",
          "type": "service",
          "impact_level": 2,
          "state": 3,
          "impact_state": 6,
          "status": 1,
          "enabled": true,
          "ko_events": 1,
          "ok_events": 1,
          "deletable": true,
          "depends_count": 2,
          "impacts_count": 1,
          "category": null
        },
        {
          "_id": "test-service-entity-get-1-nested-1",
          "name": "test-service-entity-get-1-nested-1",
          "type": "service",
          "impact_level": 1,
          "state": 2,
          "impact_state": 2,
          "status": 1,
          "enabled": true,
          "ko_events": 2,
          "ok_events": 1,
          "deletable": true,
          "depends_count": 2,
          "impacts_count": 1,
          "category": null
        },
        {
          "_id": "test-service-entity-get-1-nested-3",
          "name": "test-service-entity-get-1-nested-3",
          "type": "service",
          "impact_level": 3,
          "state": 0,
          "impact_state": 0,
          "status": 0,
          "enabled": true,
          "ko_events": 0,
          "ok_events": 1,
          "deletable": true,
          "depends_count": 2,
          "impacts_count": 1,
          "category": null
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
    When I do GET /api/v4/entityservice-dependencies?_id=test-service-entity-get-1-nested-1&with_flags=true&sort_by=impact_state&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-entity-get-1-2/test-component-entity-get-1",
          "name": "test-resource-entity-get-1-2",
          "type": "resource",
          "impact_level": 1,
          "state": 2,
          "impact_state": 2,
          "status": 1,
          "enabled": true,
          "ko_events": 1,
          "ok_events": 0,
          "deletable": false,
          "depends_count": 0,
          "impacts_count": 1,
          "category": null
        },
        {
          "_id": "test-resource-entity-get-1-1/test-component-entity-get-1",
          "name": "test-resource-entity-get-1-1",
          "type": "resource",
          "impact_level": 1,
          "state": 1,
          "impact_state": 1,
          "status": 1,
          "enabled": true,
          "ko_events": 1,
          "ok_events": 0,
          "deletable": false,
          "depends_count": 0,
          "impacts_count": 1,
          "category": null
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
    When I do GET /api/v4/entityservice-dependencies?_id=test-service-entity-get-1-nested-2&with_flags=true&sort_by=impact_state&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-entity-get-1-3/test-component-entity-get-1",
          "name": "test-resource-entity-get-1-3",
          "type": "resource",
          "impact_level": 1,
          "state": 3,
          "impact_state": 3,
          "status": 1,
          "enabled": true,
          "ko_events": 1,
          "ok_events": 0,
          "deletable": false,
          "depends_count": 0,
          "impacts_count": 1,
          "category": null
        },
        {
          "_id": "test-resource-entity-get-1-4/test-component-entity-get-1",
          "name": "test-resource-entity-get-1-4",
          "type": "resource",
          "impact_level": 1,
          "state": 0,
          "impact_state": 0,
          "status": 0,
          "enabled": true,
          "ko_events": 0,
          "ok_events": 1,
          "deletable": true,
          "depends_count": 0,
          "impacts_count": 1,
          "category": null
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
    When I do GET /api/v4/entityservice-dependencies?_id=test-service-entity-get-1-nested-3&with_flags=true&sort_by=impact_state&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-entity-get-1-5/test-component-entity-get-1",
          "name": "test-resource-entity-get-1-5",
          "type": "resource",
          "impact_level": 1,
          "state": 0,
          "impact_state": 0,
          "status": 0,
          "enabled": true,
          "ko_events": 0,
          "ok_events": 1,
          "deletable": true,
          "depends_count": 0,
          "impacts_count": 1,
          "category": null
        },
        {
          "_id": "test-resource-entity-get-1-6/test-component-entity-get-1",
          "name": "test-resource-entity-get-1-6",
          "type": "resource",
          "impact_level": 1,
          "state": 0,
          "impact_state": 0,
          "status": 0,
          "enabled": true,
          "ko_events": 0,
          "ok_events": 1,
          "deletable": true,
          "depends_count": 0,
          "impacts_count": 1,
          "category": null
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
    When I do GET /api/v4/entityservice-dependencies?_id=test-resource-entity-get-1-1/test-component-entity-get-1
    Then the response code should be 404
    When I do GET /api/v4/entityservice-impacts?_id=test-service-entity-get-1
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
    When I do GET /api/v4/entityservice-impacts?_id=test-service-entity-get-1-nested-1&with_flags=true&sort_by=impact_state&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-service-entity-get-1",
          "name": "test-service-entity-get-1",
          "type": "service",
          "impact_level": 1,
          "state": 3,
          "impact_state": 3,
          "status": 1,
          "enabled": true,
          "ko_events": 3,
          "ok_events": 1,
          "deletable": true,
          "depends_count": 3,
          "impacts_count": 0,
          "category": null
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
    When I do GET /api/v4/entityservice-impacts?_id=test-service-entity-get-1-nested-2&with_flags=true&sort_by=impact_state&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-service-entity-get-1",
          "name": "test-service-entity-get-1",
          "type": "service",
          "impact_level": 1,
          "state": 3,
          "impact_state": 3,
          "status": 1,
          "enabled": true,
          "ko_events": 3,
          "ok_events": 1,
          "deletable": true,
          "depends_count": 3,
          "impacts_count": 0,
          "category": null
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
    When I do GET /api/v4/entityservice-impacts?_id=test-service-entity-get-1-nested-3&with_flags=true&sort_by=impact_state&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-service-entity-get-1",
          "name": "test-service-entity-get-1",
          "type": "service",
          "impact_level": 1,
          "state": 3,
          "impact_state": 3,
          "status": 1,
          "enabled": true,
          "ko_events": 3,
          "ok_events": 1,
          "deletable": true,
          "depends_count": 3,
          "impacts_count": 0,
          "category": null
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
    When I do GET /api/v4/entityservice-impacts?_id=test-resource-entity-get-1-1/test-component-entity-get-1&with_flags=true&sort_by=impact_state&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-service-entity-get-1-nested-1",
          "name": "test-service-entity-get-1-nested-1",
          "type": "service",
          "impact_level": 1,
          "state": 2,
          "impact_state": 2,
          "status": 1,
          "enabled": true,
          "ko_events": 2,
          "ok_events": 1,
          "deletable": true,
          "depends_count": 2,
          "impacts_count": 1,
          "category": null
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
    When I do GET /api/v4/entityservice-impacts?_id=test-resource-entity-get-1-3/test-component-entity-get-1&with_flags=true&sort_by=impact_state&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-service-entity-get-1-nested-2",
          "name": "test-service-entity-get-1-nested-2",
          "type": "service",
          "impact_level": 2,
          "state": 3,
          "impact_state": 6,
          "status": 1,
          "enabled": true,
          "ko_events": 1,
          "ok_events": 1,
          "deletable": true,
          "depends_count": 2,
          "impacts_count": 1,
          "category": null
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
    When I do GET /api/v4/entityservice-impacts?_id=test-resource-entity-get-1-5/test-component-entity-get-1&with_flags=true&sort_by=impact_state&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-service-entity-get-1-nested-3",
          "name": "test-service-entity-get-1-nested-3",
          "type": "service",
          "impact_level": 3,
          "state": 0,
          "impact_state": 0,
          "status": 0,
          "enabled": true,
          "ko_events": 0,
          "ok_events": 1,
          "deletable": true,
          "depends_count": 2,
          "impacts_count": 1,
          "category": null
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

  Scenario: given services should be able to retrieve filtered by category dependencies and impacts
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-service-entity-get-2-1",
      "name": "test-service-entity-get-2-1",
      "output_template": "test-service-entity-get-2-1",
      "category": "test-category-to-entityservice-entity-get-2",
      "enabled": true,
      "impact_level": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-entity-get-2-1",
                "test-resource-entity-get-2-2"
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
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-service-entity-get-2-2",
      "name": "test-service-entity-get-2-2",
      "output_template": "test-service-entity-get-2-2",
      "category": "",
      "enabled": true,
      "impact_level": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-entity-get-2-2",
                "test-resource-entity-get-2-3"
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
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-entity-get-2",
      "connector_name": "test-connector-name-entity-get-2",
      "component": "test-component-entity-get-2",
      "resource": "test-resource-entity-get-2-1",
      "source_type": "resource",
      "state": 1
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-entity-get-2",
      "connector_name": "test-connector-name-entity-get-2",
      "component": "test-component-entity-get-2",
      "resource": "test-resource-entity-get-2-2",
      "source_type": "resource",
      "state": 2
    }
    """
    When I wait the end of 3 events processing
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-entity-get-2",
      "connector_name": "test-connector-name-entity-get-2",
      "component": "test-component-entity-get-2",
      "resource": "test-resource-entity-get-2-3",
      "source_type": "resource",
      "state": 3
    }
    """
    When I wait the end of 2 events processing
    When I do PUT /api/v4/entitybasics?_id=test-resource-entity-get-2-1/test-component-entity-get-2:
    """json
    {
      "category": "test-category-to-entityservice-entity-get-2",
      "enabled": true,
      "impact_level": 1,
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/entityservice-dependencies?_id=test-service-entity-get-2-1&category=test-category-to-entityservice-entity-get-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-entity-get-2-1/test-component-entity-get-2"
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
    When I do GET /api/v4/entityservice-impacts?_id=test-resource-entity-get-2-2/test-component-entity-get-2&category=test-category-to-entityservice-entity-get-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-service-entity-get-2-1"
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

  Scenario: given services should be able to search dependencies and impacts
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-service-entity-get-3-1",
      "name": "test-service-entity-get-3-1",
      "output_template": "test-service-entity-get-3-1",
      "category": "",
      "enabled": true,
      "impact_level": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-entity-get-3-1",
                "test-resource-entity-get-3-2"
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
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-service-entity-get-3-2",
      "name": "test-service-entity-get-3-2",
      "output_template": "test-service-entity-get-3-2",
      "category": "",
      "enabled": true,
      "impact_level": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-entity-get-3-2",
                "test-resource-entity-get-3-3"
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
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-entity-get-3",
      "connector_name": "test-connector-name-entity-get-3",
      "component": "test-component-entity-get-3",
      "resource": "test-resource-entity-get-3-1",
      "source_type": "resource",
      "state": 1
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-entity-get-3",
      "connector_name": "test-connector-name-entity-get-3",
      "component": "test-component-entity-get-3",
      "resource": "test-resource-entity-get-3-2",
      "source_type": "resource",
      "state": 2
    }
    """
    When I wait the end of 3 events processing
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-entity-get-3",
      "connector_name": "test-connector-name-entity-get-3",
      "component": "test-component-entity-get-3",
      "resource": "test-resource-entity-get-3-3",
      "source_type": "resource",
      "state": 3
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/entityservice-dependencies?_id=test-service-entity-get-3-1&search=test-resource-entity-get-3-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-entity-get-3-1/test-component-entity-get-3"
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
    When I do GET /api/v4/entityservice-impacts?_id=test-resource-entity-get-3-2/test-component-entity-get-3&search=test-service-entity-get-3-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-service-entity-get-3-1"
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

  Scenario: given get dependencies unauth request should not allow access
    When I do GET /api/v4/entityservice-dependencies?_id=test-entityservice-not-found
    Then the response code should be 401

  Scenario: given get dependencies request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/entityservice-dependencies?_id=test-entityservice-not-found
    Then the response code should be 403

  Scenario: given get impacts unauth request should not allow access
    When I do GET /api/v4/entityservice-impacts?_id=test-entityservice-not-found
    Then the response code should be 401

  Scenario: given get impacts request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/entityservice-impacts?_id=test-entityservice-not-found
    Then the response code should be 403
