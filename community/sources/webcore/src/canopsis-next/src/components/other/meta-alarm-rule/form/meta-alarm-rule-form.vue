<template>
  <v-stepper :value="activeStep" @change="$emit('update:active-step', $event)">
    <v-stepper-header>
      <v-stepper-step
        :complete="activeStep > 0"
        :step="0"
        :rules="[() => true]"
        class="py-0"
        editable
      >
        {{ $t('metaAlarmRule.steps.basics') }}
        <small v-if="false">{{ $t('modals.createDynamicInfo.errors.invalid') }}</small>
      </v-stepper-step>
      <v-divider />
      <v-stepper-step
        :complete="activeStep > 1"
        :step="1"
        :rules="[() => true]"
        class="py-0"
        editable
      >
        {{ $t('metaAlarmRule.steps.defineType') }}
        <small v-if="false">{{ $t('modals.createDynamicInfo.errors.invalid') }}</small>
      </v-stepper-step>
      <v-divider />
      <v-stepper-step
        :complete="activeStep > 2"
        :step="2"
        :rules="[() => true]"
        class="py-0"
        editable
      >
        {{ $t('metaAlarmRule.steps.addParameters') }}
        <small v-if="false">{{ $t('modals.createDynamicInfo.errors.invalid') }}</small>
      </v-stepper-step>
    </v-stepper-header>

    <v-stepper-items>
      <v-stepper-content
        :step="0"
        class="pa-0"
      >
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
      </v-stepper-content>
      <v-stepper-content
        :step="1"
        class="pa-0"
      >
        <v-select
          v-field="form.type"
          :items="ruleTypes"
          :label="$t('common.type')"
        />
      </v-stepper-content>
      <v-stepper-content
        :step="2"
        class="pa-0"
      >
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
      </v-stepper-content>
    </v-stepper-items>
  </v-stepper>
</template>

<script>
import { computed } from 'vue';

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
    activeStep: {
      type: Number,
      default: 0,
    },
  },
  setup(props) {
    /**
     * We are filtered 'manualgroup' because we are using in only in the alarms list widget directly
     */
    const ruleTypes = computed(() => Object.values(META_ALARMS_RULE_TYPES)
      .filter(type => !isManualGroupMetaAlarmRuleType(type)));

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
      || isTimeBasedType.value
      || isCorelType.value);
    const isCorelFormShown = computed(() => isCorelType.value);
    const withTotalEntityPattern = computed(() => isMetaAlarmRuleTypeHasTotalEntityPatterns(props.form.type));
    const withChildInactiveDelay = computed(() => isComplexType.value
      || isValueGroupType.value
      || isCorelFormShown.value);

    return {
      ruleTypes,

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
