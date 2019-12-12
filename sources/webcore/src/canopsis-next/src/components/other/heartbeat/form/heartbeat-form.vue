<template lang="pug">
  div
    v-layout(row, wrap)
      v-flex(xs3)
        v-text-field(
          v-field="form.periodValue",
          v-validate="'required'",
          :label="$t('modals.statsDateInterval.fields.periodValue')",
          :error-messages="errors.collect('periodValue')",
          type="number",
          name="periodValue"
        )
      v-flex
        v-select(
          v-field="form.periodUnit",
          v-validate="'required'",
          :items="periodUnits",
          :label="$t('modals.statsDateInterval.fields.periodUnit')",
          :error-messages="errors.collect('periodUnit')",
          name="periodUnit"
        )
    v-layout
      v-btn(@click="showEditPatternModal") {{ $t('modals.eventFilterRule.editPattern') }}
    v-layout
      v-alert(:value="errors.has('pattern')", type="error") {{ $t('modals.createHeartbeat.patternRequired') }}
</template>

<script>
import { isEmpty } from 'lodash';

import { MODALS, HEARTBEAT_DURATION_UNITS } from '@/constants';

import formMixin from '@/mixins/form';

/**
 * Modal to create widget
 */
export default {
  inject: ['$validator'],
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    periodUnits() {
      return [
        {
          text: this.$tc('common.times.minute'),
          value: HEARTBEAT_DURATION_UNITS.minute,
        },
        {
          text: this.$tc('common.times.hour'),
          value: HEARTBEAT_DURATION_UNITS.hour,
        },
      ];
    },
  },
  created() {
    this.$validator.attach({
      name: 'pattern',
      rules: 'required:true',
      getter: () => !isEmpty(this.form.pattern),
      context: () => this,
    });
  },
  methods: {
    showEditPatternModal() {
      this.$modals.show({
        name: MODALS.createEventFilterRulePattern,
        config: {
          isSimplePattern: true,
          pattern: this.form.pattern,
          action: (pattern) => {
            this.updateField('pattern', pattern);
            this.$validator.validate('pattern');
          },
        },
      });
    },
  },
};
</script>
