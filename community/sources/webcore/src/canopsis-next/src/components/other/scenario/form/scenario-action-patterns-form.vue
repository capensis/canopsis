<template>
  <c-patterns-field
    v-field="patterns"
    :name="name"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    some-required
    with-alarm
    with-entity
    both-counters
  />
</template>

<script>
import { ALARM_PATTERN_FIELDS, ENTITY_PATTERN_FIELDS } from '@/constants';

import { formValidationHeaderMixin } from '@/mixins/form';

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
    alarmAttributes() {
      return [
        {
          value: ALARM_PATTERN_FIELDS.creationDate,
        },
        {
          value: ALARM_PATTERN_FIELDS.ackAt,
        },
        {
          value: ALARM_PATTERN_FIELDS.lastUpdateDate,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.lastEventDate,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.resolved,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.activationDate,
          options: { disabled: true },
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
