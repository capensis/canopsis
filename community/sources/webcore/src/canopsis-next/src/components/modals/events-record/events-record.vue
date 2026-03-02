<template>
  <v-form>
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <events-record-events-header
          :count="count"
          :is-recording="isRecording"
          @remove="remove"
          @stop:recording="stopRecording(eventsRecordId)"
        />
        <v-tabs v-model="activeTab" centered>
          <v-tab
            :href="`#${TABS.events}`"
            :disabled="isRecording"
            class="v-tab--tooltip"
          >
            <v-tooltip :disabled="!isRecording" bottom>
              <template #activator="{ on }">
                <span v-on="on">
                  {{ $t('modals.eventsRecord.eventsTab') }}
                </span>
              </template>
              <span>{{ $t('eventsRecord.eventsAreLoading') }}</span>
            </v-tooltip>
          </v-tab>
          <v-tab :href="`#${TABS.pattern}`">
            {{ $t('modals.eventsRecord.patternTab') }}
          </v-tab>
          <v-tabs-items v-model="activeTab">
            <v-tab-item :value="TABS.events">
              <div>
                <events-record-events-filter
                  :events-record-id="eventsRecordId"
                  :count="config.eventsRecord.count"
                  :has-filter-applied="hasFilterApplied"
                  @apply:filter="applyEventFilter"
                  @reset:filter="resetFilter"
                />
                <events-record-events-list
                  :events-record-id="eventsRecordId"
                  :events="events"
                  :pending="pending"
                  :resending="isResending"
                  :resending-disabled="resendingDisabled"
                  :options="query"
                  :total-items="meta.total_count"
                  @remove="removeEvent"
                  @remove:selected="removeEvents"
                  @start:resending="startResending(eventsRecordId, $event)"
                  @stop:resending="stopResending(eventsRecordId)"
                  @update:options="updateOptions"
                />
              </div>
            </v-tab-item>
            <v-tab-item :value="TABS.pattern">
              <v-layout class="mt-4" column>
                <c-event-filter-patterns-field
                  :patterns="patternForm"
                  readonly
                />
              </v-layout>
            </v-tab-item>
          </v-tabs-items>
        </v-tabs>
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
import { computed, ref, watch, onMounted } from 'vue';

import { MODALS, EVENT_FILTER_PATTERN_FIELDS } from '@/constants';

import { pickIds } from '@/helpers/array';
import { convertDateToString } from '@/helpers/date/date';
import { patternToForm } from '@/helpers/entities/pattern/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useEventsRecord } from '@/hooks/store/modules/events-record';
import { useEventsRecordCurrent } from '@/hooks/store/modules/events-record-current';
import { useQueryOptions } from '@/hooks/query/options';

import { useEventsRecordResending } from '@/components/other/events-record/hooks/resending';
import { useEventsRecordRecording } from '@/components/other/events-record/hooks/recording';

import EventsRecordEventsHeader from '@/components/other/events-record/events-record-events-header.vue';
import EventsRecordEventsFilter from '@/components/other/events-record/events-record-events-filter.vue';
import EventsRecordEventsList from '@/components/other/events-record/events-record-events-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

const TABS = {
  events: 'events',
  pattern: 'pattern',
};

export default {
  name: MODALS.eventsRecord,
  components: {
    EventsRecordEventsHeader,
    EventsRecordEventsFilter,
    EventsRecordEventsList,
    ModalWrapper,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const events = ref([]);
    const meta = ref({});
    const activeTab = ref(TABS.events);

    const { t } = useI18n();
    const { config, close, modals } = useInnerModal(props);

    const eventsRecord = computed(() => config.value.eventsRecord);
    const eventsRecordId = computed(() => eventsRecord.value._id);

    /**
     * EVENTS RECORD STORE MODULES
     */
    const {
      removeEventsRecord,
      removeEventsRecordEvent,
      bulkRemoveEventsRecordEvent,
      fetchEventsRecordEventsListWithoutStore,
    } = useEventsRecord();

    const { recordingsById, resendingsById } = useEventsRecordCurrent();
    const { stopRecording } = useEventsRecordRecording(config.value.fetchList);

    const recording = computed(() => recordingsById.value[eventsRecordId.value]);
    const count = computed(() => recording.value?.n || config.value.eventsRecord.n || 0);
    const isRecording = computed(() => !!recording.value);
    const isResending = computed(() => !!resendingsById.value[eventsRecordId.value]);

    const patternForm = computed(() => (
      patternToForm({ event_pattern: eventsRecord.value.pattern })
    ));

    const title = computed(() => (
      t('modals.eventsRecord.title', { date: convertDateToString(eventsRecord.value.t) })
    ));

    const resendingDisabled = computed(() => isRecording.value || isResending.value);

    watch(isRecording, (newIsRecording) => {
      if (newIsRecording && activeTab.value === TABS.events) {
        activeTab.value = TABS.pattern;
      }
    }, { immediate: true });

    /**
     * QUERY
     */
    const {
      pending,
      query,
      fetchHandlerWithQuery: fetchList,
      updateQuery,
      updateQueryField,
      removeQueryField,
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

    const hasFilterApplied = computed(() => query.value.event_pattern !== undefined);

    const resetFilter = () => removeQueryField('event_pattern');

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
      TABS,
      eventsRecordId,
      events,
      pending,
      meta,
      config,
      query,
      title,
      activeTab,
      count,
      isRecording,
      hasFilterApplied,
      patternForm,
      isResending,
      resendingDisabled,

      close,
      resetFilter,
      updateOptions,
      remove,
      removeEvent,
      removeEvents,
      applyEventFilter,

      startResending,
      stopResending,
      stopRecording,
    };
  },
};
</script>
