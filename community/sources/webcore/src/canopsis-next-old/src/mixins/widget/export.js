import { EXPORT_FETCHING_INTERVAL } from '@/config';
import { EXPORT_STATUSES } from '@/constants';

import { openUrlInNewTab } from '@/helpers/url';

/**
 * @typedef {Object} ExportCsvMixin
 * @property {Object} methods
 * @property {function} methods.downloadFile
 * @property {function} methods.generateFile
 * @property {function} methods.waitGeneratingFile
 */

/**
 * Mixin creator for exporting files
 *
 * @param {string} createExport
 * @param {string} fetchExport
 * @param {number | string} completedStatus
 * @param {number | string} failedStatus
 * @returns {ExportCsvMixin}
 */
export const exportMixinCreator = ({
  createExport,
  fetchExport,
  completedStatus = EXPORT_STATUSES.completed,
  failedStatus = EXPORT_STATUSES.failed,
}) => ({
  methods: {
    async generateFile({ data, ...params } = {}) {
      const { _id: id } = await this[createExport]({ data, ...params });

      return this.waitGeneratingFile({ id, ...params });
    },

    downloadFile(url) {
      openUrlInNewTab(url);
    },

    waitGeneratingFile({ id, ...params }) {
      return new Promise((resolve, reject) => {
        setTimeout(async () => {
          try {
            const exportData = await this[fetchExport]({ id, ...params });

            if (exportData.status === completedStatus) {
              return resolve(exportData);
            }

            if (exportData.status === failedStatus) {
              return reject();
            }

            return resolve(this.waitGeneratingFile({ id, ...params }));
          } catch (err) {
            return reject(err);
          }
        }, EXPORT_FETCHING_INTERVAL);
      });
    },
  },
});
