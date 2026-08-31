<template>
  <c-calendar
    ref="calendar"
    :events="events"
    :loading="pending"
    readonly
    hide-details-menu
    @change:pagination="fetchTimespans"
  />
</template>

<script setup>
import { computed, nextTick, onMounted, ref } from 'vue';

import { COLORS } from '@/config';
import { DATETIME_FORMATS } from '@/constants';

import {
  convertDateToDateObjectByTimezone,
  convertDateToString,
  convertDateToTimestampByTimezone,
} from '@/helpers/date/date';
import { isFullDayEvent } from '@/helpers/entities/pbehavior/form';
import { pbehaviorToTimespanRequest } from '@/helpers/entities/pbehavior/timespans/form';

import { useInfo } from '@/hooks/store/modules/info';
import { usePbehaviorTimespan } from '@/hooks/store/modules/pbehavior-timespan';

const props = defineProps({
  pbehavior: {
    type: Object,
    required: true,
  },
});

const { timezone } = useInfo();
const { fetchPbehaviorTimespansListWithoutStore } = usePbehaviorTimespan();

const calendar = ref(null);
const pending = ref(false);
const timespans = ref([]);

const isTimed = computed(() => !isFullDayEvent(
    convertDateToDateObjectByTimezone(props.pbehavior.tstart, timezone.value),
    props.pbehavior.tstop &&
      convertDateToDateObjectByTimezone(props.pbehavior.tstop, timezone.value),
  ),
);

const events = computed(() => timespans.value.map((timespan, index) => {
    const fromString = convertDateToString(timespan.from, DATETIME_FORMATS.medium);
    const toString = convertDateToString(timespan.to, DATETIME_FORMATS.medium);

    return {
      id: index,
      start: convertDateToDateObjectByTimezone(timespan.from, timezone.value),
      end: convertDateToDateObjectByTimezone(timespan.to, timezone.value),
      timed: isTimed.value,
      name: `${fromString} - ${toString}`,
      color: COLORS.secondary,
    };
  }),
);

const fetchTimespans = async () => {
  const filled = calendar.value?.filled;

  if (!filled?.start || !filled?.end) {
    return;
  }

  pending.value = true;

  try {
    const from = convertDateToTimestampByTimezone(filled.start, timezone.value);
    const to = convertDateToTimestampByTimezone(filled.end, timezone.value);
    const data = pbehaviorToTimespanRequest({
      pbehavior: props.pbehavior,
      from,
      to,
    });

    timespans.value = await fetchPbehaviorTimespansListWithoutStore({ data });
  } finally {
    pending.value = false;
  }
};

onMounted(async () => {
  await nextTick();

  fetchTimespans();
});
</script>
