import {
  getAlarmWidgetMoreInfoTemplateId,
  getAlarmWidgetColumnTemplateId,
  getAlarmWidgetColumnPopupTemplateId,
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

    registerMoreInfosHandlebarsTemplates(template) {
      if (template) {
        this.registerHandlebarsTemplate(getAlarmWidgetMoreInfoTemplateId(this.widget._id), template);
      }
    },
  },
};
