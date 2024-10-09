<template>
  <c-page @refresh="fetchList">
    <v-expand-transition>
      <events-records-header
        v-if="!isCurrentEmpty"
        :current="current"
        @start:recording="startRecording"
        @stop:recording="stopRecording"
        @stop:resending="stopResending"
      />
    </v-expand-transition>
    <events-records-list
      :events-records="preparedEventsRecords"
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
import { pick, isEmpty } from 'lodash';
import { computed, ref, onMounted } from 'vue';

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
    const { current } = useEventsRecordCurrent();

    useEventRecordCurrentPolling();

    const isCurrentEmpty = computed(() => isEmpty(current.value));

    const preparedEventsRecords = computed(() => (
      current.value?._id
        ? eventsRecords.value.map(
          eventsRecord => (
            eventsRecord._id === current.value._id
              ? { ...eventsRecord, ...pick(current.value, ['is_resending', 'is_recording', 'count']) }
              : eventsRecord
          ),
        )
        : eventsRecords.value
    ));

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
      preparedEventsRecords,
      meta,
      current,
      isCurrentEmpty,
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
