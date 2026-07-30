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

import { useModals } from '@/hooks/modals';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useQueryOptions } from '@/hooks/query/options';
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

        return {
          ...eventsRecord,
          is_recording: !!recording,
          is_resending: !!resending,
          progress: recording
            ? `${recording.n ?? 0}${recording.n && eventsRecord.count ? `/${eventsRecord.count}` : ''}`
            : undefined,
        };
      });
    };

    const {
      pending,
      query,
      fetchHandlerWithQuery: fetchList,
      updateQuery,
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

    const { options, updateOptions } = useQueryOptions(query, updateQuery);
    const { startRecording, stopRecording } = useEventsRecordRecording(fetchList);
    useEventRecordCurrentPolling();

    watch([eventsRecordRecordingsById, eventsRecordResendingsById], prepareEventsRecords);

    const showRemoveEventsRecordModal = id =>
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await removeEventsRecord({ id });

            return fetchList();
          },
        },
      });

    const showEventsRecordModal = eventsRecord =>
      modals.show({
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
