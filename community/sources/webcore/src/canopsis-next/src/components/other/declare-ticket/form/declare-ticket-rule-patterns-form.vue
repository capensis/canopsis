<template>
  <c-patterns-field
    v-field="form"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :readonly="readonly"
    some-required
    with-pbehavior
    with-alarm
    with-entity
    both-counters
  />
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
      default: () => ({}),
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    alarmAttributes() {
      return [
        {
          value: ALARM_PATTERN_FIELDS.lastEventDate,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.lastUpdateDate,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.resolvedAt,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.ackAt,
          options: {
            intervalRanges: [QUICK_RANGES.custom],
          },
        },
        {
          value: ALARM_PATTERN_FIELDS.creationDate,
          options: {
            intervalRanges: [QUICK_RANGES.custom],
          },
        },
      ];
    },

    entityAttributes() {
      return [
        {
          value: ENTITY_PATTERN_FIELDS.lastEventDate,
          options: { disabled: true },
        },
      ];
    },
  },
};
</script>
