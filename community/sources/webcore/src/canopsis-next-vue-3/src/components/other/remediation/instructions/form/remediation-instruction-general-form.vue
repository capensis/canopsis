<template>
  <div>
    <c-name-field
      v-model="nameField"
      :disabled="disabledCommon"
      required
      autofocus
    />

    <c-form-block>
      <c-form-block-row :label="$t('remediation.instruction.type')">
        <c-instruction-type-field v-model="typeField" :disabled="disabled || !isNew" />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.description')">
        <c-description-field v-model="descriptionField" :disabled="disabledCommon" required />
      </c-form-block-row>

      <c-form-block-row :label="$t('remediation.instruction.timeoutAfterExecution')">
        <c-duration-field
          v-model:duration="timeoutAfterExecutionField"
          :label="$t('remediation.instruction.timeoutAfterExecution')"
          :units-label="$t('common.unit')"
          :disabled="disabled"
          :autofocus="disabledCommon"
          name="timeout_after_execution"
          required
        />
        <c-priority-field v-if="isAutoType" v-model="priorityField" :disabled="disabled" />
      </c-form-block-row>

      <c-form-block-row v-if="isAutoType" :label="$t('remediation.instruction.retryEnabled')">
        <c-enabled-field
          :model-value="form.retry_enabled"
          :label="$t('remediation.instruction.retryEnabled')"
          :disabled="disabled"
          @update:model-value="updateRetryEnabled"
        />
        <v-expand-transition>
          <c-number-field
            v-if="form.retry_enabled"
            v-model:value="retryCountField"
            :label="$t('remediation.instruction.retryCount')"
            :min="1"
            :max="100"
            name="retry_count"
            required
          />
        </v-expand-transition>
      </c-form-block-row>

      <c-form-block-row v-if="isAutoType" :label="$t('common.trigger', 2)">
        <c-triggers-field
          v-model="triggersField"
          :types="availableTriggers"
          with-additional-values
        />

        <c-enabled-field
          v-model="enabledRepeatTriggersField"
          :label="$t('remediation.instruction.enabledRepeatTrigger')"
        >
          <template #append>
            <c-help-icon
              :text="$t('remediation.instruction.tooltips.enabledRepeatTriggerTooltip')"
              icon="help"
              color="grey darken-1"
              top
            />
          </template>
        </c-enabled-field>

        <v-expand-transition>
          <c-triggers-field
            v-if="form.enabled_repeat_triggers"
            v-model="repeatTriggersField"
            :types="availableRepeatTriggers"
            :label="$t('remediation.instruction.repeatTriggers')"
            name="repeat_triggers"
            translation-key-prefix="common.repeatTriggers"
          />
        </v-expand-transition>
      </c-form-block-row>

      <c-form-block-row
        v-if="!disabledCommon"
        :label="$t('remediation.instruction.requestApproval')"
      >
        <remediation-instruction-approval-form
          v-model:approval="approvalField"
          :disabled="disabled"
          :required="requiredApprove"
        />
      </c-form-block-row>
    </c-form-block>
  </div>
</template>

<script>
import { computed } from 'vue';

import {
  REMEDIATION_AUTO_INSTRUCTION_TRIGGERS_TYPES,
  REMEDIATION_AUTO_INSTRUCTION_REPEAT_TRIGGERS_TYPES,
} from '@/constants';

import {
  isInstructionTypeAuto,
  isInstructionTypeSimpleManual,
} from '@/helpers/entities/remediation/instruction/form';

import { useModelField } from '@/hooks/form/model-field';

import RemediationInstructionApprovalForm from './remediation-instruction-approval-form.vue';

export default {
  inject: ['$validator'],
  components: {
    RemediationInstructionApprovalForm,
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    disabledCommon: {
      type: Boolean,
      default: false,
    },
    isNew: {
      type: Boolean,
      default: false,
    },
    requiredApprove: {
      type: Boolean,
      default: false,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
  emits: ['update:form'],
  setup(props, { emit }) {
    const availableTriggers = REMEDIATION_AUTO_INSTRUCTION_TRIGGERS_TYPES;
    const availableRepeatTriggers = REMEDIATION_AUTO_INSTRUCTION_REPEAT_TRIGGERS_TYPES;

    const { updateModel } = useModelField(props, emit);
    // TODO: fix that
    const createFieldModel = field => computed({
        get: () => props.form[field],
        set: value => updateModel({
            ...props.form,
            [field]: value,
          }),
      });

    const isAutoType = computed(() => isInstructionTypeAuto(props.form?.type));
    const isManualSimplified = computed(() => isInstructionTypeSimpleManual(props.form?.type));
    const nameField = createFieldModel('name');
    const typeField = createFieldModel('type');
    const descriptionField = createFieldModel('description');
    const timeoutAfterExecutionField = createFieldModel('timeout_after_execution');
    const priorityField = createFieldModel('priority');
    const retryCountField = createFieldModel('retry_count');
    const triggersField = createFieldModel('triggers');
    const enabledRepeatTriggersField = createFieldModel('enabled_repeat_triggers');
    const repeatTriggersField = createFieldModel('repeat_triggers');
    const approvalField = createFieldModel('approval');

    /**
     * Updates the retry enabled field and sets initial retry count when enabled
     *
     * @param {boolean} value - The new retry enabled state
     */
    const updateRetryEnabled = (value) => {
      const newForm = { ...props.form, retry_enabled: value };

      if (value) {
        newForm.retry_count = 1;
      }

      updateModel(newForm);
    };

    return {
      availableTriggers,
      availableRepeatTriggers,

      isAutoType,
      isManualSimplified,

      nameField,
      typeField,
      descriptionField,
      timeoutAfterExecutionField,
      priorityField,
      retryCountField,
      triggersField,
      enabledRepeatTriggersField,
      repeatTriggersField,
      approvalField,
      updateRetryEnabled,
    };
  },
};
</script>
