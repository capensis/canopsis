Feature: get service entities
  I need to be able to get service entities

  Scenario: given one dependency with maintenance pbehavior should get one entity
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-entity-1",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-1",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-1",
      "resource" : "test-resource-pbehavior-weather-service-entity-1",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-entity-1",
      "name": "test-pbehavior-weather-service-entity-1",
      "output_template": "test-pbehavior-weather-service-entity-1",
      "category": "test-category-pbehavior-weather-service-entity",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-entity-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-entity-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-entity-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services/test-pbehavior-weather-service-entity-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-pbehavior-weather-service-entity-1",
          "state": {"val": 2},
          "status": {"val": 1},
          "is_grey": true,
          "icon": "build",
          "pbehavior_info": {
            "canonical_type": "maintenance",
            "name": "test-pbehavior-weather-service-entity-1",
            "icon_name": "build"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-entity-1",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "reason": {
                "_id": "test-reason-to-engine",
                "name": "Test Engine",
                "description": "Test Engine"
              },
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "icon_name": "build",
                "name": "Engine maintenance",
                "type": "maintenance"
              }
            }
          ]
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

  Scenario: given one dependency with active pbehavior should get one entity
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-entity-2",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-2",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-2",
      "resource" : "test-resource-pbehavior-weather-service-entity-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-entity-2",
      "name": "test-pbehavior-weather-service-entity-2",
      "output_template": "test-pbehavior-weather-service-entity-2",
      "category": "test-category-pbehavior-weather-service-entity",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-entity-2"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-entity-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-active-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-entity-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do GET /api/v4/weather-services/test-pbehavior-weather-service-entity-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-pbehavior-weather-service-entity-2",
          "state": {"val": 2},
          "status": {"val": 1},
          "is_grey": false,
          "icon": "person",
          "pbehavior_info": {
            "canonical_type": "active",
            "name": "test-pbehavior-weather-service-entity-2",
            "icon_name": "brightness_3"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-entity-2",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "reason": {
                "_id": "test-reason-to-engine",
                "name": "Test Engine",
                "description": "Test Engine"
              },
              "type": {
                "_id": "test-active-type-to-engine",
                "icon_name": "brightness_3",
                "name": "Engine active",
                "type": "active"
              }
            }
          ]
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

  Scenario: given dependency with maintenance pbehavior
  and another dependency without pbehavior should get multiple entities
    Given I am admin
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-entity-3",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-3",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-3",
      "resource" : "test-resource-pbehavior-weather-service-entity-3-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-entity-3",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-3",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-3",
      "resource" : "test-resource-pbehavior-weather-service-entity-3-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-entity-3",
      "name": "test-pbehavior-weather-service-entity-3",
      "output_template": "test-pbehavior-weather-service-entity-3",
      "category": "test-category-pbehavior-weather-service-entity",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-pbehavior-weather-service-entity-3-1",
                "test-resource-pbehavior-weather-service-entity-3-2"
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
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-entity-3-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-entity-3-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services/test-pbehavior-weather-service-entity-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-pbehavior-weather-service-entity-3-1",
          "state": {"val": 3},
          "status": {"val": 1},
          "is_grey": true,
          "icon": "build",
          "pbehavior_info": {
            "canonical_type": "maintenance",
            "name": "test-pbehavior-weather-service-entity-3-1",
            "icon_name": "build"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-entity-3-1",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "reason": {
                "_id": "test-reason-to-engine",
                "name": "Test Engine",
                "description": "Test Engine"
              },
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "icon_name": "build",
                "name": "Engine maintenance",
                "type": "maintenance"
              }
            }
          ]
        },
        {
          "name": "test-resource-pbehavior-weather-service-entity-3-2",
          "state": {"val": 2},
          "status": {"val": 1},
          "is_grey": false,
          "icon": "person",
          "pbehavior_info": null,
          "pbehaviors": []
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

  Scenario: given service with maintenance pbehavior should not get gray flag
    Given I am admin
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-entity-4",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-4",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-4",
      "resource" : "test-resource-pbehavior-weather-service-entity-4",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-entity-4",
      "name": "test-pbehavior-weather-service-entity-4",
      "output_template": "test-pbehavior-weather-service-entity-4",
      "category": "test-category-pbehavior-weather-service-entity",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-entity-4"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-entity-4",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-entity-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do GET /api/v4/weather-services/test-pbehavior-weather-service-entity-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-pbehavior-weather-service-entity-4",
          "state": {"val": 3},
          "status": {"val": 1},
          "is_grey": false,
          "icon": "wb_cloudy",
          "pbehavior_info": null,
          "pbehaviors": []
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

  Scenario: given service with maintenance pbehavior without alarm should not get gray flag
    Given I am admin
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-entity-5",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-5",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-5",
      "resource" : "test-resource-pbehavior-weather-service-entity-5",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-entity-5",
      "name": "test-pbehavior-weather-service-entity-5",
      "output_template": "test-pbehavior-weather-service-entity-5",
      "category": "test-category-pbehavior-weather-service-entity",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-entity-5"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-entity-5",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-entity-5"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do GET /api/v4/weather-services/test-pbehavior-weather-service-entity-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-pbehavior-weather-service-entity-5",
          "state": {"val": 0},
          "status": {"val": 0},
          "is_grey": false,
          "icon": "wb_sunny",
          "pbehavior_info": null,
          "pbehaviors": []
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

  Scenario: given dependency with maintenance pbehavior without alarm should get gray flag
    Given I am admin
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-entity-6",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-6",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-6",
      "resource" : "test-resource-pbehavior-weather-service-entity-6",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-entity-6",
      "name": "test-pbehavior-weather-service-entity-6",
      "output_template": "test-pbehavior-weather-service-entity-6",
      "category": "test-category-pbehavior-weather-service-entity",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-entity-6"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-entity-6",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-entity-6"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services/test-pbehavior-weather-service-entity-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-pbehavior-weather-service-entity-6",
          "state": {"val": 0},
          "status": {"val": 0},
          "is_grey": true,
          "icon": "build",
          "pbehavior_info": {
            "canonical_type": "maintenance",
            "name": "test-pbehavior-weather-service-entity-6",
            "icon_name": "build"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-entity-6",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "reason": {
                "_id": "test-reason-to-engine",
                "name": "Test Engine",
                "description": "Test Engine"
              },
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "icon_name": "build",
                "name": "Engine maintenance",
                "type": "maintenance"
              }
            }
          ]
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

  Scenario: given dependency with origin pbehavior should get origin icon
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-entity-7",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-7",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-7",
      "resource" : "test-resource-pbehavior-weather-service-entity-7",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-entity-7",
      "name": "test-pbehavior-weather-service-entity-7",
      "output_template": "test-pbehavior-weather-service-entity-7",
      "category": "test-category-pbehavior-weather-service-entity",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-entity-7"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/bulk/entity-pbehaviors:
    """json
    [
      {
        "entity": "test-resource-pbehavior-weather-service-entity-7/test-component-pbehavior-weather-service-entity-7",
        "origin": "pbehavior-weather-service-entity-7",
        "name": "test-pbehavior-weather-service-entity-7",
        "tstart": {{ now }},
        "tstop": {{ nowAdd "1h" }},
        "color": "#FFFFFF",
        "type": "test-maintenance-type-to-engine",
        "reason": "test-reason-to-engine"
      }
    ]
    """
    Then the response code should be 207
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services/test-pbehavior-weather-service-entity-7?pbh_origin=pbehavior-weather-service-entity-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-pbehavior-weather-service-entity-7",
          "state": {"val": 2},
          "status": {"val": 1},
          "is_grey": true,
          "icon": "build",
          "pbh_origin_icon": "build",
          "pbehavior_info": {
            "canonical_type": "maintenance",
            "name": "test-pbehavior-weather-service-entity-7",
            "icon_name": "build"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-entity-7",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "reason": {
                "_id": "test-reason-to-engine",
                "name": "Test Engine",
                "description": "Test Engine"
              },
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "icon_name": "build",
                "name": "Engine maintenance",
                "type": "maintenance"
              }
            }
          ]
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

  Scenario: given dependency with origin pause pbehavior and maintenance pbehavior should get origin icon
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-entity-8",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-8",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-8",
      "resource" : "test-resource-pbehavior-weather-service-entity-8",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-entity-8",
      "name": "test-pbehavior-weather-service-entity-8",
      "output_template": "test-pbehavior-weather-service-entity-8",
      "category": "test-category-pbehavior-weather-service-entity",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-entity-8"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "name": "test-pbehavior-weather-service-entity-8-1",
      "enabled": true,
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-pause-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-entity-8"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/bulk/entity-pbehaviors:
    """json
    [
      {
        "entity": "test-resource-pbehavior-weather-service-entity-8/test-component-pbehavior-weather-service-entity-8",
        "origin": "pbehavior-weather-service-entity-8",
        "name": "test-pbehavior-weather-service-entity-8-2",
        "tstart": {{ now }},
        "tstop": {{ nowAdd "1h" }},
        "color": "#FFFFFF",
        "type": "test-maintenance-type-to-engine",
        "reason": "test-reason-to-engine"
      }
    ]
    """
    Then the response code should be 207
    When I do GET /api/v4/weather-services/test-pbehavior-weather-service-entity-8?pbh_origin=pbehavior-weather-service-entity-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-pbehavior-weather-service-entity-8",
          "state": {"val": 2},
          "status": {"val": 1},
          "is_grey": true,
          "icon": "pause",
          "pbh_origin_icon": "build",
          "pbehavior_info": {
            "canonical_type": "pause",
            "name": "test-pbehavior-weather-service-entity-8-1",
            "icon_name": "pause"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-entity-8-1",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "reason": {
                "_id": "test-reason-to-engine",
                "name": "Test Engine",
                "description": "Test Engine"
              },
              "type": {
                "_id": "test-pause-type-to-engine",
                "icon_name": "pause",
                "name": "Engine pause",
                "type": "pause"
              }
            }
          ]
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
