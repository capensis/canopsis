<template lang="pug">
  div
    v-stepper.state-setting-form(v-model="stepper")
      v-stepper-header
        v-stepper-step(
          :complete="stepper > steps.BASICS",
          :step="steps.BASICS",
          :rules="[() => !hasBasicsFormAnyError]",
          editable
        ) {{ $t('stateSetting.steps.basics') }}
        v-divider
        v-stepper-step(
          :complete="stepper > steps.RULE_PATTERNS",
          :step="steps.RULE_PATTERNS",
          :rules="[() => !hasRulePatternsFormAnyError]",
          editable
        ) {{ $t('stateSetting.steps.rulePatterns') }}
        v-divider
        v-stepper-step(
          :complete="stepper > steps.CONDITIONS",
          :step="steps.CONDITIONS",
          :rules="[() => !hasConditionsFormAnyError]",
          editable
        ) {{ $t('stateSetting.steps.conditions') }}
      v-stepper-items
        v-stepper-content(:step="steps.BASICS")
          state-setting-basics-step(
            ref="basicsForm",
            v-field="form"
          )
        v-stepper-content(:step="steps.RULE_PATTERNS")
          state-setting-rule-patterns-step(
            ref="rulePatternsForm",
            v-field="form.rule_patterns"
          )
        v-stepper-content(:step="steps.CONDITIONS")
          state-setting-impacting-patterns-step(
            ref="conditionsForm",
            v-if="isInheritedMethod",
            v-field="form.impacting_patterns"
          )
          state-setting-conditions-step(
            ref="conditionsForm",
            v-else,
            v-field="form.conditions"
          )
</template>

<script>
import { STATE_SETTING_METHODS } from '@/constants';

import { formMixin } from '@/mixins/form';

import StateSettingBasicsStep from './steps/state-setting-basics-step.vue';
import StateSettingRulePatternsStep from './steps/state-setting-rule-patterns-step.vue';
import StateSettingImpactingPatternsStep from './steps/state-setting-impacting-patterns-step.vue';
import StateSettingConditionsStep from './steps/state-setting-conditions-step.vue';

export default {
  components: {
    StateSettingBasicsStep,
    StateSettingRulePatternsStep,
    StateSettingImpactingPatternsStep,
    StateSettingConditionsStep,
  },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      stepper: 1,
      hasBasicsFormAnyError: false,
      hasRulePatternsFormAnyError: false,
      hasConditionsFormAnyError: false,
    };
  },
  computed: {
    steps() {
      return {
        BASICS: 1,
        RULE_PATTERNS: 2,
        CONDITIONS: 3,
      };
    },

    isInheritedMethod() {
      return this.form.method === STATE_SETTING_METHODS.inherited;
    },
  },

  mounted() {
    this.$watch(() => this.$refs.basicsForm.hasAnyError, (value) => {
      this.hasBasicsFormAnyError = value;
    });

    this.$watch(() => this.$refs.rulePatternsForm.hasAnyError, (value) => {
      this.hasRulePatternsFormAnyError = value;
    });

    this.$watch(() => this.$refs.conditionsForm.hasAnyError, (value) => {
      this.hasConditionsFormAnyError = value;
    });
  },
};
</script>

<style lang="scss">
.state-setting-form {
  background-color: transparent !important;

  .v-stepper__wrapper {
    overflow: unset;
  }
}
</style>
