<template lang="pug">
  c-patterns-field(
    v-field="form",
    :disabled="disabled",
    :some-required="flapping",
    :alarm-attributes="alarmAttributes",
    :alarm-excluded-attributes="alarmExcludedAttributes",
    :entity-excluded-items="entityExcludedItems",
    with-alarm,
    with-entity
  )
</template>

<script>
import { ALARM_PATTERN_FIELDS, ENTITY_PATTERN_FIELDS, QUICK_RANGES } from '@/constants';

export default {
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    flapping: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
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
      ];
    },

    alarmExcludedAttributes() {
      return [
        ALARM_PATTERN_FIELDS.lastUpdateDate,
        ALARM_PATTERN_FIELDS.lastEventDate,
        ALARM_PATTERN_FIELDS.resolvedAt,
      ];
    },

    entityExcludedItems() {
      return [
        ENTITY_PATTERN_FIELDS.lastEventDate,
      ];
    },
  },
};
</script>
