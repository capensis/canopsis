<template>
  <c-patterns-field
    v-field="form"
    :readonly="readonly"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :with-total-entity="withTotalEntity"
    :some-required="someRequired"
    with-alarm
    with-entity
    both-counters
  />
</template>

<script>
import { ALARM_PATTERN_FIELDS, ENTITY_PATTERN_FIELDS } from '@/constants';

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
    readonly: {
      type: Boolean,
      default: false,
    },
    withTotalEntity: {
      type: Boolean,
      default: false,
    },
    someRequired: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    entityAttributes() {
      return [
        {
          value: ENTITY_PATTERN_FIELDS.lastEventDate,
          options: { disabled: true },
        },
      ];
    },

    alarmAttributes() {
      return [
        {
          value: ALARM_PATTERN_FIELDS.creationDate,
        },
        {
          value: ALARM_PATTERN_FIELDS.ackAt,
        },
        {
          value: ALARM_PATTERN_FIELDS.resolved,
          options: { disabled: true },
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
          value: ALARM_PATTERN_FIELDS.activationDate,
          options: { disabled: true },
        },
      ];
    },
  },
};
</script>
