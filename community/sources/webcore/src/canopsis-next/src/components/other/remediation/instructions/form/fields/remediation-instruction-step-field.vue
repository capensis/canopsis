<template lang="pug">
  .step-field
    v-layout
      v-flex.mt-3(xs1)
        c-draggable-step-number(
          :disabled="disabled",
          :color="hasChildrenError ? 'error' : 'primary'",
          drag-class="step-drag-handler"
        ) {{ stepNumber }}
      v-flex(xs11)
        v-layout(row)
          v-layout
            c-expand-btn.step-expand(
              v-model="expanded",
              :color="!expanded && hasChildrenError ? 'error' : 'grey darken-3'"
            )
            v-layout(column)
              v-layout(row)
                v-flex(xs8)
                  v-text-field(
                    v-field="step.name",
                    v-validate="'required'",
                    :label="$t('common.name')",
                    :error-messages="errors.collect(nameFieldName)",
                    :name="nameFieldName",
                    :disabled="disabled",
                    box
                  )
                v-flex.pl-2(xs3)
                  v-text-field.step-time-complete-unit(
                    :value="timeToComplete | duration('refreshFieldFormat')",
                    :label="$t('remediation.instruction.timeToComplete')",
                    readonly
                  )
                v-flex.mt-1(xs1)
                  v-layout(justify-center)
                    c-action-btn(v-if="!disabled", type="delete", @click="$emit('remove')")
              v-expand-transition(mode="out-in")
                v-layout(v-if="expanded", column)
                  c-workflow-field(
                    v-field="step.stop_on_fail",
                    :label="$t('remediation.instruction.workflow')",
                    :continue-label="$t('remediation.instruction.remainingStep')",
                    :disabled="disabled"
                  )
                  remediation-instruction-step-endpoint-field(v-field="step.endpoint", :disabled="disabled")
                  remediation-instruction-operations-form(
                    v-field="step.operations",
                    :name="operationFieldName",
                    :step-number="stepNumber",
                    :disabled="disabled"
                  )
</template>

<script>
import { formMixin, validationChildrenMixin } from '@/mixins/form';

import { remediationInstructionStepToForm } from '@/helpers/forms/remediation-instruction';
import { isOmitEqual } from '@/helpers/equal';
import { toSeconds } from '@/helpers/date/duration';

import { confirmableFormMixinCreator } from '@/mixins/confirmable-form';

import RemediationInstructionOperationsForm from '../remediation-instruction-operations-form.vue';

import RemediationInstructionStepEndpointField from './remediation-instruction-step-endpoint-field.vue';

export default {
  inject: ['$validator'],
  components: {
    RemediationInstructionOperationsForm,
    RemediationInstructionStepEndpointField,
  },
  mixins: [
    formMixin,
    validationChildrenMixin,
    confirmableFormMixinCreator({
      field: 'step',
      method: 'remove',
      comparator(step) {
        const emptyStep = remediationInstructionStepToForm();
        const paths = [
          'key',
          step.operations.length ? ['operations', 0, 'key'] : 'operations',
        ];

        return isOmitEqual(step, emptyStep, paths);
      },
    }),
  ],
  model: {
    prop: 'step',
    event: 'input',
  },
  props: {
    step: {
      type: Object,
      required: true,
    },
    stepNumber: {
      type: [Number, String],
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      expanded: true,
    };
  },
  computed: {
    fieldSuffix() {
      return this.step.key ? `-${this.step.key}` : '';
    },

    nameFieldName() {
      return `name${this.fieldSuffix}`;
    },

    operationFieldName() {
      return `operations${this.fieldSuffix}`;
    },

    timeToComplete() {
      return this.step.operations.reduce((acc, operation) => {
        const { time_to_complete: { value, unit } } = operation;

        return acc + toSeconds(value, unit);
      }, 0);
    },
  },
  methods: {
    remove() {
      this.$emit('remove');
    },
  },
};
</script>

<style lang="scss" scoped>
  .step-field {
    & /deep/ .step-expand {
      margin: 24px 2px 0 2px !important;
      width: 20px !important;
      height: 20px !important;
    }

    & /deep/ .step-time-complete-unit .v-input__slot {
      &:before, &:after {
        content: none !important;
      }
    }
  }
</style>
