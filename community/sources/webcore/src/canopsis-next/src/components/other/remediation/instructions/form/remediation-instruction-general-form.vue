<template>
  <div>
    <c-name-field
      v-field="form.name"
      :disabled="disabledCommon"
      required
      autofocus
    />

    <c-form-block>
      <c-form-block-row :label="$t('remediation.instruction.type')" indented>
        <c-instruction-type-field
          v-field="form.type"
          :disabled="disabled || !isNew"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.description')">
        <c-description-field
          v-field="form.description"
          :disabled="disabledCommon"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('remediation.instruction.timeoutAfterExecution')">
        <c-duration-field
          v-field="form.timeout_after_execution"
          :label="$t('remediation.instruction.timeoutAfterExecution')"
          :units-label="$t('common.unit')"
          :disabled="disabled"
          :autofocus="disabledCommon"
          name="timeout_after_execution"
          required
        />
        <c-priority-field
          v-if="isAutoType"
          v-field="form.priority"
          :disabled="disabled"
        />
      </c-form-block-row>

      <c-form-block-row v-if="isAutoType" :label="$t('remediation.instruction.retryEnabled')">
        <c-enabled-field
          :value="form.retry_enabled"
          :label="$t('remediation.instruction.retryEnabled')"
          :disabled="disabled"
          @input="updateRetryEnabled"
        />
        <v-expand-transition>
          <c-number-field
            v-if="form.retry_enabled"
            v-field="form.retry_count"
            :label="$t('remediation.instruction.retryCount')"
            :min="1"
            :max="100"
            name="retry_count"
            required
          />
        </v-expand-transition>
      </c-form-block-row>

      <c-form-block-row v-if="isAutoType" :label="$tc('common.trigger', 2)">
        <c-triggers-field
          v-field="form.triggers"
          :types="availableTriggers"
          with-additional-values
        />

        <c-enabled-field
          v-field="form.enabled_repeat_triggers"
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
            v-field="form.repeat_triggers"
            :types="availableRepeatTriggers"
            :label="$t('remediation.instruction.repeatTriggers')"
            name="repeat_triggers"
            translation-key-prefix="common.repeatTriggers"
          />
        </v-expand-transition>
      </c-form-block-row>

      <c-form-block-row v-if="!disabledCommon" :label="$t('remediation.instruction.requestApproval')">
        <remediation-instruction-approval-form
          v-field="form.approval"
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

import { isInstructionTypeAuto, isInstructionTypeSimpleManual } from '@/helpers/entities/remediation/instruction/form';

import { useModelField } from '@/hooks/form/model-field';

import RemediationInstructionApprovalForm from './remediation-instruction-approval-form.vue';

export default {
  inject: ['$validator'],
  components: {
    RemediationInstructionApprovalForm,
  },
  model: {
    prop: 'form',
    event: 'input',
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
  setup(props, { emit }) {
    const availableTriggers = REMEDIATION_AUTO_INSTRUCTION_TRIGGERS_TYPES;
    const availableRepeatTriggers = REMEDIATION_AUTO_INSTRUCTION_REPEAT_TRIGGERS_TYPES;

    const { updateModel } = useModelField(props, emit);

    const isAutoType = computed(() => isInstructionTypeAuto(props.form?.type));
    const isManualSimplified = computed(() => isInstructionTypeSimpleManual(props.form?.type));

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

      updateRetryEnabled,
    };
  },
};
</script>
