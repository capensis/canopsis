<template lang="pug">
  v-layout
    v-flex.mt-3(xs1)
      c-draggable-step-number(
        drag-class="operation-drag-handler",
        :disabled="disabled",
        :color="hasChildrenError ? 'error' : 'primary'"
      ) {{ operationNumber }}
    v-flex(xs11)
      v-layout(row, justify-space-between)
        v-flex(xs11)
          v-layout(row)
            c-expand-btn.operation-expand(
              v-if="!disabled",
              v-model="expanded",
              :color="!expanded && hasChildrenError ? 'error' : ''"
            )
            v-layout(column)
              v-text-field(
                v-field="operation.name",
                v-validate="'required'",
                :label="$t('common.name')",
                :error-messages="errors.collect(nameFieldName)",
                :name="nameFieldName",
                :disabled="disabled",
                box
              )
              v-expand-transition(mode="out-in")
                v-layout(v-if="expanded", column)
                  remediation-instruction-time-to-complete-field(
                    v-field="operation.time_to_complete",
                    :disabled="disabled",
                    :name="timeToCompleteFieldName"
                  )
                  text-editor-blurred(
                    v-if="disabled",
                    :value="operation.description",
                    :label="$t('common.description')",
                    :disabled="disabled",
                    hide-details
                  )
                  text-editor-field(
                    v-else,
                    v-field="operation.description",
                    v-validate="'required'",
                    :label="$t('common.description')",
                    :error-messages="errors.collect(descriptionFieldName)",
                    :name="descriptionFieldName"
                  )
                  jobs-chips(
                    v-if="disabled &&  operation.jobs && operation.jobs.length",
                    :jobs="operation.jobs"
                  )
                  jobs-select(v-if="!disabled", v-field="operation.jobs")
        span
          c-action-btn.mt-1(v-if="!disabled", type="delete", @click="$emit('remove')")
</template>

<script>
import { isOmitEqual } from '@/helpers/equal';

import { remediationInstructionStepOperationToForm } from '@/helpers/forms/remediation-instruction';

import { formMixin, validationChildrenMixin } from '@/mixins/form';
import { confirmableFormMixinCreator } from '@/mixins/confirmable-form';

import TextEditorField from '@/components/forms/fields/text-editor-field.vue';
import JobsChips from '@/components/other/remediation/instructions/partials/jobs-chips.vue';
import JobsSelect from '@/components/other/remediation/instructions/partials/jobs-select.vue';
import TextEditorBlurred from '@/components/common/text-editor/text-editor-blurred.vue';

import RemediationInstructionTimeToCompleteField from './remediation-instruction-time-to-complete-field.vue';

export default {
  inject: ['$validator'],
  components: {
    TextEditorBlurred,
    RemediationInstructionTimeToCompleteField,
    TextEditorField,
    JobsChips,
    JobsSelect,
  },
  mixins: [
    formMixin,
    validationChildrenMixin,
    confirmableFormMixinCreator({
      field: 'operation',
      method: 'remove',
      comparator(operation) {
        const emptyOperation = remediationInstructionStepOperationToForm();
        const paths = ['key'];

        return isOmitEqual(operation, emptyOperation, paths);
      },
    }),
  ],
  model: {
    prop: 'operation',
    event: 'input',
  },
  props: {
    operation: {
      type: Object,
      default: () => ({}),
    },
    operationNumber: {
      type: [Number, String],
      default: 0,
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
    fieldName() {
      return this.operation.key ? this.operation.key : '';
    },

    nameFieldName() {
      return `${this.fieldName}.name`;
    },

    timeToCompleteFieldName() {
      return `${this.fieldName}.timeToComplete`;
    },

    descriptionFieldName() {
      return `${this.fieldName}.description`;
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
  .operation-expand {
    margin: 24px 2px 0 2px !important;
    width: 20px;
    height: 20px;
  }
</style>
