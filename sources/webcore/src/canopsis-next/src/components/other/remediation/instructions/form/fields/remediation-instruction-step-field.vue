<template lang="pug">
  v-layout(row)
    v-flex(xs10)
      v-layout
        v-text-field(
          v-field="step.name",
          v-validate="'required'",
          :label="$t('common.name')",
          :disabled="step.saved",
          :error-messages="nameErrorMessages",
          :name="name",
          box,
          @keyup.stop.enter="saveName"
        )
      v-layout(v-if="timeToComplete > 0")
        v-text-field(
          :value="timeToComplete | duration(undefined,  'refreshFieldFormat')",
          :label="$t('remediationInstructions.timeToComplete')",
          disabled,
          box
        )
      v-layout
        remediation-instruction-steps-workflow-field(v-field="step.workflow", :disabled="step.saved")
      v-layout(v-if="!step.saved", justify-end)
        v-btn.mt-0(depressed, flat, @click="cancelChangeName") {{ $t('common.cancel') }}
        v-btn.mt-0.mr-0.primary(@click="saveName") {{ $t('common.save') }}
    v-flex.mt-3(v-if="step.saved && !hideActions", xs2)
      v-layout(justify-start)
        v-btn.ma-0.ml-2(icon, small, @click="editName")
          v-icon edit
        v-btn.ma-0.ml-1(icon, small, @click.prevent="$emit('remove')")
          v-icon(color="error") delete
</template>

<script>
import formMixin from '@/mixins/form';

import { getUnitValueFromOtherUnit } from '@/helpers/time';

import RemediationInstructionStepsWorkflowField from './remediation-instruction-steps-workflow-field.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: { RemediationInstructionStepsWorkflowField },
  mixins: [formMixin],
  model: {
    prop: 'step',
    event: 'input',
  },
  props: {
    step: {
      type: Object,
      required: true,
    },
    hideActions: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      oldStep: null,
    };
  },
  computed: {
    fieldSuffix() {
      return this.step.key ? `-${this.step.key}` : '';
    },

    name() {
      return `name${this.fieldSuffix}`;
    },

    nameErrorMessages() {
      return this.errors.collect(this.name).map(error => error.replace(this.fieldSuffix, ''));
    },

    timeToComplete() {
      return this.step.operations.reduce((acc, operation) => {
        const { time_to_complete: { interval, unit } } = operation;

        return acc + getUnitValueFromOtherUnit(interval, unit);
      }, 0);
    },
  },
  methods: {
    editName() {
      this.oldStep = this.step;

      this.updateField('saved', false);
    },

    cancelChangeName() {
      if (this.oldStep) {
        this.updateModel({
          ...this.oldStep,
          saved: true,
        });
      } else {
        this.$emit('remove');
      }
    },

    async saveName() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        this.oldStep = null;

        this.updateField('saved', true);
      }
    },
  },
};
</script>
