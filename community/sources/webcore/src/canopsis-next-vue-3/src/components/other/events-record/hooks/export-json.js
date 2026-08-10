import { ref, unref } from 'vue';

import { getEventsRecordFileUrl } from '@/helpers/entities/events-record/url';
import { openUrlInNewTab } from '@/helpers/url';
import { useFilePolling } from '@/hooks/polling';
import { useEventsRecord } from '@/hooks/store/modules/events-record';

export const useEventsRecordExportJson = ({ eventsRecordId, eventIds = [], eventPattern = [] }) => {
  const downloading = ref(false);

  const { createEventsRecordExport, fetchEventsRecordExport } = useEventsRecord();

  const { poll: generateFile } = useFilePolling({
    createHandler: createEventsRecordExport,
    fetchHandler: fetchEventsRecordExport,
    endHandler: fileData => openUrlInNewTab(getEventsRecordFileUrl(fileData?._id ?? fileData?.id)),
  });

  const exportJson = async () => {
    try {
      downloading.value = true;

      await generateFile({
        id: unref(eventsRecordId),
        eventIds: unref(eventIds),
        eventPattern: unref(eventPattern),
      });
    } finally {
      downloading.value = false;
    }
  };

  return {
    downloading,
    exportJson,
  };
};
