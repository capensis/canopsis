<template lang="pug">
  v-layout(column)
    draggable(v-field="steps", :options="draggableOptions")
      v-card.my-2(v-for="(step, index) in steps", :key="step.key")
        v-card-text
          v-layout(row, wrap)
            v-flex.mt-3(xs1)
              draggable-step-number(drag-class="step-drag-handler") {{ index + 1 }}
            v-flex(xs11)
              remediation-instruction-step-field(
                v-field="steps[index]",
                :index="index",
                @remove="removeStep(index)"
              )
    v-layout
      v-btn.ml-0(
        outline,
        color="primary",
        @click="addStep"
      ) {{ $t('remediationInstructions.addStep') }}
</template>

<script>
import Draggable from 'vuedraggable';

import { MODALS } from '@/constants';
import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { generateRemediationInstructionStep } from '@/helpers/entities';

import formArrayMixin from '@/mixins/form/array';

import DraggableStepNumber from '../partials/draggable-step-number.vue';

import RemediationInstructionStepField from './fields/remediation-instruction-step-field.vue';

export default {
  components: {
    Draggable,
    DraggableStepNumber,
    RemediationInstructionStepField,
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
