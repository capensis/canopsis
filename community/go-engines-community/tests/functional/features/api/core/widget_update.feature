Feature: Update a widget
  I need to be able to update a widget
  Only admin should be able to update a widget

  @concurrent
  Scenario: given update request should update widget
    When I am admin
    Then I do PUT /api/v4/widgets/test-widget-to-update:
    """json
    {
      "title": "test-widget-to-update-title-updated",
      "type": "test-widget-to-update-type",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 0}
      },
      "parameters": {
        "mainFilter": "test-widgetfilter-to-widget-update-1-3",
        "test-widget-to-update-param-str": "teststr",
        "test-widget-to-update-param-int": 2,
        "test-widget-to-update-param-bool": true,
        "test-widget-to-update-param-arr": ["teststr1", "teststr2"],
        "test-widget-to-update-param-map": {"testkey": "teststr"}
      },
      "filters": [
        {
          "_id": "test-widgetfilter-to-widget-update-1-3",
          "title": "test-widgetfilter-to-widget-update-1-3-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-widget-update-1-3-pattern"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-widget-update-1-3-pattern"
                }
              }
            ]
          ],
          "pbehavior_pattern": [
            [
              {
                "field": "pbehavior_info.type",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-widget-update-1-3-pattern"
                }
              }
            ]
          ]
        },
        {
          "_id": "test-widgetfilter-to-widget-update-1-2",
          "title": "test-widgetfilter-to-widget-update-1-2-title",
          "corporate_alarm_pattern": "test-pattern-to-widget-edit-1"
        }
      ]
    }
    """
    Then the response code should be 200
    When I save response mainFilterID={{ (index .lastResponse.filters 0)._id }}
    Then the response key "parameters.mainFilter" should not be "test-widgetfilter-to-widget-update-1-3"
    Then the response body should contain:
    """json
    {
      "_id": "test-widget-to-update",
      "title": "test-widget-to-update-title-updated",
      "type": "test-widget-to-update-type",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 0}
      },
      "parameters": {
        "mainFilter": "{{ .mainFilterID }}",
        "test-widget-to-update-param-str": "teststr",
        "test-widget-to-update-param-int": 2,
        "test-widget-to-update-param-bool": true,
        "test-widget-to-update-param-arr": ["teststr1", "teststr2"],
        "test-widget-to-update-param-map": {"testkey": "teststr"}
      },
      "author": {
        "_id": "root",
        "name": "root"
      },
      "is_private": false,
      "filters": [
        {
          "title": "test-widgetfilter-to-widget-update-1-3-title",
          "is_private": false,
          "is_user_preference": false,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-widget-update-1-3-pattern"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-widget-update-1-3-pattern"
                }
              }
            ]
          ],
          "pbehavior_pattern": [
            [
              {
                "field": "pbehavior_info.type",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-widget-update-1-3-pattern"
                }
              }
            ]
          ]
        },
        {
          "_id": "test-widgetfilter-to-widget-update-1-2",
          "title": "test-widgetfilter-to-widget-update-1-2-title",
          "is_private": false,
          "is_user_preference": false,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "corporate_alarm_pattern": "test-pattern-to-widget-edit-1",
          "corporate_alarm_pattern_title": "test-pattern-to-widget-edit-1-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-widget-edit-1-pattern"
                }
              }
            ]
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/widget-filters?widget=test-widget-to-update
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "test-widgetfilter-to-widget-update-1-3-title",
          "is_private": false,
          "is_user_preference": false,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-widget-update-1-3-pattern"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-widget-update-1-3-pattern"
                }
              }
            ]
          ],
          "pbehavior_pattern": [
            [
              {
                "field": "pbehavior_info.type",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-widget-update-1-3-pattern"
                }
              }
            ]
          ]
        },
        {
          "_id": "test-widgetfilter-to-widget-update-1-2",
          "title": "test-widgetfilter-to-widget-update-1-2-title",
          "is_private": false,
          "is_user_preference": false,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "corporate_alarm_pattern": "test-pattern-to-widget-edit-1",
          "corporate_alarm_pattern_title": "test-pattern-to-widget-edit-1-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-widget-edit-1-pattern"
                }
              }
            ]
          ]
        },
        {
          "_id": "test-widgetfilter-to-widget-update-1-4",
          "title": "test-widgetfilter-to-widget-update-1-4-title",
          "is_user_preference": true,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-widget-update-1-4-pattern"
                }
              }
            ]
          ]
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
  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/widgets/test-widget-to-update
    Then the response code should be 401

  @concurrent
  Scenario: given update request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/widgets/test-widget-to-update
    Then the response code should be 403

  @concurrent
  Scenario: given update request and auth user without view permission should not allow access
    When I am admin
    When I do PUT /api/v4/widgets/test-widget-to-check-access:
    """json
    {
      "title": "test-widget-to-check-access-title",
      "type": "test-widget-to-check-access-type"
    }
    """
    Then the response code should be 403

  @concurrent
  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/widgets/test-widget-to-update:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "type": "Type is missing."
      }
    }
    """

  @concurrent
  Scenario: given update request with not exist id should return not allow access error
    When I am admin
    When I do PUT /api/v4/widgets/test-widget-not-found:
    """json
    {
      "title": "test-widget-not-found-title",
      "type": "test-widget-not-found-type"
    }
    """
    Then the response code should be 404

  @concurrent
  Scenario: given update owned private widget request should update widget
    When I am admin
    Then I do PUT /api/v4/widgets/test-private-widget-to-update-1:
    """json
    {
      "title": "test-private-widget-to-update-title-updated",
      "type": "test-private-widget-to-update-type",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 0}
      },
      "parameters": {
        "mainFilter": "test-private-widgetfilter-to-private-widget-update-1-3",
        "test-private-widget-to-update-param-str": "teststr",
        "test-private-widget-to-update-param-int": 2,
        "test-private-widget-to-update-param-bool": true,
        "test-private-widget-to-update-param-arr": ["teststr1", "teststr2"],
        "test-private-widget-to-update-param-map": {"testkey": "teststr"}
      },
      "filters": [
        {
          "_id": "test-private-widgetfilter-to-private-widget-update-1-3",
          "title": "test-private-widgetfilter-to-private-widget-update-1-3-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-private-widgetfilter-to-private-widget-update-1-3-pattern"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-private-widgetfilter-to-private-widget-update-1-3-pattern"
                }
              }
            ]
          ],
          "pbehavior_pattern": [
            [
              {
                "field": "pbehavior_info.type",
                "cond": {
                  "type": "eq",
                  "value": "test-private-widgetfilter-to-private-widget-update-1-3-pattern"
                }
              }
            ]
          ]
        },
        {
          "_id": "test-private-widgetfilter-to-private-widget-update-1-2",
          "title": "test-private-widgetfilter-to-private-widget-update-1-2-title",
          "corporate_alarm_pattern": "test-pattern-to-widget-edit-1"
        }
      ]
    }
    """
    Then the response code should be 200
    When I save response mainFilterID={{ (index .lastResponse.filters 0)._id }}
    Then the response key "parameters.mainFilter" should not be "test-private-widgetfilter-to-private-widget-update-1-3"
    Then the response body should contain:
    """json
    {
      "_id": "test-private-widget-to-update-1",
      "title": "test-private-widget-to-update-title-updated",
      "type": "test-private-widget-to-update-type",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 0}
      },
      "parameters": {
        "mainFilter": "{{ .mainFilterID }}",
        "test-private-widget-to-update-param-str": "teststr",
        "test-private-widget-to-update-param-int": 2,
        "test-private-widget-to-update-param-bool": true,
        "test-private-widget-to-update-param-arr": ["teststr1", "teststr2"],
        "test-private-widget-to-update-param-map": {"testkey": "teststr"}
      },
      "author": {
        "_id": "root",
        "name": "root"
      },
      "is_private": true,
      "filters": [
        {
          "title": "test-private-widgetfilter-to-private-widget-update-1-3-title",
          "is_private": true,
          "is_user_preference": false,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-private-widgetfilter-to-private-widget-update-1-3-pattern"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-private-widgetfilter-to-private-widget-update-1-3-pattern"
                }
              }
            ]
          ],
          "pbehavior_pattern": [
            [
              {
                "field": "pbehavior_info.type",
                "cond": {
                  "type": "eq",
                  "value": "test-private-widgetfilter-to-private-widget-update-1-3-pattern"
                }
              }
            ]
          ]
        },
        {
          "_id": "test-private-widgetfilter-to-private-widget-update-1-2",
          "title": "test-private-widgetfilter-to-private-widget-update-1-2-title",
          "is_private": true,
          "is_user_preference": false,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "corporate_alarm_pattern": "test-pattern-to-widget-edit-1",
          "corporate_alarm_pattern_title": "test-pattern-to-widget-edit-1-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-widget-edit-1-pattern"
                }
              }
            ]
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given update not owned private widget request should return error
    When I am admin
    Then I do PUT /api/v4/widgets/test-private-widget-to-update-2:
    """json
    {
      "title": "test-private-widget-to-update-title-updated",
      "type": "test-private-widget-to-update-type",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 0}
      },
      "parameters": {
        "test-private-widget-to-update-param-str": "teststr",
        "test-private-widget-to-update-param-int": 2,
        "test-private-widget-to-update-param-bool": true,
        "test-private-widget-to-update-param-arr": ["teststr1", "teststr2"],
        "test-private-widget-to-update-param-map": {"testkey": "teststr"}
      }
    }
    """
    Then the response code should be 403

  @concurrent
  Scenario: given update owned private widget request with api_private_view_groups
    but without api_view permissions should return filters should update private widget
    When I am test-role-to-private-views-without-view-perm
    Then I do PUT /api/v4/widgets/test-private-widget-to-update-3:
    """json
    {
      "title": "test-private-widget-to-update-title-updated",
      "type": "test-private-widget-to-update-type",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 0}
      },
      "parameters": {
        "mainFilter": "test-private-widgetfilter-to-private-widget-update-3-3",
        "test-private-widget-to-update-param-str": "teststr",
        "test-private-widget-to-update-param-int": 2,
        "test-private-widget-to-update-param-bool": true,
        "test-private-widget-to-update-param-arr": ["teststr1", "teststr2"],
        "test-private-widget-to-update-param-map": {"testkey": "teststr"}
      },
      "filters": [
        {
          "_id": "test-private-widgetfilter-to-private-widget-update-3-3",
          "title": "test-private-widgetfilter-to-private-widget-update-3-3-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-private-widgetfilter-to-private-widget-update-3-3-pattern"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-private-widgetfilter-to-private-widget-update-3-3-pattern"
                }
              }
            ]
          ],
          "pbehavior_pattern": [
            [
              {
                "field": "pbehavior_info.type",
                "cond": {
                  "type": "eq",
                  "value": "test-private-widgetfilter-to-private-widget-update-3-3-pattern"
                }
              }
            ]
          ]
        }
      ]
    }
    """
    Then the response code should be 200
    When I save response mainFilterID={{ (index .lastResponse.filters 0)._id }}
    Then the response key "parameters.mainFilter" should not be "test-private-widgetfilter-to-private-widget-update-3-3"
    Then the response body should contain:
    """json
    {
      "_id": "test-private-widget-to-update-3",
      "title": "test-private-widget-to-update-title-updated",
      "type": "test-private-widget-to-update-type",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 0}
      },
      "parameters": {
        "mainFilter": "{{ .mainFilterID }}",
        "test-private-widget-to-update-param-str": "teststr",
        "test-private-widget-to-update-param-int": 2,
        "test-private-widget-to-update-param-bool": true,
        "test-private-widget-to-update-param-arr": ["teststr1", "teststr2"],
        "test-private-widget-to-update-param-map": {"testkey": "teststr"}
      },
      "author": {
        "_id": "test-user-to-private-views-without-view-perm",
        "name": "test-user-to-private-views-without-view-perm"
      },
      "is_private": true,
      "filters": [
        {
          "title": "test-private-widgetfilter-to-private-widget-update-3-3-title",
          "is_private": true,
          "is_user_preference": false,
          "author": {
            "_id": "test-user-to-private-views-without-view-perm",
            "name": "test-user-to-private-views-without-view-perm"
          },
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-private-widgetfilter-to-private-widget-update-3-3-pattern"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-private-widgetfilter-to-private-widget-update-3-3-pattern"
                }
              }
            ]
          ],
          "pbehavior_pattern": [
            [
              {
                "field": "pbehavior_info.type",
                "cond": {
                  "type": "eq",
                  "value": "test-private-widgetfilter-to-private-widget-update-3-3-pattern"
                }
              }
            ]
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given update public widget request with api_private_view_groups
    but without api_view permissions should return filters should not allow access
    When I am test-role-to-private-views-without-view-perm
    Then I do PUT /api/v4/widgets/test-widget-to-update:
    """json
    {
      "title": "test-widget-to-update-title-updated",
      "type": "test-widget-to-update-type",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 0}
      },
      "parameters": {
        "mainFilter": "test-widgetfilter-to-widget-update-3-3",
        "test-widget-to-update-param-str": "teststr",
        "test-widget-to-update-param-int": 2,
        "test-widget-to-update-param-bool": true,
        "test-widget-to-update-param-arr": ["teststr1", "teststr2"],
        "test-widget-to-update-param-map": {"testkey": "teststr"}
      },
      "filters": []
    }
    """
    Then the response code should be 403
