<template>
  <c-card-iterator-item
    :item-number="operationNumber"
    :removable="!disabled"
    drag-handle-class="operation-drag-handler"
    class="remediation-instruction-operation-field"
    default-expanded
    static-delete
    @remove="remove"
  >
    <template #header>
      <c-name-field
        :model-value="operation.name"
        :name="nameFieldName"
        :disabled="disabled"
        required
        @update:model-value="updateField('name', $event)"
      />
    </template>

    <v-row>
      <v-col cols="12">
        <c-duration-field
          :duration="operation.time_to_complete"
          :label="t('remediation.instruction.timeToComplete')"
          :units-label="t('common.unit')"
          :units="timeToCompleteUnits"
          :disabled="disabled"
          :name="timeToCompleteFieldName"
          required
          @update:duration="updateField('time_to_complete', $event)"
        />
      </v-col>
      <v-col v-if="disabled" cols="12">
        <text-editor-blurred
          :value="operation.description"
          :label="t('common.description')"
          :disabled="disabled"
          hide-details
        />
      </v-col>
      <v-col v-else cols="12">
        <text-editor-field
          :model-value="operation.description"
          :label="t('common.description')"
          :name="descriptionFieldName"
          :variables="templateVars.operation"
          required
          @update:model-value="updateField('description', $event)"
        />
      </v-col>
      <v-col v-if="disabled && operation.jobs && operation.jobs.length" cols="12">
        <jobs-chips :jobs="operation.jobs" />
      </v-col>
      <v-col v-if="!disabled" cols="12">
        <jobs-select
          :jobs="operation.jobs"
          @update:jobs="updateField('jobs', $event)"
        />
      </v-col>
    </v-row>
  </c-card-iterator-item>
</template>

<script>
import { computed, toRef } from 'vue';

import { SHORT_AVAILABLE_TIME_UNITS } from '@/constants';

import { isOmitEqual } from '@/helpers/collection';
import { remediationInstructionStepOperationToForm } from '@/helpers/entities/remediation/instruction/form';

import { useConfirmableForm } from '@/hooks/confirmable-form';
import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

import TextEditorBlurred from '@/components/common/text-editor/text-editor-blurred.vue';
import TextEditorField from '@/components/forms/fields/text-editor-field.vue';
import JobsChips from '@/components/other/remediation/instructions/partials/jobs-chips.vue';
import JobsSelect from '@/components/other/remediation/instructions/partials/jobs-select.vue';

export default {
  components: {
    TextEditorBlurred,
    TextEditorField,
    JobsChips,
    JobsSelect,
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
  emits: ['update:operation', 'remove'],
  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);

    const fieldName = computed(() => (props.operation.key ? props.operation.key : ''));
    const nameFieldName = computed(() => `${fieldName.value}.name`);
    const timeToCompleteFieldName = computed(() => `${fieldName.value}.timeToComplete`);
    const descriptionFieldName = computed(() => `${fieldName.value}.description`);
    const timeToCompleteUnits = Object.values(SHORT_AVAILABLE_TIME_UNITS);

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
      t,
      timeToCompleteUnits,
      nameFieldName,
      timeToCompleteFieldName,
      descriptionFieldName,
      updateField,
      remove,
    };
  },
};
</script>

<style lang="scss" scoped>
.remediation-instruction-operation-field {
  :deep(.c-card-iterator-item__row) {
    align-items: center;
  }

  .theme--light & {
    background-color: #f9f9f9;
  }

  .theme--dark & {
    background-color: #2d2d2d;
  }
}
</style>
