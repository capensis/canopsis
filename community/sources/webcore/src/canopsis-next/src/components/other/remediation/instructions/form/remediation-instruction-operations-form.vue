<template lang="pug">
  v-layout.mt-2(column)
    v-layout(row)
      v-flex(v-if="!operations.length", xs11)
        v-alert(:value="true", type="info") {{ $t('remediation.instruction.emptyOperations') }}
    draggable(
      v-field="operations",
      :options="draggableOptions",
      :class="{ 'grey lighten-1': isDragging }",
      @start="startDragging",
      @end="endDragging"
    )
      remediation-instruction-operation-field.py-1(
        v-for="(operation, index) in operations",
        v-field="operations[index]",
        :key="operation.key",
        :index="index",
        :operation-number="getOperationNumber(index)",
        :disabled="disabled",
        @remove="removeOperation(index)"
      )
    v-layout(row, align-center)
      v-btn.ml-0(
        :color="hasOperationsErrors ? 'error' : 'primary'",
        :disabled="disabled",
        outline,
        @click="addOperation"
      ) {{ $t('remediation.instruction.addOperation') }}
      span.error--text(v-show="hasOperationsErrors") {{ $t('remediation.instruction.errors.operationRequired') }}
</template>

<script>
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { remediationInstructionStepOperationToForm } from '@/helpers/forms/remediation-instruction';
import { getLetterByIndex } from '@/helpers/string';

import { formArrayMixin } from '@/mixins/form';

import RemediationInstructionOperationField from './fields/remediation-instruction-operation-field.vue';

export default {
  inject: ['$validator'],
  components: {
    Draggable,
    RemediationInstructionOperationField,
  },
  mixins: [formArrayMixin],
  model: {
    prop: 'operations',
    event: 'input',
  },
  props: {
    name: {
      type: String,
      default: 'operations',
    },
    operations: {
      type: Array,
      default: () => ([]),
    },
    stepNumber: {
      type: [String, Number],
      default: '',
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
    hasOperationsErrors() {
      return this.errors.has(this.name);
    },

    draggableOptions() {
      return {
        disabled: this.disabled,
        animation: VUETIFY_ANIMATION_DELAY,
        handle: '.operation-drag-handler',
        ghostClass: 'grey',
        group: {
          name: 'remediation-instruction-operations',
          pull: false,
          put: false,
        },
      };
    },
  },
  watch: {
    operations() {
      this.$validator.validate(this.name);
    },
  },
  created() {
    this.$validator.attach({
      name: this.name,
      rules: 'min_value:1',
      getter: () => this.operations.length,
      vm: this,
    });
  },
  beforeDestroy() {
    this.$validator.detach(this.name);
  },
  methods: {
    getOperationNumber(index) {
      return `${this.stepNumber}${getLetterByIndex(index)}`;
    },

    addOperation() {
      this.addItemIntoArray(remediationInstructionStepOperationToForm());
    },

    removeOperation(index) {
      this.removeItemFromArray(index);
    },

    startDragging() {
      this.isDragging = true;
    },

    endDragging() {
      this.isDragging = false;
    },
  },
};
</script>
