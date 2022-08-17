Feature: Get alarms
  I need to be able to get a alarms

  Scenario: given get search request should return alarms only
  with string in connector, connector_name, component or resource fields
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-get
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-4",
          "assigned_instructions": [
            {
              "_id": "test-instruction-with-entity-pattern-1",
              "execution": null,
              "name": "test-instruction-with-entity-pattern-1-name"
            },
            {
              "_id": "test-instruction-with-entity-pattern-2",
              "execution": null,
              "name": "test-instruction-with-entity-pattern-2-name"
            },
            {
              "_id": "test-instruction-with-patterns-combined",
              "execution": null,
              "name": "test-instruction-with-patterns-combined-name"
            }
          ],
          "entity": {
            "_id": "test-alarm-get-resource-4/test-alarm-get-component",
            "category": null,
            "component": "test-alarm-get-component",
            "depends": [
              "test-alarm-get-connector/test-alarm-get-connectorname"
            ],
            "enabled": true,
            "impact": [
              "test-alarm-get-component"
            ],
            "impact_level": 1,
            "infos": {},
            "measurements": null,
            "name": "test-alarm-get-resource-4",
            "type": "resource"
          },
          "impact_state": 1,
          "infos": {},
          "t": 1597030222,
          "v": {
            "ack": {
              "_t": "ack",
              "a": "root",
              "m": "",
              "t": 1597030351,
              "val": 0
            },
            "canceled": {
              "_t": "cancel",
              "a": "root",
              "m": "Test",
              "t": 1597030366,
              "val": 4
            },
            "children": [],
            "component": "test-alarm-get-component",
            "connector": "test-alarm-get-connector",
            "connector_name": "test-alarm-get-connectorname",
            "creation_date": 1597030222,
            "display_name": "PK-SL-XK",
            "infos": {},
            "infos_rule_version": {},
            "initial_long_output": "",
            "initial_output": "test-alarm-get alarm",
            "last_event_date": 1597030222,
            "last_update_date": 1597030366,
            "long_output": "",
            "long_output_history": [
              ""
            ],
            "output": "test-alarm-get criticité",
            "parents": [
              "test-alarm-get-entity-meta-1/metaalarm"
            ],
            "resolved": 1597034023,
            "duration": 3801,
            "current_state_duration": 3801,
            "resource": "test-alarm-get-resource-4",
            "state": {
              "_t": "stateinc",
              "a": "test-alarm-get-connector.test-alarm-get-connectorname",
              "m": "test-alarm-get alarm",
              "t": 1597030222,
              "val": 1
            },
            "status": {
              "_t": "cancel",
              "a": "root",
              "m": "Test",
              "t": 1597030366,
              "val": 4
            },
            "tags": [],
            "total_state_changes": 1
          }
        },
        {
          "_id": "test-alarm-get-3",
          "assigned_instructions": [
            {
              "_id": "test-instruction-with-entity-pattern-1",
              "execution": null,
              "name": "test-instruction-with-entity-pattern-1-name"
            },
            {
              "_id": "test-instruction-with-entity-pattern-2",
              "execution": null,
              "name": "test-instruction-with-entity-pattern-2-name"
            },
            {
              "_id": "test-instruction-with-patterns-combined",
              "execution": null,
              "name": "test-instruction-with-patterns-combined-name"
            }
          ],
          "entity": {
            "_id": "test-alarm-get-resource-3/test-alarm-get-component",
            "category": null,
            "component": "test-alarm-get-component",
            "depends": [
              "test-alarm-get-connector/test-alarm-get-connectorname"
            ],
            "enabled": true,
            "impact": [
              "test-alarm-get-component"
            ],
            "impact_level": 1,
            "infos": {},
            "measurements": null,
            "name": "test-alarm-get-resource-3",
            "type": "resource"
          },
          "impact_state": 1,
          "infos": {},
          "t": 1597030220,
          "v": {
            "children": [],
            "component": "test-alarm-get-component",
            "connector": "test-alarm-get-connector",
            "connector_name": "test-alarm-get-connectorname",
            "creation_date": 1597030220,
            "display_name": "ZD-SY-RM",
            "infos": {},
            "infos_rule_version": {},
            "initial_long_output": "",
            "initial_output": "test-alarm-get alarm",
            "last_event_date": 1597030220,
            "last_update_date": 1597030220,
            "long_output": "",
            "long_output_history": [
              ""
            ],
            "output": "test-alarm-get alarm",
            "parents": [
              "test-alarm-get-entity-meta-1/metaalarm"
            ],
            "resource": "test-alarm-get-resource-3",
            "state": {
              "_t": "stateinc",
              "a": "test-alarm-get-connector.test-alarm-get-connectorname",
              "m": "test-alarm-get alarm",
              "t": 1597030220,
              "val": 1
            },
            "status": {
              "_t": "statusinc",
              "a": "test-alarm-get-connector.test-alarm-get-connectorname",
              "m": "test-alarm-get alarm",
              "t": 1597030220,
              "val": 1
            },
            "tags": [],
            "total_state_changes": 1
          }
        },
        {
          "_id": "test-alarm-get-2",
          "assigned_instructions": [
            {
              "_id": "test-instruction-with-entity-pattern-1",
              "execution": null,
              "name": "test-instruction-with-entity-pattern-1-name"
            },
            {
              "_id": "test-instruction-with-entity-pattern-2",
              "execution": null,
              "name": "test-instruction-with-entity-pattern-2-name"
            },
            {
              "_id": "test-instruction-with-patterns-combined",
              "execution": null,
              "name": "test-instruction-with-patterns-combined-name"
            }
          ],
          "entity": {
            "_id": "test-alarm-get-resource-2/test-alarm-get-component",
            "category": null,
            "component": "test-alarm-get-component",
            "depends": [
              "test-alarm-get-connector/test-alarm-get-connectorname"
            ],
            "enabled": true,
            "impact": [
              "test-alarm-get-component"
            ],
            "impact_level": 1,
            "infos": {},
            "measurements": null,
            "name": "test-alarm-get-resource-2",
            "type": "resource"
          },
          "impact_state": 1,
          "infos": {},
          "t": 1597030219,
          "v": {
            "children": [],
            "component": "test-alarm-get-component",
            "connector": "test-alarm-get-connector",
            "connector_name": "test-alarm-get-connectorname",
            "creation_date": 1597030219,
            "display_name": "PU-YA-QB",
            "infos": {},
            "infos_rule_version": {},
            "initial_long_output": "",
            "initial_output": "test-alarm-get alarm",
            "last_event_date": 1597030219,
            "last_update_date": 1597030219,
            "long_output": "test-alarm-get-correlation-search",
            "long_output_history": [
              ""
            ],
            "output": "test-alarm-get alarm",
            "parents": [
              "test-alarm-get-entity-meta-1/metaalarm",
              "test-alarm-get-entity-meta-2/metaalarm"
            ],
            "resource": "test-alarm-get-resource-2",
            "state": {
              "_t": "stateinc",
              "a": "test-alarm-get-connector.test-alarm-get-connectorname",
              "m": "test-alarm-get alarm",
              "t": 1597030219,
              "val": 1
            },
            "status": {
              "_t": "statusinc",
              "a": "test-alarm-get-connector.test-alarm-get-connectorname",
              "m": "test-alarm-get alarm",
              "t": 1597030219,
              "val": 1
            },
            "tags": [],
            "total_state_changes": 1
          }
        },
        {
          "_id": "test-alarm-get-1",
          "assigned_instructions": [
            {
              "_id": "test-instruction-with-entity-pattern-1",
              "execution": null,
              "name": "test-instruction-with-entity-pattern-1-name"
            },
            {
              "_id": "test-instruction-with-entity-pattern-2",
              "execution": null,
              "name": "test-instruction-with-entity-pattern-2-name"
            },
            {
              "_id": "test-instruction-with-patterns-combined",
              "execution": null,
              "name": "test-instruction-with-patterns-combined-name"
            }
          ],
          "entity": {
            "_id": "test-alarm-get-resource-1/test-alarm-get-component",
            "category": {
              "_id": "test-category-to-alarm-get-1",
              "name": "test-category-to-alarm-get-1-name",
              "author": "test-category-to-alarm-get-1-author",
              "created": 1592215337,
              "updated": 1592215337
            },
            "component": "test-alarm-get-component",
            "depends": [
              "test-alarm-get-connector/test-alarm-get-connectorname"
            ],
            "enabled": true,
            "impact": [
              "test-alarm-get-component"
            ],
            "impact_level": 1,
            "infos": {
              "test-alarm-get-resource-1-info-1": {
                "name": "test-alarm-get-resource-1-info-1-name",
                "description": "test-alarm-get-resource-1-info-1-description",
                "value": "test-alarm-get-resource-1-info-1-value"
              },
              "test-alarm-get-resource-1-info-2": {
                "name": "test-alarm-get-resource-1-info-2-name",
                "description": "test-alarm-get-resource-1-info-2-description",
                "value": false
              },
              "test-alarm-get-resource-1-info-3": {
                "name": "test-alarm-get-resource-1-info-3-name",
                "description": "test-alarm-get-resource-1-info-3-description",
                "value": 1022
              },
              "test-alarm-get-resource-1-info-4": {
                "name": "test-alarm-get-resource-1-info-4-name",
                "description": "test-alarm-get-resource-1-info-4-description",
                "value": 10.45
              },
              "test-alarm-get-resource-1-info-5": {
                "name": "test-alarm-get-resource-1-info-5-name",
                "description": "test-alarm-get-resource-1-info-5-description",
                "value": null
              },
              "test-alarm-get-resource-1-info-6": {
                "name": "test-alarm-get-resource-1-info-6-name",
                "description": "test-alarm-get-resource-1-info-6-description",
                "value": ["test-alarm-get-resource-1-info-6-value", false, 1022, 10.45, null]
              },
              "test-alarm-get-resource-1-info-7": {
                "name": "test-alarm-get-resource-1-info-7",
                "description": "test-alarm-get-resource-1-info-7-description",
                "value": "test-alarm-get-resource-1-info-7-value"
              }
            },
            "measurements": null,
            "name": "test-alarm-get-resource-1",
            "type": "resource"
          },
          "impact_state": 3,
          "infos": {},
          "t": 1596942720,
          "v": {
            "children": [],
            "component": "test-alarm-get-component",
            "connector": "test-alarm-get-connector",
            "connector_name": "test-alarm-get-connectorname",
            "creation_date": 1596942720,
            "display_name": "RC-KC_tW",
            "infos": {},
            "infos_rule_version": {},
            "initial_long_output": "",
            "initial_output": "test-alarm-get alarm",
            "last_comment": {
              "_t": "comment",
              "a": "system",
              "m": "comment test-alarm-get",
              "t": 1597033041,
              "val": 0
            },
            "last_event_date": 1596942720,
            "last_update_date": 1596942720,
            "long_output": "test-alarm-get-correlation-search",
            "long_output_history": [
              ""
            ],
            "output": "test-alarm-get alarm",
            "parents": [],
            "resource": "test-alarm-get-resource-1",
            "state": {
              "_t": "stateinc",
              "a": "test-alarm-get-connector.test-alarm-get-connectorname",
              "m": "test-alarm-get alarm",
              "t": 1596942720,
              "val": 3
            },
            "status": {
              "_t": "statusinc",
              "a": "test-alarm-get-connector.test-alarm-get-connectorname",
              "m": "test-alarm-get alarm",
              "t": 1596942720,
              "val": 1
            },
            "tags": [],
            "total_state_changes": 1
          }
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

  Scenario: given get sort request should return sorted alarms
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-get&sort_key=status&sort_dir=desc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-4"
        },
        {
          "_id": "test-alarm-get-1"
        },
        {
          "_id": "test-alarm-get-2"
        },
        {
          "_id": "test-alarm-get-3"
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

  Scenario: given get sort request should return sorted alarms by duration asc
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-get&sort_key=v.duration&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-4"
        },
        {
          "_id": "test-alarm-get-3"
        },
        {
          "_id": "test-alarm-get-2"
        },
        {
          "_id": "test-alarm-get-1"
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

  Scenario: given get sort request should return sorted alarms by duration desc
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-get&sort_key=v.duration&sort_dir=desc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-1"
        },
        {
          "_id": "test-alarm-get-2"
        },
        {
          "_id": "test-alarm-get-3"
        },
        {
          "_id": "test-alarm-get-4"
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

  Scenario: given get multi sort request should return sorted alarms by t and last_event_date
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-multi-sort-get&multi_sort[]=t,asc&multi_sort[]=v.last_event_date,desc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-multi-sort-get-1",
          "t": 1000000000,
          "v": {
            "last_event_date": 1000000000
          }
        },
        {
          "_id": "test-alarm-multi-sort-get-2",
          "t": 1000000001,
          "v": {
            "last_event_date": 1000000003
          }
        },
        {
          "_id": "test-alarm-multi-sort-get-3",
          "t": 1000000001,
          "v": {
            "last_event_date": 1000000002
          }
        },
        {
          "_id": "test-alarm-multi-sort-get-4",
          "t": 1000000002,
          "v": {
            "last_event_date": 1000000004
          }
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
    When I do GET /api/v4/alarms?search=test-alarm-multi-sort-get&multi_sort[]=t,asc&multi_sort[]=v.last_event_date,asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-multi-sort-get-1",
          "t": 1000000000,
          "v": {
            "last_event_date": 1000000000
          }
        },
        {
          "_id": "test-alarm-multi-sort-get-3",
          "t": 1000000001,
          "v": {
            "last_event_date": 1000000002
          }
        },
        {
          "_id": "test-alarm-multi-sort-get-2",
          "t": 1000000001,
          "v": {
            "last_event_date": 1000000003
          }
        },
        {
          "_id": "test-alarm-multi-sort-get-4",
          "t": 1000000002,
          "v": {
            "last_event_date": 1000000004
          }
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

  Scenario: given get multi sort with simple sort is not allowed
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-multi-sort-get&multi_sort[]=t,asc&multi_sort[]=v.last_event_date,desc&sort_key=v.duration&sort_dir=desc
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "sort_key": "Can't be present both SortBy and MultiSort."
      }
    }
    """

  Scenario: given get with invalid multi sort
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-multi-sort-get&multi_sort[]=t,asc&multi_sort[]=v.last_event_date
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "multi_sort[]": "Invalid multi_sort value."
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-alarm-multi-sort-get&multi_sort[]=t,asc&multi_sort[]=v.last_event_date,bad
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "multi_sort[]": "Invalid multi_sort value."
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-alarm-multi-sort-get&multi_sort[]=t,asc&multi_sort[]=v.last_event_date,desc,extra
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "multi_sort[]": "Invalid multi_sort value."
      }
    }
    """

  Scenario: given get time inverval request should return alarms which were created
  in this time interval.
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-get&tstart=1596931200&tstop=1597017600
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-1"
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

  Scenario: given get opened request should return only opened alarms
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-get&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-3"
        },
        {
          "_id": "test-alarm-get-2"
        },
        {
          "_id": "test-alarm-get-1"
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

  Scenario: given get recent resolved request should return only recent resolved alarms and opened
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-get
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-4"
        },
        {
          "_id": "test-alarm-get-3"
        },
        {
          "_id": "test-alarm-get-2"
        },
        {
          "_id": "test-alarm-get-1"
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

  Scenario: given get resolved request should return only resolved alarms
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-get&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-resolved-collection"
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

  Scenario: given get filter request should return alarms which are matched to the filter
    When I am admin
    When I do GET /api/v4/alarms?filter={"$or":[{"uid":"test-alarm-get-2"},{"uid":"test-alarm-get-4"}]}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-4"
        },
        {
          "_id": "test-alarm-get-2"
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

  Scenario: given get with_steps request should return alarms with steps
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-get&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-4",
          "v": {
            "steps": [
              {
                  "_t": "stateinc",
                  "t": 1597030222,
                  "a": "test-alarm-get-connector.test-alarm-get-connectorname",
                  "m": "test-alarm-get alarm",
                  "val": 1
              },
              {
                  "_t": "statusinc",
                  "t": 1597030222,
                  "a": "test-alarm-get-connector.test-alarm-get-connectorname",
                  "m": "test-alarm-get alarm",
                  "val": 1
              },
              {
                  "_t": "ack",
                  "t": 1597030351,
                  "a": "root",
                  "m": "",
                  "val": 0
              },
              {
                  "_t": "cancel",
                  "t": 1597030366,
                  "a": "root",
                  "m": "Test",
                  "val": 0
              },
              {
                  "_t": "statusinc",
                  "t": 1597030366,
                  "a": "test-alarm-get-connector.test-alarm-get-connectorname",
                  "m": "Test",
                  "val": 4
              }
            ]
          }
        },
        {
          "_id": "test-alarm-get-3",
          "v": {
            "steps": [
              {
                  "_t": "stateinc",
                  "t": 1597030220,
                  "a": "test-alarm-get-connector.test-alarm-get-connectorname",
                  "m": "test-alarm-get alarm",
                  "val": 1
              },
              {
                  "_t": "statusinc",
                  "t": 1597030220,
                  "a": "test-alarm-get-connector.test-alarm-get-connectorname",
                  "m": "test-alarm-get alarm",
                  "val": 1
              }
            ]
          }
        },
        {
          "_id": "test-alarm-get-2",
          "v": {
            "steps": [
              {
                  "_t": "stateinc",
                  "t": 1597030219,
                  "a": "test-alarm-get-connector.test-alarm-get-connectorname",
                  "m": "test-alarm-get alarm",
                  "val": 1
              },
              {
                  "_t": "statusinc",
                  "t": 1597030219,
                  "a": "test-alarm-get-connector.test-alarm-get-connectorname",
                  "m": "test-alarm-get alarm",
                  "val": 1
              }
            ]
          }
        },
        {
          "_id": "test-alarm-get-1",
          "v": {
            "steps": [
              {
                  "_t": "stateinc",
                  "t": 1596942720,
                  "a": "test-alarm-get-connector.test-alarm-get-connectorname",
                  "m": "test-alarm-get alarm",
                  "val": 3
              },
              {
                  "_t": "statusinc",
                  "t": 1596942720,
                  "a": "test-alarm-get-connector.test-alarm-get-connectorname",
                  "m": "test-alarm-get alarm",
                  "val": 1
              },
              {
                  "_t": "comment",
                  "t": 1597033041,
                  "a": "system",
                  "m": "comment test-alarm-get",
                  "val": 0
              }
            ]
          }
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

  Scenario: given get search by infos.*.value field request should return alarms
  only with string in entity.infos.*.value field
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-get-resource-1-info-1&active_columns[]=infos.test-alarm-get-resource-1-info-1.value
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-1"
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

  Scenario: given get correlation request should return meta alarms or alarms without parent
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-get&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-meta-1",
          "metaalarm": true,
          "consequences": {
              "total": 3
          },
          "rule": {
            "id": "test-alarm-get-metaalarm-rule-1",
            "name": "Test alarm get"
          }
        },
        {
          "_id": "test-alarm-get-meta-2",
          "metaalarm": true,
          "consequences": {
              "total": 1
          },
          "rule": {
            "id": "test-alarm-get-metaalarm-rule-2",
            "name": "Test alarm get"
          }
        },
        {
          "_id": "test-alarm-get-meta-manual-1",
          "metaalarm": true,
          "consequences": {
              "total": 0
          },
          "rule": {
            "id": "test-alarm-get-metaalarm-rule-manual-1",
            "name": "Test manual 1"
          }
        },
        {
          "_id": "test-alarm-get-meta-manual-2",
          "metaalarm": true,
          "consequences": {
              "total": 0
          },
          "rule": {
            "id": "test-alarm-get-metaalarm-rule-manual-1",
            "name": "Test manual 1"
          }
        },
        {
          "_id": "test-alarm-get-1",
          "metaalarm": false
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

  Scenario: given get correlation request should return only manual meta alarms
    When I am admin
    When I do GET /api/v4/alarms?correlation=true&manual=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-meta-manual-1",
          "metaalarm": true,
          "consequences": {
              "total": 0
          },
          "rule": {
            "id": "test-alarm-get-metaalarm-rule-manual-1",
            "name": "Test manual 1"
          }
        },
        {
          "_id": "test-alarm-get-meta-manual-2",
          "metaalarm": true,
          "consequences": {
              "total": 0
          },
          "rule": {
            "id": "test-alarm-get-metaalarm-rule-manual-1",
            "name": "Test manual 1"
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

  Scenario: given get correlation with_consequences request should return
  meta alarms with children or alarms without parent
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-get&correlation=true&with_consequences=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-meta-1",
          "metaalarm": true,
          "consequences": {
              "data": [
                {
                  "_id": "test-alarm-get-4"
                },
                {
                  "_id": "test-alarm-get-3"
                },
                {
                  "_id": "test-alarm-get-2"
                }
              ],
              "total": 3
          },
          "rule": {
            "id": "test-alarm-get-metaalarm-rule-1",
            "name": "Test alarm get"
          }
        },
        {
          "_id": "test-alarm-get-meta-2",
          "metaalarm": true,
          "consequences": {
            "data": [
              {
                "_id": "test-alarm-get-2"
              }
            ],
            "total": 1
          },
          "rule": {
            "id": "test-alarm-get-metaalarm-rule-2",
            "name": "Test alarm get"
          }
        },
        {
          "_id": "test-alarm-get-meta-manual-1",
          "metaalarm": true,
          "consequences": {
              "total": 0
          },
          "rule": {
            "id": "test-alarm-get-metaalarm-rule-manual-1",
            "name": "Test manual 1"
          }
        },
        {
          "_id": "test-alarm-get-meta-manual-2",
          "metaalarm": true,
          "consequences": {
            "total": 0
          },
          "rule": {
            "id": "test-alarm-get-metaalarm-rule-manual-1",
            "name": "Test manual 1"
          }
        },
        {
          "_id": "test-alarm-get-1",
          "metaalarm": false
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

  Scenario: given get correlation search by v.long_output request should return
  meta alarms with string in one of children v.long_output field or meta alarms
  and alarms without parents with string in v.long_output field
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-get-correlation-search&correlation=true&active_columns[]=v.long_output
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-meta-1",
          "metaalarm": true,
          "consequences": {
              "total": 3
          },
          "filtered_children": [
            "test-alarm-get-2"
          ]
        },
        {
          "_id": "test-alarm-get-meta-2",
          "metaalarm": true,
          "consequences": {
              "total": 1
          },
          "filtered_children": [
            "test-alarm-get-2"
          ]
        },
        {
          "_id": "test-alarm-get-1",
          "metaalarm": false
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

  Scenario: given get search expression request should return alarms which are matched
  to expression filter
    When I am admin
    When I do GET /api/v4/alarms?search=entity._id="test-alarm-get-resource-1/test-alarm-get-component"
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-1",
          "entity": {
            "_id": "test-alarm-get-resource-1/test-alarm-get-component"
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

  Scenario: given get search expression request should return alarms which are matched
  to expression filter
    When I am admin
    When I do GET /api/v4/alarms?search=resource="test-alarm-get-resource-1"
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-1",
          "entity": {
            "_id": "test-alarm-get-resource-1/test-alarm-get-component"
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

  Scenario: given get search expression request should return alarms which are matched
  to expression filter
    When I am admin
    When I do GET /api/v4/alarms?search=v.output%20LIKE%20"criticité"
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-4",
          "v": {
            "output": "test-alarm-get criticité"
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

  Scenario: given get search expression request should return alarms which are matched
  to expression filter
    When I am admin
    When I do GET /api/v4/alarms?search=entity_id%20LIKE%20"1/test-alarm-get-component$"
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-1"
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

  Scenario: given get search expression request should return alarms which are matched
  to expression filter
    When I am admin
    When I do GET /api/v4/alarms?search=v.parents%20CONTAINS%20"test-alarm-get-entity-meta-1/metaalarm"
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-4"
        },
        {
          "_id": "test-alarm-get-3"
        },
        {
          "_id": "test-alarm-get-2"
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

  Scenario: given get search request should return assigned instruction for the alarm
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-with-instruction-resource-1&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-with-instruction-1",
          "assigned_instructions": [
            {
                "_id": "test-instruction-to-assign",
                "name": "test-instruction-to-assign-name",
                "execution": null
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

  Scenario: given get search request should return assigned instruction, which have an execution for the alarm
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-with-instruction-resource-2&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-with-instruction-2",
          "assigned_instructions": [
            {
              "_id": "test-instruction-to-assign-with-execution",
              "name": "test-instruction-to-assign-with-execution-name",
              "execution": {
                "_id": "execution-for-test-alarm-with-instruction-resource-2",
                "status": 0
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

  Scenario: given get search request should return assigned instruction, which have several executions for the alarm
  where some of them is not executed
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-with-instruction-resource-3&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-with-instruction-3",
          "assigned_instructions": [
            {
              "_id": "test-instruction-to-assign-with-execution-2",
              "execution": {
                "_id": "execution-for-test-alarm-with-instruction-resource-3",
                "status": 0
              },
              "name": "test-instruction-to-assign-with-execution-2-name"
            },
            {
              "_id": "test-instruction-to-assign-without-execution",
              "execution": null,
              "name": "test-instruction-to-assign-without-execution-name"
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

  Scenario: given get search request should return alarms with assigned instructions depending from exclude or include instructions fields
    When I am admin
    When I do GET /api/v4/alarms?include_instruction_types[]=0&include_instruction_types[]=1&search=test-alarm-with-instruction
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-with-instruction-1"
        },
        {
          "_id": "test-alarm-with-instruction-2"
        },
        {
          "_id": "test-alarm-with-instruction-3"
        },
        {
          "_id": "test-alarm-with-instruction-4"
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
    When I do GET /api/v4/alarms?include_instruction_types[]=0&&search=test-alarm-with-instruction
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-with-instruction-1"
        },
        {
          "_id": "test-alarm-with-instruction-2"
        },
        {
          "_id": "test-alarm-with-instruction-3"
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
    When I do GET /api/v4/alarms?include_instruction_types[]=1&search=test-alarm-with-instruction
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-with-instruction-4"
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
    When I do GET /api/v4/alarms?exclude_instruction_types[]=0&exclude_instruction_types[]=1&search=test-alarm-with-instruction
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
    When I do GET /api/v4/alarms?exclude_instruction_types[]=0&search=test-alarm-with-instruction
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-with-instruction-4"
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
    When I do GET /api/v4/alarms?exclude_instruction_types[]=1&search=test-alarm-with-instruction
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-with-instruction-1"
        },
        {
          "_id": "test-alarm-with-instruction-2"
        },
        {
          "_id": "test-alarm-with-instruction-3"
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
    When I do GET /api/v4/alarms?exclude_instruction_types[]=0&exclude_instruction_types[]=1&search=instruction-not-exists
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
    When I do GET /api/v4/alarms?include_instructions[]=test-instruction-to-assign&search=test-alarm-with-instruction
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-with-instruction-1"
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
    When I do GET /api/v4/alarms?include_instructions[]=test-instruction-to-assign&include_instructions[]=test-instruction-to-assign-with-execution&search=test-alarm-with-instruction
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-with-instruction-1"
        },
        {
          "_id": "test-alarm-with-instruction-2"
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
    When I do GET /api/v4/alarms?include_instruction_types[]=0&include_instruction_types[]=1&exclude_instructions[]=test-instruction-to-assign&search=test-alarm-with-instruction
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-with-instruction-2"
        },
        {
          "_id": "test-alarm-with-instruction-3"
        },
        {
          "_id": "test-alarm-with-instruction-4"
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
    When I do GET /api/v4/alarms?include_instruction_types[]=0&include_instruction_types[]=1&exclude_instructions[]=test-instruction-to-assign&exclude_instructions[]=test-instruction-to-assign-with-execution&exclude_instructions[]=test-instruction-to-auto-instruction-filter&search=test-alarm-with-instruction
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-with-instruction-3"
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
    When I do GET /api/v4/alarms?include_instructions[]=test-instruction-with-entity-pattern-1
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-4"
        },
        {
          "_id": "test-alarm-get-3"
        },
        {
          "_id": "test-alarm-get-2"
        },
        {
          "_id": "test-alarm-get-1"
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
    When I do GET /api/v4/alarms?include_instructions[]=test-instruction-with-entity-pattern-2
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-4"
        },
        {
          "_id": "test-alarm-get-3"
        },
        {
          "_id": "test-alarm-get-2"
        },
        {
          "_id": "test-alarm-get-1"
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
    When I do GET /api/v4/alarms?include_instructions[]=test-instruction-with-patterns-combined
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-get-4"
        },
        {
          "_id": "test-alarm-get-3"
        },
        {
          "_id": "test-alarm-get-2"
        },
        {
          "_id": "test-alarm-get-1"
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

  Scenario: get alarm with idle_since
    When I am admin
    When I do GET /api/v4/alarms?search=test-idle-since
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-idle-since-1",
          "entity": {
            "_id": "test-idle-since-resource-1/test-idle-since-component",
            "idle_since": 123
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

  Scenario: get alarm with events_count
    When I am admin
    When I do GET /api/v4/alarms?search=test-events-count
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-events-count-1",
          "entity": {
            "_id": "test-events-count-resource-1/test-events-count-component"
          },
          "v": {
            "events_count": 5
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

  Scenario: given get correlation with_instructions request should return
  meta alarms with children, children should have assigned instruction if they have it, the corresponding
  metaalarm should have a mark about it.
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-meta-children-with-instructions&correlation=true&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-meta-children-with-instructions-1",
          "metaalarm": true,
          "children_instructions": true,
          "consequences": {
            "data": [
              {
                "_id": "test-alarm-meta-children-with-instructions-alarm-1-1",
                "assigned_instructions": [
                  {
                    "_id": "test-alarm-meta-children-with-instructions"
                  }
                ]
              }
            ],
            "total": 1
          }
        },
        {
          "_id": "test-alarm-meta-children-with-instructions-2",
          "metaalarm": true,
          "children_instructions": false,
          "consequences": {
            "data": [
              {
                "_id": "test-alarm-meta-children-with-instructions-alarm-2-1"
              }
            ],
            "total": 1
          }
        },
        {
          "_id": "test-alarm-meta-children-with-instructions-3",
          "metaalarm": true,
          "children_instructions": true,
          "consequences": {
            "total": 2
          }
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

  Scenario: given get correlation without with_instructions request should return
  meta alarms without mark that children has assigned instructions even if they have it.
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-meta-children-with-instructions&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-meta-children-with-instructions-1",
          "metaalarm": true,
          "children_instructions": false
        },
        {
          "_id": "test-alarm-meta-children-with-instructions-2",
          "metaalarm": true,
          "children_instructions": false
        },
        {
          "_id": "test-alarm-meta-children-with-instructions-3",
          "metaalarm": true,
          "children_instructions": false
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

  Scenario: given get time inverval request should return alarms by default t field
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&tstart=2000000000&tstop=2000000010&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-time-fields-search-get-1"
        },
        {
          "_id": "test-alarm-time-fields-search-get-2"
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
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&tstart=2000000010&tstop=2000000020&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-time-fields-search-get-2"
        },
        {
          "_id": "test-alarm-time-fields-search-get-3"
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

  Scenario: given get time inverval request should return alarms by time_field
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&tstart=1900000000&tstop=1900000010&sort_key=v.resource&sort_dir=asc
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
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&time_field=t&tstart=2000000000&tstop=2000000010&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-time-fields-search-get-1"
        },
        {
          "_id": "test-alarm-time-fields-search-get-2"
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
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&time_field=t&tstart=2000000010&tstop=2000000020&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-time-fields-search-get-2"
        },
        {
          "_id": "test-alarm-time-fields-search-get-3"
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
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&time_field=creation_date&tstart=1900000000&tstop=1900000010&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-time-fields-search-get-1"
        },
        {
          "_id": "test-alarm-time-fields-search-get-2"
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
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&opened=true&time_field=creation_date&tstart=1900000010&tstop=1900000020&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-time-fields-search-get-2"
        },
        {
          "_id": "test-alarm-time-fields-search-get-3"
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
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&opened=true&tstart=1800000000&tstop=1800000010&sort_key=v.resource&sort_dir=asc
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
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&opened=true&time_field=creation_date&tstart=1800000000&tstop=1800000010&sort_key=v.resource&sort_dir=asc
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
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&opened=true&time_field=last_update_date&tstart=1800000000&tstop=1800000010&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-time-fields-search-get-1"
        },
        {
          "_id": "test-alarm-time-fields-search-get-2"
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
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&opened=true&time_field=last_update_date&tstart=1800000010&tstop=1800000020&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-time-fields-search-get-2"
        },
        {
          "_id": "test-alarm-time-fields-search-get-3"
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
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&opened=true&tstart=1700000000&tstop=1700000010&sort_key=v.resource&sort_dir=asc
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
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&opened=true&time_field=creation_date&tstart=1700000000&tstop=1700000010&sort_key=v.resource&sort_dir=asc
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
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&opened=true&time_field=last_update_date&tstart=1700000000&tstop=1700000010&sort_key=v.resource&sort_dir=asc
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
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&opened=true&time_field=last_event_date&tstart=1700000000&tstop=1700000010&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-time-fields-search-get-1"
        },
        {
          "_id": "test-alarm-time-fields-search-get-2"
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
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&opened=true&time_field=last_event_date&tstart=1700000010&tstop=1700000020&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-time-fields-search-get-2"
        },
        {
          "_id": "test-alarm-time-fields-search-get-3"
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
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&time_field=resolved&tstart=2100000000&tstop=2100000010&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-time-fields-search-get-4"
        },
        {
          "_id": "test-alarm-time-fields-search-get-5"
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
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&time_field=resolved&tstart=2100000010&tstop=2100000020&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-time-fields-search-get-5"
        },
        {
          "_id": "test-alarm-time-fields-search-get-6"
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

  Scenario: given get time inverval request for closed alarms should return alarms by default resolved field
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&opened=false&tstart=2200000000&tstop=2200000010&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-time-fields-search-get-resolved-1"
        },
        {
          "_id": "test-alarm-time-fields-search-get-resolved-2"
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
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&opened=false&tstart=2200000010&tstop=2200000020&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-time-fields-search-get-resolved-2"
        },
        {
          "_id": "test-alarm-time-fields-search-get-resolved-3"
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

  Scenario: given get time inverval request with bad time_field should return error
    When I am admin
    When I do GET /api/v4/alarms?search=test-alarm-time-fields-search-get&time_field=test&tstart=2200000010&tstop=2200000020&sort_key=v.resource&sort_dir=asc
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "time_field": "TimeField must be one of [t creation_date resolved last_update_date last_event_date] or empty."
      }
    }
    """

  Scenario: given get correlation alarms with the same children, but different alarms shouldn't have alarms of each other
    When I am admin
    When I do GET /api/v4/alarms?sort_key=t&sort_dir=desc&correlation=true&with_steps=true&with_consequences=true&search=test-api-correlation
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-alarm-api-correlation-metaalarm-1",
          "metaalarm": true,
          "consequences": {
            "data": [
              {
                "_id": "test-api-correlation-get-1"
              },
              {
                "_id": "test-api-correlation-get-2"
              }
            ],
            "total": 2
          }
        },
        {
          "_id": "test-alarm-api-correlation-metaalarm-2",
          "metaalarm": true,
          "consequences": {
            "data": [
              {
                "_id": "test-api-correlation-get-3"
              },
              {
                "_id": "test-api-correlation-get-4"
              }
            ],
            "total": 2
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

  Scenario: given get unauth request should not allow access
    When I do GET /api/v4/alarms
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/alarms
    Then the response code should be 403
