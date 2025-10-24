<template>
  <c-patterns-field
    v-field="form"
    :readonly="readonly"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    with-alarm
    with-entity
    both-counters
    some-required
  />
</template>

<script>
import { ALARM_PATTERN_FIELDS, ENTITY_PATTERN_FIELDS } from '@/constants';

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
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const alarmAttributes = [
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

    const entityAttributes = [
      {
        value: ENTITY_PATTERN_FIELDS.lastEventDate,
        options: { disabled: true },
      },
    ];

    return {
      alarmAttributes,
      entityAttributes,
    };
  },
};
</script>
