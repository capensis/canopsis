<template>
  <v-layout column>
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
      class="mb-3"
    />
    <meta-alarm-rule-corel-form
      v-if="isCorelFormShown"
      v-field="form.config"
      :variables="variables"
    />
    <meta-alarm-rule-patterns-form
      v-field="form.patterns"
      :with-total-entity="withTotalEntityPattern"
      :some-required="isAttributeType"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import {
  isAttributeMetaAlarmRuleType,
  isComplexMetaAlarmRuleType,
  isCorelMetaAlarmRuleType,
  isMetaAlarmRuleTypeHasTotalEntityPatterns,
  isTimebasedMetaAlarmRuleType,
  isValueGroupMetaAlarmRuleType,
} from '@/helpers/entities/meta-alarm/rule/form';

import MetaAlarmRuleThresholdForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-threshold-form.vue';
import MetaAlarmRuleValuePathsForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-value-paths-form.vue';
import MetaAlarmRulePatternsForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-patterns-form.vue';
import MetaAlarmRuleCorelForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-corel-form.vue';
import MetaAlarmRuleTimeBasedForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-time-based-form.vue';

export default {
  components: {
    MetaAlarmRuleTimeBasedForm,
    MetaAlarmRuleCorelForm,
    MetaAlarmRulePatternsForm,
    MetaAlarmRuleValuePathsForm,
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
    variables: {
      type: Array,
      required: false,
    },
  },
  setup(props) {
    /**
     * Rule types
     */
    const isAttributeType = computed(() => isAttributeMetaAlarmRuleType(props.form.type));
    const isTimeBasedType = computed(() => isTimebasedMetaAlarmRuleType(props.form.type));
    const isComplexType = computed(() => isComplexMetaAlarmRuleType(props.form.type));
    const isValueGroupType = computed(() => isValueGroupMetaAlarmRuleType(props.form.type));
    const isCorelType = computed(() => isCorelMetaAlarmRuleType(props.form.type));

    /**
     * Conditions for forms showing
     */
    const isThresholdFormShown = computed(() => isComplexType.value || isValueGroupType.value);
    const isValuePathsFormShown = computed(() => isValueGroupType.value);
    const isTimeBasedFormShown = computed(() => isComplexType.value
      || isValueGroupType.value
      || isTimeBasedType.value);
    const isCorelFormShown = computed(() => isCorelType.value);
    const withTotalEntityPattern = computed(() => isMetaAlarmRuleTypeHasTotalEntityPatterns(props.form.type));
    const withChildInactiveDelay = computed(() => isComplexType.value || isValueGroupType.value);

    return {
      isAttributeType,

      isThresholdFormShown,
      isValuePathsFormShown,
      isTimeBasedFormShown,
      isCorelFormShown,
      withTotalEntityPattern,
      withChildInactiveDelay,
    };
  },
};
</script>
