import { createNamespacedHelpers } from 'vuex';

import { EXPORT_STATUSES } from '@/constants';

const { mapGetters, mapActions } = createNamespacedHelpers('entity');

/**
 * @mixin Helpers' for context store
 */
export default {
  computed: {
    ...mapGetters({
      getContextEntitiesListByWidgetId: 'getListByWidgetId',
      getContextEntitiesMetaByWidgetId: 'getMetaByWidgetId',
      getContextEntitiesPendingByWidgetId: 'getPendingByWidgetId',
      getContextEntitiesFetchingParamsByWidgetId: 'getFetchingParamsByWidgetId',
      getContextExportByWidgetId: 'getExportByWidgetId',
    }),

    contextEntities() {
      return this.getContextEntitiesListByWidgetId(this.widget._id);
    },
    contextEntitiesMeta() {
      return this.getContextEntitiesMetaByWidgetId(this.widget._id);
    },
    contextEntitiesPending() {
      return this.getContextEntitiesPendingByWidgetId(this.widget._id);
    },
    contextExportPending() {
      const exportData = this.getContextExportByWidgetId(this.widget._id);

      return exportData && exportData.status === EXPORT_STATUSES.running;
    },
  },
  methods: {
    ...mapActions({
      fetchContextEntitiesList: 'fetchList',
      removeContextEntity: 'remove',
      updateContextEntity: 'update',
      createContextEntity: 'create',
      refreshContextEntitiesLists: 'refreshLists',
      createContextExport: 'createContextExport',
      fetchContextExport: 'fetchContextExport',
      fetchContextCsvFile: 'fetchContextCsvFile',
    }),

    async updateContextEntityWithPopup(data) {
      await this.updateContextEntity({ data });
      this.$popups.success({ text: this.$t('modals.createEntity.success.edit') });
    },

    async createContextEntityWithPopup(data) {
      await this.createContextEntity({ data });
      this.$popups.success({ text: this.$t('modals.createEntity.success.create') });
    },

    async duplicateContextEntityWithPopup(data) {
      await this.createContextEntity({ data });
      this.$popups.success({ text: this.$t('modals.createEntity.success.duplicate') });
    },
  },
};
