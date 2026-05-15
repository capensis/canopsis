<template>
  <v-layout class="gap-3" column>
    <div class="mb-2">
      <div class="text-subtitle-2 mb-2">
        {{ $t(`metaAlarmRule.types.${form.type}.text`) }}
      </div>
      <div class="text-body-2">
        {{ $t(`metaAlarmRule.types.${form.type}.helpText`) }}
      </div>
    </div>

    <c-form-block v-if="isThresholdFormShown || isTimeBasedFormShown || isValuePathsFormShown || isCorelFormShown">
      <meta-alarm-rule-threshold-field
        v-if="isThresholdFormShown"
        v-field="form.config"
      />

      <meta-alarm-rule-time-based-field
        v-if="isTimeBasedFormShown"
        v-field="form.config.time_interval"
      />

      <meta-alarm-rule-child-inactive-delay-field
        v-if="isChildInactiveDelayFormShown"
        v-field="form.config.child_inactive_delay"
      />

      <meta-alarm-rule-value-paths-field
        v-if="isValuePathsFormShown"
        v-field="form.config.value_paths"
        required
      />

      <meta-alarm-rule-corel-form
        v-if="isCorelFormShown"
        v-field="form.config"
        :template-vars="templateVars"
      />
    </c-form-block>

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

import MetaAlarmRuleThresholdField from '@/components/other/meta-alarm-rule/form/fields/meta-alarm-rule-threshold-field.vue';
import MetaAlarmRuleValuePathsField from '@/components/other/meta-alarm-rule/form/fields/meta-alarm-value-paths-field.vue';
import MetaAlarmRuleTimeBasedField from '@/components/other/meta-alarm-rule/form/fields/meta-alarm-rule-time-based-field.vue';
import MetaAlarmRuleChildInactiveDelayField from '@/components/other/meta-alarm-rule/form/fields/meta-alarm-rule-child-inactive-delay-field.vue';
import MetaAlarmRuleCorelForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-corel-form.vue';
import MetaAlarmRulePatternsForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-patterns-form.vue';

export default {
  components: {
    MetaAlarmRuleTimeBasedField,
    MetaAlarmRuleValuePathsField,
    MetaAlarmRuleThresholdField,
    MetaAlarmRuleChildInactiveDelayField,
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
    templateVars: {
      type: Object,
      default: () => ({}),
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
    const isChildInactiveDelayFormShown = computed(() => isValueGroupType.value);
    const isCorelFormShown = computed(() => isCorelType.value);
    const withTotalEntityPattern = computed(() => isMetaAlarmRuleTypeHasTotalEntityPatterns(props.form.type));

    return {
      isAttributeType,

      isThresholdFormShown,
      isValuePathsFormShown,
      isTimeBasedFormShown,
      isChildInactiveDelayFormShown,
      isCorelFormShown,
      withTotalEntityPattern,
    };
  },
};
</script>
