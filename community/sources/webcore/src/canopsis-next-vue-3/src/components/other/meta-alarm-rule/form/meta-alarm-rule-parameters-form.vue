<template>
  <c-layout gap="3" column>
    <div class="mb-2">
      <div class="text-subtitle-2 mb-2">
        {{ t(`metaAlarmRule.types.${typeField}.text`) }}
      </div>
      <div class="text-body-2">
        {{ t(`metaAlarmRule.types.${typeField}.helpText`) }}
      </div>
    </div>

    <c-form-block v-if="isFormBlockShown">
      <meta-alarm-rule-threshold-field
        v-if="isThresholdFormShown"
        v-model:config="configField"
      />

      <meta-alarm-rule-time-based-field
        v-if="isTimeBasedFormShown"
        v-model="timeIntervalField"
      />

      <meta-alarm-rule-child-inactive-delay-field
        v-if="isChildInactiveDelayFormShown"
        v-model="childInactiveDelayField"
      />

      <meta-alarm-value-paths-field
        v-if="isValuePathsFormShown"
        v-model:items="valuePathsField"
        name="config.value_paths"
        required
      />

      <meta-alarm-rule-corel-form
        v-if="isCorelFormShown"
        v-model:config="configField"
        :template-vars="templateVars"
      />
    </c-form-block>

    <meta-alarm-rule-patterns-form
      v-model:form="patternsField"
      :with-total-entity="withTotalEntityPattern"
      :some-required="isAttributeType"
    />
  </c-layout>
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

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import MetaAlarmRuleChildInactiveDelayField from './fields/meta-alarm-rule-child-inactive-delay-field.vue';
import MetaAlarmRuleThresholdField from './fields/meta-alarm-rule-threshold-field.vue';
import MetaAlarmRuleTimeBasedField from './fields/meta-alarm-rule-time-based-field.vue';
import MetaAlarmValuePathsField from './fields/meta-alarm-value-paths-field.vue';
import MetaAlarmRuleCorelForm from './meta-alarm-rule-corel-form.vue';
import MetaAlarmRulePatternsForm from './meta-alarm-rule-patterns-form.vue';

export default {
  components: {
    MetaAlarmRuleTimeBasedField,
    MetaAlarmValuePathsField,
    MetaAlarmRuleThresholdField,
    MetaAlarmRuleChildInactiveDelayField,
    MetaAlarmRuleCorelForm,
    MetaAlarmRulePatternsForm,
  },

  props: {
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },

  setup() {
    const { t } = useI18n();
    const {
      typeField,
      configField,
      timeIntervalField,
      childInactiveDelayField,
      valuePathsField,
      patternsField,
    } = useFormModels({
      typeField: 'type',
      configField: 'config',
      timeIntervalField: 'config.time_interval',
      childInactiveDelayField: 'config.child_inactive_delay',
      valuePathsField: 'config.value_paths',
      patternsField: 'patterns',
    });

    const isAttributeType = computed(() => isAttributeMetaAlarmRuleType(typeField.value));
    const isTimeBasedType = computed(() => isTimebasedMetaAlarmRuleType(typeField.value));
    const isComplexType = computed(() => isComplexMetaAlarmRuleType(typeField.value));
    const isValueGroupType = computed(() => isValueGroupMetaAlarmRuleType(typeField.value));
    const isCorelType = computed(() => isCorelMetaAlarmRuleType(typeField.value));
    const isThresholdFormShown = computed(() => isComplexType.value || isValueGroupType.value);
    const isValuePathsFormShown = computed(() => isValueGroupType.value);
    const isTimeBasedFormShown = computed(
      () => isComplexType.value || isValueGroupType.value || isTimeBasedType.value,
    );
    const isChildInactiveDelayFormShown = computed(
      () => isComplexType.value || isValueGroupType.value,
    );
    const isCorelFormShown = computed(() => isCorelType.value);
    const isFormBlockShown = computed(
      () =>
        isThresholdFormShown.value ||
        isTimeBasedFormShown.value ||
        isChildInactiveDelayFormShown.value ||
        isValuePathsFormShown.value ||
        isCorelFormShown.value,
    );
    const withTotalEntityPattern = computed(() =>
      isMetaAlarmRuleTypeHasTotalEntityPatterns(typeField.value),
    );

    return {
      t,
      typeField,
      configField,
      timeIntervalField,
      childInactiveDelayField,
      valuePathsField,
      patternsField,
      isAttributeType,
      isThresholdFormShown,
      isValuePathsFormShown,
      isTimeBasedFormShown,
      isChildInactiveDelayFormShown,
      isCorelFormShown,
      isFormBlockShown,
      withTotalEntityPattern,
    };
  },
};
</script>
