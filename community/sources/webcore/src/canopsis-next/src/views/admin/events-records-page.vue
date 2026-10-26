<template>
  <c-page @refresh="fetchList">
    <v-expand-transition>
      <events-records-header
        :recordings="eventsRecordRecordings"
        :resendings="eventsRecordResendings"
        :limit="eventsRecordLimit"
        @start:recording="startRecording"
      />
    </v-expand-transition>

    <events-records-list
      :events-records="eventsRecords"
      :pending="pending"
      :options="options"
      :total-items="meta.total_count"
      @stop="stopRecording"
      @show="showEventsRecordModal"
      @remove="showRemoveEventsRecordModal"
      @update:options="updateOptions"
    />
  </c-page>
</template>

<script>
import { ref, watch } from 'vue';

import { MODALS } from '@/constants';

import { convertQueryToRequest } from '@/helpers/query';

import { useModals } from '@/hooks/modals';
import { useQueryOptions } from '@/hooks/query/options';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useEventsRecord } from '@/hooks/store/modules/events-record';

import { useEventRecordCurrentPolling } from '@/components/other/events-record/hooks/current-polling';
import { useEventsRecordRecording } from '@/components/other/events-record/hooks/recording';

import EventsRecordsHeader from '@/components/other/events-record/events-records-header.vue';
import EventsRecordsList from '@/components/other/events-record/events-records-list.vue';

export default {
  components: {
    EventsRecordsHeader,
    EventsRecordsList,
  },

  setup() {
    const eventsRecords = ref([]);
    const meta = ref({});
    const modals = useModals();

    const {
      removeEventsRecord,
      fetchEventsRecordsListWithoutStore,
      eventsRecordRecordings,
      eventsRecordResendings,
      eventsRecordLimit,
      eventsRecordRecordingsById,
      eventsRecordResendingsById,
    } = useEventsRecord();

    const prepareEventsRecords = () => {
      eventsRecords.value = eventsRecords.value.map((eventsRecord) => {
        const recording = eventsRecordRecordingsById.value[eventsRecord._id];
        const resending = eventsRecordResendingsById.value[eventsRecord._id];
        const liveCount = recording?.n;

        return {
          ...eventsRecord,
          count: liveCount ?? eventsRecord.count ?? 0,
          is_recording: !!recording,
          is_resending: !!resending,
          progress: recording ? String(liveCount ?? 0) : undefined,
        };
      });
    };

    const { fetchCurrentEventsRecord } = useEventRecordCurrentPolling();

    const {
      pending,
      query,
      fetchHandlerWithQuery: fetchList,
      updateQuery,
    } = usePendingWithLocalQuery({
      fetchHandler: async (fetchQuery) => {
        const response = await fetchEventsRecordsListWithoutStore({
          params: convertQueryToRequest(fetchQuery),
        });

        await fetchCurrentEventsRecord();

        eventsRecords.value = response.data;
        meta.value = response.meta;
        prepareEventsRecords();
      },
    });

    const { options, updateOptions } = useQueryOptions(query, updateQuery);
    const { startRecording, stopRecording } = useEventsRecordRecording(fetchList);

    watch([eventsRecordRecordingsById, eventsRecordResendingsById], prepareEventsRecords);

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

    fetchList();

    return {
      eventsRecords,
      meta,
      pending,
      options,
      eventsRecordRecordings,
      eventsRecordResendings,
      eventsRecordLimit,
      updateOptions,
      fetchList,
      startRecording,
      stopRecording,
      showEventsRecordModal,
      showRemoveEventsRecordModal,
    };
  },
};
</script>
