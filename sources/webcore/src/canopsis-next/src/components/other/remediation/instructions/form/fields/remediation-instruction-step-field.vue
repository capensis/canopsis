<template lang="pug">
  v-layout(row)
    v-layout
      expand-button.step-expand(v-model="expanded")
      v-layout(column)
        v-layout(row)
          v-flex(xs8)
            v-text-field(
              v-field="step.name",
              v-validate="'required'",
              :label="$t('common.name')",
              :error-messages="nameErrorMessages",
              :name="nameFieldName",
              box
            )
          v-flex.pl-2(xs3)
            v-text-field.step-time-complete-unit(
              :value="timeToComplete | duration(undefined, 'refreshFieldFormat')",
              :label="$t('remediationInstructions.timeToComplete')",
              readonly
            )
          v-flex.mt-3(xs1)
            v-layout(justify-center)
              v-btn.ma-0(icon, small, @click.prevent="$emit('remove')")
                v-icon(color="error") delete
        v-expand-transition(mode="out-in")
          v-layout(v-show="expanded", column)
            remediation-instruction-steps-workflow-field(v-field="step.stop_on_fail")
            remediation-instruction-operations-form(
              v-field="step.operations",
              :name="operationFieldName",
              :step-number="index + 1"
            )
</template>

<script>
import formMixin from '@/mixins/form';

import { getUnitValueFromOtherUnit } from '@/helpers/time';

import ExpandButton from '@/components/other/buttons/expand-button.vue';

import RemediationInstructionOperationsForm from '../remediation-instruction-operations-form.vue';

import RemediationInstructionStepsWorkflowField from './remediation-instruction-steps-workflow-field.vue';

export default {
  inject: ['$validator'],
  components: {
    ExpandButton,
    RemediationInstructionStepsWorkflowField,
    RemediationInstructionOperationsForm,
  },
  mixins: [formMixin],
  model: {
    prop: 'step',
    event: 'input',
  },
  props: {
    step: {
      type: Object,
      required: true,
    },
    index: {
      type: Number,
      required: true,
    },
  },
  data() {
    return {
      expanded: true,
    };
  },
  computed: {
    fieldSuffix() {
      return this.step.key ? `-${this.step.key}` : '';
    },

    nameFieldName() {
      return `name${this.fieldSuffix}`;
    },

    operationFieldName() {
      return `operations${this.fieldSuffix}`;
    },

    nameErrorMessages() {
      return this.errors.collect(this.nameFieldName).map(error => error.replace(this.fieldSuffix, ''));
    },

    timeToComplete() {
      return this.step.operations.reduce((acc, operation) => {
        const { time_to_complete: { interval, unit } } = operation;

        return acc + getUnitValueFromOtherUnit(interval, unit);
      }, 0);
    },
  },
};
</script>

<style lang="scss">
  .step-expand {
    margin: 24px 2px 0 2px !important;
    width: 20px !important;
    height: 20px !important;
  }
  .step-time-complete-unit .v-input__slot {
    &:before, &:after {
      content: none !important;
    }
  }
</style>
