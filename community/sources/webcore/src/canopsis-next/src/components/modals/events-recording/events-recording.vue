<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.launchEventsRecording.title') }}</span>
      </template>
      <template #text="">
        <events-recording-events-list
          :events="events"
          :pending="pending"
          :options="query"
          :total-items="meta.total_count"
          @resend="resend"
          @export="exportJson"
          @remove="remove"
          @update:options="updateQuery"
        />
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="close"
        >
          {{ $t('common.cancel') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { ref, onMounted } from 'vue';

import { MODALS } from '@/constants';
import { PAGINATION_LIMIT } from '@/config';

import { useInnerModal } from '@/hooks/modals';
import { useEventsRecording } from '@/hooks/store/modules/events-recording';
import { usePendingHandler } from '@/hooks/query/pending';
import { useLocalQuery } from '@/hooks/query/local-query';

import EventsRecordingEventsList from '@/components/other/events-recording/events-recording-events-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.eventsRecording,
  components: { EventsRecordingEventsList, ModalWrapper },
  setup(props) {
    const events = ref([]);
    const meta = ref({});

    const { config, close } = useInnerModal(props);

    /**
     * EVENTS RECORDING STORE MODULE
     */
    const { fetchEventsRecordingEventsListWithoutStore } = useEventsRecording();

    /**
     * PENDING
     */
    const {
      pending,
      handler: fetchList,
    } = usePendingHandler(async (fetchQuery) => {
      const response = await fetchEventsRecordingEventsListWithoutStore({
        params: {
          limit: fetchQuery.itemsPerPage,
          page: fetchQuery.page,
        },
      });

      events.value = response.data;
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

    const resend = () => {};
    const exportJson = () => {};
    const remove = () => {};

    onMounted(() => fetchList(query.value));

    return {
      events,
      pending,
      meta,
      config,
      query,

      close,
      updateQuery,
      resend,
      exportJson,
      remove,
    };
  },
};
</script>
