<template>
  <v-layout
    class="mt-2"
    column
  >
    <v-layout>
      <v-flex xs11>
        <v-alert
          :value="!operations.length"
          type="info"
        >
          {{ $t('remediation.instruction.emptyOperations') }}
        </v-alert>
      </v-flex>
    </v-layout>
    <c-draggable-list-field
      v-field="operations"
      :disabled="disabled"
      :class="{ 'grey lighten-1': isDragging }"
      :group="draggableGroup"
      ghost-class="grey"
      handle=".operation-drag-handler"
      @start="startDragging"
      @end="endDragging"
    >
      <remediation-instruction-operation-field
        v-field="operations[index]"
        v-for="(operation, index) in operations"
        :key="operation.key"
        :index="index"
        :operation-number="getOperationNumber(index)"
        :disabled="disabled"
        class="py-1"
        @remove="removeOperation(index)"
      />
    </c-draggable-list-field>
    <v-layout align-center>
      <v-btn
        :color="hasOperationsErrors ? 'error' : 'primary'"
        :disabled="disabled"
        class="mr-2"
        outlined
        @click="addOperation"
      >
        {{ $t('remediation.instruction.addOperation') }}
      </v-btn>
      <span
        v-show="hasOperationsErrors"
        class="error--text"
      >
        {{ $t('remediation.instruction.errors.operationRequired') }}
      </span>
    </v-layout>
  </v-layout>
</template>

<script>
import { remediationInstructionStepOperationToForm } from '@/helpers/entities/remediation/instruction/form';
import { getLetterByIndex } from '@/helpers/string';

import { formArrayMixin } from '@/mixins/form';

import RemediationInstructionOperationField from './fields/remediation-instruction-operation-field.vue';

export default {
  inject: ['$validator'],
  components: {
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

    draggableGroup() {
      return {
        name: 'remediation-instruction-operations',
        pull: false,
        put: false,
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
