<template>
  <c-page @refresh="fetchList">
    <events-recordings-header
      :progress="meta.inProgress"
      @launch="launch"
      @stop="stop"
    />
    <events-recordings-list
      :events-recording="eventsRecordings"
      :pending="pending"
      :options.sync="options"
      :total-items="meta.total_count"
      @show="showEventsRecordingModal"
      @remove="showRemoveEventsRecordingModal"
    />
  </c-page>
</template>

<script>
import { ref, onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS } from '@/constants';

import { useModals } from '@/hooks/modals';
import { useEventsRecording } from '@/hooks/store/modules/events-recording';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useQueryOptions } from '@/hooks/query/options';

import EventsRecordingsHeader from '@/components/other/events-recording/events-recordings-header.vue';
import EventsRecordingsList from '@/components/other/events-recording/events-recordings-list.vue';

export default {
  components: { EventsRecordingsHeader, EventsRecordingsList },
  setup() {
    const eventsRecordings = ref([]);
    const meta = ref({});
    const modals = useModals();

    /**
     * STORE
     */
    const {
      launchEventsRecording,
      stopEventsRecording,
      removeEventsRecording,
      fetchEventsRecordingsListWithoutStore,
    } = useEventsRecording();

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
        const response = await fetchEventsRecordingsListWithoutStore({
          params: {
            limit: fetchQuery.itemsPerPage,
            page: fetchQuery.page,
          },
        });

        eventsRecordings.value = response.data;
        meta.value = response.meta;
      },
    });

    const { options } = useQueryOptions(query, updateQuery);

    /**
     * METHODS
     */
    const launch = () => modals.show({
      name: MODALS.launchEventsRecording,
      config: {
        action: async (pattern) => {
          await launchEventsRecording(pattern);

          return fetchList();
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

          return fetchList();
        },
      },
    });

    const downloadEventsRecording = () => {};

    onMounted(() => fetchList(query.value));

    return {
      eventsRecordings,
      meta,
      pending,
      options,

      launch,
      stop,
      showEventsRecordingModal,
      showRemoveEventsRecordingModal,
      downloadEventsRecording,
      updateQuery,
      fetchList,
    };
  },
};
</script>
