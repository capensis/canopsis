<template>
  <c-patterns-field
    v-field="form"
    :with-alarm="!isEntityType"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :readonly="readonly"
    :entity-counters-type="isEntityType"
    :both-counters="!isEntityType"
    some-required
    with-entity
  />
</template>

<script>
import { ALARM_PATTERN_FIELDS, ENTITY_PATTERN_FIELDS, QUICK_RANGES } from '@/constants';

import { formValidationHeaderMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [
    formValidationHeaderMixin,
  ],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    isEntityType: {
      type: Boolean,
      default: false,
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
          options: {
            intervalRanges: [QUICK_RANGES.custom],
          },
        },
        {
          value: ALARM_PATTERN_FIELDS.ackAt,
          options: {
            intervalRanges: [QUICK_RANGES.custom],
          },
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
