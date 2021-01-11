<template lang="pug">
  .step-field
    v-layout
      v-flex.mt-3(xs1)
        draggable-step-number(
          drag-class="step-drag-handler",
          :color="hasChildrenError ? 'error' : 'primary'"
        ) {{ stepNumber }}
      v-flex(xs11)
        v-layout(row)
          v-layout
            expand-button.step-expand(
              v-model="expanded",
              :color="!expanded && hasChildrenError ? 'error' : 'grey darken-3'"
            )
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
                    :value="timeToComplete | duration('refreshFieldFormat')",
                    :label="$t('remediationInstructions.timeToComplete')",
                    readonly
                  )
                v-flex.mt-3(xs1)
                  v-layout(justify-center)
                    v-btn.ma-0(icon, small, @click.prevent="remove")
                      v-icon(color="error") delete
              v-expand-transition(mode="out-in")
                v-layout(v-if="expanded", column)
                  remediation-instruction-step-workflow-field(v-field="step.stop_on_fail")
                  remediation-instruction-step-endpoint-field(v-field="step.endpoint")
                  remediation-instruction-operations-form(
                    v-field="step.operations",
                    :name="operationFieldName",
                    :step-number="stepNumber"
                  )
</template>

<script>
import formMixin from '@/mixins/form';
import validationChildrenMixin from '@/mixins/form/validation-children';

import { isOmitEqual } from '@/helpers/is-omit-equal';
import { toSeconds } from '@/helpers/duration';
import { generateRemediationInstructionStep } from '@/helpers/entities';

import confirmableFormMixin from '@/mixins/confirmable-form';

import ExpandButton from '@/components/common/buttons/expand-button.vue';

import DraggableStepNumber from '../../partials/draggable-step-number.vue';

import RemediationInstructionOperationsForm from '../remediation-instruction-operations-form.vue';

import RemediationInstructionStepWorkflowField from './remediation-instruction-step-workflow-field.vue';
import RemediationInstructionStepEndpointField from './remediation-instruction-step-endpoint-field.vue';

export default {
  inject: ['$validator'],
  components: {
    DraggableStepNumber,
    ExpandButton,
    RemediationInstructionOperationsForm,
    RemediationInstructionStepWorkflowField,
    RemediationInstructionStepEndpointField,
  },
  mixins: [
    formMixin,
    validationChildrenMixin,
    confirmableFormMixin({
      field: 'step',
      method: 'remove',
      comparator(step) {
        const emptyStep = generateRemediationInstructionStep();
        const paths = [
          'key',
          step.operations.length ? ['operations', 0, 'key'] : 'operations',
        ];

        return isOmitEqual(step, emptyStep, paths);
      },
    }),
  ],
  model: {
    prop: 'step',
    event: 'input',
  },
  props: {
    step: {
      type: Object,
      required: true,
    },
    stepNumber: {
      type: [Number, String],
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
        const { time_to_complete: { value, unit } } = operation;

        return acc + toSeconds(value, unit);
      }, 0);
    },
  },
  methods: {
    remove() {
      this.$emit('remove');
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
