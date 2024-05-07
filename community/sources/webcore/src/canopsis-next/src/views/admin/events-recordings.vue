<template>
  <div>
    <c-page-header />
    <v-card class="ma-4 mt-0">
      <events-recordings-header
        :progress="meta.inProgress"
        @launch="launch"
        @stop="stop"
      />
      <events-recordings-list
        :events-recording="eventsRecordings"
        :pending="pending"
        :options="query"
        :total-items="meta.total_count"
        @show="showEventsRecordingModal"
        @remove="showRemoveEventsRecordingModal"
        @update:options="updateQuery"
      />
    </v-card>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS } from '@/constants';

import { useLocalQuery } from '@/hooks/query/local-query';
import { useModals } from '@/hooks/modals';
import { useEventsRecording } from '@/hooks/store/modules/events-recording';
import { usePendingHandler } from '@/hooks/query/pending';

import EventsRecordingsHeader from '@/components/other/events-recording/events-recordings-header.vue';
import EventsRecordingsList from '@/components/other/events-recording/events-recordings-list.vue';

export default {
  components: { EventsRecordingsHeader, EventsRecordingsList },
  setup() {
    const eventsRecordings = ref([]);
    const meta = ref({});
    const modals = useModals();

    /**
     * EVENTS RECORDING STORE MODULE
     */
    const {
      launchEventsRecording,
      stopEventsRecording,
      removeEventsRecording,
      fetchEventsRecordingsListWithoutStore,
    } = useEventsRecording();

    /**
     * PENDING
     */
    const {
      pending,
      handler: fetchList,
    } = usePendingHandler(async (fetchQuery) => {
      const response = await fetchEventsRecordingsListWithoutStore({
        params: {
          limit: fetchQuery.itemsPerPage,
          page: fetchQuery.page,
        },
      });

      eventsRecordings.value = response.data;
      meta.value = response.meta;
    });

    /**
     * QUERY
     */
    const {
      query,
      updateQuery,
    } = useLocalQuery({
      initialQuery: { page: 1, itemsPerPage: PAGINATION_LIMIT },
      onUpdate: fetchList,
    });
    const fetchListWithPreviousParams = () => fetchList(query.value);

    /**
     * METHODS
     */
    const launch = () => modals.show({
      name: MODALS.launchEventsRecording,
      config: {
        action: async (pattern) => {
          await launchEventsRecording(pattern);

          return fetchListWithPreviousParams();
        },
      },
    });
    const stop = async () => stopEventsRecording();
    const showEventsRecordingModal = eventsRecording => modals.show({
      name: MODALS.eventsRecording,
      config: {
        eventsRecording,
      },
    });

    const showRemoveEventsRecordingModal = id => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await removeEventsRecording({ id });

          return fetchListWithPreviousParams();
        },
      },
    });

    const downloadEventsRecording = () => {};

    onMounted(() => fetchList(query.value));

    return {
      eventsRecordings,
      pending,
      meta,
      query,

      launch,
      stop,
      showEventsRecordingModal,
      showRemoveEventsRecordingModal,
      downloadEventsRecording,
      updateQuery,
    };
  },
};
</script>
