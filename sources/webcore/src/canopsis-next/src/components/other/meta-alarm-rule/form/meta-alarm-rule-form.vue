<template lang="pug">
  div
    v-layout(align-center)
      v-text-field(
        v-field="form._id",
        :label="$t('metaAlarmRule.id')",
        :disabled="isDisabledIdField",
        :readonly="isDisabledIdField"
      )
        v-tooltip(v-show="!isDisabledIdField", slot="append", left)
          v-icon(slot="activator") help
          span {{ $t('metaAlarmRule.idHelp') }}
    v-text-field(
      v-validate="'required'",
      v-field="form.name",
      :error-messages="errors.collect('name')",
      :label="$t('common.name')",
      name="name"
    )
    v-select(v-field="form.type", :items="ruleTypes", :label="$t('common.type')")
    v-text-field(
      v-if="isValueGroupType",
      v-field="form.config.value_path",
      v-validate="'required'",
      :label="$t('metaAlarmRule.fields.valuePath')",
      :error-messages="errors.collect('valuePath')",
      name="valuePath"
    )
    meta-alarm-rule-threshold-form(v-if="isThresholdFormShown", v-field="form.config")
    meta-alarm-rule-threshold-count-form(v-if="isThresholdCountFormShown", v-field="form.config")
    meta-alarm-rule-time-based-form(v-if="isTimeBasedFormShown", v-field="form.config")
    meta-alarm-rule-patterns-form(v-if="isPatternsFormShown", v-field="form.config")
</template>

<script>
import { META_ALARMS_RULE_TYPES } from '@/constants';

import MetaAlarmRuleThresholdForm from './meta-alarm-rule-threshold-form.vue';
import MetaAlarmRulePatternsForm from './meta-alarm-rule-patterns-form.vue';
import MetaAlarmRuleTimeBasedForm from './meta-alarm-rule-time-based-form.vue';
import MetaAlarmRuleThresholdCountForm from './meta-alarm-rule-threshold-count-form.vue';

export default {
  inject: ['$validator'],
  components: {
    MetaAlarmRuleThresholdCountForm,
    MetaAlarmRuleTimeBasedForm,
    MetaAlarmRulePatternsForm,
    MetaAlarmRuleThresholdForm,
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
    isDisabledIdField: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    ruleTypes() {
      return Object.values(META_ALARMS_RULE_TYPES);
    },

    /**
     * Conditions for forms showing
     */
    isThresholdFormShown() {
      return this.isComplexType;
    },

    isThresholdCountFormShown() {
      return this.isValueGroupType;
    },

    isTimeBasedFormShown() {
      return this.isComplexType || this.isValueGroupType || this.isTimeBasedType;
    },

    isPatternsFormShown() {
      return this.isComplexType || this.isValueGroupType || this.isPatternsType;
    },

    /**
     * Rule types
     */
    isPatternsType() {
      return this.form.type === META_ALARMS_RULE_TYPES.attribute;
    },

    isTimeBasedType() {
      return this.form.type === META_ALARMS_RULE_TYPES.timebased;
    },

    isComplexType() {
      return this.form.type === META_ALARMS_RULE_TYPES.complex;
    },

    isValueGroupType() {
      return this.form.type === META_ALARMS_RULE_TYPES.valuegroup;
    },
  },
};
</script>

