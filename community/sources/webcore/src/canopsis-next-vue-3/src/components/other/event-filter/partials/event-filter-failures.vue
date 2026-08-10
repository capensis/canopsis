<template>
  <c-layout wrap>
    <div :style="gridFlexStyle(3)">
      <event-filter-failure-type-field
        :value="query.type"
        @update:model-value="updateQueryField('type', $event)"
      />
    </div>
    <div v-if="eventFilter.unread_failures_count" :style="gridFlexStyle(12)">
      <v-btn
        :loading="markAsReadPending"
        class="v-btn-legacy-m--y"
        color="primary"
        @click="markNewFailuresAsRead"
      >
        <v-icon left>
          done_all
        </v-icon>
        <span>{{ $t('eventFilter.markAsRead') }}</span>
      </v-btn>
    </div>
    <div :style="gridFlexStyle(12)">
      <event-filter-failures-list
        :failures="eventFilterFailures"
        :options="options"
        :total-items="eventFilterFailuresMeta.total_count"
        :pending="eventFilterFailuresPending"
        @update:options="updateOptions"
      />
    </div>
  </c-layout>
</template>

<script>
import { onMounted, ref, watch } from 'vue';

import { gridFlexStyle } from '@/helpers/layout';
import { convertQueryToRequest } from '@/helpers/query';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';
import { useEventFilter } from '@/hooks/store/modules/event-filter';

import EventFilterFailureTypeField from '../form/fields/event-filter-failure-type-field.vue';

import EventFilterFailuresList from './event-filter-failures-list.vue';

export default {
  components: { EventFilterFailureTypeField, EventFilterFailuresList },
  emits: ['refresh'],
  props: {
    eventFilter: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const { fetchEventFilterFailuresListWithoutStore, markNewEventFilterFailuresAsRead } =
      useEventFilter();

    const eventFilterFailures = ref([]);
    const eventFilterFailuresMeta = ref({});
    const eventFilterFailuresPending = ref(false);
    const markAsReadPending = ref(false);

    const {
      query,
      updateQueryField,
      options,
      updateOptions,
      handler: fetchList,
    } = useLocalQueryWithOptions({
      initialQuery: {
        type: undefined,
      },
      onUpdate: async (fetchQuery) => {
        eventFilterFailuresPending.value = true;

        try {
          const { data, meta } = await fetchEventFilterFailuresListWithoutStore({
            id: props.eventFilter._id,
            params: convertQueryToRequest(fetchQuery),
          });

          eventFilterFailures.value = data;
          eventFilterFailuresMeta.value = meta;
        } catch (err) {
          console.error(err);
        } finally {
          eventFilterFailuresPending.value = false;
        }
      },
    });

    const markNewFailuresAsRead = async () => {
      markAsReadPending.value = true;

      try {
        await markNewEventFilterFailuresAsRead({ id: props.eventFilter._id });

        fetchList();
        emit('refresh');
      } catch (err) {
        console.error(err);
      } finally {
        markAsReadPending.value = false;
      }
    };

    watch(() => props.eventFilter.failures_count, fetchList);
    onMounted(fetchList);

    return {
      gridFlexStyle,
      query,
      options,
      updateOptions,
      updateQueryField,
      eventFilterFailures,
      eventFilterFailuresMeta,
      eventFilterFailuresPending,
      markAsReadPending,
      markNewFailuresAsRead,
    };
  },
};
</script>
