<template>
  <c-page @refresh="fetchList">
    <v-expand-transition>
      <events-records-header
        :recordings="recordings"
        :resendings="resendings"
        :limit="limit"
        @start:recording="startRecording"
      />
    </v-expand-transition>
    <events-records-list
      :events-records="eventsRecords"
      :pending="pending"
      :options.sync="options"
      :total-items="meta.total_count"
      @stop="stopRecording"
      @show="showEventsRecordModal"
      @remove="showRemoveEventsRecordModal"
    />
  </c-page>
</template>

<script>
import { computed, ref, watch, onMounted } from 'vue';

import { MODALS } from '@/constants';

import { useModals } from '@/hooks/modals';
import { useEventsRecord } from '@/hooks/store/modules/events-record';
import { useEventsRecordCurrent } from '@/hooks/store/modules/events-record-current';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useQueryOptions } from '@/hooks/query/options';

import { useEventsRecordRecording } from '@/components/other/events-record/hooks/recording';
import { useEventsRecordResending } from '@/components/other/events-record/hooks/resending';
import { useEventRecordCurrentPolling } from '@/components/other/events-record/hooks/current-polling';

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
    const { removeEventsRecord, fetchEventsRecordsListWithoutStore } = useEventsRecord();
    const { resendings, recordings, limit, recordingsById, resendingsById } = useEventsRecordCurrent();

    const prepareEventsRecords = () => eventsRecords.value.forEach((eventRecord, index) => {
      const recording = recordingsById.value[eventRecord._id];
      const resending = resendingsById.value[eventRecord._id];
      const isRecording = !!recording;
      const isResending = !!resending;

      if (eventRecord.is_recording !== isRecording || eventRecord.is_resending !== isResending) {
        eventsRecords.value[index] = {

          ...eventRecord,

          count: recording?.count || 0,
          is_recording: isRecording,
          is_resending: isResending,
        };
      }
    });

    useEventRecordCurrentPolling();

    const inProgressCount = computed(() => recordings.value.length);

    watch(() => [resendings, recordings], prepareEventsRecords);

    /**
     * QUERY
     */
    const {
      query,
      pending,
      updateQuery,
      handler: fetchList,
    } = usePendingWithLocalQuery({
      fetchHandler: async (fetchQuery) => {
        const response = await fetchEventsRecordsListWithoutStore({
          params: {
            limit: fetchQuery.itemsPerPage,
            page: fetchQuery.page,
          },
        });

        eventsRecords.value = response.data;
        meta.value = response.meta;

        prepareEventsRecords();
      },
    });

    const { options } = useQueryOptions(query, updateQuery);

    /**
     * RECORDING AND RESENDING
     */
    const { startRecording, stopRecording } = useEventsRecordRecording(fetchList);
    const { stopResending } = useEventsRecordResending();

    /**
     * METHODS
     */
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

        fetchList,
        removeAction: showRemoveEventsRecordModal,
      },
    });

    onMounted(() => fetchList(query.value));

    return {
      eventsRecords,
      meta,
      recordings,
      resendings,
      limit,
      inProgressCount,
      pending,
      options,

      fetchList,
      startRecording,
      stopRecording,
      stopResending,
      showEventsRecordModal,
      showRemoveEventsRecordModal,
    };
  },
};
</script>
