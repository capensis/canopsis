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
        const widgetId = this.widget._id;

        const { _id: id } = await this.createContextExport({ params, widgetId });

        await this.waitGeneratingContextFile({ id, widgetId });

        const csvFile = await this.fetchContextCsvFile({ id, widgetId });

        saveCsvFile(csvFile, name);
      } catch (err) {
        this.$popups.error({ text: err.error || this.$t('errors.default') });
      }
    },

    waitGeneratingContextFile({ id, widgetId }) {
      return new Promise((resolve, reject) => {
        const interval = setInterval(async () => {
          try {
            const exportContextData = await this.fetchExportContext({ id, widgetId });

            if (exportContextData.status === EXPORT_STATUSES.completed) {
              resolve(exportContextData);
            }

            if (exportContextData.status === EXPORT_STATUSES.failed) {
              reject();
            }

            if (exportContextData.status !== EXPORT_STATUSES.running) {
              clearInterval(interval);
            }
          } catch (err) {
            clearInterval(interval);
            reject(err);
          }
        }, EXPORT_FETCHING_INTERVAL);
      });
    },
  },
};
