<template>
  <div>
    <c-patterns-field
      v-field="form"
      :alarm-attributes="alarmAttributes"
      :entity-attributes="entityAttributes"
      with-alarm
      with-entity
      both-counters
    />
    <c-collapse-panel
      :title="$t('remediation.pattern.tabs.pbehaviorTypes.title')"
      class="mt-3"
    >
      <remediation-patterns-pbehavior-types-form v-field="form" />
    </c-collapse-panel>
  </div>
</template>

<script>
import { ALARM_PATTERN_FIELDS, ENTITY_PATTERN_FIELDS, QUICK_RANGES } from '@/constants';

import { formValidationHeaderMixin } from '@/mixins/form';

import RemediationPatternsPbehaviorTypesForm from '@/components/other/remediation/patterns/form/remediation-patterns-pbehavior-types-form.vue';

export default {
  inject: ['$validator'],
  components: {
    RemediationPatternsPbehaviorTypesForm,
  },
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
