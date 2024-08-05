<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <events-record-events-header :events-record="config.eventsRecord" />
        <events-record-events-list
          :events-record-id="config.eventsRecord._id"
          :events="events"
          :pending="pending"
          :downloading="downloading"
          :options="query"
          :total-items="meta.total_count"
          @resend="resendEvent"
          @export="exportJsonEvent"
          @remove="removeEvent"
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
import { computed, ref, onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS, TIME_UNITS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useEventsRecord } from '@/hooks/store/modules/events-record';
import { usePendingHandler } from '@/hooks/query/pending';
import { useLocalQuery } from '@/hooks/query/local-query';

import { useExportJson } from '@/components/other/events-record/hooks/export-json';
import EventsRecordEventsHeader from '@/components/other/events-record/events-record-events-header.vue';
import EventsRecordEventsList from '@/components/other/events-record/events-record-events-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.eventsRecord,
  components: { EventsRecordEventsHeader, EventsRecordEventsList, ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const events = ref([]);
    const meta = ref({});

    const resending = ref(false);

    const { t, tc } = useI18n();
    const { config, close, modals } = useInnerModal(props);

    const eventsRecord = computed(() => config.value.eventsRecord);
    const eventsRecordId = computed(() => eventsRecord.value._id);

    /**
     * EVENTS RECORD STORE MODULE
     */
    const {
      playbackEventsRecordEvents,
      /* stopPlaybackEventsRecordEvents,
      removeEventsRecordEvent, */
      fetchEventsRecordEventsListWithoutStore,
    } = useEventsRecord();

    /**
     * PENDING
     */
    const {
      pending,
      handler: fetchList,
    } = usePendingHandler(async (fetchQuery) => {
      const response = await fetchEventsRecordEventsListWithoutStore({
        id: config.value.eventsRecord._id,
        params: {
          limit: fetchQuery.itemsPerPage,
          page: fetchQuery.page,
        },
      });

      events.value = response.data;
      meta.value = response.meta;
    });

    /**
     * EXPORT JSON
     */
    const { downloadingsById, exportJson: exportJsonMethod } = useExportJson();

    const downloading = computed(() => downloadingsById.value[eventsRecordId.value]);

    const exportJson = () => exportJsonMethod(eventsRecordId.value, []);

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

    const resendEvent = (eventIds = []) => modals.show({
      name: MODALS.duration,
      config: {
        title: tc('eventsRecord.resendEvents', 1),
        label: t('eventsRecord.delayBetweenEvents'),
        units: [
          { value: TIME_UNITS.millisecond, text: 'common.times.millisecond' },
          { value: TIME_UNITS.second, text: 'common.times.second' },
        ],
        action: async (delay) => {
          resending.value = true;

          await playbackEventsRecordEvents({ id: eventsRecordId.value, data: { delay, event_ids: eventIds } });

          resending.value = false;
        },
      },
    });

    /* const stopResendEvent = async () => {
      await stopPlaybackEventsRecordEvents({ id: eventsRecordId });

      resending.value = false;
    }; */

    const remove = () => config.value.removeAction?.();

    const exportJsonEvent = () => {};
    const removeEvent = async () => {
      // await removeEventsRecordEvent();

      // return config.value.afterSubmit?.();
    };
    const applyEventFilter = () => {};

    const title = computed(() => (
      t('modals.eventsRecord.title', { date: convertDateToString(eventsRecord.value.t) })
    ));

    onMounted(() => fetchList(query.value));

    return {
      events,
      pending,
      meta,
      config,
      query,
      title,
      downloading,

      close,
      updateQuery,
      resendEvent,
      exportJsonEvent,
      removeEvent,
      remove,
      applyEventFilter,
      exportJson,
    };
  },
};
</script>
