import { createNamespacedHelpers } from 'vuex';

import { EXPORT_STATUSES } from '@/constants';
import { EXPORT_FETCHING_INTERVAL } from '@/config';

import { saveCsvFile } from '@/helpers/files';

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
    contextExportData() {
      return this.getContextExportByWidgetId(this.widget._id);
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
      fetchExportContext: 'fetchExportContext',
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

    async exportContext({ params, name } = {}) {
      try {
        await this.createContextExport({ params, widgetId: this.widget._id });

        this.startFetchExportContextData({ id: this.contextExportData._id, widgetId: this.widget._id, name });
      } catch (err) {
        this.$popups.error({ text: err.error || this.$t('errors.default') });
      }
    },

    startFetchExportContextData({ id, widgetId, name }) {
      setTimeout(async () => {
        try {
          const exportContextData = await this.fetchExportContext({ id });

          switch (exportContextData.status) {
            case EXPORT_STATUSES.running:
              this.startFetchExportContextData({ id });
              break;
            case EXPORT_STATUSES.completed: {
              const csvFile = await this.fetchContextCsvFile({ id, widgetId });

              saveCsvFile(csvFile, name);
              break;
            }
          }
        } catch (err) {
          this.$popups.error({ text: err.error || this.$t('errors.default') });
        }
      }, EXPORT_FETCHING_INTERVAL);
    },
  },
};
