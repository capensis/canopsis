<template>
  <v-form>
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <events-record-events-header
          :count="count"
          :is-recording="isRecording"
          @remove="remove"
          @stop:recording="stopRecording(eventsRecordId)"
        />

        <v-tabs v-model="activeTab" align-tabs="center">
          <v-tab :value="TABS.events" :disabled="isRecording" class="v-tab--tooltip">
            <v-tooltip :disabled="!isRecording" location="bottom">
              <template #activator="{ props: tooltipProps }">
                <span v-bind="tooltipProps">
                  {{ t('modals.eventsRecord.eventsTab') }}
                </span>
              </template>
              <span>{{ t('eventsRecord.eventsAreLoading') }}</span>
            </v-tooltip>
          </v-tab>
          <v-tab :value="TABS.pattern">
            {{ t('modals.eventsRecord.patternTab') }}
          </v-tab>
        </v-tabs>

        <v-window v-model="activeTab">
          <v-window-item :value="TABS.events">
            <div>
              <events-record-events-filter
                :events-record-id="eventsRecordId"
                :query="query"
                @apply:filter="applyEventFilter"
                @reset:filter="resetFilter"
              />
              <events-record-events-list
                :events-record-id="eventsRecordId"
                :events="events"
                :pending="pending"
                :resending="isResending"
                :resending-disabled="resendingDisabled"
                :options="options"
                :total-items="meta.total_count"
                @remove="removeEvent"
                @remove:selected="removeEvents"
                @start:resending="startResending(eventsRecordId, $event)"
                @stop:resending="stopResending(eventsRecordId)"
                @update:options="updateOptions"
              />
            </div>
          </v-window-item>
          <v-window-item :value="TABS.pattern">
            <c-layout class="mt-4">
              <c-event-filter-patterns-field :patterns="patternForm" readonly />
            </c-layout>
          </v-window-item>
        </v-window>
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed, onMounted, ref, watch } from 'vue'

import { MODALS } from '@/constants'

import { pickIds } from '@/helpers/array'
import { convertDateToString } from '@/helpers/date/date'
import { patternToForm } from '@/helpers/entities/pattern/form'

import { useI18n } from '@/hooks/i18n'
import { useModals } from '@/hooks/modals'
import { usePendingWithLocalQuery } from '@/hooks/query/shared'
import { useQueryOptions } from '@/hooks/query/options'
import { useEventsRecord } from '@/hooks/store/modules/events-record'

import { useEventsRecordRecording } from '@/components/other/events-record/hooks/recording'
import { useEventsRecordResending } from '@/components/other/events-record/hooks/resending'
import EventsRecordEventsFilter from '@/components/other/events-record/events-record-events-filter.vue'
import EventsRecordEventsHeader from '@/components/other/events-record/events-record-events-header.vue'
import EventsRecordEventsList from '@/components/other/events-record/events-record-events-list.vue'

const TABS = {
  events: 'events',
  pattern: 'pattern',
}

export default {
  name: MODALS.eventsRecord,

  components: {
    EventsRecordEventsHeader,
    EventsRecordEventsFilter,
    EventsRecordEventsList,
  },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n()
    const modals = useModals()
    const config = computed(() => props.modal.config ?? {})
    const events = ref([])
    const meta = ref({})
    const activeTab = ref(TABS.events)

    const eventsRecord = computed(() => config.value.eventsRecord ?? {})
    const eventsRecordId = computed(() => eventsRecord.value._id)

    const {
      removeEventsRecord,
      removeEventsRecordEvent,
      bulkRemoveEventsRecordEvents,
      fetchEventsRecordEventsListWithoutStore,
      eventsRecordRecordingsById,
      eventsRecordResendingsById,
    } = useEventsRecord()

    const { stopRecording } = useEventsRecordRecording(config.value.fetchList)
    const { startResending, stopResending } = useEventsRecordResending()

    const recording = computed(() => eventsRecordRecordingsById.value[eventsRecordId.value])
    const count = computed(() => recording.value?.n ?? eventsRecord.value.count ?? 0)
    const isRecording = computed(() => !!recording.value)
    const isResending = computed(() => !!eventsRecordResendingsById.value[eventsRecordId.value])
    const resendingDisabled = computed(() => isRecording.value || isResending.value)
    const patternForm = computed(() => patternToForm({ event_pattern: eventsRecord.value.pattern }))
    const title = computed(() =>
      t('modals.eventsRecord.title', { date: convertDateToString(eventsRecord.value.t) }),
    )

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
          id: eventsRecordId.value,
          params: {
            page: fetchQuery.page,
            limit: fetchQuery.itemsPerPage,
            event_pattern: JSON.stringify(fetchQuery.event_pattern),
          },
        })

        events.value = response.data
        meta.value = response.meta
      },
    })

    const { options, updateOptions } = useQueryOptions(query, updateQuery)

    const closeModal = () => modals.hide({ id: props.modal.id })
    const resetFilter = () => removeQueryField('event_pattern')

    const remove = () =>
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await removeEventsRecord({ id: eventsRecordId.value })
            await config.value.fetchList?.()

            return closeModal()
          },
        },
      })

    const removeEvent = (event) =>
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await removeEventsRecordEvent({ id: event._id })

            return fetchList()
          },
        },
      })

    const removeEvents = (eventsForRemove) =>
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await bulkRemoveEventsRecordEvents({ data: pickIds(eventsForRemove) })

            return fetchList()
          },
        },
      })

    const applyEventFilter = () =>
      modals.show({
        name: MODALS.applyEventFilter,
        config: {
          title: t('eventsRecord.applyEventFilter'),
          eventsRecordId: eventsRecordId.value,
          eventPattern: query.value.event_pattern,
          action: (eventPattern) => {
            updateQueryField('event_pattern', eventPattern)

            return fetchList()
          },
        },
      })

    watch(
      isRecording,
      (newIsRecording) => {
        if (newIsRecording && activeTab.value === TABS.events) {
          activeTab.value = TABS.pattern
        }
      },
      { immediate: true },
    )

    onMounted(() => fetchList(query.value))

    return {
      TABS,
      t,
      eventsRecordId,
      events,
      pending,
      meta,
      query,
      options,
      title,
      activeTab,
      count,
      isRecording,
      patternForm,
      isResending,
      resendingDisabled,
      closeModal,
      resetFilter,
      updateOptions,
      remove,
      removeEvent,
      removeEvents,
      applyEventFilter,
      startResending,
      stopResending,
      stopRecording,
    }
  },
}
</script>
