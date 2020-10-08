<template lang="pug">
  v-layout.mt-2(column)
    v-layout.py-1(v-for="(operation, index) in operations", :key="operation.key")
      v-flex.mt-3(xs1)
        draggable-step-number(drag-class="operation-drag-handler") {{ getStepLabel(index) }}
      v-flex(xs11)
        remediation-instruction-operation-field(v-field="operations[index]", @remove="removeOperation(index)")
    v-layout(row)
      div
        v-btn.ml-0.accent.darken-1(@click="addOperation") {{ $t('remediationInstructions.addOperation') }}
        div.error--text(v-show="errors.has(fieldName)") {{ $t('remediationInstructions.errors.operationRequired') }}
</template>

<script>
import { generateRemediationInstructionStepOperation } from '@/helpers/entities';

import { FIRST_LETTER_ALPHABET_CHAR_CODE, MODALS } from '@/constants';

import formArrayMixin from '@/mixins/form/array';

import DraggableStepNumber from '../partials/draggable-step-number.vue';

import RemediationInstructionOperationField from './fields/remediation-instruction-operation-field.vue';

export default {
  components: {
    RemediationInstructionOperationField,
    DraggableStepNumber,
  },
  inject: ['$validator'],
  mixins: [formArrayMixin],
  model: {
    prop: 'operations',
    event: 'input',
  },
  props: {
    step: {
      type: Object,
      default: () => ({}),
    },
    operations: {
      type: Array,
      default: () => ([]),
    },
    stepNumber: {
      type: [String, Number],
      default: '',
    },
  },
  data() {
    return {
      isDragging: false,
    };
  },
  computed: {
    fieldName() {
      return `operations${this.step.key ? this.step.key : ''}`;
    },
  },
  watch: {
    operations() {
      this.$validator.validate(this.fieldName);
    },
  },
  created() {
    this.$validator.attach({
      name: this.fieldName,
      rules: 'min_value:1',
      getter: () => this.operations.length,
      context: () => this,
    });
  },
  beforeDestroy() {
    this.$validator.detach(this.fieldName);
  },
  methods: {
    getStepLabel(index) {
      return `${this.stepNumber}${this.getCharByIndex(index)}`;
    },

    getCharByIndex(index) {
      return String.fromCharCode(FIRST_LETTER_ALPHABET_CHAR_CODE + index);
    },

    addOperation() {
      this.addItemIntoArray(generateRemediationInstructionStepOperation());
    },

    removeOperation(index) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => {
            this.removeItemFromArray(index);
          },
        },
      });
    },
  },
};
</script>
