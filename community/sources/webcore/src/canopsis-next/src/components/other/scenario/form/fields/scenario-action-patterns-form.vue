<template lang="pug">
  c-patterns-field(
    v-field="patterns",
    :name="name",
    :alarm-attributes="alarmAttributes",
    :alarm-excluded-attributes="alarmExcludedAttributes",
    :entity-excluded-items="entityExcludedItems",
    some-required,
    with-alarm,
    with-entity
  )
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form';
import { ALARM_PATTERN_FIELDS, ENTITY_PATTERN_FIELDS, QUICK_RANGES } from '@/constants';

export default {
  inject: ['$validator'],
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'patterns',
    event: 'input',
  },
  props: {
    patterns: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'patterns',
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
        ENTITY_PATTERN_FIELDS.impact,
        ENTITY_PATTERN_FIELDS.depends,
      ];
    },
  },
};
</script>
