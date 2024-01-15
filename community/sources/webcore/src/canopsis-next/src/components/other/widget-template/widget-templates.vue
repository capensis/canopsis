<template>
  <widget-templates-list
    :options.sync="options"
    :widget-templates="widgetTemplates"
    :pending="widgetTemplatesPending"
    :total-items="widgetTemplatesMeta.total_count"
    :updatable="hasUpdateAnyWidgetTemplateAccess"
    :removable="hasDeleteAnyWidgetTemplateAccess"
    @edit="showEditWidgetTemplateModal"
    @remove="showRemoveWidgetTemplateModal"
  />
</template>

<script>
import { MODALS } from '@/constants';

import { localQueryMixin } from '@/mixins/query-local/query';
import { entitiesWidgetTemplatesMixin } from '@/mixins/entities/widget-template';
import { permissionsTechnicalWidgetTemplateMixin } from '@/mixins/permissions/technical/widget-templates';

import WidgetTemplatesList from './widget-templates-list.vue';

export default {
  components: { WidgetTemplatesList },
  mixins: [
    localQueryMixin,
    entitiesWidgetTemplatesMixin,
    permissionsTechnicalWidgetTemplateMixin,
  ],
  mounted() {
    this.fetchList();
  },
  methods: {
    showEditWidgetTemplateModal(widgetTemplate) {
      this.$modals.show({
        name: MODALS.createWidgetTemplate,
        config: {
          widgetTemplate,

          title: this.$t('modals.createWidgetTemplate.edit.title'),
          action: async (newWidgetTemplate) => {
            await this.updateWidgetTemplate({ id: widgetTemplate._id, data: newWidgetTemplate });

            return this.fetchList();
          },
        },
      });
    },

    showRemoveWidgetTemplateModal(widgetTemplateId) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removeWidgetTemplate({ id: widgetTemplateId });

            return this.fetchList();
          },
        },
      });
    },

    fetchList() {
      return this.fetchWidgetTemplatesList({ params: this.getQuery() });
    },
  },
};
</script>
