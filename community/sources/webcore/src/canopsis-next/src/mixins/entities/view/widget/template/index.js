import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('view/widget/template');

export const entitiesWidgetTemplatesMixin = {
  computed: {
    ...mapGetters({
      widgetTemplates: 'items',
      widgetTemplatesMeta: 'meta',
      widgetTemplatesPending: 'pending',
    }),
  },
  methods: {
    ...mapActions({
      fetchWidgetTemplatesList: 'fetchList',
      fetchWidgetTemplatesListWithPreviousParams: 'fetchListWithPreviousParams',
      createWidgetTemplate: 'create',
      updateWidgetTemplate: 'update',
      removeWidgetTemplate: 'remove',
    }),
  },
};
