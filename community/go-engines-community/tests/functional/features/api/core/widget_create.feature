Feature: Create a widget
  I need to be able to create a widget
  Only admin should be able to create a widget

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-create-1-title",
      "tab": "test-tab-to-widget-edit",
      "type": "test-widget-to-create-1-type",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 0}
      },
      "parameters": {
        "mainFilter": "test-widgetfilter-to-widget-create-1-2",
        "test-widget-to-create-1-param-str": "teststr",
        "test-widget-to-create-1-param-int": 2,
        "test-widget-to-create-1-param-bool": true,
        "test-widget-to-create-1-param-arr": ["teststr1", "teststr2"],
        "test-widget-to-create-1-param-map": {"testkey": "teststr"}
      },
      "filters": [
        {
          "_id": "test-widgetfilter-to-widget-create-1-1",
          "title": "test-widgetfilter-to-widget-create-1-1-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-widget-create-1-1-pattern"
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
                  "value": "test-widgetfilter-to-widget-create-1-1-pattern"
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
                  "value": "test-widgetfilter-to-widget-create-1-1-pattern"
                }
              }
            ]
          ]
        },
        {
          "_id": "test-widgetfilter-to-widget-create-1-2",
          "title": "test-widgetfilter-to-widget-create-1-2-title",
          "corporate_alarm_pattern": "test-pattern-to-widget-edit-1"
        }
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "title": "test-widget-to-create-1-title",
      "type": "test-widget-to-create-1-type",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 0}
      },
      "parameters": {
        "test-widget-to-create-1-param-str": "teststr",
        "test-widget-to-create-1-param-int": 2,
        "test-widget-to-create-1-param-bool": true,
        "test-widget-to-create-1-param-arr": ["teststr1", "teststr2"],
        "test-widget-to-create-1-param-map": {"testkey": "teststr"}
      },
      "author": {
        "_id": "root",
        "name": "root"
      },
      "filters": [
        {
          "title": "test-widgetfilter-to-widget-create-1-1-title",
          "is_private": false,
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
                  "value": "test-widgetfilter-to-widget-create-1-1-pattern"
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
                  "value": "test-widgetfilter-to-widget-create-1-1-pattern"
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
                  "value": "test-widgetfilter-to-widget-create-1-1-pattern"
                }
              }
            ]
          ]
        },
        {
          "title": "test-widgetfilter-to-widget-create-1-2-title",
          "is_private": false,
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
    When I save response widgetID={{ .lastResponse._id }}
    When I save response mainFilterID={{ (index .lastResponse.filters 1)._id }}
    Then the response key "parameters.mainFilter" should not be "test-widgetfilter-to-widget-create-1-2"
    When I do GET /api/v4/widgets/{{ .widgetID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "test-widget-to-create-1-title",
      "type": "test-widget-to-create-1-type",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 0}
      },
      "parameters": {
        "mainFilter": "{{ .mainFilterID }}",
        "test-widget-to-create-1-param-str": "teststr",
        "test-widget-to-create-1-param-int": 2,
        "test-widget-to-create-1-param-bool": true,
        "test-widget-to-create-1-param-arr": ["teststr1", "teststr2"],
        "test-widget-to-create-1-param-map": {"testkey": "teststr"}
      },
      "author": {
        "_id": "root",
        "name": "root"
      },
      "filters": [
        {
          "title": "test-widgetfilter-to-widget-create-1-1-title",
          "is_private": false,
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
                  "value": "test-widgetfilter-to-widget-create-1-1-pattern"
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
                  "value": "test-widgetfilter-to-widget-create-1-1-pattern"
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
                  "value": "test-widgetfilter-to-widget-create-1-1-pattern"
                }
              }
            ]
          ]
        },
        {
          "title": "test-widgetfilter-to-widget-create-1-2-title",
          "is_private": false,
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

  Scenario: given create AlarmsList request with widget templates should return ok
    When I am admin
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-create-2-title",
      "tab": "test-tab-to-widget-edit",
      "type": "AlarmsList",
      "parameters": {
        "widgetColumnsTemplate": "test-widgettemplate-to-widget-edit-1",
        "widgetGroupColumnsTemplate": "test-widgettemplate-to-widget-edit-1",
        "serviceDependenciesColumnsTemplate": "test-widgettemplate-to-widget-edit-2"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "title": "test-widget-to-create-2-title",
      "type": "AlarmsList",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "parameters": {
        "widgetColumnsTemplate": "test-widgettemplate-to-widget-edit-1",
        "widgetColumnsTemplateTitle": "test-widgettemplate-to-widget-edit-1-title",
        "widgetColumns": [
          {
            "value": "v.resource"
          },
          {
            "value": "v.component"
          },
          {
            "value": "extra_details"
          }
        ],
        "widgetGroupColumnsTemplate": "test-widgettemplate-to-widget-edit-1",
        "widgetGroupColumnsTemplateTitle": "test-widgettemplate-to-widget-edit-1-title",
        "widgetGroupColumns": [
          {
            "value": "v.resource"
          },
          {
            "value": "v.component"
          },
          {
            "value": "extra_details"
          }
        ],
        "serviceDependenciesColumnsTemplate": "test-widgettemplate-to-widget-edit-2",
        "serviceDependenciesColumnsTemplateTitle": "test-widgettemplate-to-widget-edit-2-title",
        "serviceDependenciesColumns": [
          {
            "value": "_id"
          },
          {
            "value": "type"
          }
        ]
      }
    }
    """
    When I save response widgetID={{ .lastResponse._id }}
    When I do GET /api/v4/widgets/{{ .widgetID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "test-widget-to-create-2-title",
      "type": "AlarmsList",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "parameters": {
        "widgetColumnsTemplate": "test-widgettemplate-to-widget-edit-1",
        "widgetColumnsTemplateTitle": "test-widgettemplate-to-widget-edit-1-title",
        "widgetColumns": [
          {
            "value": "v.resource"
          },
          {
            "value": "v.component"
          },
          {
            "value": "extra_details"
          }
        ],
        "widgetGroupColumnsTemplate": "test-widgettemplate-to-widget-edit-1",
        "widgetGroupColumnsTemplateTitle": "test-widgettemplate-to-widget-edit-1-title",
        "widgetGroupColumns": [
          {
            "value": "v.resource"
          },
          {
            "value": "v.component"
          },
          {
            "value": "extra_details"
          }
        ],
        "serviceDependenciesColumnsTemplate": "test-widgettemplate-to-widget-edit-2",
        "serviceDependenciesColumnsTemplateTitle": "test-widgettemplate-to-widget-edit-2-title",
        "serviceDependenciesColumns": [
          {
            "value": "_id"
          },
          {
            "value": "type"
          }
        ]
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/widgets
    Then the response code should be 401

  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/widgets
    Then the response code should be 403

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/widgets:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "type": "Type is missing.",
        "tab": "Tab is missing."
      }
    }
    """
    When I do POST /api/v4/widgets:
    """json
    {
      "filters": [
        {}
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "type": "Type is missing.",
        "tab": "Tab is missing.",
        "filters.0.alarm_pattern": "AlarmPattern is missing.",
        "filters.0.corporate_alarm_pattern": "CorporateAlarmPattern is missing.",
        "filters.0.corporate_entity_pattern": "CorporateEntityPattern is missing.",
        "filters.0.corporate_pbehavior_pattern": "CorporatePbehaviorPattern is missing.",
        "filters.0.entity_pattern": "EntityPattern is missing.",
        "filters.0.pbehavior_pattern": "PbehaviorPattern is missing.",
        "filters.0.title": "Title is missing.",
        "filters.0.weather_service_pattern": "WeatherServicePattern is missing."
      }
    }
    """
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-create-2-title",
      "tab": "test-tab-to-widget-edit",
      "type": "test-widget-to-create-2-type",
      "filters": [
        {
          "title": "test-widgetfilter-to-widget-create-2-1-title",
          "corporate_entity_pattern": "test-pattern-to-widget-edit-1"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "filters.0.corporate_entity_pattern": "CorporateEntityPattern doesn't exist."
      }
    }
    """

  Scenario: given Junit invalid create request should return errors
    When I am admin
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-create-2-title",
      "tab": "test-tab-to-widget-edit",
      "type": "Junit"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameters.directory": "Directory is missing."
      }
    }
    """
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-create-2-title",
      "tab": "test-tab-to-widget-edit",
      "type": "Junit",
      "parameters": {
        "is_api": true,
        "directory": "testdirectory",
        "screenshot_directories": ["testdirectory"],
        "video_directories": ["testdirectory"]
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameters.directory": "Directory is not empty.",
        "parameters.screenshot_directories": "ScreenshotDirectories is not empty.",
        "parameters.video_directories": "VideoDirectories is not empty."
      }
    }
    """
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-create-2-title",
      "tab": "test-tab-to-widget-edit",
      "type": "Junit",
      "parameters": {
        "directory": "testdirectory",
        "video_filemask": "test",
        "screenshot_filemask": "test",
        "report_fileregexp": "(.*)"
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameters.video_filemask": "VideoFilemask is not a valid file mask.",
        "parameters.screenshot_filemask": "ScreenshotFilemask is not a valid file mask.",
        "parameters.report_fileregexp": "ReportFileRegexp is invalid regexp."
      }
    }
    """

  Scenario: given Map invalid create request should return errors
    When I am admin
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-create-2-title",
      "tab": "test-tab-to-widget-edit",
      "type": "Map"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameters.map": "Map is missing."
      }
    }
    """

  Scenario: given create request and auth user without view permission should not allow access
    When I am admin
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-create-2-title",
      "type": "test-widget-to-create-2-type",
      "tab": "test-tab-to-widget-check-access"
    }
    """
    Then the response code should be 403

  Scenario: given create request with not exist tab should return not allow access error
    When I am admin
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-create-2-title",
      "type": "test-widget-to-create-2-type",
      "tab": "test-tab-not-found"
    }
    """
    Then the response code should be 403

  Scenario: given create request with invalid columns should return bad request error
    When I am admin
    When I do POST /api/v4/widgets:
    """json
    {
      "type": "AlarmsList",
      "parameters": {
        "widgetColumns": [
          {}
        ]
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "parameters.widgetColumns.0.value": "Value is missing."
      }
    }
    """
    When I do POST /api/v4/widgets:
    """json
    {
      "type": "AlarmsList",
      "parameters": {
        "widgetColumns": [
          {
            "value": "unknown"
          }
        ]
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "parameters.widgetColumns.0.value": "Value is invalid."
      }
    }
    """
    When I do POST /api/v4/widgets:
    """json
    {
      "type": "Context",
      "parameters": {
        "serviceDependenciesColumns": [
          {
            "value": "unknown"
          }
        ]
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "parameters.serviceDependenciesColumns.0.value": "Value is invalid."
      }
    }
    """
    When I do POST /api/v4/widgets:
    """json
    {
      "type": "ServiceWeather",
      "parameters": {
        "alarmsList": {
          "widgetColumns": [
            {
              "value": "unknown"
            }
          ]
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "parameters.alarmsList.widgetColumns.0.value": "Value is invalid."
      }
    }
    """
    When I do POST /api/v4/widgets:
    """json
    {
      "type": "AlarmsList",
      "tab": "test-tab-to-widget-edit",
      "parameters": {
        "widgetColumnsTemplate": "not-exist"
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "parameters.widgetColumnsTemplate": "Template doesn't exist."
      }
    }
    """
    When I do POST /api/v4/widgets:
    """json
    {
      "type": "ServiceWeather",
      "tab": "test-tab-to-widget-edit",
      "parameters": {
        "alarmsList": {
          "widgetColumnsTemplate": "not-exist"
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "parameters.alarmsList.widgetColumnsTemplate": "Template doesn't exist."
      }
    }
    """
    When I do POST /api/v4/widgets:
    """json
    {
      "type": "AlarmsList",
      "tab": "test-tab-to-widget-edit",
      "parameters": {
        "serviceDependenciesColumnsTemplate": "not-exist"
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "parameters.serviceDependenciesColumnsTemplate": "Template doesn't exist."
      }
    }
    """
