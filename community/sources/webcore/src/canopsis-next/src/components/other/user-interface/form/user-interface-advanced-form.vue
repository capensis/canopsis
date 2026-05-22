<template>
  <c-form-block>
    <c-form-block-row :label="$t('userInterface.infoPopupTimeout')">
      <c-duration-field
        v-field="form.popup_timeout.info"
        :label="$t('userInterface.infoPopupTimeout')"
        :units-label="$t('common.unit')"
        :disabled="disabled"
        :min="1"
        name="popup_timeout.info"
      />
    </c-form-block-row>

    <c-form-block-row :label="$t('userInterface.errorPopupTimeout')">
      <c-duration-field
        v-field="form.popup_timeout.error"
        :label="$t('userInterface.errorPopupTimeout')"
        :units-label="$t('common.unit')"
        :disabled="disabled"
        :min="1"
        name="popup_timeout.error"
      />
    </c-form-block-row>

    <c-form-block-row :label="$t('userInterface.matchedItems')">
      <v-layout class="gap-4">
        <v-flex xs6>
          <c-number-field
            v-field="form.max_matched_items"
            :label="$t('userInterface.maxMatchedItems')"
            :disabled="disabled"
            :min="1"
            name="max_matched_items"
          >
            <template #append>
              <c-help-icon
                :text="$t('userInterface.tooltips.maxMatchedItems')"
                color="grey darken-1"
                icon="help"
                left
              />
            </template>
          </c-number-field>
        </v-flex>
        <v-flex xs6>
          <c-number-field
            v-field="form.check_count_request_timeout"
            :label="$t('userInterface.checkCountRequestTimeout')"
            :disabled="disabled"
            :min="1"
            name="check_count_request_timeout"
          >
            <template #append>
              <c-help-icon
                :text="$t('userInterface.tooltips.checkCountRequestTimeout')"
                color="grey darken-1"
                icon="help"
                left
              />
            </template>
          </c-number-field>
        </v-flex>
      </v-layout>
    </c-form-block-row>

    <c-form-block-row
      :label="$t('settings.kioskMode')"
      align-center
    >
      <c-enabled-field
        v-field="form.show_header_on_kiosk_mode"
        :label="$t('userInterface.showHeaderOnKioskMode')"
        :disabled="disabled"
        hide-details
        no-margin
      >
        <template #append>
          <c-help-icon
            :text="$t('userInterface.showHeaderOnKioskModeTooltip')"
            color="grey darken-1"
            icon="help"
            top
          />
        </template>
      </c-enabled-field>
    </c-form-block-row>

    <c-form-block-row
      :label="$t('userInterface.animations')"
      align-center
    >
      <c-enabled-field
        v-field="form.disabled_transitions"
        :label="$t('userInterface.disabledTransitions')"
        :disabled="disabled"
        hide-details
        no-margin
      >
        <template #append>
          <c-help-icon
            :text="$t('userInterface.disabledTransitionsTooltip')"
            color="grey darken-1"
            icon="help"
            top
          />
        </template>
      </c-enabled-field>
    </c-form-block-row>

    <c-form-block-row
      :label="$t('common.state')"
      align-center
    >
      <c-enabled-field
        v-field="form.allow_change_severity_to_info"
        :label="$t('userInterface.allowChangeSeverityToInfo')"
        :disabled="disabled"
        hide-details
        no-margin
      />
    </c-form-block-row>

    <c-form-block-row
      :label="$t('userInterface.instructionApprovement')"
      align-center
    >
      <c-enabled-field
        v-field="form.required_instruction_approve"
        :label="$t('userInterface.requiredInstructionApprove')"
        :disabled="disabled"
        hide-details
        no-margin
      />
    </c-form-block-row>

    <c-form-block-row
      :label="$t('userInterface.autoSuggestPbehaviorName')"
      align-center
    >
      <c-enabled-field
        v-field="form.auto_suggest_pbehavior_name"
        :label="$t('userInterface.autoSuggestPbehaviorName')"
        :disabled="disabled"
        hide-details
        no-margin
      />
    </c-form-block-row>
  </c-form-block>
</template>

<script>
export default {
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
  },
};
</script>
