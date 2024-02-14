<template>
  <v-stepper
    v-model="stepper"
    class="state-setting-form"
  >
    <v-stepper-header>
      <v-stepper-step
        :complete="stepper > steps.BASICS"
        :step="steps.BASICS"
        :rules="[() => !hasBasicsFormAnyError]"
        editable
      >
        {{ $t('stateSetting.steps.basics') }}
      </v-stepper-step>
      <v-divider />
      <v-stepper-step
        :complete="stepper > steps.ENTITY_PATTERN"
        :step="steps.ENTITY_PATTERN"
        :rules="[() => !hasEntityPatternFormAnyError]"
        editable
      >
        {{ $t('stateSetting.steps.rulePatterns') }}
      </v-stepper-step>
      <v-divider />
      <v-stepper-step
        :complete="stepper > steps.THRESHOLDS"
        :step="steps.THRESHOLDS"
        :rules="[() => !hasThresholdsFormAnyError]"
        editable
      >
        {{ $t('stateSetting.steps.conditions') }}
      </v-stepper-step>
    </v-stepper-header>
    <v-stepper-items>
      <v-stepper-content :step="steps.BASICS">
        <state-setting-basics-step
          v-field="form"
          ref="basicsForm"
        />
      </v-stepper-content>
      <v-stepper-content :step="steps.ENTITY_PATTERN">
        <c-alert
          class="mb-4"
          type="info"
        >
          {{ methodMessage }}
        </c-alert>
        <state-setting-entity-pattern-step
          v-field="form.entity_pattern"
          ref="entityPatternForm"
          :entity-types="patternEntityTypes"
        />
      </v-stepper-content>
      <v-stepper-content :step="steps.THRESHOLDS">
        <c-alert
          class="mb-4"
          type="info"
        >
          {{ methodMessage }}
        </c-alert>
        <state-setting-inherited-entity-pattern-step
          v-if="isInheritedMethod"
          v-field="form.inherited_entity_pattern"
          ref="thresholdsForm"
        />
        <state-setting-thresholds-step
          v-else
          v-field="form.state_thresholds"
          ref="thresholdsForm"
        />
      </v-stepper-content>
    </v-stepper-items>
  </v-stepper>
</template>

<script>
import { STATE_SETTING_METHODS } from '@/constants';

import { formMixin } from '@/mixins/form';

import StateSettingBasicsStep from './steps/state-setting-basics-step.vue';
import StateSettingEntityPatternStep from './steps/state-setting-entity-pattern-step.vue';
import StateSettingInheritedEntityPatternStep from './steps/state-setting-inherited-entity-pattern-step.vue';
import StateSettingThresholdsStep from './steps/state-setting-thresholds-step.vue';

export default {
  components: {
    StateSettingBasicsStep,
    StateSettingEntityPatternStep,
    StateSettingInheritedEntityPatternStep,
    StateSettingThresholdsStep,
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
      hasEntityPatternFormAnyError: false,
      hasThresholdsFormAnyError: false,
    };
  },
  computed: {
    steps() {
      return {
        BASICS: 1,
        ENTITY_PATTERN: 2,
        THRESHOLDS: 3,
      };
    },

    isInheritedMethod() {
      return this.form.method === STATE_SETTING_METHODS.inherited;
    },

    methodMessage() {
      return this.$t(`stateSetting.methods.${this.form.method}.stepTitle`);
    },

    patternEntityTypes() {
      return [this.form.type];
    },
  },

  mounted() {
    this.$watch(() => this.$refs.basicsForm.hasAnyError, (value) => {
      this.hasBasicsFormAnyError = value;
    });

    this.$watch(() => this.$refs.entityPatternForm.hasAnyError, (value) => {
      this.hasEntityPatternFormAnyError = value;
    });

    this.$watch(() => this.$refs.thresholdsForm.hasAnyError, (value) => {
      this.hasThresholdsFormAnyError = value;
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
