<template lang="pug">
  .step-field
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
              :name="name",
              box
            )
          v-flex.pl-2(xs3)
            v-text-field.step-time-complete-unit(
              :value="timeToComplete | duration('refreshFieldFormat')",
              :label="$t('remediationInstructions.timeToComplete')",
              readonly
            )
          v-flex.mt-3(xs1)
            v-layout(justify-center)
              v-btn.ma-0(icon, @click.prevent="$emit('remove')")
                v-icon(color="error") delete
        v-expand-transition(mode="out-in")
          v-layout(v-show="expanded", column)
            remediation-instruction-step-workflow-field(v-field="step.stop_on_fail")
            remediation-instruction-step-endpoint-field(v-field="step.endpoint")
            remediation-instruction-operations-form(
              v-field="step.operations",
              :step="step",
              :step-number="index + 1"
            )
</template>

<script>
import formMixin from '@/mixins/form';

import { getUnitValueFromOtherUnit } from '@/helpers/time';

import ExpandButton from '@/components/other/buttons/expand-button.vue';


import RemediationInstructionOperationsForm from '../remediation-instruction-operations-form.vue';

import RemediationInstructionStepWorkflowField from './remediation-instruction-step-workflow-field.vue';
import RemediationInstructionStepEndpointField from './remediation-instruction-step-endpoint-field.vue';

export default {
  inject: ['$validator'],
  components: {
    ExpandButton,
    RemediationInstructionStepWorkflowField,
    RemediationInstructionOperationsForm,
    RemediationInstructionStepEndpointField,
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

    name() {
      return `name${this.fieldSuffix}`;
    },

    nameErrorMessages() {
      return this.errors.collect(this.name).map(error => error.replace(this.fieldSuffix, ''));
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

<style lang="scss" scoped>
  .step-field {
    & /deep/ .step-expand {
      margin: 24px 2px 0 2px !important;
      width: 20px !important;
      height: 20px !important;
    }


    & /deep/ .step-time-complete-unit .v-input__slot {
      &:before, &:after {
        content: none !important;
      }
    }
  }
</style>
