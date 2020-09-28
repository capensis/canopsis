<template lang="pug">
  v-layout(column)
    v-layout.my-1(v-for="(step, index) in steps", :key="step.key", xs-10)
      v-layout(row, wrap)
        v-flex.mt-3(xs1)
          draggable-step-number(:draggable="everySaved") {{ index + 1 }}
        v-flex.pl-3(xs11)
          remediation-instruction-step-field(
            v-field="steps[index]",
            :hide-actions="!everySaved",
            @remove="removeStep(index)"
          )
          remediation-instruction-operations-form(v-show="step.saved", v-field="steps[index].operations")
          remediation-instruction-step-actions(
            v-show="step.saved",
            :has-operations="!!step.operations.length",
            :add-disabled="!hasSteps || !everySaved",
            @add-step="addStepBetween(index)",
            @add-operation="addStepOperation(index)",
            @add-endpoint="addEndpoint(index)"
          )
    v-layout
      v-btn.ml-0.primary(v-if="!hasSteps", @click="addStep") {{ $t('remediationInstructions.addStep') }}
</template>

<script>
import { MODALS } from '@/constants';

import { generateRemediationInstructionStep, generateRemediationInstructionStepOperation } from '@/helpers/entities';

import formMixin from '@/mixins/form';
import formArrayMixin from '@/mixins/form/array';

import RemediationInstructionStepField from './remediation-instruction-step-field.vue';
import RemediationInstructionStepsWorkflowField from './remediation-instruction-steps-workflow-field.vue';
import RemediationInstructionOperationsForm from './remediation-instruction-operations-form.vue';
import DraggableStepNumber from './draggable-step-number.vue';
import RemediationInstructionStepActions from './remediation-instruction-step-actions.vue';

export default {
  components: {
    DraggableStepNumber,
    RemediationInstructionStepField,
    RemediationInstructionStepActions,
    RemediationInstructionOperationsForm,
    RemediationInstructionStepsWorkflowField,
  },
  inject: ['$validator'],
  mixins: [formMixin, formArrayMixin],
  model: {
    prop: 'steps',
    event: 'input',
  },
  props: {
    steps: {
      type: Array,
      default: () => ([]),
    },
  },
  computed: {
    everySaved() {
      return this.steps.every(({ saved }) => saved);
    },

    hasSteps() {
      return !!this.steps.length;
    },
  },
  methods: {
    addStep() {
      this.addItemIntoArray(generateRemediationInstructionStep());
    },

    addStepBetween(stepIndex) {
      const steps = [...this.steps];

      steps.splice(stepIndex + 1, 0, generateRemediationInstructionStep());

      this.updateModel(steps);
    },

    addStepOperation(stepIndex) {
      const step = this.steps[stepIndex];

      this.updateItemInArray(stepIndex, {
        ...step,
        operations: [...step.operations, generateRemediationInstructionStepOperation()],
      });
    },

    addEndpoint() {},

    removeStep(index) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => this.removeItemFromArray(index),
        },
      });
    },
  },
};
</script>
