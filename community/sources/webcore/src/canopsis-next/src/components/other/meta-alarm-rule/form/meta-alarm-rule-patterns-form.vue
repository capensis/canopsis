<template lang="pug">
  c-patterns-field(
    v-field="form",
    :alarm-attributes="alarmAttributes",
    :alarm-excluded-attributes="alarmExcludedAttributes",
    name="config",
    with-alarm,
    with-entity,
    with-event
  )
</template>

<script>
import { ALARM_PATTERN_FIELDS, QUICK_RANGES } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    intervalOptions() {
      return {
        intervalRanges: [QUICK_RANGES.custom],
      };
    },

    alarmAttributes() {
      return [
        {
          value: ALARM_PATTERN_FIELDS.creationDate,
          options: this.intervalOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.ackAt,
          options: this.intervalOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.resolvedAt,
          options: this.intervalOptions,
        },
      ];
    },

    alarmExcludedAttributes() {
      return [
        ALARM_PATTERN_FIELDS.lastUpdateDate,
        ALARM_PATTERN_FIELDS.lastEventDate,
      ];
    },
  },
};
</script>
