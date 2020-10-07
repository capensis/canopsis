<template lang="pug">
  v-layout.mt-2(column)
    v-layout(v-for="(operation, index) in operations", :key="operation.key")
      v-flex.mt-3(xs1)
        draggable-step-number
          span {{ stepNumber }}
          span {{ getCharByIndex(index) }}
      v-flex.pl-3(xs9)
        remediation-instruction-operation-field(
          v-field="operations[index]",
          @remove="removeOperation(index)"
        )
      v-flex(xs2)
    v-layout(v-show="!hideActions", row)
      div
        v-btn.ml-0.accent.darken-1(@click="addOperation") {{ $t('remediationInstructions.addOperation') }}
        div.error--text(v-show="errors.has(fieldName)") {{ $t('remediationInstructions.errors.operationRequired') }}
      v-btn.accent.darken-1(
        v-if="!!operations.length",
        @click="addEndpoint"
      ) {{ $t('remediationInstructions.addEndpoint') }}
</template>

<script>
import { generateRemediationInstructionStepOperation } from '@/helpers/entities';

import { FIRST_LETTER_ALPHABET_CHAR_CODE, MODALS } from '@/constants';

import formArrayMixin from '@/mixins/form/array';

import DraggableStepNumber from '../partials/draggable-step-number.vue';

import RemediationInstructionOperationField from './fields/remediation-instruction-operation-field.vue';

export default {
  components: { RemediationInstructionOperationField, DraggableStepNumber },
  inject: ['$validator'],
  mixins: [formArrayMixin],
  model: {
    prop: 'operations',
    event: 'input',
  },
  props: {
    operations: {
      type: Array,
      default: () => ([]),
    },
    stepNumber: {
      type: [String, Number],
      default: '',
    },
    hideActions: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    fieldName() {
      return `operations${this.stepNumber}`;
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
    getCharByIndex(index) {
      return String.fromCharCode(FIRST_LETTER_ALPHABET_CHAR_CODE + index);
    },

    addOperation() {
      this.addItemIntoArray(generateRemediationInstructionStepOperation());
      this.$nextTick(() => this.$validator.validate(this.fieldName));
    },

    removeOperation(index) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => {
            this.removeItemFromArray(index);
            this.$nextTick(() => this.$validator.validate(this.fieldName));
          },
        },
      });
    },

    addEndpoint() {},
  },
};
</script>
