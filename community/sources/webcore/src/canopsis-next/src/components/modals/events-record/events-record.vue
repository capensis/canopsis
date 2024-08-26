<template>
  <v-form>
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <events-record-events-header
          :events-record-id="eventsRecordId"
          :count="config.eventsRecord.count"
          @remove="remove"
          @apply:filter="applyEventFilter"
        />
        <events-record-events-list
          :events-record-id="eventsRecordId"
          :events="events"
          :pending="pending"
          :resending="resending"
          :resending-disabled="resendingDisabled"
          :options="query"
          :total-items="meta.total_count"
          @remove="removeEvent"
          @remove:selected="removeEvents"
          @start:resending="startResending(eventsRecordId, $event)"
          @stop:resending="stopResending"
          @update:options="updateOptions"
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

import { MODALS, EVENT_FILTER_PATTERN_FIELDS } from '@/constants';

import { pickIds } from '@/helpers/array';
import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useEventsRecord } from '@/hooks/store/modules/events-record';
import { useEventsRecordCurrent } from '@/hooks/store/modules/events-record-current';
import { useQueryOptions } from '@/hooks/query/options';

import { useEventsRecordResending } from '@/components/other/events-record/hooks/resending';

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

    const { t } = useI18n();
    const { config, close, modals } = useInnerModal(props);

    const eventsRecord = computed(() => config.value.eventsRecord);
    const eventsRecordId = computed(() => eventsRecord.value._id);

    const title = computed(() => (
      t('modals.eventsRecord.title', { date: convertDateToString(eventsRecord.value.t) })
    ));

    /**
     * EVENTS RECORD STORE MODULES
     */
    const {
      removeEventsRecord,
      removeEventsRecordEvent,
      bulkRemoveEventsRecordEvent,
      fetchEventsRecordEventsListWithoutStore,
    } = useEventsRecord();

    const { current } = useEventsRecordCurrent();

    const resending = computed(() => current.value.is_resending && current.value._id === eventsRecordId.value);
    const resendingDisabled = computed(() => current.value.is_recording || current.value.is_resending);

    /**
     * QUERY
     */
    const {
      pending,
      query,
      fetchHandlerWithQuery: fetchList,
      updateQuery,
      updateQueryField,
    } = usePendingWithLocalQuery({
      fetchHandler: async (fetchQuery) => {
        const response = await fetchEventsRecordEventsListWithoutStore({
          id: config.value.eventsRecord._id,
          params: {
            page: fetchQuery.page,
            limit: fetchQuery.itemsPerPage,
            event_pattern: JSON.stringify(fetchQuery.event_pattern),
          },
        });

        events.value = response.data;
        meta.value = response.meta;
      },
    });

    const { updateOptions } = useQueryOptions(query, updateQuery);

    /**
     * RESEND
     */
    const { startResending, stopResending } = useEventsRecordResending();

    /**
     * Remove events record
     *
     * @return {*}
     */
    const remove = () => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await removeEventsRecord({ id: eventsRecordId.value });
          await config.value.fetchList?.();

          return close();
        },
      },
    });

    /**
     * Remove event from events record
     *
     * @param {Object} event
     * @return {*}
     */
    const removeEvent = event => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await removeEventsRecordEvent({ id: event._id });

          return fetchList();
        },
      },
    });

    /**
     * Remove several events from events record
     *
     * @param {Object[]} eventsForRemove
     * @return {*}
     */
    const removeEvents = eventsForRemove => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await bulkRemoveEventsRecordEvent({ data: pickIds(eventsForRemove) });

          return fetchList();
        },
      },
    });

    /**
     * Apply event filter for events record events
     *
     * @return {*}
     */
    const applyEventFilter = () => modals.show({
      name: MODALS.applyEventFilter,
      config: {
        title: t('eventsRecord.applyEventFilter'),
        eventPattern: query.value.event_pattern,
        excludedAttributes: [
          { value: EVENT_FILTER_PATTERN_FIELDS.eventType },
          { value: EVENT_FILTER_PATTERN_FIELDS.state },
          { value: EVENT_FILTER_PATTERN_FIELDS.sourceType },
          { value: EVENT_FILTER_PATTERN_FIELDS.extraInfos },
          { value: EVENT_FILTER_PATTERN_FIELDS.longOutput },
          { value: EVENT_FILTER_PATTERN_FIELDS.author },
          { value: EVENT_FILTER_PATTERN_FIELDS.initiator },
        ],
        action: eventPattern => updateQueryField('event_pattern', eventPattern),
      },
    });

    onMounted(() => fetchList(query.value));

    return {
      eventsRecordId,
      events,
      pending,
      meta,
      config,
      query,
      title,
      resending,
      resendingDisabled,

      close,
      updateOptions,
      remove,
      removeEvent,
      removeEvents,
      applyEventFilter,

      startResending,
      stopResending,
    };
  },
};
</script>
