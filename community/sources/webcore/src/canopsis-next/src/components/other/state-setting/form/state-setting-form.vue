<template>
  <c-layout gap="3" column>
    <c-enabled-field v-model="enabled" with-background />

    <c-form-general-patterns-tabs
      :form="formValues"
      :patterns-label="t('common.condition', 2)"
      :additional-label="t('stateSetting.steps.targetEntities')"
      :ai-chat-expand-tab-keys="aiChatExpandTabKeys"
    >
      <template #general="{ setRef }">
        <state-setting-general-form :ref="setRef" />
      </template>

      <template #additional="{ setRef }">
        <c-layout gap="4" column>
          <c-alert type="info">
            {{ methodMessage }}
          </c-alert>
          <state-setting-entity-patterns-form :ref="setRef" :entity-types="patternEntityTypes" />
        </c-layout>
      </template>

      <template #patterns="{ setRef }">
        <c-layout gap="4" column>
          <c-alert type="info">
            {{ methodMessage }}
          </c-alert>
          <state-setting-inherited-entity-pattern-step v-if="isInheritedMethod" :ref="setRef" />
          <state-setting-thresholds-step v-else :ref="setRef" />
        </c-layout>
      </template>
    </c-form-general-patterns-tabs>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import {
  FORM_GENERAL_PATTERNS_TABS,
  PATTERNS_FIELDS,
  STATE_SETTING_METHODS,
  STATE_SETTINGS_INHERITED_ENTITY_PATTERN_FIELD,
} from '@/constants';

import { useFormContext, useFormModel } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import StateSettingEntityPatternsForm from './state-setting-entity-patterns-form.vue';
import StateSettingGeneralForm from './state-setting-general-form.vue';
import StateSettingInheritedEntityPatternStep from './steps/state-setting-inherited-entity-pattern-step.vue';
import StateSettingThresholdsStep from './steps/state-setting-thresholds-step.vue';

const AI_CHAT_EXPAND_TAB_KEYS = {
  [PATTERNS_FIELDS.entity]: FORM_GENERAL_PATTERNS_TABS.patterns,
  [STATE_SETTINGS_INHERITED_ENTITY_PATTERN_FIELD]: FORM_GENERAL_PATTERNS_TABS.additional,
};

export default {
  components: {
    StateSettingEntityPatternsForm,
    StateSettingGeneralForm,
    StateSettingInheritedEntityPatternStep,
    StateSettingThresholdsStep,
  },

  setup() {
    const { t } = useI18n();
    const form = useFormContext();
    const enabled = useFormModel('enabled');
    const method = useFormModel('method');
    const type = useFormModel('type');

    const isInheritedMethod = computed(() => method.value === STATE_SETTING_METHODS.inherited);
    const methodMessage = computed(() => t(`stateSetting.methods.${method.value}.stepTitle`));
    const patternEntityTypes = computed(() => [type.value]);

    return {
      t,
      formValues: form.values,
      enabled,
      aiChatExpandTabKeys: AI_CHAT_EXPAND_TAB_KEYS,
      isInheritedMethod,
      methodMessage,
      patternEntityTypes,
    };
  },
};
</script>
