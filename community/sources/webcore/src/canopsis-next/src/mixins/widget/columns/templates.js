import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT, ENTITIES_TYPES } from '@/constants';

import { widgetColumnsToForm } from '@/helpers/forms/shared/widget-column';

const { mapActions } = createNamespacedHelpers('view/widget/template');

export const widgetColumnsTemplatesMixin = {
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
      return this.widgetTemplates.map(template => ({
        ...template,

        columns: widgetColumnsToForm(template.columns),
      }));
    },

    alarmTypeTemplates() { // TODO: May be move this logic to component?
      return this.preparedWidgetTemplates.filter(({ type }) => type === ENTITIES_TYPES.alarm);
    },

    entityTypeTemplates() {
      return this.preparedWidgetTemplates.filter(({ type }) => type === ENTITIES_TYPES.entity);
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
