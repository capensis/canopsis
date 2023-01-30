import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT } from '@/constants';

import { widgetColumnsToForm } from '@/helpers/forms/shared/widget-column';

const { mapActions } = createNamespacedHelpers('widgetTemplate');

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
