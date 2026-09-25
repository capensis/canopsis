<template>
  <div>
    <c-name-field
      v-model="nameField"
      :disabled="disabledCommon"
      required
      autofocus
    />

    <c-form-block>
      <c-form-block-row :label="t('remediation.instruction.type')">
        <c-instruction-type-field v-model="typeField" :disabled="disabled || !isNew" />
      </c-form-block-row>

      <c-form-block-row :label="t('common.description')" indented>
        <c-description-field v-model="descriptionField" :disabled="disabledCommon" required />
      </c-form-block-row>

      <c-form-block-row :label="t('remediation.instruction.timeoutAfterExecution')">
        <c-duration-field
          v-model:duration="timeoutAfterExecutionField"
          :label="t('remediation.instruction.timeoutAfterExecution')"
          :units-label="t('common.unit')"
          :disabled="disabled"
          :autofocus="disabledCommon"
          name="timeout_after_execution"
          required
        />
        <c-priority-field v-if="isAutoType" v-model="priorityField" :disabled="disabled" />
      </c-form-block-row>

      <c-form-block-row v-if="isAutoType" :label="t('remediation.instruction.retryEnabled')">
        <c-enabled-field
          :model-value="retryEnabledField"
          :label="t('remediation.instruction.retryEnabled')"
          :disabled="disabled"
          @update:model-value="updateRetryEnabled"
        />
        <v-expand-transition>
          <c-number-field
            v-if="retryEnabledField"
            v-model:value="retryCountField"
            :label="t('remediation.instruction.retryCount')"
            :min="1"
            :max="100"
            name="retry_count"
            required
          />
        </v-expand-transition>
      </c-form-block-row>

      <c-form-block-row v-if="isAutoType" :label="t('common.trigger', 2)">
        <c-triggers-field
          v-model="triggersField"
          :types="availableTriggers"
          required
          with-additional-values
        />

        <c-enabled-field
          v-model="enabledRepeatTriggersField"
          :label="t('remediation.instruction.enabledRepeatTrigger')"
        >
          <template #append>
            <c-help-icon
              :text="t('remediation.instruction.tooltips.enabledRepeatTriggerTooltip')"
              icon="help"
              color="grey darken-1"
              top
            />
          </template>
        </c-enabled-field>

        <v-expand-transition>
          <c-triggers-field
            v-if="enabledRepeatTriggersField"
            v-model="repeatTriggersField"
            :types="availableRepeatTriggers"
            :label="t('remediation.instruction.repeatTriggers')"
            name="repeat_triggers"
            translation-key-prefix="common.repeatTriggers"
            required
          />
        </v-expand-transition>
      </c-form-block-row>

      <c-form-block-row
        v-if="!disabledCommon"
        :label="t('remediation.instruction.requestApproval')"
        align-center
      >
        <remediation-instruction-approval-form
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

import { isInstructionTypeAuto } from '@/helpers/entities/remediation/instruction/form';

import { useFormContext, useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import RemediationInstructionApprovalForm from './remediation-instruction-approval-form.vue';

export default {
  components: {
    RemediationInstructionApprovalForm,
  },
  props: {
    form: {
      type: Object,
      required: false,
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
  setup(props) {
    const availableTriggers = REMEDIATION_AUTO_INSTRUCTION_TRIGGERS_TYPES;
    const availableRepeatTriggers = REMEDIATION_AUTO_INSTRUCTION_REPEAT_TRIGGERS_TYPES;

    const { t } = useI18n();
    const formContext = useFormContext();

    if (props.form) {
      formContext.reset(props.form);
    }

    const {
      nameField,
      typeField,
      descriptionField,
      timeoutAfterExecutionField,
      priorityField,
      retryEnabledField,
      retryCountField,
      triggersField,
      enabledRepeatTriggersField,
      repeatTriggersField,
    } = useFormModels({
      nameField: 'name',
      typeField: 'type',
      descriptionField: 'description',
      timeoutAfterExecutionField: 'timeout_after_execution',
      priorityField: 'priority',
      retryEnabledField: 'retry_enabled',
      retryCountField: 'retry_count',
      triggersField: 'triggers',
      enabledRepeatTriggersField: 'enabled_repeat_triggers',
      repeatTriggersField: 'repeat_triggers',
    });

    const isAutoType = computed(() => isInstructionTypeAuto(typeField.value));

    /**
     * Updates the retry enabled field and sets initial retry count when enabled
     *
     * @param {boolean} value - The new retry enabled state
     */
    const updateRetryEnabled = (value) => {
      retryEnabledField.value = value;

      if (value) {
        retryCountField.value = 1;
      }
    };

    return {
      t,
      availableTriggers,
      availableRepeatTriggers,

      isAutoType,

      nameField,
      typeField,
      descriptionField,
      timeoutAfterExecutionField,
      priorityField,
      retryEnabledField,
      retryCountField,
      triggersField,
      enabledRepeatTriggersField,
      repeatTriggersField,
      updateRetryEnabled,
    };
  },
};
</script>
