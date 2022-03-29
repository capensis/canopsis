<template lang="pug">
  c-patterns-field(
    v-field="form",
    :disabled="disabled",
    :alarm-attributes="alarmAttributes",
    :alarm-excluded-attributes="alarmExcludedAttributes",
    :entity-excluded-items="entityExcludedItems",
    with-alarm,
    with-entity,
    some-required
  )
</template>

<script>
import { ALARM_PATTERN_FIELDS, ENTITY_PATTERN_FIELDS, QUICK_RANGES } from '@/constants';

import { formValidationHeaderMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
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

    entityExcludedItems() {
      return [
        ENTITY_PATTERN_FIELDS.lastEventDate,
      ];
    },
  },
};
</script>
