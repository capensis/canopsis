<template lang="pug">
  v-layout(column)
    v-flex(v-if="!steps.length", xs12)
      v-alert(:value="true", type="info") {{ $t('remediationInstructions.emptySteps') }}
    draggable(v-field="steps", :options="draggableOptions")
      v-card.my-2(v-for="(step, index) in steps", :key="step.key")
        v-card-text
          remediation-instruction-step-field(
            v-field="steps[index]",
            :step-number="index + 1",
            @remove="removeStep(index)"
          )
    v-layout(row, align-center)
      v-btn.ml-0(
        :color="hasStepsErrors ? 'error' : 'primary'",
        outline,
        @click="addStep"
      ) {{ $t('remediationInstructions.addStep') }}
      span.error--text(v-show="hasStepsErrors") {{ $t('remediationInstructions.errors.stepRequired') }}
</template>

<script>
import Draggable from 'vuedraggable';

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
    name: {
      type: String,
      default: 'steps',
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
      context: () => this,
      vm: this,
    });
  },
  beforeDestroy() {
    this.$validator.detach(this.name);
  },
  methods: {
    addStep() {
      this.addItemIntoArray(generateRemediationInstructionStep());
    },

    removeStep(index) {
      this.removeItemFromArray(index);
    },
  },
};
</script>
