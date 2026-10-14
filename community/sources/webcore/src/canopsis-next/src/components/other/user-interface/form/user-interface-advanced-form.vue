<template>
  <c-form-block>
    <c-form-block-row :label="t('userInterface.infoPopupTimeout')">
      <c-duration-field
        v-model:duration="popupTimeoutInfoField"
        :label="t('userInterface.infoPopupTimeout')"
        :units-label="t('common.unit')"
        :disabled="disabled"
        :min="1"
        name="popup_timeout.info"
        same-width
      />
    </c-form-block-row>

    <c-form-block-row :label="t('userInterface.errorPopupTimeout')">
      <c-duration-field
        v-model:duration="popupTimeoutErrorField"
        :label="t('userInterface.errorPopupTimeout')"
        :units-label="t('common.unit')"
        :disabled="disabled"
        :min="1"
        name="popup_timeout.error"
        same-width
      />
    </c-form-block-row>

    <c-form-block-row :label="t('userInterface.matchedItems')">
      <v-row density="compact">
        <v-col cols="6">
          <c-number-field
            v-model:value="maxMatchedItemsField"
            :label="t('userInterface.maxMatchedItems')"
            :disabled="disabled"
            :min="1"
            name="max_matched_items"
          >
            <template #append>
              <c-help-icon
                :text="t('userInterface.tooltips.maxMatchedItems')"
                color="grey-darken-1"
                icon="help"
                left
              />
            </template>
          </c-number-field>
        </v-col>
        <v-col cols="6">
          <c-number-field
            v-model:value="checkCountRequestTimeoutField"
            :label="t('userInterface.checkCountRequestTimeout')"
            :disabled="disabled"
            :min="1"
            name="check_count_request_timeout"
          >
            <template #append>
              <c-help-icon
                :text="t('userInterface.tooltips.checkCountRequestTimeout')"
                color="grey-darken-1"
                icon="help"
                left
              />
            </template>
          </c-number-field>
        </v-col>
      </v-row>
    </c-form-block-row>

    <c-form-block-row :label="t('settings.kioskMode')" align-center>
      <c-enabled-field
        v-model="showHeaderOnKioskModeField"
        :label="t('userInterface.showHeaderOnKioskMode')"
        :disabled="disabled"
        hide-details
        no-margin
      >
        <template #append>
          <c-help-icon
            :text="t('userInterface.showHeaderOnKioskModeTooltip')"
            color="grey-darken-1"
            icon="help"
            top
          />
        </template>
      </c-enabled-field>
    </c-form-block-row>

    <c-form-block-row :label="t('userInterface.animations')" align-center>
      <c-enabled-field
        v-model="disabledTransitionsField"
        :label="t('userInterface.disabledTransitions')"
        :disabled="disabled"
        hide-details
        no-margin
      >
        <template #append>
          <c-help-icon
            :text="t('userInterface.disabledTransitionsTooltip')"
            color="grey-darken-1"
            icon="help"
            top
          />
        </template>
      </c-enabled-field>
    </c-form-block-row>

    <c-form-block-row :label="t('common.state')" align-center>
      <c-enabled-field
        v-model="allowChangeSeverityToInfoField"
        :label="t('userInterface.allowChangeSeverityToInfo')"
        :disabled="disabled"
        hide-details
        no-margin
      />
    </c-form-block-row>

    <c-form-block-row :label="t('userInterface.instructionApprovement')" align-center>
      <c-enabled-field
        v-model="requiredInstructionApproveField"
        :label="t('userInterface.requiredInstructionApprove')"
        :disabled="disabled"
        hide-details
        no-margin
      />
    </c-form-block-row>

    <c-form-block-row :label="t('userInterface.autoSuggestPbehaviorName')" align-center>
      <c-enabled-field
        v-model="autoSuggestPbehaviorNameField"
        :label="t('userInterface.autoSuggestPbehaviorName')"
        :disabled="disabled"
        hide-details
        no-margin
      />
    </c-form-block-row>
  </c-form-block>
</template>

<script>
import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { t } = useI18n();
    const {
      popupTimeoutInfoField,
      popupTimeoutErrorField,
      maxMatchedItemsField,
      checkCountRequestTimeoutField,
      showHeaderOnKioskModeField,
      disabledTransitionsField,
      allowChangeSeverityToInfoField,
      requiredInstructionApproveField,
      autoSuggestPbehaviorNameField,
    } = useFormModels({
      popupTimeoutInfoField: 'popup_timeout.info',
      popupTimeoutErrorField: 'popup_timeout.error',
      maxMatchedItemsField: 'max_matched_items',
      checkCountRequestTimeoutField: 'check_count_request_timeout',
      showHeaderOnKioskModeField: 'show_header_on_kiosk_mode',
      disabledTransitionsField: 'disabled_transitions',
      allowChangeSeverityToInfoField: 'allow_change_severity_to_info',
      requiredInstructionApproveField: 'required_instruction_approve',
      autoSuggestPbehaviorNameField: 'auto_suggest_pbehavior_name',
    });

    return {
      t,
      popupTimeoutInfoField,
      popupTimeoutErrorField,
      maxMatchedItemsField,
      checkCountRequestTimeoutField,
      showHeaderOnKioskModeField,
      disabledTransitionsField,
      allowChangeSeverityToInfoField,
      requiredInstructionApproveField,
      autoSuggestPbehaviorNameField,
    };
  },
};
</script>
