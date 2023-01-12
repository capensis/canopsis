Feature: Update a widget
  I need to be able to update a widget
  Only admin should be able to update a widget

  Scenario: given updated or deleted widget template request should return updated AlarmList widget
    When I am admin
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-widget-widgettemplate-1-title",
      "tab": "test-tab-to-widget-edit",
      "type": "AlarmsList",
      "parameters": {
        "widgetColumnsTemplate": "test-widgettemplate-to-widget-widgettemplate-1",
        "serviceDependenciesColumnsTemplate": "test-widgettemplate-to-widget-widgettemplate-2"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "parameters": {
        "widgetColumnsTemplate": "test-widgettemplate-to-widget-widgettemplate-1",
        "widgetColumnsTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-1-title",
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
        "serviceDependenciesColumnsTemplate": "test-widgettemplate-to-widget-widgettemplate-2",
        "serviceDependenciesColumnsTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-2-title",
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
    When I save response widgetId={{ .lastResponse._id }}
    When I do PUT /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-1:
    """json
    {
      "title": "test-widgettemplate-to-widget-widgettemplate-1-title-updated",
      "type": "alarm",
      "columns": [
        {
          "value": "v.resource"
        },
        {
          "value": "v.component"
        },
        {
          "value": "v.display_name"
        },
        {
          "value": "extra_details"
        }
      ]
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-2:
    """json
    {
      "title": "test-widgettemplate-to-widget-widgettemplate-2-title-updated",
      "type": "entity",
      "columns": [
        {
          "value": "_id"
        },
        {
          "value": "type"
        },
        {
          "value": "component"
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/widgets/{{ .widgetId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "parameters": {
        "widgetColumnsTemplate": "test-widgettemplate-to-widget-widgettemplate-1",
        "widgetColumnsTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-1-title-updated",
        "widgetColumns": [
          {
            "value": "v.resource"
          },
          {
            "value": "v.component"
          },
          {
            "value": "v.display_name"
          },
          {
            "value": "extra_details"
          }
        ],
        "serviceDependenciesColumnsTemplate": "test-widgettemplate-to-widget-widgettemplate-2",
        "serviceDependenciesColumnsTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-2-title-updated",
        "serviceDependenciesColumns": [
          {
            "value": "_id"
          },
          {
            "value": "type"
          },
          {
            "value": "component"
          }
        ]
      }
    }
    """
    When I do DELETE /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-1
    Then the response code should be 204
    When I do DELETE /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-2
    Then the response code should be 204
    When I do GET /api/v4/widgets/{{ .widgetId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "parameters": {
        "widgetColumns": [
          {
            "value": "v.resource"
          },
          {
            "value": "v.component"
          },
          {
            "value": "v.display_name"
          },
          {
            "value": "extra_details"
          }
        ],
        "serviceDependenciesColumns": [
          {
            "value": "_id"
          },
          {
            "value": "type"
          },
          {
            "value": "component"
          }
        ]
      }
    }
    """
    Then the response key "parameters.widgetColumnsTemplate" should not exist
    Then the response key "parameters.widgetColumnsTemplateTitle" should not exist
    Then the response key "parameters.serviceDependenciesColumnsTemplate" should not exist
    Then the response key "parameters.serviceDependenciesColumnsTemplateTitle" should not exist

  Scenario: given updated or deleted widget template request should return updated ServiceWeather widget
    When I am admin
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-widget-widgettemplate-2-title",
      "tab": "test-tab-to-widget-edit",
      "type": "ServiceWeather",
      "parameters": {
        "alarmsList": {
          "itemsPerPage": 10,
          "widgetColumnsTemplate": "test-widgettemplate-to-widget-widgettemplate-3"
        },
        "serviceDependenciesColumnsTemplate": "test-widgettemplate-to-widget-widgettemplate-4"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "parameters": {
        "alarmsList": {
          "itemsPerPage": 10,
          "widgetColumnsTemplate": "test-widgettemplate-to-widget-widgettemplate-3",
          "widgetColumnsTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-3-title",
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
          ]
        },
        "serviceDependenciesColumnsTemplate": "test-widgettemplate-to-widget-widgettemplate-4",
        "serviceDependenciesColumnsTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-4-title",
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
    When I save response widgetId={{ .lastResponse._id }}
    When I do PUT /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-3:
    """json
    {
      "title": "test-widgettemplate-to-widget-widgettemplate-3-title-updated",
      "type": "alarm",
      "columns": [
        {
          "value": "v.resource"
        },
        {
          "value": "v.component"
        },
        {
          "value": "v.display_name"
        },
        {
          "value": "extra_details"
        }
      ]
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-4:
    """json
    {
      "title": "test-widgettemplate-to-widget-widgettemplate-4-title-updated",
      "type": "entity",
      "columns": [
        {
          "value": "_id"
        },
        {
          "value": "type"
        },
        {
          "value": "component"
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/widgets/{{ .widgetId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "parameters": {
        "alarmsList": {
          "itemsPerPage": 10,
          "widgetColumnsTemplate": "test-widgettemplate-to-widget-widgettemplate-3",
          "widgetColumnsTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-3-title-updated",
          "widgetColumns": [
            {
              "value": "v.resource"
            },
            {
              "value": "v.component"
            },
            {
              "value": "v.display_name"
            },
            {
              "value": "extra_details"
            }
          ]
        },
        "serviceDependenciesColumnsTemplate": "test-widgettemplate-to-widget-widgettemplate-4",
        "serviceDependenciesColumnsTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-4-title-updated",
        "serviceDependenciesColumns": [
          {
            "value": "_id"
          },
          {
            "value": "type"
          },
          {
            "value": "component"
          }
        ]
      }
    }
    """
    When I do DELETE /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-3
    Then the response code should be 204
    When I do DELETE /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-4
    Then the response code should be 204
    When I do GET /api/v4/widgets/{{ .widgetId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "parameters": {
        "alarmsList": {
          "itemsPerPage": 10,
          "widgetColumns": [
            {
              "value": "v.resource"
            },
            {
              "value": "v.component"
            },
            {
              "value": "v.display_name"
            },
            {
              "value": "extra_details"
            }
          ]
        },
        "serviceDependenciesColumns": [
          {
            "value": "_id"
          },
          {
            "value": "type"
          },
          {
            "value": "component"
          }
        ]
      }
    }
    """
    Then the response key "parameters.alarmsList.widgetColumnsTemplate" should not exist
    Then the response key "parameters.alarmsList.widgetColumnsTemplateTitle" should not exist
    Then the response key "parameters.serviceDependenciesColumnsTemplate" should not exist
    Then the response key "parameters.serviceDependenciesColumnsTemplateTitle" should not exist
