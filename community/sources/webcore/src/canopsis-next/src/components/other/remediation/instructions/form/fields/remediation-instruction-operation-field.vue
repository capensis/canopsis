<template lang="pug">
  v-layout
    v-flex.mt-3(xs1)
      c-draggable-step-number(
        drag-class="operation-drag-handler",
        :color="hasChildrenError ? 'error' : 'primary'"
      ) {{ operationNumber }}
    v-flex(xs11)
      v-layout(row)
        v-flex.pr-1(xs11)
          v-layout(row)
            c-expand-btn.operation-expand(
              v-model="expanded",
              :color="!expanded && hasChildrenError ? 'error' : 'grey darken-3'"
            )
            v-layout(column)
              v-text-field(
                v-field="operation.name",
                v-validate="'required'",
                :label="$t('common.name')",
                :error-messages="errors.collect(nameFieldName)",
                :name="nameFieldName",
                box
              )
              v-expand-transition(mode="out-in")
                v-layout(v-if="expanded", column)
                  remediation-instruction-time-to-complete-field(
                    v-field="operation.time_to_complete",
                    :name="timeToCompleteFieldName"
                  )
                  text-editor-field(
                    v-field="operation.description",
                    v-validate="'required'",
                    :label="$t('common.description')",
                    :error-messages="errors.collect(descriptionFieldName)",
                    :name="descriptionFieldName"
                  )
                  jobs-select(v-field="operation.jobs")
        v-flex.mt-3(xs1)
          v-layout(justify-center)
            v-btn.ma-0(icon, small, @click.prevent="remove")
              v-icon(color="error") delete
</template>

<script>
import { isOmitEqual } from '@/helpers/validators/is-omit-equal';

import { remediationInstructionStepOperationToForm } from '@/helpers/forms/remediation-instruction';

import formMixin from '@/mixins/form';
import validationChildrenMixin from '@/mixins/form/validation-children';
import confirmableFormMixin from '@/mixins/confirmable-form';

import TextEditorField from '@/components/forms/fields/text-editor-field.vue';
import JobsSelect from '@/components/other/remediation/instructions/partials/jobs-select.vue';

import RemediationInstructionTimeToCompleteField from './remediation-instruction-time-to-complete-field.vue';

export default {
  inject: ['$validator'],
  components: {
    RemediationInstructionTimeToCompleteField,
    TextEditorField,
    JobsSelect,
  },
  mixins: [
    formMixin,
    validationChildrenMixin,
    confirmableFormMixin({
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
