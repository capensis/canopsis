Feature: Get alarms
  I need to be able to get a alarms

  Scenario: given get opened and recently resolved alarms request should return alarms
    When I am admin
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-3",
          "entity": {
            "_id": "test-resource-to-alarm-get-3/test-component-to-alarm-get",
            "category": null,
            "component": "test-component-to-alarm-get",
            "depends": [
              "test-connector-default/test-connector-default-name"
            ],
            "enabled": true,
            "impact": [
              "test-component-to-alarm-get"
            ],
            "impact_level": 1,
            "infos": {},
            "measurements": null,
            "name": "test-resource-to-alarm-get-3",
            "type": "resource"
          },
          "impact_state": 0,
          "infos": {},
          "t": 1597030221,
          "v": {
            "children": [],
            "component": "test-component-to-alarm-get",
            "connector": "test-connector-default",
            "connector_name": "test-connector-default-name",
            "creation_date": 1597030221,
            "display_name": "PU-YA-QB",
            "infos": {},
            "infos_rule_version": {},
            "initial_long_output": "",
            "initial_output": "test-alarm-to-get-3-output",
            "last_event_date": 1597030221,
            "last_update_date": 1597030221,
            "long_output": "",
            "long_output_history": [""],
            "output": "test-alarm-to-get-3-output",
            "parents": [],
            "resource": "test-resource-to-alarm-get-3",
            "duration": 20,
            "current_state_duration": 0,
            "active_duration": 20,
            "pbh_inactive_duration": 0,
            "snooze_duration": 0,
            "resolved": 1597030241,
            "state": {
              "_t": "statedec",
              "a": "test-connector-default.test-connector-default-name",
              "m": "test-alarm-to-get-3-output",
              "t": 1597030241,
              "initiator": "external",
              "val": 0
            },
            "status": {
              "_t": "statusdec",
              "a": "test-connector-default.test-connector-default-name",
              "m": "test-alarm-to-get-3-output",
              "t": 1597030241,
              "initiator": "external",
              "val": 0
            },
            "total_state_changes": 1
          }
        },
        {
          "_id": "test-alarm-to-get-2",
          "entity": {
            "_id": "test-resource-to-alarm-get-2/test-component-to-alarm-get",
            "category": {
              "_id": "test-category-to-alarm-get-2",
              "name": "test-category-to-alarm-get-2-name",
              "author": "root",
              "created": 1592215337,
              "updated": 1592215337
            },
            "component": "test-component-to-alarm-get",
            "depends": [
              "test-connector-default/test-connector-default-name"
            ],
            "enabled": true,
            "impact": [
              "test-component-to-alarm-get"
            ],
            "impact_level": 1,
            "infos": {},
            "measurements": null,
            "name": "test-resource-to-alarm-get-2",
            "type": "resource"
          },
          "impact_state": 1,
          "infos": {},
          "t": 1597030220,
          "v": {
            "children": [],
            "component": "test-component-to-alarm-get",
            "connector": "test-connector-default",
            "connector_name": "test-connector-default-name",
            "creation_date": 1597030220,
            "display_name": "PU-YA-QB",
            "infos": {},
            "infos_rule_version": {},
            "initial_long_output": "test-alarm-to-get-2-output-long",
            "initial_output": "test-alarm-to-get-2-output",
            "last_event_date": 1597030220,
            "last_update_date": 1597030220,
            "long_output": "test-alarm-to-get-2-output-long",
            "long_output_history": [
              "test-alarm-to-get-2-output-long"
            ],
            "output": "test-alarm-to-get-2-output",
            "parents": [],
            "resource": "test-resource-to-alarm-get-2",
            "pbh_inactive_duration": 0,
            "snooze_duration": 0,
            "state": {
              "_t": "stateinc",
              "a": "test-connector-default.test-connector-default-name",
              "m": "test-alarm-to-get-2-output",
              "t": 1597030220,
              "initiator": "external",
              "val": 1
            },
            "status": {
              "_t": "statusinc",
              "a": "test-connector-default.test-connector-default-name",
              "m": "test-alarm-to-get-2-output",
              "t": 1597030220,
              "initiator": "external",
              "val": 1
            },
            "total_state_changes": 1
          }
        },
        {
          "_id": "test-alarm-to-get-1",
          "entity": {
            "_id": "test-resource-to-alarm-get-1/test-component-to-alarm-get",
            "category": {
              "_id": "test-category-to-alarm-get-1",
              "name": "test-category-to-alarm-get-1-name",
              "author": "root",
              "created": 1592215337,
              "updated": 1592215337
            },
            "component": "test-component-to-alarm-get",
            "depends": [
              "test-connector-default/test-connector-default-name"
            ],
            "enabled": true,
            "impact": [
              "test-component-to-alarm-get"
            ],
            "impact_level": 1,
            "infos": {
              "test-resource-to-alarm-get-1-info-1": {
                "name": "test-resource-to-alarm-get-1-info-1-name",
                "description": "test-resource-to-alarm-get-1-info-1-description",
                "value": "test-resource-to-alarm-get-1-info-1-value"
              },
              "test-resource-to-alarm-get-1-info-2": {
                "name": "test-resource-to-alarm-get-1-info-2-name",
                "description": "test-resource-to-alarm-get-1-info-2-description",
                "value": false
              },
              "test-resource-to-alarm-get-1-info-3": {
                "name": "test-resource-to-alarm-get-1-info-3-name",
                "description": "test-resource-to-alarm-get-1-info-3-description",
                "value": 1022
              },
              "test-resource-to-alarm-get-1-info-4": {
                "name": "test-resource-to-alarm-get-1-info-4-name",
                "description": "test-resource-to-alarm-get-1-info-4-description",
                "value": 10.45
              },
              "test-resource-to-alarm-get-1-info-5": {
                "name": "test-resource-to-alarm-get-1-info-5-name",
                "description": "test-resource-to-alarm-get-1-info-5-description",
                "value": null
              },
              "test-resource-to-alarm-get-1-info-6": {
                "name": "test-resource-to-alarm-get-1-info-6-name",
                "description": "test-resource-to-alarm-get-1-info-6-description",
                "value": ["test-resource-to-alarm-get-1-info-6-value"]
              }
            },
            "measurements": null,
            "name": "test-resource-to-alarm-get-1",
            "type": "resource"
          },
          "impact_state": 3,
          "infos": {},
          "t": 1597030219,
          "v": {
            "children": [],
            "component": "test-component-to-alarm-get",
            "connector": "test-connector-default",
            "connector_name": "test-connector-default-name",
            "creation_date": 1597030219,
            "display_name": "RC-KC_tW",
            "infos": {},
            "infos_rule_version": {},
            "initial_long_output": "",
            "initial_output": "test-alarm-to-get-1-output",
            "last_comment": {
              "_t": "comment",
              "a": "root",
              "m": "test-alarm-to-get-1-comment-2",
              "t": 1597030221,
              "initiator": "user",
              "val": 0
            },
            "last_event_date": 1597030250,
            "last_update_date": 1597030219,
            "long_output": "",
            "long_output_history": [
              ""
            ],
            "output": "test-alarm-to-get-1-output",
            "parents": [],
            "resource": "test-resource-to-alarm-get-1",
            "pbh_inactive_duration": 0,
            "snooze_duration": 0,
            "state": {
              "_t": "stateinc",
              "a": "test-connector-default.test-connector-default-name",
              "m": "test-alarm-to-get-1-output",
              "t": 1597030219,
              "initiator": "external",
              "val": 3
            },
            "status": {
              "_t": "statusinc",
              "a": "test-connector-default.test-connector-default-name",
              "m": "test-alarm-to-get-1-output",
              "t": 1597030219,
              "initiator": "external",
              "val": 1
            },
            "total_state_changes": 1
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

  Scenario: given get opened alarms request should return alarms
    When I am admin
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-2"
        },
        {
          "_id": "test-alarm-to-get-1"
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

  Scenario: given get resolved alarms request should return alarms
    When I am admin
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-3"
        },
        {
          "_id": "test-alarm-to-get-4"
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

  Scenario: given time interval get request should return alarms
    When I am admin
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get&tstart=1597030220&tstop=1597030320
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-3"
        },
        {
          "_id": "test-alarm-to-get-2"
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
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get&opened=true&tstart=1597030220&tstop=1597030320
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-2"
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
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get&opened=false&tstart=1597030220&tstop=1597030320
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-3"
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
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get&time_field=v.last_event_date&tstart=1597030250&tstop=1597030320
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-1"
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
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get&tstart=1597030320&tstop=1597030350
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

  Scenario: given category get request should return alarms
    When I am admin
    When I do GET /api/v4/alarms?category=test-category-to-alarm-get-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-2"
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

  Scenario: given filter get request should return alarms
    When I am admin
    When I do GET /api/v4/alarms?filters[]=test-widgetfilter-to-alarm-get-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-2"
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
    When I do GET /api/v4/alarms?filters[]=test-widgetfilter-to-alarm-get-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-2"
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

  Scenario: given search expression get request should return alarms
    When I am admin
    When I do GET /api/v4/alarms?search=resource%20LIKE%20"to-alarm-get-2"
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-2"
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
    When I do GET /api/v4/alarms?search=entity.name="test-resource-to-alarm-get-2"
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-2"
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

  Scenario: given get sort request should return sorted alarms
    When I am admin
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get&sort_by=v.state.val&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-1"
        },
        {
          "_id": "test-alarm-to-get-2"
        },
        {
          "_id": "test-alarm-to-get-3"
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
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get&sort_by=v.duration&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-3"
        },
        {
          "_id": "test-alarm-to-get-2"
        },
        {
          "_id": "test-alarm-to-get-1"
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
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get&sort_by=v.duration&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-1"
        },
        {
          "_id": "test-alarm-to-get-2"
        },
        {
          "_id": "test-alarm-to-get-3"
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

  Scenario: given get multi sort request should return sorted alarms
    When I am admin
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-multi-sort-get&multi_sort[]=t,asc&multi_sort[]=v.last_event_date,desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-multi-sort-get-2"
        },
        {
          "_id": "test-alarm-to-multi-sort-get-1"
        },
        {
          "_id": "test-alarm-to-multi-sort-get-3"
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
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-multi-sort-get&multi_sort[]=t,asc&multi_sort[]=v.last_event_date,asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-multi-sort-get-2"
        },
        {
          "_id": "test-alarm-to-multi-sort-get-3"
        },
        {
          "_id": "test-alarm-to-multi-sort-get-1"
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

  Scenario: given invalid get request should return error
    When I am admin
    When I do GET /api/v4/alarms?filters[]=not-exist
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "filter": "Filter doesn't exist."
      }
    }
    """
    When I do GET /api/v4/alarms?time_field=unknown
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "time_field": "TimeField must be one of [t v.creation_date v.resolved v.last_update_date v.last_event_date] or empty."
      }
    }
    """
    When I do GET /api/v4/alarms?multi_sort[]=t,asc&multi_sort[]=v.last_event_date,desc&sort_by=v.duration&sort=desc
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "sort_by": "Can't be present both SortBy and MultiSort."
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-alarm-multi-sort-get&multi_sort[]=t,asc&multi_sort[]=v.last_event_date
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "multi_sort": "Invalid multi_sort value."
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-alarm-multi-sort-get&multi_sort[]=t,asc&multi_sort[]=v.last_event_date,bad
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "multi_sort": "Invalid multi_sort value."
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-alarm-multi-sort-get&multi_sort[]=t,asc&multi_sort[]=v.last_event_date,desc,extra
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "multi_sort": "Invalid multi_sort value."
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

  Scenario: given get one request should return alarm
    When I am admin
    When I do GET /api/v4/alarms/test-alarm-to-get-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-alarm-to-get-1",
      "entity": {
        "_id": "test-resource-to-alarm-get-1/test-component-to-alarm-get",
        "category": {
          "_id": "test-category-to-alarm-get-1",
          "name": "test-category-to-alarm-get-1-name",
          "author": "root",
          "created": 1592215337,
          "updated": 1592215337
        },
        "component": "test-component-to-alarm-get",
        "depends": [
          "test-connector-default/test-connector-default-name"
        ],
        "enabled": true,
        "impact": [
          "test-component-to-alarm-get"
        ],
        "impact_level": 1,
        "infos": {
          "test-resource-to-alarm-get-1-info-1": {
            "name": "test-resource-to-alarm-get-1-info-1-name",
            "description": "test-resource-to-alarm-get-1-info-1-description",
            "value": "test-resource-to-alarm-get-1-info-1-value"
          },
          "test-resource-to-alarm-get-1-info-2": {
            "name": "test-resource-to-alarm-get-1-info-2-name",
            "description": "test-resource-to-alarm-get-1-info-2-description",
            "value": false
          },
          "test-resource-to-alarm-get-1-info-3": {
            "name": "test-resource-to-alarm-get-1-info-3-name",
            "description": "test-resource-to-alarm-get-1-info-3-description",
            "value": 1022
          },
          "test-resource-to-alarm-get-1-info-4": {
            "name": "test-resource-to-alarm-get-1-info-4-name",
            "description": "test-resource-to-alarm-get-1-info-4-description",
            "value": 10.45
          },
          "test-resource-to-alarm-get-1-info-5": {
            "name": "test-resource-to-alarm-get-1-info-5-name",
            "description": "test-resource-to-alarm-get-1-info-5-description",
            "value": null
          },
          "test-resource-to-alarm-get-1-info-6": {
            "name": "test-resource-to-alarm-get-1-info-6-name",
            "description": "test-resource-to-alarm-get-1-info-6-description",
            "value": ["test-resource-to-alarm-get-1-info-6-value"]
          }
        },
        "measurements": null,
        "name": "test-resource-to-alarm-get-1",
        "type": "resource"
      },
      "impact_state": 3,
      "infos": {},
      "t": 1597030219,
      "v": {
        "children": [],
        "component": "test-component-to-alarm-get",
        "connector": "test-connector-default",
        "connector_name": "test-connector-default-name",
        "creation_date": 1597030219,
        "display_name": "RC-KC_tW",
        "infos": {},
        "infos_rule_version": {},
        "initial_long_output": "",
        "initial_output": "test-alarm-to-get-1-output",
        "last_comment": {
          "_t": "comment",
          "a": "root",
          "m": "test-alarm-to-get-1-comment-2",
          "t": 1597030221,
          "initiator": "user",
          "val": 0
        },
        "last_event_date": 1597030250,
        "last_update_date": 1597030219,
        "long_output": "",
        "long_output_history": [
          ""
        ],
        "output": "test-alarm-to-get-1-output",
        "parents": [],
        "resource": "test-resource-to-alarm-get-1",
        "pbh_inactive_duration": 0,
        "snooze_duration": 0,
        "state": {
          "_t": "stateinc",
          "a": "test-connector-default.test-connector-default-name",
          "m": "test-alarm-to-get-1-output",
          "t": 1597030219,
          "initiator": "external",
          "val": 3
        },
        "status": {
          "_t": "statusinc",
          "a": "test-connector-default.test-connector-default-name",
          "m": "test-alarm-to-get-1-output",
          "t": 1597030219,
          "initiator": "external",
          "val": 1
        },
        "total_state_changes": 1
      }
    }
    """
    When I do GET /api/v4/alarms/test-alarm-to-get-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-alarm-to-get-4",
      "entity": {
        "_id": "test-resource-to-alarm-get-3/test-component-to-alarm-get",
        "category": null,
        "component": "test-component-to-alarm-get",
        "depends": [
          "test-connector-default/test-connector-default-name"
        ],
        "enabled": true,
        "impact": [
          "test-component-to-alarm-get"
        ],
        "impact_level": 1,
        "infos": {},
        "measurements": null,
        "name": "test-resource-to-alarm-get-3",
        "type": "resource"
      },
      "impact_state": 0,
      "infos": {},
      "t": 1597030121,
      "v": {
        "children": [],
        "component": "test-component-to-alarm-get",
        "connector": "test-connector-default",
        "connector_name": "test-connector-default-name",
        "creation_date": 1597030121,
        "display_name": "PU-YA-QB",
        "infos": {},
        "infos_rule_version": {},
        "initial_long_output": "",
        "initial_output": "test-alarm-to-get-4-output",
        "last_event_date": 1597030121,
        "last_update_date": 1597030121,
        "resolved": 1597030141,
        "long_output": "",
        "long_output_history": [
          ""
        ],
        "output": "test-alarm-to-get-4-output",
        "parents": [],
        "resource": "test-resource-to-alarm-get-3",
        "duration": 20,
        "current_state_duration": 0,
        "active_duration": 20,
        "pbh_inactive_duration": 0,
        "snooze_duration": 0,
        "state": {
          "_t": "statedec",
          "a": "test-connector-default.test-connector-default-name",
          "m": "test-alarm-to-get-4-output",
          "t": 1597030141,
          "initiator": "external",
          "val": 0
        },
        "status": {
          "_t": "statusdec",
          "a": "test-connector-default.test-connector-default-name",
          "m": "test-alarm-to-get-4-output",
          "t": 1597030141,
          "initiator": "external",
          "val": 0
        },
        "total_state_changes": 1
      }
    }
    """

  Scenario: given get one request should return not found error
    When I am admin
    When I do GET /api/v4/alarms/not-exist
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given get one unauth request should not allow access
    When I do GET /api/v4/alarms/not-exist
    Then the response code should be 401

  Scenario: given get one request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/alarms/not-exist
    Then the response code should be 403

  Scenario: given get details request should return alarm
    When I am admin
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "test-alarm-to-get-1",
        "opened": true,
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "test-alarm-to-get-1",
        "steps": {
          "page": 2,
          "limit": 2
        }
      },
      {
        "_id": "test-alarm-to-get-1",
        "steps": {
          "limit": 5
        }
      },
      {
        "_id": "test-alarm-to-get-3",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "test-alarm-to-get-4",
        "opened": false,
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "test-alarm-to-get-1",
        "opened": false,
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "test-alarm-to-get-3",
        "opened": true,
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "not-exist",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "not-exist"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "_id": "test-alarm-to-get-1",
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-1-output",
                "t": 1597030219,
                "initiator": "external",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-1-output",
                "t": 1597030219,
                "initiator": "external",
                "val": 1
              },
              {
                "_t": "comment",
                "a": "root",
                "m": "test-alarm-to-get-1-comment-1",
                "t": 1597030220,
                "initiator": "user"
              },
              {
                "_t": "comment",
                "a": "root",
                "m": "test-alarm-to-get-1-comment-2",
                "t": 1597030221,
                "initiator": "user"
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 4
            }
          }
        }
      },
      {
        "status": 200,
        "_id": "test-alarm-to-get-1",
        "data": {
          "steps": {
            "data": [
              {
                "_t": "comment",
                "a": "root",
                "m": "test-alarm-to-get-1-comment-1",
                "t": 1597030220,
                "initiator": "user"
              },
              {
                "_t": "comment",
                "a": "root",
                "m": "test-alarm-to-get-1-comment-2",
                "t": 1597030221,
                "initiator": "user"
              }
            ],
            "meta": {
              "page": 2,
              "page_count": 2,
              "per_page": 2,
              "total_count": 4
            }
          }
        }
      },
      {
        "status": 200,
        "_id": "test-alarm-to-get-1",
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-1-output",
                "t": 1597030219,
                "initiator": "external",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-1-output",
                "t": 1597030219,
                "initiator": "external",
                "val": 1
              },
              {
                "_t": "comment",
                "a": "root",
                "m": "test-alarm-to-get-1-comment-1",
                "t": 1597030220,
                "initiator": "user"
              },
              {
                "_t": "comment",
                "a": "root",
                "m": "test-alarm-to-get-1-comment-2",
                "t": 1597030221,
                "initiator": "user"
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 5,
              "total_count": 4
            }
          }
        }
      },
      {
        "status": 200,
        "_id": "test-alarm-to-get-3",
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-3-output",
                "t": 1597030221,
                "initiator": "external",
                "val": 1
              },
              {
                "_t": "statusinc",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-3-output",
                "t": 1597030221,
                "initiator": "external",
                "val": 1
              },
              {
                "_t": "statedec",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-3-output",
                "t": 1597030241,
                "initiator": "external",
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-3-output",
                "t": 1597030241,
                "initiator": "external",
                "val": 0
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 4
            }
          }
        }
      },
      {
        "status": 200,
        "_id": "test-alarm-to-get-4",
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-4-output",
                "t": 1597030121,
                "initiator": "external",
                "val": 1
              },
              {
                "_t": "statusinc",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-4-output",
                "t": 1597030121,
                "initiator": "external",
                "val": 1
              },
              {
                "_t": "statedec",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-4-output",
                "t": 1597030141,
                "initiator": "external",
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-4-output",
                "t": 1597030141,
                "initiator": "external",
                "val": 0
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 4
            }
          }
        }
      },
      {
        "status": 404,
        "_id": "test-alarm-to-get-1",
        "error": "Not found"
      },
      {
        "status": 404,
        "_id": "test-alarm-to-get-3",
        "error": "Not found"
      },
      {
        "status": 404,
        "_id": "not-exist",
        "error": "Not found"
      },
      {
        "status": 400,
        "_id": "not-exist",
        "errors": {
          "steps": "Steps is missing.",
          "children": "Children is missing."
        }
      }
    ]
    """

  Scenario: given get details unauth request should not allow access
    When I do POST /api/v4/alarm-details
    Then the response code should be 401

  Scenario: given get details request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/alarm-details
    Then the response code should be 403

  Scenario: given get by component request should return opened resource alarms
    When I am admin
    When I do GET /api/v4/component-alarms?_id=test-component-to-alarm-get&sort_by=v.resource&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-2"
        },
        {
          "_id": "test-alarm-to-get-1"
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

  Scenario: given get by component request should return validation error
    When I am admin
    When I do GET /api/v4/component-alarms
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "_id": "ID is missing."
      }
    }
    """

  Scenario: given get by component request should return not found error
    When I am admin
    When I do GET /api/v4/component-alarms?_id=not-exist
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given get by component unauth request should not allow access
    When I do GET /api/v4/component-alarms
    Then the response code should be 401

  Scenario: given get by component and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/component-alarms
    Then the response code should be 403

  Scenario: given get resolved by entity request should return resolved entity alarms
    When I am admin
    When I do GET /api/v4/resolved-alarms?_id=test-resource-to-alarm-get-3/test-component-to-alarm-get&sort_by=_id&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-4"
        },
        {
          "_id": "test-alarm-to-get-3"
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
    When I do GET /api/v4/resolved-alarms?_id=test-resource-to-alarm-get-3/test-component-to-alarm-get&tstart=1597030241
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-3"
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
    When I do GET /api/v4/resolved-alarms?_id=test-resource-to-alarm-get-3/test-component-to-alarm-get&tstop=1597030141
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-4"
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

  Scenario: given get resolved by entity request should return validation error
    When I am admin
    When I do GET /api/v4/resolved-alarms
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "_id": "ID is missing."
      }
    }
    """

  Scenario: given get resolved by entity request should return not found error
    When I am admin
    When I do GET /api/v4/resolved-alarms?_id=not-exist
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given get resolved by entity unauth request should not allow access
    When I do GET /api/v4/resolved-alarms
    Then the response code should be 401

  Scenario: given get resolved by entity and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/resolved-alarms
    Then the response code should be 403
