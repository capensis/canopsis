<template>
  <c-page @refresh="fetchList">
    <events-records-header
      :progress="status.isRecording"
      @start="start"
      @stop="stop"
    />
    <events-records-list
      :events-records="eventsRecords"
      :pending="pending"
      :options.sync="options"
      :total-items="meta.total_count"
      @show="showEventsRecordModal"
      @export="downloadEventsRecord"
      @remove="showRemoveEventsRecordModal"
    />
  </c-page>
</template>

<script>
import { ref, onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS } from '@/constants';

import { getEventsRecordFileUrl } from '@/helpers/entities/events-record/url';

import { useModals } from '@/hooks/modals';
import { useEventsRecord } from '@/hooks/store/modules/events-record';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useQueryOptions } from '@/hooks/query/options';
import { useExportFile } from '@/hooks/export-file';

import EventsRecordsHeader from '@/components/other/events-record/events-records-header.vue';
import EventsRecordsList from '@/components/other/events-record/events-records-list.vue';

export default {
  components: { EventsRecordsHeader, EventsRecordsList },
  setup() {
    const eventsRecords = ref([]);
    const meta = ref({});
    const status = ref({});
    const modals = useModals();

    /**
     * STORE
     */
    const {
      createEventsRecordExport,
      fetchEventsRecordExport,
      startEventsRecord,
      stopEventsRecord,
      removeEventsRecord,
      fetchEventsRecordsListWithoutStore,
    } = useEventsRecord();

    /**
     * QUERY
     */
    const {
      query,
      pending,
      updateQuery,
      handler: fetchList,
    } = usePendingWithLocalQuery({
      initialQuery: { page: 1, itemsPerPage: PAGINATION_LIMIT },
      fetchHandler: async (fetchQuery) => {
        const response = await fetchEventsRecordsListWithoutStore({
          params: {
            limit: fetchQuery.itemsPerPage,
            page: fetchQuery.page,
          },
        });

        eventsRecords.value = response.data;
        meta.value = response.meta;
        status.value = response.status;
      },
    });

    const { options } = useQueryOptions(query, updateQuery);

    /**
     * METHODS
     */
    const start = () => modals.show({
      name: MODALS.startEventsRecord,
      config: {
        action: async (data) => {
          await startEventsRecord({ data });

          return fetchList();
        },
      },
    });
    const stop = async () => stopEventsRecord();
    const showEventsRecordModal = eventsRecord => modals.show({
      name: MODALS.eventsRecord,
      config: {
        eventsRecord,
      },
    });

    const showRemoveEventsRecordModal = id => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await removeEventsRecord({ id });

          return fetchList();
        },
      },
    });

    const {
      generateFile,
      downloadFile,
    } = useExportFile({
      createHandler: createEventsRecordExport,
      fetchHandler: fetchEventsRecordExport,
    });

    const downloadEventsRecord = async (eventsRecord) => {
      const fileData = await generateFile({ id: eventsRecord._id });

      downloadFile(getEventsRecordFileUrl(fileData._id));
    };

    onMounted(() => fetchList(query.value));

    return {
      eventsRecords,
      meta,
      status,
      pending,
      options,

      start,
      stop,
      showEventsRecordModal,
      showRemoveEventsRecordModal,
      downloadEventsRecord,
      updateQuery,
      fetchList,
    };
  },
};
</script>
