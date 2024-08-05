<template>
  <c-page @refresh="fetchList">
    <v-expand-transition>
      <events-records-header
        v-if="status"
        :progress="status.isRecording"
        @start="start"
        @stop="stop"
      />
    </v-expand-transition>
    <events-records-list
      :events-records="preparedEventsRecords"
      :pending="pending"
      :options.sync="options"
      :total-items="meta.total_count"
      @show="showEventsRecordModal"
      @remove="showRemoveEventsRecordModal"
    />
  </c-page>
</template>

<script>
import { computed, ref, onMounted } from 'vue';

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
    const status = ref();

    const modals = useModals();

    const preparedEventsRecords = computed(() => (
      status.value?.isRecording
        ? eventsRecords.value.map(
          (eventsRecord, index) => (index ? eventsRecord : { ...eventsRecord, isRecording: true }),
        )
        : eventsRecords.value
    ));

    /**
     * STORE
     */
    const {
      startEventsRecordCurrent,
      stopEventsRecordCurrent,
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
          await startEventsRecordCurrent({ data });

          return fetchList();
        },
      },
    });
    const stop = async () => {
      await stopEventsRecordCurrent();

      return fetchList();
    };

    const showRemoveEventsRecordModal = id => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await removeEventsRecord({ id });

          return fetchList();
        },
      },
    });

    const showEventsRecordModal = eventsRecord => modals.show({
      name: MODALS.eventsRecord,
      config: {
        eventsRecord,

        removeAction: showRemoveEventsRecordModal,
      },
    });

    onMounted(() => fetchList(query.value));

    return {
      preparedEventsRecords,
      meta,
      status,
      pending,
      options,

      start,
      stop,
      showEventsRecordModal,
      showRemoveEventsRecordModal,
      updateQuery,
      fetchList,
    };
  },
};
</script>
