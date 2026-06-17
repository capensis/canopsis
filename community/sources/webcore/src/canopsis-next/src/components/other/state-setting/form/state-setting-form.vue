<template>
  <v-layout class="gap-2" column>
    <c-enabled-field v-field="form.enabled" with-background />
    <v-stepper
      v-model="stepper"
      class="state-setting-form"
    >
      <v-stepper-header>
        <v-stepper-step
          :complete="stepper > STATE_SETTING_FORM_STEPS.basics"
          :step="STATE_SETTING_FORM_STEPS.basics"
          :rules="[() => !hasBasicsFormAnyError]"
          editable
        >
          {{ $t('stateSetting.steps.basics') }}
        </v-stepper-step>
        <v-divider />
        <v-stepper-step
          :complete="stepper > STATE_SETTING_FORM_STEPS.entityPattern"
          :step="STATE_SETTING_FORM_STEPS.entityPattern"
          :rules="[() => !hasEntityPatternFormAnyError]"
          editable
        >
          {{ $t('stateSetting.steps.rulePatterns') }}
        </v-stepper-step>
        <v-divider />
        <v-stepper-step
          :complete="stepper > STATE_SETTING_FORM_STEPS.thresholds"
          :step="STATE_SETTING_FORM_STEPS.thresholds"
          :rules="[() => !hasThresholdsFormAnyError]"
          editable
        >
          {{ $t('stateSetting.steps.conditions') }}
        </v-stepper-step>
      </v-stepper-header>
      <v-stepper-items>
        <v-stepper-content :step="STATE_SETTING_FORM_STEPS.basics">
          <state-setting-basics-step
            v-field="form"
            ref="basicsFormElement"
          />
        </v-stepper-content>
        <v-stepper-content :step="STATE_SETTING_FORM_STEPS.entityPattern">
          <c-alert
            class="mb-4"
            type="info"
          >
            {{ methodMessage }}
          </c-alert>
          <state-setting-entity-pattern-step
            v-field="form.entity_pattern"
            ref="entityPatternFormElement"
            :entity-types="patternEntityTypes"
          />
        </v-stepper-content>
        <v-stepper-content :step="STATE_SETTING_FORM_STEPS.thresholds">
          <c-alert
            class="mb-4"
            type="info"
          >
            {{ methodMessage }}
          </c-alert>
          <state-setting-inherited-entity-pattern-step
            v-if="isInheritedMethod"
            v-field="form.inherited_entity_pattern"
            ref="thresholdsFormElement"
          />
          <state-setting-thresholds-step
            v-else
            v-field="form.state_thresholds"
            ref="thresholdsFormElement"
          />
        </v-stepper-content>
      </v-stepper-items>
    </v-stepper>
  </v-layout>
</template>

<script>
import { computed, inject, ref, watch } from 'vue';

import { STATE_SETTING_METHODS, PATTERNS_FIELDS } from '@/constants';

import { useAiChatExpand } from '@/hooks/ai/ai-chat-form';
import { useI18n } from '@/hooks/i18n';

import StateSettingBasicsStep from './steps/state-setting-basics-step.vue';
import StateSettingEntityPatternStep from './steps/state-setting-entity-pattern-step.vue';
import StateSettingInheritedEntityPatternStep from './steps/state-setting-inherited-entity-pattern-step.vue';
import StateSettingThresholdsStep from './steps/state-setting-thresholds-step.vue';

const STATE_SETTING_FORM_STEPS = {
  basics: 1,
  entityPattern: 2,
  thresholds: 3,
};

export default {
  components: {
    StateSettingBasicsStep,
    StateSettingEntityPatternStep,
    StateSettingInheritedEntityPatternStep,
    StateSettingThresholdsStep,
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
  },
  setup(props) {
    inject('$validator');

    const { t } = useI18n();

    const stepper = ref(STATE_SETTING_FORM_STEPS.basics);
    const hasBasicsFormAnyError = ref(false);
    const hasEntityPatternFormAnyError = ref(false);
    const hasThresholdsFormAnyError = ref(false);

    const basicsFormElement = ref(null);
    const entityPatternFormElement = ref(null);
    const thresholdsFormElement = ref(null);

    const isInheritedMethod = computed(() => props.form.method === STATE_SETTING_METHODS.inherited);

    const methodMessage = computed(() => t(`stateSetting.methods.${props.form.method}.stepTitle`));

    const patternEntityTypes = computed(() => [props.form.type]);

    useAiChatExpand({
      activeTab: stepper,
      neededTab: {
        [PATTERNS_FIELDS.entity]: STATE_SETTING_FORM_STEPS.entityPattern,
        inherited_entity_pattern: STATE_SETTING_FORM_STEPS.thresholds,
      },
    });

    watch(() => basicsFormElement.value?.hasAnyError, (value) => {
      hasBasicsFormAnyError.value = value ?? false;
    });

    watch(() => entityPatternFormElement.value?.hasAnyError, (value) => {
      hasEntityPatternFormAnyError.value = value ?? false;
    });

    watch(() => thresholdsFormElement.value?.hasAnyError, (value) => {
      hasThresholdsFormAnyError.value = value ?? false;
    });

    return {
      STATE_SETTING_FORM_STEPS,
      stepper,
      hasBasicsFormAnyError,
      hasEntityPatternFormAnyError,
      hasThresholdsFormAnyError,
      basicsFormElement,
      entityPatternFormElement,
      thresholdsFormElement,
      isInheritedMethod,
      methodMessage,
      patternEntityTypes,
    };
  },
};
</script>

<style lang="scss">
.state-setting-form {
  .v-stepper__wrapper {
    overflow: unset;
  }
}
</style>
