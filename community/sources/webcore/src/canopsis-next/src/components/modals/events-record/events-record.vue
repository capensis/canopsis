<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <v-layout class="gap-3 my-4" justify-center align-center>
          <span class="text-subtitle-2">
            {{ $t('modals.eventsRecord.subtitle', { count: config.eventsRecord.count }) }}
          </span>
          <c-action-btn
            :tooltip="$t('modals.eventsRecord.buttonTooltip')"
            type="delete"
            @click="remove"
          />
        </v-layout>
        <v-layout class="gap-3 my-4">
          <v-btn
            color="primary"
            @click="applyEventFilter"
          >
            {{ $t('eventsRecord.applyEventFilter') }}
          </v-btn>
          <v-btn
            color="primary"
            @click="exportJson"
          >
            <v-icon class="mr-2" color="white">
              file_download
            </v-icon>
            <span>{{ $t('common.exportToJson') }}</span>
          </v-btn>
        </v-layout>
        <events-record-events-list
          :events="events"
          :pending="pending"
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

import EventsRecordEventsList from '@/components/other/events-record/events-record-events-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.eventsRecord,
  components: { EventsRecordEventsList, ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const events = ref([]);
    const meta = ref({});

    const { t, tc } = useI18n();
    const { config, close, modals } = useInnerModal(props);

    /**
     * EVENTS RECORD STORE MODULE
     */
    const { fetchEventsRecordEventsListWithoutStore } = useEventsRecord();

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
     * QUERY
     */
    const {
      query,
      updateQuery,
    } = useLocalQuery({
      initialQuery: { page: 1, itemsPerPage: PAGINATION_LIMIT },
      onUpdate: fetchList,
    });

    const resendEvent = () => modals.show({
      name: MODALS.duration,
      config: {
        title: tc('eventsRecord.resendEvents', 1),
        label: t('eventsRecord.delayBetweenEvents'),
        units: [
          { value: TIME_UNITS.millisecond, text: 'common.times.millisecond' },
          { value: TIME_UNITS.second, text: 'common.times.second' },
        ],
        action: () => {},
      },
    });
    const exportJsonEvent = () => {};
    const removeEvent = () => {};
    const remove = () => {};
    const applyEventFilter = () => {};
    const exportJson = () => {};

    const title = computed(() => (
      t('modals.eventsRecord.title', { date: convertDateToString(config.created) })
    ));

    onMounted(() => fetchList(query.value));

    return {
      events,
      pending,
      meta,
      config,
      query,
      title,

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
