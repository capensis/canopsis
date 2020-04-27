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
    meta-alarm-rule-threshhold-form(v-if="isComplexType", v-field="form.config")
    meta-alarm-rule-timebased-form(v-if="isComplexType || isTimebasedType", v-field="form.config")
    meta-alarm-rule-patterns-form(v-if="isComplexType || isPatternsType", v-field="form.config")
</template>

<script>
import { META_ALARMS_RULE_TYPES } from '@/constants';

import MetaAlarmRuleThreshholdForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-threshhold-form.vue';
import MetaAlarmRulePatternsForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-patterns-form.vue';
import MetaAlarmRuleTimebasedForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-timebased-form.vue';

export default {
  inject: ['$validator'],
  components: {
    MetaAlarmRuleTimebasedForm,
    MetaAlarmRulePatternsForm,
    MetaAlarmRuleThreshholdForm,
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
    isPatternsType() {
      return this.form.type === META_ALARMS_RULE_TYPES.attribute;
    },
    isComplexType() {
      return this.form.type === META_ALARMS_RULE_TYPES.complex;
    },
    isTimebasedType() {
      return this.form.type === META_ALARMS_RULE_TYPES.timebased;
    },
  },
};
</script>

