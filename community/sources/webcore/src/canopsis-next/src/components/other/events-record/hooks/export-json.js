import { set, ref } from 'vue';

import { getEventsRecordFileUrl } from '@/helpers/entities/events-record/url';

import { useEventsRecord } from '@/hooks/store/modules/events-record';
import { useExportFile } from '@/hooks/export-file';

export const useExportJson = () => {
  const downloadingsById = ref({});

  const {
    createEventsRecordExport,
    fetchEventsRecordExport,
  } = useEventsRecord();

  const { generateFile } = useExportFile({
    createHandler: createEventsRecordExport,
    fetchHandler: fetchEventsRecordExport,
    urlPreparer: fileData => getEventsRecordFileUrl(fileData?._id),
  });

  const exportJson = async (eventsRecordId, eventIds = []) => {
    set(downloadingsById.value, eventsRecordId, true);

    await generateFile({ id: eventsRecordId, eventIds });

    set(downloadingsById.value, eventsRecordId, false);
  };

  return {
    downloadingsById,

    exportJson,
  };
};
