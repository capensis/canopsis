import { EXPORT_FETCHING_INTERVAL } from '@/config';
import { EXPORT_STATUSES } from '@/constants';

import { openUrlInNewTab } from '@/helpers/url';

export const useExportFile = ({
  createExport,
  fetchExport,
  completedStatus = EXPORT_STATUSES.completed,
  failedStatus = EXPORT_STATUSES.failed,
  interval = EXPORT_FETCHING_INTERVAL,
} = {}) => {
  const waitGeneratingFile = ({ id, ...params } = {}) => new Promise((resolve, reject) => {
    setTimeout(async () => {
      try {
        const exportData = await fetchExport({ id, ...params });

        if (exportData.status === completedStatus) {
          return resolve(exportData);
        }

        if (exportData.status === failedStatus) {
          return reject();
        }

        return resolve(waitGeneratingFile({ id, ...params }));
      } catch (err) {
        return reject(err);
      }
    }, interval);
  });

  const generateFile = async ({ data, ...params } = {}) => {
    const { _id: id } = await createExport(data, ...params);

    return waitGeneratingFile({ id, ...params });
  };

  const downloadFile = url => openUrlInNewTab(url);

  return {
    generateFile,
    downloadFile,
  };
};
