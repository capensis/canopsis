<template>
  <div>
    <c-id-field
      v-field="form._id"
      :disabled="isDisabledIdField"
      :help-text="$t('metaAlarmRule.idHelp')"
    />
    <c-name-field
      v-field="form.name"
      required
    />
    <c-description-field
      v-field="form.output_template"
      :label="$t('metaAlarmRule.outputTemplate')"
      :help-text="$t('metaAlarmRule.outputTemplateHelp')"
    />
    <c-enabled-field
      v-field="form.auto_resolve"
      :label="$t('metaAlarmRule.autoResolve')"
    />
    <v-select
      v-field="form.type"
      :items="ruleTypes"
      :label="$t('common.type')"
    />
    <meta-alarm-rule-corel-form
      v-if="isCorelFormShown"
      v-field="form.config"
    />
    <meta-alarm-rule-threshold-form
      v-if="isThresholdFormShown"
      v-field="form.config"
    />
    <meta-alarm-rule-time-based-form
      v-if="isTimeBasedFormShown"
      v-field="form.config"
      :with-child-inactive-delay="withChildInactiveDelay"
    />
    <meta-alarm-rule-value-paths-form
      v-if="isValuePathsFormShown"
      v-field="form.config"
      class="mb-2"
    />
    <meta-alarm-rule-patterns-form
      v-field="form.patterns"
      :with-total-entity="withTotalEntityPattern"
      :some-required="isAttributeType"
    />
  </div>
</template>

<script>
import { META_ALARMS_RULE_TYPES } from '@/constants';

import {
  isAttributeMetaAlarmRuleType,
  isComplexMetaAlarmRuleType,
  isCorelMetaAlarmRuleType,
  isManualGroupMetaAlarmRuleType,
  isMetaAlarmRuleTypeHasTotalEntityPatterns,
  isTimebasedMetaAlarmRuleType,
  isValueGroupMetaAlarmRuleType,
} from '@/helpers/entities/meta-alarm/rule/form';

import MetaAlarmRuleThresholdForm from './meta-alarm-rule-threshold-form.vue';
import MetaAlarmRuleTimeBasedForm from './meta-alarm-rule-time-based-form.vue';
import MetaAlarmRuleValuePathsForm from './meta-alarm-rule-value-paths-form.vue';
import MetaAlarmRuleCorelForm from './meta-alarm-rule-corel-form.vue';
import MetaAlarmRulePatternsForm from './meta-alarm-rule-patterns-form.vue';

export default {
  inject: ['$validator'],
  components: {
    MetaAlarmRuleTimeBasedForm,
    MetaAlarmRuleThresholdForm,
    MetaAlarmRuleValuePathsForm,
    MetaAlarmRuleCorelForm,
    MetaAlarmRulePatternsForm,
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
        .filter(type => !isManualGroupMetaAlarmRuleType(type));
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

    isCorelFormShown() {
      return this.isCorelType;
    },

    withTotalEntityPattern() {
      return isMetaAlarmRuleTypeHasTotalEntityPatterns(this.form.type);
    },

    withChildInactiveDelay() {
      return this.isComplexType
        || this.isValueGroupType
        || this.isCorelFormShown;
    },

    /**
     * Rule types
     */
    isAttributeType() {
      return isAttributeMetaAlarmRuleType(this.form.type);
    },

    isTimeBasedType() {
      return isTimebasedMetaAlarmRuleType(this.form.type);
    },

    isComplexType() {
      return isComplexMetaAlarmRuleType(this.form.type);
    },

    isValueGroupType() {
      return isValueGroupMetaAlarmRuleType(this.form.type);
    },

    isCorelType() {
      return isCorelMetaAlarmRuleType(this.form.type);
    },
  },
};
</script>
