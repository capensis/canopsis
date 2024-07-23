<template>
  <c-page @refresh="fetchList">
    <events-records-header
      :progress="meta.inProgress"
      @start="start"
      @stop="stop"
    />
    <events-records-list
      :events-records="eventsRecords"
      :pending="pending"
      :options.sync="options"
      :total-items="meta.total_count"
      @show="showEventsRecordModal"
      @remove="showRemoveEventsRecordModal"
    />
  </c-page>
</template>

<script>
import { ref, onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS } from '@/constants';

import { useModals } from '@/hooks/modals';
import { useEventsRecord } from '@/hooks/store/modules/events-record';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useQueryOptions } from '@/hooks/query/options';

import EventsRecordsHeader from '@/components/other/events-record/events-records-header.vue';
import EventsRecordsList from '@/components/other/events-record/events-records-list.vue';

export default {
  components: { EventsRecordsHeader, EventsRecordsList },
  setup() {
    const eventsRecords = ref([]);
    const meta = ref({});
    const modals = useModals();

    /**
     * STORE
     */
    const {
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

        eventsRecords.value = response.records;
        meta.value = {};
      },
    });

    const { options } = useQueryOptions(query, updateQuery);

    /**
     * METHODS
     */
    const start = () => modals.show({
      name: MODALS.startEventsRecord,
      config: {
        action: async (pattern) => {
          await startEventsRecord(pattern);

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

    const downloadEventsRecord = () => {};

    onMounted(() => fetchList(query.value));

    return {
      eventsRecords,
      meta,
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
