import { EXPORT_FETCHING_INTERVAL } from '@/config';
import { EXPORT_STATUSES } from '@/constants';

import { saveCsvFile } from '@/helpers/file/files';

/**
 * @typedef {Object} ExportMixin
 * @property {Object} methods
 * @property {function} methods.exportWidgetAsCsv
 * @property {function} methods.generateWidgetFile
 * @property {function} methods.waitGeneratingCsvFile
 */

/**
 * Mixin creator for exporting widget data
 *
 * @param {string} createExport
 * @param {string} fetchExportFile
 * @param {string} fetchExport
 * @returns {ExportMixin}
 */
export default ({ createExport, fetchExport, fetchExportFile }) => ({
  methods: {
    async generateWidgetFile({ params } = {}) {
      const widgetId = this.widget._id;

      const { _id: id } = await this[createExport]({ params, widgetId });

      await this.waitGeneratingCsvFile({ id, widgetId });

      return this[fetchExportFile]({ id, widgetId });
    },

    async exportWidgetAsCsv({ params, name } = {}) {
      try {
        const file = await this.generateWidgetFile({ params });

        saveCsvFile(file, name);
      } catch (err) {
        this.$popups.error({ text: err.error || this.$t('errors.default') });
      }
    },

    waitGeneratingCsvFile({ id, widgetId }) {
      return new Promise((resolve, reject) => {
        setTimeout(async () => {
          try {
            const exportData = await this[fetchExport]({ id, widgetId });

            if (exportData.status === EXPORT_STATUSES.completed) {
              return resolve(exportData);
            }

            if (exportData.status === EXPORT_STATUSES.failed) {
              return reject();
            }

            return resolve(this.waitGeneratingCsvFile({ id, widgetId }));
          } catch (err) {
            return reject(err);
          }
        }, EXPORT_FETCHING_INTERVAL);
      });
    },
  },
});
