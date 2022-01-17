<template lang="pug">
  div
    v-layout(align-center)
      v-text-field(
        v-validate,
        v-field="form._id",
        :label="$t('common.id')",
        :disabled="isDisabledIdField",
        :readonly="isDisabledIdField",
        :error-messages="errors.collect('_id')",
        name="_id"
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
    v-textarea(
      v-field="form.output_template",
      :label="$t('metaAlarmRule.outputTemplate')"
    )
      v-tooltip(slot="append", left)
        v-icon(slot="activator") help
        div(v-html="$t('metaAlarmRule.outputTemplateHelp')")
    v-switch(
      v-field="form.auto_resolve",
      :label="$t('metaAlarmRule.autoResolve')",
      color="primary"
    )
    v-select(v-field="form.type", :items="ruleTypes", :label="$t('common.type')")
    meta-alarm-rule-corel-form(v-if="isCorelFormShown", v-field="form.config")
    meta-alarm-rule-threshold-form(v-if="isThresholdFormShown", v-field="form.config")
    meta-alarm-rule-time-based-form(v-if="isTimeBasedFormShown", v-field="form.config")
    meta-alarm-rule-value-paths-form(v-if="isValuePathsFormShown", v-field="form.config")
    c-patterns-field(
      v-if="isPatternsFormShown",
      v-field="form.config",
      :total-entity="withTotalEntityPatterns",
      name="config",
      alarm,
      entity,
      event
    )
</template>

<script>
import { META_ALARMS_RULE_TYPES } from '@/constants';

import MetaAlarmRuleThresholdForm from './meta-alarm-rule-threshold-form.vue';
import MetaAlarmRuleTimeBasedForm from './meta-alarm-rule-time-based-form.vue';
import MetaAlarmRuleValuePathsForm from './meta-alarm-rule-value-paths-form.vue';
import MetaAlarmRuleCorelForm from './meta-alarm-rule-corel-form.vue';

export default {
  inject: ['$validator'],
  components: {
    MetaAlarmRuleTimeBasedForm,
    MetaAlarmRuleThresholdForm,
    MetaAlarmRuleValuePathsForm,
    MetaAlarmRuleCorelForm,
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
    /**
     * We are filtered 'manualgroup' because we are using in only in the alarms list widget directly
     */
    ruleTypes() {
      return Object.values(META_ALARMS_RULE_TYPES)
        .filter(type => type !== META_ALARMS_RULE_TYPES.manualgroup);
    },

    /**
     * Conditions for forms showing
     */
    isThresholdFormShown() {
      return this.isComplexType || this.isValueGroupType;
    },

    isValuePathsFormShown() {
      return this.isValueGroupType;
    },

    isTimeBasedFormShown() {
      return this.isComplexType
        || this.isValueGroupType
        || this.isTimeBasedType
        || this.isCorelFormShown;
    },

    isPatternsFormShown() {
      return this.isComplexType
        || this.isValueGroupType
        || this.isPatternsType
        || this.isCorelFormShown;
    },

    withTotalEntityPatterns() {
      return this.isComplexType || this.isValueGroupType;
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

    isCorelFormShown() {
      return this.form.type === META_ALARMS_RULE_TYPES.corel;
    },
  },
};
</script>
