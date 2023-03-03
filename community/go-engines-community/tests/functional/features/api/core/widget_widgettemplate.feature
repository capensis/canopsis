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
        "serviceDependenciesColumnsTemplate": "test-widgettemplate-to-widget-widgettemplate-2",
        "moreInfoTemplateTemplate": "test-widgettemplate-to-widget-widgettemplate-3"
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
        ],
        "moreInfoTemplateTemplate": "test-widgettemplate-to-widget-widgettemplate-3",
        "moreInfoTemplateTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-3-title",
        "moreInfoTemplate": "{{ `{{ alarm.v.display_name }}` }}"
      }
    }
    """
    When I save response widgetId={{ .lastResponse._id }}
    When I do PUT /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-1:
    """json
    {
      "title": "test-widgettemplate-to-widget-widgettemplate-1-title-updated",
      "type": "alarm_columns",
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
      "type": "entity_columns",
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
    When I do PUT /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-3:
    """json
    {
      "title": "test-widgettemplate-to-widget-widgettemplate-3-title-updated",
      "type": "alarm_more_infos",
      "content": "updated {{ `{{ alarm.v.display_name }}` }}"
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
        ],
        "moreInfoTemplateTemplate": "test-widgettemplate-to-widget-widgettemplate-3",
        "moreInfoTemplateTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-3-title-updated",
        "moreInfoTemplate": "updated {{ `{{ alarm.v.display_name }}` }}"
      }
    }
    """
    When I do DELETE /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-1
    Then the response code should be 204
    When I do DELETE /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-2
    Then the response code should be 204
    When I do DELETE /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-3
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
        ],
        "moreInfoTemplate": "updated {{ `{{ alarm.v.display_name }}` }}"
      }
    }
    """
    Then the response key "parameters.widgetColumnsTemplate" should not exist
    Then the response key "parameters.widgetColumnsTemplateTitle" should not exist
    Then the response key "parameters.serviceDependenciesColumnsTemplate" should not exist
    Then the response key "parameters.serviceDependenciesColumnsTemplateTitle" should not exist
    Then the response key "parameters.moreInfoTemplateTemplate" should not exist
    Then the response key "parameters.moreInfoTemplateTemplateTitle" should not exist

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
          "widgetColumnsTemplate": "test-widgettemplate-to-widget-widgettemplate-4"
        },
        "serviceDependenciesColumnsTemplate": "test-widgettemplate-to-widget-widgettemplate-5",
        "blockTemplateTemplate": "test-widgettemplate-to-widget-widgettemplate-6",
        "modalTemplateTemplate": "test-widgettemplate-to-widget-widgettemplate-7",
        "entityTemplateTemplate": "test-widgettemplate-to-widget-widgettemplate-8"
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
          "widgetColumnsTemplate": "test-widgettemplate-to-widget-widgettemplate-4",
          "widgetColumnsTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-4-title",
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
        "serviceDependenciesColumnsTemplate": "test-widgettemplate-to-widget-widgettemplate-5",
        "serviceDependenciesColumnsTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-5-title",
        "serviceDependenciesColumns": [
          {
            "value": "_id"
          },
          {
            "value": "type"
          }
        ],
        "blockTemplateTemplate": "test-widgettemplate-to-widget-widgettemplate-6",
        "blockTemplateTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-6-title",
        "blockTemplate": "{{ `{{ entity.name }}` }}",
        "modalTemplateTemplate": "test-widgettemplate-to-widget-widgettemplate-7",
        "modalTemplateTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-7-title",
        "modalTemplate": "{{ `{{ entities name=\"entity._id\" }}` }}",
        "entityTemplateTemplate": "test-widgettemplate-to-widget-widgettemplate-8",
        "entityTemplateTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-8-title",
        "entityTemplate": "{{ `{{ entity.infos.name1.value }}` }}"
      }
    }
    """
    When I save response widgetId={{ .lastResponse._id }}
    When I do PUT /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-4:
    """json
    {
      "title": "test-widgettemplate-to-widget-widgettemplate-4-title-updated",
      "type": "alarm_columns",
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
    When I do PUT /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-5:
    """json
    {
      "title": "test-widgettemplate-to-widget-widgettemplate-5-title-updated",
      "type": "entity_columns",
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
    When I do PUT /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-6:
    """json
    {
      "title": "test-widgettemplate-to-widget-widgettemplate-6-title-updated",
      "type": "weather_item",
      "content": "updated {{ `{{ entity.name }}` }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-7:
    """json
    {
      "title": "test-widgettemplate-to-widget-widgettemplate-7-title-updated",
      "type": "weather_modal",
      "content": "updated {{ `{{ entities name=\"entity._id\" }}` }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-8:
    """json
    {
      "title": "test-widgettemplate-to-widget-widgettemplate-8-title-updated",
      "type": "weather_entity",
      "content": "updated {{ `{{ entity.infos.name1.value }}` }}"
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
          "widgetColumnsTemplate": "test-widgettemplate-to-widget-widgettemplate-4",
          "widgetColumnsTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-4-title-updated",
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
        "serviceDependenciesColumnsTemplate": "test-widgettemplate-to-widget-widgettemplate-5",
        "serviceDependenciesColumnsTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-5-title-updated",
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
        ],
        "blockTemplateTemplate": "test-widgettemplate-to-widget-widgettemplate-6",
        "blockTemplateTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-6-title-updated",
        "blockTemplate": "updated {{ `{{ entity.name }}` }}",
        "modalTemplateTemplate": "test-widgettemplate-to-widget-widgettemplate-7",
        "modalTemplateTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-7-title-updated",
        "modalTemplate": "updated {{ `{{ entities name=\"entity._id\" }}` }}",
        "entityTemplateTemplate": "test-widgettemplate-to-widget-widgettemplate-8",
        "entityTemplateTemplateTitle": "test-widgettemplate-to-widget-widgettemplate-8-title-updated",
        "entityTemplate": "updated {{ `{{ entity.infos.name1.value }}` }}"
      }
    }
    """
    When I do DELETE /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-4
    Then the response code should be 204
    When I do DELETE /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-5
    Then the response code should be 204
    When I do DELETE /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-6
    Then the response code should be 204
    When I do DELETE /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-7
    Then the response code should be 204
    When I do DELETE /api/v4/widget-templates/test-widgettemplate-to-widget-widgettemplate-8
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
        ],
        "blockTemplate": "updated {{ `{{ entity.name }}` }}",
        "modalTemplate": "updated {{ `{{ entities name=\"entity._id\" }}` }}",
        "entityTemplate": "updated {{ `{{ entity.infos.name1.value }}` }}"
      }
    }
    """
    Then the response key "parameters.alarmsList.widgetColumnsTemplate" should not exist
    Then the response key "parameters.alarmsList.widgetColumnsTemplateTitle" should not exist
    Then the response key "parameters.serviceDependenciesColumnsTemplate" should not exist
    Then the response key "parameters.serviceDependenciesColumnsTemplateTitle" should not exist
    Then the response key "parameters.blockTemplateTemplate" should not exist
    Then the response key "parameters.blockTemplateTemplateTitle" should not exist
    Then the response key "parameters.modalTemplateTemplate" should not exist
    Then the response key "parameters.modalTemplateTemplateTitle" should not exist
    Then the response key "parameters.entityTemplateTemplate" should not exist
    Then the response key "parameters.entityTemplateTemplateTitle" should not exist
