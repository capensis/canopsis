<template>
  <v-layout column>
    <c-enabled-field v-field="form.enabled" with-background />
    <v-tabs
      v-model="activeTab"
      slider-color="primary"
      centered
    >
      <v-tab :class="{ 'error--text': hasGeneralError }">
        {{ $t('common.general') }}
      </v-tab>
      <v-tab :class="{ 'error--text': hasPatternsError }">
        {{ $tc('common.pattern') }}
      </v-tab>

      <v-tab-item eager>
        <idle-rule-general-form
          v-field="form"
          ref="generalFormElement"
          :is-entity-type="isEntityType"
        />
      </v-tab-item>
      <v-tab-item eager>
        <idle-rule-patterns-form
          v-field="form.patterns"
          ref="patternsFormElement"
          :is-entity-type="isEntityType"
          class="mt-2"
        />
      </v-tab-item>
    </v-tabs>
  </v-layout>
</template>

<script>
import { computed, ref, watch } from 'vue';

import { isIdleRuleEntityType } from '@/helpers/entities/idle-rule/form';

import { useAiChatExpand } from '@/hooks/ai/ai-chat-form';

import IdleRuleGeneralForm from './idle-rule-general-form.vue';
import IdleRulePatternsForm from './idle-rule-patterns-form.vue';

const IDLE_RULE_FORM_TABS = {
  general: 0,
  patterns: 1,
};

export default {
  inject: ['$validator'],
  components: {
    IdleRuleGeneralForm,
    IdleRulePatternsForm,
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
    const activeTab = ref(IDLE_RULE_FORM_TABS.general);
    const hasGeneralError = ref(false);
    const hasPatternsError = ref(false);

    const generalFormElement = ref(null);
    const patternsFormElement = ref(null);

    const isEntityType = computed(() => isIdleRuleEntityType(props.form.type));

    useAiChatExpand({ activeTab, neededTab: IDLE_RULE_FORM_TABS.patterns });

    watch(() => generalFormElement.value?.hasAnyError, (value) => {
      hasGeneralError.value = value ?? false;
    });

    watch(() => patternsFormElement.value?.hasAnyError, (value) => {
      hasPatternsError.value = value ?? false;
    });

    return {
      activeTab,
      hasGeneralError,
      hasPatternsError,
      generalFormElement,
      patternsFormElement,
      isEntityType,
    };
  },
};
</script>
