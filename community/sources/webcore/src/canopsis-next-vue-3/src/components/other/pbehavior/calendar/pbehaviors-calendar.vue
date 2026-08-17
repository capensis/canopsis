<template>
  <c-calendar
    ref="calendar"
    :events="events"
    :loading="pending"
    :timezone="timezone"
    readonly
    hide-details-menu
    @update:timezone="timezone = $event"
    @change:pagination="fetchEvents"
  />
</template>

<script setup>
import { nextTick, onMounted, ref } from 'vue';

import { DATETIME_FORMATS } from '@/constants';

import { getMostReadableTextColor } from '@/helpers/color';
import {
  convertDateToDateObjectByTimezone,
  convertDateToString,
  convertDateToTimestampByTimezone,
  getLocalTimezone,
} from '@/helpers/date/date';
import { getPbehaviorColor, isFullDayEvent } from '@/helpers/entities/pbehavior/form';

import { usePbehavior } from '@/hooks/store/modules/pbehavior';

const props = defineProps({
  entityId: {
    type: String,
    default: '',
  },
});

const { fetchPbehaviorsCalendarWithoutStore, fetchEntitiesPbehaviorsCalendarWithoutStore } = usePbehavior();

const calendar = ref(null);
const pending = ref(false);
const events = ref([]);
const timezone = ref(getLocalTimezone());

const convertPbehaviorsCalendarToEvents = (pbehaviors = []) =>
  pbehaviors.map((pbehavior, index) => {
    const start = convertDateToDateObjectByTimezone(pbehavior.from, pbehavior.timezone);
    const end =
      pbehavior.to && convertDateToDateObjectByTimezone(pbehavior.to, pbehavior.timezone);
    const isTimed = !isFullDayEvent(start, end);
    const fromString = convertDateToString(pbehavior.from, DATETIME_FORMATS.medium);
    const toString = pbehavior.to
      ? convertDateToString(pbehavior.to, DATETIME_FORMATS.medium)
      : '';
    const color = getPbehaviorColor(pbehavior);
    const iconColor = getMostReadableTextColor(color, { level: 'AA', size: 'large' });

    return {
      id: index,
      color,
      iconColor,
      start,
      end,
      icon: pbehavior.type?.icon_name,
      name: `${fromString}${toString ? ` - ${toString}` : ''} ${pbehavior.title}`,
      timed: isTimed,
    };
  });

const fetchPbehaviorsCalendar = () => {
  const filled = calendar.value?.filled;

  if (!filled?.start || !filled?.end) {
    return [];
  }

  const params = {
    from: convertDateToTimestampByTimezone(filled.start, timezone.value),
    to: convertDateToTimestampByTimezone(filled.end, timezone.value),
  };

  if (props.entityId) {
    return fetchEntitiesPbehaviorsCalendarWithoutStore({
      params: {
        ...params,
        _id: props.entityId,
      },
    });
  }

  return fetchPbehaviorsCalendarWithoutStore({ params });
};

const fetchEvents = async () => {
  pending.value = true;

  try {
    const pbehaviorsCalendar = await fetchPbehaviorsCalendar();

    events.value = convertPbehaviorsCalendarToEvents(pbehaviorsCalendar);
  } finally {
    pending.value = false;
  }
};

onMounted(async () => {
  await nextTick();

  fetchEvents();
});
</script>
