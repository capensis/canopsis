<template lang="pug">
  v-layout(column)
    v-flex(v-if="!steps.length", xs12)
      v-alert(:value="true", type="info") {{ $t('remediation.instruction.emptySteps') }}
    c-card-iterator.mb-2(
      v-field="steps",
      item-key="key",
      :disabled="disabled",
      :draggable-group="draggableGroup"
    )
      template(#item="{ index }")
        remediation-instruction-step-field(
          v-field="steps[index]",
          :step-number="index + 1",
          :disabled="disabled",
          @remove="removeStep(index)"
        )
    v-layout(row, align-center)
      v-btn.ml-0(
        :color="hasStepsErrors ? 'error' : 'primary'",
        :disabled="disabled",
        outline,
        @click="addStep"
      ) {{ $t('remediation.instruction.addStep') }}
      span.error--text(v-show="hasStepsErrors") {{ $t('remediation.instruction.errors.stepRequired') }}
</template>

<script>
import { remediationInstructionStepToForm } from '@/helpers/forms/remediation-instruction';

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
  watch: {
    steps() {
      this.$validator.validate(this.name);
    },
  },
  created() {
    this.$validator.attach({
      name: this.name,
      rules: 'min_value:1',
      getter: () => this.steps.length,
      vm: this,
    });
  },
  beforeDestroy() {
    this.$validator.detach(this.name);
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
