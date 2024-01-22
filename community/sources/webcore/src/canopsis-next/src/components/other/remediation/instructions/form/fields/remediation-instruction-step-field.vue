<template>
  <c-card-iterator-item
    class="remediation-instruction-step-field"
    :item-number="stepNumber"
    offset-left
    @remove="remove"
  >
    <template #header="">
      <v-layout class="mt-3">
        <v-flex xs9>
          <v-text-field
            v-field="step.name"
            v-validate="'required'"
            :label="$t('common.name')"
            :error-messages="errors.collect(nameFieldName)"
            :name="nameFieldName"
            :disabled="disabled"
            filled
          />
        </v-flex>
        <v-flex
          class="pl-2"
          xs3
        >
          <v-text-field
            class="remediation-instruction-step-field__time-to-complete"
            :value="timeToComplete | duration('refreshFieldFormat')"
            :label="$t('remediation.instruction.timeToComplete')"
            readonly
          />
        </v-flex>
      </v-layout>
    </template>
    <c-workflow-field
      v-field="step.stop_on_fail"
      :label="$t('remediation.instruction.workflow')"
      :continue-label="$t('remediation.instruction.remainingStep')"
      :disabled="disabled"
    />
    <remediation-instruction-step-endpoint-field
      v-field="step.endpoint"
      :disabled="disabled"
    />
    <remediation-instruction-operations-form
      v-field="step.operations"
      :name="operationFieldName"
      :step-number="stepNumber"
      :disabled="disabled"
    />
  </c-card-iterator-item>
</template>

<script>
import { remediationInstructionStepToForm } from '@/helpers/entities/remediation/instruction/form';
import { isOmitEqual } from '@/helpers/collection';
import { toSeconds } from '@/helpers/date/duration';

import { formMixin, validationChildrenMixin } from '@/mixins/form';
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

<style lang="scss">
.remediation-instruction-step-field {
  &__time-to-complete .v-input__slot {
    &:before, &:after {
      content: none !important;
    }
  }
}
</style>
