<template>
  <c-card-iterator-form
    v-field="steps"
    :disabled="disabled"
    :draggable-group="draggableGroup"
    :required-error-message="$t('remediation.instruction.errors.stepRequired')"
    :empty-message="$t('remediation.instruction.emptySteps')"
    :add-button-label="$t('remediation.instruction.addStep')"
    :name="name"
    item-key="key"
    iterator-class="mb-2"
    required
    @add="addStep"
  >
    <template #item="{ index }">
      <remediation-instruction-step-field
        v-field="steps[index]"
        :step-number="index + 1"
        :disabled="disabled"
        :template-vars="templateVars"
        @remove="removeStep(index)"
      />
    </template>
  </c-card-iterator-form>
</template>

<script>
import { remediationInstructionStepToForm } from '@/helpers/entities/remediation/instruction/form';

import { formArrayMixin } from '@/mixins/form';

import RemediationInstructionStepField from './fields/remediation-instruction-step-field.vue';

export default {
  inject: ['$validator'],
  components: {
    RemediationInstructionStepField,
  },
  mixins: [formArrayMixin],
  model: {
    prop: 'steps',
    event: 'input',
  },
  props: {
    steps: {
      type: Array,
      default: () => ([]),
    },
    name: {
      type: String,
      default: 'steps',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      isDragging: false,
    };
  },
  computed: {
    hasStepsErrors() {
      return this.errors.has(this.name);
    },

    draggableGroup() {
      return {
        name: 'remediation-instruction-steps',
      };
    },
  },
  methods: {
    addStep() {
      this.addItemIntoArray(remediationInstructionStepToForm());
    },

    removeStep(index) {
      this.removeItemFromArray(index);
    },
  },
};
</script>
