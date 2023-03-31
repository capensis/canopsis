Feature: Get alarms
  I need to be able to get a alarms

  @concurrent
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
            "connector": "test-connector-default/test-connector-default-name",
            "component": "test-component-to-alarm-get",
            "enabled": true,
            "impact_level": 1,
            "infos": {},
            "name": "test-resource-to-alarm-get-3",
            "type": "resource",
            "depends_count": 0,
            "impacts_count": 0
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
              "user_id": "",
              "m": "test-alarm-to-get-3-output",
              "t": 1597030241,
              "initiator": "external",
              "val": 0
            },
            "status": {
              "_t": "statusdec",
              "a": "test-connector-default.test-connector-default-name",
              "user_id": "",
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
              "name": "test-category-to-alarm-get-2-name"
            },
            "connector": "test-connector-default/test-connector-default-name",
            "component": "test-component-to-alarm-get",
            "enabled": true,
            "impact_level": 1,
            "infos": {},
            "name": "test-resource-to-alarm-get-2",
            "type": "resource",
            "depends_count": 0,
            "impacts_count": 0
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
              "user_id": "",
              "m": "test-alarm-to-get-2-output",
              "t": 1597030220,
              "initiator": "external",
              "val": 1
            },
            "status": {
              "_t": "statusinc",
              "a": "test-connector-default.test-connector-default-name",
              "user_id": "",
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
              "name": "test-category-to-alarm-get-1-name"
            },
            "connector": "test-connector-default/test-connector-default-name",
            "component": "test-component-to-alarm-get",
            "enabled": true,
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
            "name": "test-resource-to-alarm-get-1",
            "type": "resource",
            "depends_count": 0,
            "impacts_count": 0
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
              "user_id": "",
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
              "user_id": "",
              "m": "test-alarm-to-get-1-output",
              "t": 1597030219,
              "initiator": "external",
              "val": 3
            },
            "status": {
              "_t": "statusinc",
              "a": "test-connector-default.test-connector-default-name",
              "user_id": "",
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

  @concurrent
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

  @concurrent
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

  @concurrent
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

  @concurrent
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

  @concurrent
  Scenario: given tags get request should return alarms
    When I am admin
    When I do GET /api/v4/alarms?tag=test-tag-to-alarm-get-1
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

  @concurrent
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

  @concurrent
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

  @concurrent
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

  @concurrent
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

  @concurrent
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

  @concurrent
  Scenario: given get unauth request should not allow access
    When I do GET /api/v4/alarms
    Then the response code should be 401

  @concurrent
  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/alarms
    Then the response code should be 403

  @concurrent
  Scenario: given tags filter with has_not condition should filter alarms properly
    When I am admin
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-test-tag-filter&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-test-tag-filter-1"
        },
        {
          "_id": "test-alarm-to-test-tag-filter-2"
        },
        {
          "_id": "test-alarm-to-test-tag-filter-3"
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
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-test-tag-filter&filters[]=test-widgetfilter-to-alarm-test-tag-filter&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-test-tag-filter-1"
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

  @concurrent
  Scenario: given get opened alarms request should return alarms with links
    When I am admin
    When I do GET /api/v4/alarms?with_links=true&search=test-resource-to-alarm-link-get&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-link-get-1",
          "links": {
            "test-category-to-alarm-link-get-1": [
              {
                "rule_id": "test-link-rule-to-alarm-link-get-1",
                "label": "test-link-rule-to-alarm-link-get-1-link-1-label",
                "icon_name": "test-link-rule-to-alarm-link-get-1-link-1-icon",
                "url": "http://test-link-rule-to-alarm-link-get-1-link-1-url.com?resources[]=test-resource-to-alarm-link-get-1/test-component-default|test-link-mongo-data-regexp-1-status&"
              },
              {
                "rule_id": "test-link-rule-to-alarm-link-get-1",
                "label": "test-link-rule-to-alarm-link-get-1-link-3-label",
                "icon_name": "test-link-rule-to-alarm-link-get-1-link-3-icon",
                "url": "http://test-link-rule-to-alarm-link-get-1-link-3-url.com?resources[]=test-resource-to-alarm-link-get-1/test-component-default&",
                "single": true
              },
              {
                "rule_id": "test-link-rule-to-alarm-link-get-2",
                "label": "test-link-rule-to-alarm-link-get-2-link-1-label",
                "icon_name": "test-link-rule-to-alarm-link-get-2-link-1-icon",
                "url": "http://test-link-rule-to-alarm-link-get-2-link-1-url.com?info[]=test-resource-to-alarm-link-get-1-info-1-val&"
              },
              {
                "rule_id": "test-link-rule-to-alarm-link-get-3",
                "label": "test-link-rule-to-alarm-link-get-3-link-1-label",
                "icon_name": "test-link-rule-to-alarm-link-get-3-link-1-icon",
                "url": "http://test-link-rule-to-alarm-link-get-3-link-1-url.com?resources[]=test-resource-to-alarm-link-get-1/test-component-default&"
              },
              {
                "rule_id": "test-link-rule-to-alarm-link-get-4",
                "label": "test-link-rule-to-alarm-link-get-4-link-1-label",
                "icon_name": "test-link-rule-to-alarm-link-get-4-link-1-icon",
                "url": "http://test-link-rule-to-alarm-link-get-4-link-1-url.com?info[]=test-resource-to-alarm-link-get-1-info-1-val|test-link-mongo-data-1-status&"
              }
            ],
            "test-category-to-alarm-link-get-2": [
              {
                "rule_id": "test-link-rule-to-alarm-link-get-1",
                "label": "test-link-rule-to-alarm-link-get-1-link-2-label",
                "icon_name": "test-link-rule-to-alarm-link-get-1-link-2-icon",
                "url": "http://test-link-rule-to-alarm-link-get-1-link-2-url.com?resources[]=test-resource-to-alarm-link-get-1/test-component-default&"
              }
            ],
            "test-category-to-alarm-link-get-3": [
              {
                "rule_id": "test-link-rule-to-alarm-link-get-2",
                "label": "test-link-rule-to-alarm-link-get-2-link-2-label",
                "icon_name": "test-link-rule-to-alarm-link-get-2-link-2-icon",
                "url": "http://test-link-rule-to-alarm-link-get-2-link-2-url.com?info[]=test-resource-to-alarm-link-get-1-info-1-val&"
              }
            ],
            "test-category-to-alarm-link-get-4": [
              {
                "rule_id": "test-link-rule-to-alarm-link-get-3",
                "label": "test-link-rule-to-alarm-link-get-3-link-2-label",
                "icon_name": "test-link-rule-to-alarm-link-get-3-link-2-icon",
                "url": "http://test-link-rule-to-alarm-link-get-3-link-2-url.com?resources[]=test-resource-to-alarm-link-get-1/test-component-default&"
              }
            ],
            "test-category-to-alarm-link-get-5": [
              {
                "rule_id": "test-link-rule-to-alarm-link-get-4",
                "label": "test-link-rule-to-alarm-link-get-4-link-2-label",
                "icon_name": "test-link-rule-to-alarm-link-get-4-link-2-icon",
                "url": "http://test-link-rule-to-alarm-link-get-4-link-2-url.com?info[]=test-resource-to-alarm-link-get-1-info-1-val|test-link-mongo-data-1-status&"
              }
            ]
          }
        },
        {
          "_id": "test-alarm-to-link-get-2",
          "links": {
            "test-category-to-alarm-link-get-1": [
              {
                "rule_id": "test-link-rule-to-alarm-link-get-2",
                "label": "test-link-rule-to-alarm-link-get-2-link-1-label",
                "icon_name": "test-link-rule-to-alarm-link-get-2-link-1-icon",
                "url": "http://test-link-rule-to-alarm-link-get-2-link-1-url.com?info[]=test-resource-to-alarm-link-get-2-info-1-val&"
              },
              {
                "rule_id": "test-link-rule-to-alarm-link-get-4",
                "label": "test-link-rule-to-alarm-link-get-4-link-1-label",
                "icon_name": "test-link-rule-to-alarm-link-get-4-link-1-icon",
                "url": "http://test-link-rule-to-alarm-link-get-4-link-1-url.com?info[]=test-resource-to-alarm-link-get-2-info-1-val|test-link-mongo-data-2-status&"
              }
            ],
            "test-category-to-alarm-link-get-3": [
              {
                "rule_id": "test-link-rule-to-alarm-link-get-2",
                "label": "test-link-rule-to-alarm-link-get-2-link-2-label",
                "icon_name": "test-link-rule-to-alarm-link-get-2-link-2-icon",
                "url": "http://test-link-rule-to-alarm-link-get-2-link-2-url.com?info[]=test-resource-to-alarm-link-get-2-info-1-val&"
              }
            ],
            "test-category-to-alarm-link-get-5": [
              {
                "rule_id": "test-link-rule-to-alarm-link-get-4",
                "label": "test-link-rule-to-alarm-link-get-4-link-2-label",
                "icon_name": "test-link-rule-to-alarm-link-get-4-link-2-icon",
                "url": "http://test-link-rule-to-alarm-link-get-4-link-2-url.com?info[]=test-resource-to-alarm-link-get-2-info-1-val|test-link-mongo-data-2-status&"
              }
            ]
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
