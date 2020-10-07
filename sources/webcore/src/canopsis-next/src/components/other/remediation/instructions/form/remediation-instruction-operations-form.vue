<template lang="pug">
  v-layout(column)
    v-layout(v-if="hasSavedOperations")
      v-flex(xs10)
        v-layout(justify-end)
          v-btn.mx-0.secondary(
            v-if="allCollapsed",
            @click="expandAllItems"
          ) {{ $t('remediationInstructions.expandAll') }}
          v-btn.mx-0.primary(
            v-else,
            @click="collapseAllItems"
          ) {{ $t('remediationInstructions.hideAll') }}
    draggable(
      v-field="operations",
      :options="draggableOptions",
      :class="{ 'grey lighten-2': isDragging }",
      @start="startDragging",
      @end="endDragging"
    )
      v-layout.py-1(v-for="(operation, index) in operations", :key="operation.key")
        v-flex.mt-3(xs1)
          draggable-step-number(
            drag-class="operation-drag-handler",
            :draggable="!hideActions"
          ) {{ getStepLabel(index) }}
        v-flex(xs11)
          remediation-instruction-operation-field(
            v-field="operations[index]",
            :expanded="isExpandedOperation(operation)",
            :hideActions="hideActions",
            @expand="expandHandler(operation)",
            @remove="removeOperation(index)"
          )
    v-layout(v-show="!hideActions", row)
      div
        v-btn.ml-0.accent.darken-1(@click="addOperation") {{ $t('remediationInstructions.addOperation') }}
        div.error--text(v-show="errors.has(fieldName)") {{ $t('remediationInstructions.errors.operationRequired') }}
      v-btn.accent.darken-1(v-if="hasOperations", @click="addEndpoint") {{ $t('remediationInstructions.addEndpoint') }}
</template>

<script>
import Draggable from 'vuedraggable';

import { generateRemediationInstructionStepOperation } from '@/helpers/entities';
import { VUETIFY_ANIMATION_DELAY } from '@/config';
import { FIRST_LETTER_ALPHABET_CHAR_CODE, MODALS } from '@/constants';

import formArrayMixin from '@/mixins/form/array';

import DraggableStepNumber from '../partials/draggable-step-number.vue';

import RemediationInstructionOperationField from './fields/remediation-instruction-operation-field.vue';

export default {
  components: {
    Draggable,
    DraggableStepNumber,
    RemediationInstructionOperationField,
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
    hideActions: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      collapsedItems: [],
      isDragging: false,
    };
  },
  computed: {
    fieldName() {
      return `operations${this.step.key ? this.step.key : ''}`;
    },

    savedOperation() {
      return this.operations.filter(operation => operation.saved);
    },

    hasOperations() {
      return !!this.operations.length;
    },

    hasSavedOperations() {
      return !!this.savedOperation.length;
    },

    allCollapsed() {
      return this.savedOperation.length === this.collapsedItems.length;
    },

    draggableOptions() {
      return {
        animation: VUETIFY_ANIMATION_DELAY,
        handle: '.operation-drag-handler',
        ghostClass: 'white',
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

    addEndpoint() {},

    isExpandedOperation(operation) {
      return !this.collapsedItems.includes(operation.key);
    },

    expandHandler(operation) {
      if (this.isExpandedOperation(operation)) {
        this.collapsedItems.push(operation.key);
      } else {
        this.collapsedItems = this.collapsedItems.filter(id => id !== operation.key);
      }
    },

    collapseAllItems() {
      this.collapsedItems = this.operations.reduce((acc, operation) => {
        if (operation.saved) {
          acc.push(operation.key);
        }

        return acc;
      }, []);
    },

    expandAllItems() {
      this.collapsedItems = [];
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
