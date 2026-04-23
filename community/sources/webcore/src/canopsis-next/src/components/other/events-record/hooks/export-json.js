import { ref, unref } from 'vue';

import { openUrlInNewTab } from '@/helpers/url';
import { getEventsRecordFileUrl } from '@/helpers/entities/events-record/url';

import { useEventsRecord } from '@/hooks/store/modules/events-record';
import { useFilePolling } from '@/hooks/polling';

/**
 * Function to use events record export JSON
 *
 * @param {Object} options - Options for events record export JSON
 * @param {string} options.eventsRecordId - ID of the events record
 * @param {Array} [options.eventIds = []] - Array of event IDs (default: empty array)
 * @param {Object} [options.params = {}] - Object of parameters (default: empty object)
 * @returns {Object} Object containing downloading status and function to export JSON
 */
export const useEventsRecordExportJson = ({ eventsRecordId, eventIds = [], params = {} }) => {
  const downloading = ref(false);

  const {
    createEventsRecordExport,
    fetchEventsRecordExport,
  } = useEventsRecord();

  const { poll: generateFile } = useFilePolling({
    createHandler: createEventsRecordExport,
    fetchHandler: fetchEventsRecordExport,
    endHandler: fileData => openUrlInNewTab(getEventsRecordFileUrl(fileData?._id)),
  });

  /**
   * Function to export JSON file
   *
   * @returns {Promise} Promise that resolves when the export process is completed
   */
  const exportJson = async () => {
    try {
      downloading.value = true;

      await generateFile({ id: unref(eventsRecordId), eventIds: unref(eventIds), params: unref(params) });
    } finally {
      downloading.value = false;
    }
  };

  return {
    downloading,

    exportJson,
  };
};
