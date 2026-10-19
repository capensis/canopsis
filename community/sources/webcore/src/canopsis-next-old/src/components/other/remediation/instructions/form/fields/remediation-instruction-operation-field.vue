<template>
  <c-form-block>
    <v-layout class="remediation-instruction-operation-field pa-3" align-start>
      <c-draggable-step-number
        :disabled="disabled"
        :color="hasChildrenError ? 'error' : 'primary'"
        class="operation-number mt-4"
        drag-class="operation-drag-handler"
      >
        {{ operationNumber }}
      </c-draggable-step-number>
      <v-flex xs11>
        <v-layout>
          <c-expand-btn
            v-if="!disabled"
            v-model="expanded"
            :color="!expanded && hasChildrenError ? 'error' : ''"
            class="operation-expand"
          />
          <v-layout column>
            <v-text-field
              v-field="operation.name"
              v-validate="'required'"
              :label="$t('common.name')"
              :error-messages="errors.collect(nameFieldName)"
              :name="nameFieldName"
              :disabled="disabled"
              filled
            />
            <v-expand-transition mode="out-in">
              <v-layout
                v-if="expanded"
                column
              >
                <remediation-instruction-time-to-complete-field
                  v-field="operation.time_to_complete"
                  :disabled="disabled"
                  :name="timeToCompleteFieldName"
                />
                <text-editor-blurred
                  v-if="disabled"
                  :value="operation.description"
                  :label="$t('common.description')"
                  :disabled="disabled"
                  hide-details
                />
                <text-editor-field
                  v-else
                  v-field="operation.description"
                  v-validate="'required'"
                  :label="$t('common.description')"
                  :error-messages="errors.collect(descriptionFieldName)"
                  :name="descriptionFieldName"
                  :variables="templateVars.operation"
                />
                <jobs-chips
                  v-if="disabled && operation.jobs && operation.jobs.length"
                  :jobs="operation.jobs"
                />
                <jobs-select
                  v-if="!disabled"
                  v-field="operation.jobs"
                />
              </v-layout>
            </v-expand-transition>
          </v-layout>
        </v-layout>
      </v-flex>
      <c-action-btn
        v-if="!disabled"
        class="mt-2"
        type="delete"
        @click="remove"
      />
    </v-layout>
  </c-form-block>
</template>

<script>
import { computed, ref, toRef } from 'vue';

import { isOmitEqual } from '@/helpers/collection';
import { remediationInstructionStepOperationToForm } from '@/helpers/entities/remediation/instruction/form';

import { useConfirmableForm } from '@/hooks/confirmable-form';
import { useValidationChildren } from '@/hooks/validator/validation-children';

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
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const expanded = ref(true);
    const { hasChildrenError } = useValidationChildren();

    const fieldName = computed(() => (props.operation.key ? props.operation.key : ''));
    const nameFieldName = computed(() => `${fieldName.value}.name`);
    const timeToCompleteFieldName = computed(() => `${fieldName.value}.timeToComplete`);
    const descriptionFieldName = computed(() => `${fieldName.value}.description`);

    const { confirmAction: remove } = useConfirmableForm({
      form: toRef(props, 'operation'),
      action: () => emit('remove'),
      comparator: (operation) => {
        const emptyOperation = remediationInstructionStepOperationToForm();
        const paths = ['key'];

        return isOmitEqual(operation, emptyOperation, paths);
      },
    });

    return {
      expanded,
      hasChildrenError,
      nameFieldName,
      timeToCompleteFieldName,
      descriptionFieldName,
      remove,
    };
  },
};
</script>

<style lang="scss" scoped>
  .remediation-instruction-operation-field {
    .theme--light & {
      background-color: #F9F9F9;
    }
    .theme--dark & {
      background-color: #2D2D2D;
    }
  }

  .operation-number {
    min-width: 70px;
  }

  .operation-expand {
    margin: 24px 8px 0 2px !important;
    width: 20px;
    height: 20px;
  }
</style>
