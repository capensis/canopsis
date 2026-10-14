import {
  getAlarmWidgetMoreInfoTemplateId,
  getAlarmWidgetColumnTemplateId,
  getAlarmWidgetColumnPopupTemplateId,
  getAlarmWidgetGroupColumnTemplateId,
} from '@/helpers/entities/alarm/list';
import { registerTemplate, unregisterTemplate } from '@/helpers/handlebars';

export const widgetAlarmsListPrecompileHandlebarsTemplatesMixin = {
  beforeCreate() {
    this.registeredHandlebarsTemplate = [];
  },
  beforeDestroy() {
    this.unregisterAllRegisteredHandlebarsTemplated();
  },
  watch: {
    'widget.parameters.widgetColumns': {
      handler: 'registerWidgetColumnsHandlebarsTemplates',
      immediate: true,
    },
    'widget.parameters.moreInfoTemplate': {
      handler: 'registerMoreInfosHandlebarsTemplates',
      immediate: true,
    },
    'widget.parameters.widgetGroupColumns': {
      handler: 'registerWidgetGroupColumnsHandlebarsTemplates',
      immediate: true,
    },
  },
  methods: {
    unregisterAllRegisteredHandlebarsTemplated() {
      this.registeredHandlebarsTemplate.forEach(unregisterTemplate);
    },

    registerHandlebarsTemplate(id, template) {
      registerTemplate(id, template);

      this.registeredHandlebarsTemplate.push(id);
    },

    registerWidgetColumnsHandlebarsTemplates(columns) {
      columns.forEach(({ value, template, popupTemplate }) => {
        if (template) {
          this.registerHandlebarsTemplate(getAlarmWidgetColumnTemplateId(this.widget._id, value), template);
        }

        if (popupTemplate) {
          this.registerHandlebarsTemplate(getAlarmWidgetColumnPopupTemplateId(this.widget._id, value), popupTemplate);
        }
      });
    },

    registerWidgetGroupColumnsHandlebarsTemplates(columns) {
      columns.forEach(({ value, template }) => {
        if (template) {
          this.registerHandlebarsTemplate(getAlarmWidgetGroupColumnTemplateId(this.widget._id, value), template);
        }
      });
    },

    registerMoreInfosHandlebarsTemplates(template) {
      if (template) {
        this.registerHandlebarsTemplate(getAlarmWidgetMoreInfoTemplateId(this.widget._id), template);
      }
    },
  },
};
