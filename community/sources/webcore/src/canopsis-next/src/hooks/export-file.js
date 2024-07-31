import { unref } from 'vue';

import { EXPORT_FETCHING_INTERVAL } from '@/config';
import { EXPORT_STATUSES } from '@/constants';

import { openUrlInNewTab } from '@/helpers/url';

/**
 * Function to handle exporting files
 *
 * @param {Object} options - Options for exporting files
 * @param {Function} options.createHandler - Function to create a file
 * @param {Function} options.fetchHandler - Function to fetch file status
 * @param {number} [options.completedStatus = EXPORT_STATUSES.completed] - Status code for completed export
 * @param {number} [options.failedStatus = EXPORT_STATUSES.failed] - Status code for failed export
 * @param {number} [options.interval = EXPORT_FETCHING_INTERVAL] - Interval for fetching export status
 * @returns {Object} Object with methods to generate and download files
 */
export const useExportFile = ({
  createHandler,
  fetchHandler,
  completedStatus = EXPORT_STATUSES.completed,
  failedStatus = EXPORT_STATUSES.failed,
  interval = EXPORT_FETCHING_INTERVAL,
}) => {
  /**
   * Function to wait for file generation
   *
   * @param {Object} options - Options for generating file
   * @returns {Promise} Promise that resolves when file generation is complete
   */
  const waitGeneratingFile = options => new Promise((resolve, reject) => {
    setTimeout(async () => {
      try {
        const exportData = await fetchHandler(options);

        if (exportData.status === unref(completedStatus)) {
          return resolve(exportData);
        }

        if (exportData.status === failedStatus) {
          return reject();
        }

        return resolve(waitGeneratingFile(options));
      } catch (err) {
        return reject(err);
      }
    }, unref(interval));
  });

  /**
   * Generate a file by creating it using the provided data and parameters
   *
   * @param {Object} options - Options for generating the file
   * @param {Object} options.data - Data to create the file
   * @param {...any} options.params - Additional parameters for file creation
   * @returns {Promise} Promise that resolves when the file is generated
   */
  const generateFile = async ({ data, ...params } = {}) => {
    const { _id: id } = await createHandler({ data, ...params });

    return waitGeneratingFile({ id, ...params });
  };

  /**
   * Download a file by opening the URL in a new tab
   *
   * @param {string} url - The URL of the file to download
   * @return {WindowProxy} The window object representing the new tab
   */
  const downloadFile = url => openUrlInNewTab(url);

  return {
    generateFile,
    downloadFile,
  };
};
