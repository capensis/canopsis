import { unref } from 'vue';

import { EXPORT_FETCHING_INTERVAL } from '@/config';
import { EXPORT_STATUSES } from '@/constants';

import { usePolling } from './polling';

/**
 * Function to handle exporting a file
 *
 * @param {Object} options - Options for exporting a file
 * @param {Function} options.createHandler - Function to create the file
 * @param {Function} options.fetchHandler - Function to fetch the file
 * @param {Function} [options.endHandler] - Function to handle end of polling
 * @param {number} [options.completedStatus = EXPORT_STATUSES.completed] - Status code for completed export (default: 1)
 * @param {number} [options.failedStatus = EXPORT_STATUSES.failed] - Status code for failed export (default: 2)
 * @param {number} [options.interval = EXPORT_FETCHING_INTERVAL] - Interval in milliseconds for polling (default: 2000)
 * @returns {Object} Object containing the function to generate the file
 */
export const useExportFile = ({
  createHandler,
  fetchHandler,
  endHandler = () => {},
  completedStatus = EXPORT_STATUSES.completed,
  failedStatus = EXPORT_STATUSES.failed,
  interval = EXPORT_FETCHING_INTERVAL,
}) => {
  /**
   * Function to process the export file
   *
   * @param {Object} options - Options for the export file
   * @param {string} options._id - ID of the file
   * @param {Object} rest - Additional data for the export
   * @param {Function} resolve - Function to resolve the export process
   * @param {Function} reject - Function to reject the export process
   * @returns {Promise} Promise that resolves when the export process is completed
   */
  const processHandler = async ({ _id: id, ...rest }, resolve, reject) => {
    const exportData = await fetchHandler({ id, ...rest });

    if (exportData.status === unref(completedStatus)) {
      return resolve(exportData);
    }

    if (exportData.status === failedStatus) {
      return reject();
    }

    return exportData;
  };

  const {
    poll,
  } = usePolling({
    interval,
    endHandler,
    processHandler,
    startHandler: createHandler,
  });

  return {
    generateFile: poll,
  };
};
