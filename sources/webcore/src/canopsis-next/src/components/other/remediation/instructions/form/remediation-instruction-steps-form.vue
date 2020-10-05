<template lang="pug">
  v-layout(column)
    draggable(
      :value="steps",
      :options="draggableOptions",
      :class="{ 'grey lighten-2': isDragging }",
      @change="changeStepsOrdering",
      @start="startDragging",
      @end="endDragging"
    )
      v-layout.my-1(v-for="(step, index) in steps", :key="step.key", row, wrap)
        v-flex.mt-3(xs1)
          draggable-step-number(drag-class="step-drag-handler", :draggable="allSaved") {{ index + 1 }}
        v-flex.pl-3(xs11)
          remediation-instruction-step-field(
            v-field="steps[index]",
            :hide-actions="!allSaved",
            @remove="removeStep(index)"
          )
          remediation-instruction-operations-form(
            v-field="steps[index].operations",
            :step="step",
            :hide-actions="!allSaved",
            :stepNumber="index + 1"
          )
    v-layout
      v-btn.ml-0.primary(v-if="allSaved", @click="addStep") {{ $t('remediationInstructions.addStep') }}
</template>

<script>
import Draggable from 'vuedraggable';

import { MODALS } from '@/constants';
import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { generateRemediationInstructionStep } from '@/helpers/entities';
import { dragDropChangePositionHandler } from '@/helpers/dragdrop';

import formArrayMixin from '@/mixins/form/array';

import DraggableStepNumber from '../partials/draggable-step-number.vue';

import RemediationInstructionStepField from './fields/remediation-instruction-step-field.vue';
import RemediationInstructionOperationsForm from './remediation-instruction-operations-form.vue';

export default {
  components: {
    Draggable,
    DraggableStepNumber,
    RemediationInstructionStepField,
    RemediationInstructionOperationsForm,
  },
  inject: ['$validator'],
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
  },
  data() {
    return {
      isDragging: false,
    };
  },
  computed: {
    allSaved() {
      return this.everyStepsSaved && this.everyOperationsSaved;
    },

    everyStepsSaved() {
      return this.steps.every(step => step.saved);
    },

    everyOperationsSaved() {
      return this.steps.every(step => step.operations.every(operation => operation.saved));
    },

    draggableOptions() {
      return {
        animation: VUETIFY_ANIMATION_DELAY,
        handle: '.step-drag-handler',
        ghostClass: 'white',
        group: {
          name: 'remediation-instruction-steps',
        },
      };
    },
  },
  methods: {
    addStep() {
      this.addItemIntoArray(generateRemediationInstructionStep());
    },

    changeStepsOrdering(event) {
      this.updateModel(dragDropChangePositionHandler(this.steps, event));
    },

    startDragging() {
      this.isDragging = true;
    },

    endDragging() {
      this.isDragging = false;
    },

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
