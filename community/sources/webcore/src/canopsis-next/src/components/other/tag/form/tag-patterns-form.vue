<template>
  <c-patterns-field
    class="mt-2"
    v-field="form"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :readonly="readonly"
    with-alarm
    with-entity
    some-required
  />
</template>

<script>
import { ALARM_PATTERN_FIELDS, ENTITY_PATTERN_FIELDS } from '@/constants';

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
          value: ALARM_PATTERN_FIELDS.creationDate,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.duration,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.lastEventDate,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.lastUpdateDate,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.activationDate,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.ackAt,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.ackBy,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.ackMessage,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.ackInitiator,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.resolved,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.canceled,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.longOutput,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.initialOutput,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.initialLongOutput,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.totalStateChanges,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.tags,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.activated,
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
