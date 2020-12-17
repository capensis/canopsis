<template lang="pug">
  div
    v-layout(row, wrap)
      v-text-field(
        v-field="form.name",
        v-validate="'required'",
        :label="$t('common.name')",
        :error-messages="errors.collect('name')",
        name="name"
      )
    v-layout(row, wrap)
      time-interval-field(
        v-field="form.expectedInterval",
        :intervalLabel="$t('modals.statsDateInterval.fields.periodValue')",
        :unitLabel="$t('modals.statsDateInterval.fields.periodUnit')",
        :units="periodUnits"
      )
    v-layout(row, wrap)
      v-textarea(
        v-field="form.description",
        v-validate="'required'",
        :label="$t('common.description')",
        :error-messages="errors.collect('description')",
        name="description"
      )
    v-layout(row, wrap)
      v-text-field(
        v-field="form.output",
        :label="$t('common.output')"
      )
    v-layout
      v-btn.ml-0(@click="showEditPatternModal") {{ $t('modals.eventFilterRule.editPattern') }}
    v-layout
      v-alert(:value="errors.has('pattern')", type="error") {{ $t('modals.createHeartbeat.patternRequired') }}
</template>

<script>
import { isEmpty } from 'lodash';

import { MODALS, HEARTBEAT_DURATION_UNITS } from '@/constants';

import formMixin from '@/mixins/form';

import TimeIntervalField from '@/components/forms/fields/time-interval.vue';

/**
 * Modal to create widget
 */
export default {
  inject: ['$validator'],
  components: { TimeIntervalField },
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
        name: MODALS.createPattern,
        config: {
          onlySimpleRule: true,
          pattern: this.form.pattern,
          action: (pattern) => {
            this.updateField('pattern', pattern);
            this.$nextTick(() => this.$validator.validate('pattern'));
          },
        },
      });
    },
  },
};
</script>
