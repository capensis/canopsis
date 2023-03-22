import { EXPORT_FETCHING_INTERVAL } from '@/config';
import { EXPORT_STATUSES } from '@/constants';

import { saveCsvFile } from '@/helpers/file/files';

/**
 * @typedef {Object} ExportCsvMixin
 * @property {Object} methods
 * @property {function} methods.exportAsCsv
 * @property {function} methods.generateFile
 * @property {function} methods.waitGeneratingCsvFile
 */

/**
 * Mixin creator for exporting widget data
 *
 * @param {string} createExport
 * @param {string} fetchExportFile
 * @param {string} fetchExport
 * @returns {ExportCsvMixin}
 */
export const exportCsvMixinCreator = ({ createExport, fetchExport, fetchExportFile }) => ({
  methods: {
    async generateFile({ data, ...params }) {
      const { _id: id } = await this[createExport]({ data, ...params });

      await this.waitGeneratingCsvFile({ id, ...params });

      return this[fetchExportFile]({ id, ...params });
    },

    async exportAsCsv({ data, name, ...params }) {
      const file = await this.generateFile({ data, ...params });

      saveCsvFile(file, name);
    },

    waitGeneratingCsvFile({ id, ...params }) {
      return new Promise((resolve, reject) => {
        setTimeout(async () => {
          try {
            const exportData = await this[fetchExport]({ id, ...params });

            if (exportData.status === EXPORT_STATUSES.completed) {
              return resolve(exportData);
            }

            if (exportData.status === EXPORT_STATUSES.failed) {
              return reject();
            }

            return resolve(this.waitGeneratingCsvFile({ id, ...params }));
          } catch (err) {
            return reject(err);
          }
        }, EXPORT_FETCHING_INTERVAL);
      });
    },
  },
});
