<template>
  <v-layout class="gap-3" column>
    <c-enabled-field v-field="form.enabled" with-background />
    <c-form-general-patterns-tabs
      :form="form"
      :patterns-label="$tc('common.condition', 2)"
      :additional-label="$t('stateSetting.steps.targetEntities')"
      :ai-chat-expand-tab-keys="aiChatExpandTabKeys"
    >
      <template #general="{ setRef }">
        <state-setting-general-form
          v-field="form"
          :ref="setRef"
        />
      </template>

      <template #additional="{ setRef }">
        <v-layout class="gap-4" column>
          <c-alert type="info">
            {{ methodMessage }}
          </c-alert>
          <state-setting-entity-patterns-form
            v-field="form.entity_pattern"
            :ref="setRef"
            :entity-types="patternEntityTypes"
          />
        </v-layout>
      </template>

      <template #patterns="{ setRef }">
        <v-layout class="gap-4" column>
          <c-alert type="info">
            {{ methodMessage }}
          </c-alert>
          <state-setting-inherited-entity-pattern-step
            v-if="isInheritedMethod"
            v-field="form.inherited_entity_pattern"
            :ref="setRef"
          />
          <state-setting-thresholds-step
            v-else
            v-field="form.state_thresholds"
            :ref="setRef"
          />
        </v-layout>
      </template>
    </c-form-general-patterns-tabs>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import {
  FORM_GENERAL_PATTERNS_TABS,
  STATE_SETTING_METHODS,
  PATTERNS_FIELDS,
  STATE_SETTINGS_INHERITED_ENTITY_PATTERN_FIELD,
} from '@/constants';

import { useI18n } from '@/hooks/i18n';

import StateSettingGeneralForm from './state-setting-general-form.vue';
import StateSettingEntityPatternsForm from './state-setting-entity-patterns-form.vue';
import StateSettingInheritedEntityPatternStep from './steps/state-setting-inherited-entity-pattern-step.vue';
import StateSettingThresholdsStep from './steps/state-setting-thresholds-step.vue';

const AI_CHAT_EXPAND_TAB_KEYS = {
  [PATTERNS_FIELDS.entity]: FORM_GENERAL_PATTERNS_TABS.patterns,
  [STATE_SETTINGS_INHERITED_ENTITY_PATTERN_FIELD]: FORM_GENERAL_PATTERNS_TABS.additional,
};

export default {
  components: {
    StateSettingGeneralForm,
    StateSettingEntityPatternsForm,
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
    const { t } = useI18n();

    const isInheritedMethod = computed(() => props.form.method === STATE_SETTING_METHODS.inherited);

    const methodMessage = computed(() => t(`stateSetting.methods.${props.form.method}.stepTitle`));

    const patternEntityTypes = computed(() => [props.form.type]);

    return {
      aiChatExpandTabKeys: AI_CHAT_EXPAND_TAB_KEYS,
      isInheritedMethod,
      methodMessage,
      patternEntityTypes,
    };
  },
};
</script>
