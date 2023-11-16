<template>
  <v-layout column>
    <c-alert
      :value="!steps.length"
      type="info"
    >
      {{ $t('remediation.instruction.emptySteps') }}
    </c-alert>
    <c-card-iterator-field
      class="mb-2"
      v-field="steps"
      item-key="key"
      :disabled="disabled"
      :draggable-group="draggableGroup"
    >
      <template #item="{ index }">
        <remediation-instruction-step-field
          v-field="steps[index]"
          :step-number="index + 1"
          :disabled="disabled"
          @remove="removeStep(index)"
        />
      </template>
    </c-card-iterator-field>
    <v-layout align-center>
      <v-btn
        class="mr-2"
        :color="hasStepsErrors ? 'error' : 'primary'"
        :disabled="disabled"
        outlined
        @click="addStep"
      >
        {{ $t('remediation.instruction.addStep') }}
      </v-btn>
      <span
        class="error--text"
        v-show="hasStepsErrors"
      >
        {{ $t('remediation.instruction.errors.stepRequired') }}
      </span>
    </v-layout>
  </v-layout>
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
