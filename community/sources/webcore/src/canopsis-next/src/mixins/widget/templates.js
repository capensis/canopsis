import { filter } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { COLUMNS_WIDGET_TEMPLATES_TYPES, MAX_LIMIT, WIDGET_TEMPLATES_TYPES } from '@/constants';

import { widgetColumnsToForm } from '@/helpers/forms/shared/widget-column';

const { mapActions } = createNamespacedHelpers('widgetTemplate');

export const widgetTemplatesMixin = {
  data() {
    return {
      widgetTemplates: [],
      widgetTemplatesPending: false,
    };
  },
  mounted() {
    this.fetchList();
  },
  computed: {
    preparedWidgetTemplates() {
      return this.widgetTemplates.map(template => (
        COLUMNS_WIDGET_TEMPLATES_TYPES.includes(template.type)
          ? { ...template, columns: widgetColumnsToForm(template.columns) }
          : template
      ));
    },

    alarmColumnsWidgetTemplates() {
      return filter(this.preparedWidgetTemplates, { type: WIDGET_TEMPLATES_TYPES.alarmColumns });
    },

    entityColumnsWidgetTemplates() {
      return filter(this.preparedWidgetTemplates, { type: WIDGET_TEMPLATES_TYPES.entityColumns });
    },

    alarmMoreInfosWidgetTemplates() {
      return filter(this.preparedWidgetTemplates, { type: WIDGET_TEMPLATES_TYPES.alarmMoreInfos });
    },

    weatherItemWidgetTemplates() {
      return filter(this.preparedWidgetTemplates, { type: WIDGET_TEMPLATES_TYPES.weatherItem });
    },

    weatherModalWidgetTemplates() {
      return filter(this.preparedWidgetTemplates, { type: WIDGET_TEMPLATES_TYPES.weatherModal });
    },

    weatherEntityWidgetTemplates() {
      return filter(this.preparedWidgetTemplates, { type: WIDGET_TEMPLATES_TYPES.weatherEntity });
    },
  },
  methods: {
    ...mapActions({ fetchWidgetTemplatesListWithoutStore: 'fetchListWithoutStore' }),

    async fetchList() {
      try {
        this.widgetTemplatesPending = true;

        const { data } = await this.fetchWidgetTemplatesListWithoutStore({ params: { limit: MAX_LIMIT } });

        this.widgetTemplates = data;
      } catch (err) {
        console.warn(err);
      } finally {
        this.widgetTemplatesPending = false;
      }
    },
  },
};
