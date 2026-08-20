<template>
  <v-layout class="gap-2" column>
    <c-enabled-field
      v-if="!noEnabled"
      v-field="form.enabled"
      hide-details
      with-background
    />
    <pbehavior-general-form
      v-if="noPattern"
      v-field="form"
      :no-enabled="noEnabled"
      :no-comments="noComments"
      :no-timezone="noTimezone"
      :with-start-on-trigger="withStartOnTrigger"
      :with-inherited="withInherited"
      :name-label="nameLabel"
      :name-tooltip="nameTooltip"
    />

    <v-tabs
      v-else
      v-model="activeTab"
      slider-color="primary"
      centered
    >
      <v-tab :class="{ 'error--text': hasGeneralError }">
        {{ $t('common.general') }}
      </v-tab>
      <v-tab :class="{ 'error--text': hasPatternsError }">
        {{ $tc('common.pattern', 2) }}
      </v-tab>

      <v-tab-item eager>
        <pbehavior-general-form
          v-field="form"
          ref="generalFormElement"
          :no-enabled="noEnabled"
          :no-comments="noComments"
          :no-timezone="noTimezone"
          :with-start-on-trigger="withStartOnTrigger"
          :with-inherited="withInherited"
          :name-label="nameLabel"
          :name-tooltip="nameTooltip"
        />
      </v-tab-item>
      <v-tab-item eager>
        <pbehavior-patterns-form
          v-field="form.patterns"
          ref="patternsFormElement"
          :pbehavior-id="pbehaviorId"
          :pbehavior-counter-type="pbehaviorCounterType"
          class="mt-4"
        />
      </v-tab-item>
    </v-tabs>
  </v-layout>
</template>

<script>
import { nextTick, ref, watch } from 'vue';

import { useAiChatExpand } from '@/hooks/ai/ai-chat-form';

import PbehaviorGeneralForm from './pbehavior-general-form.vue';
import PbehaviorPatternsForm from './pbehavior-patterns-form.vue';

const PBEHAVIOR_FORM_TABS = {
  general: 0,
  patterns: 1,
};

export default {
  inject: ['$validator'],
  components: {
    PbehaviorGeneralForm,
    PbehaviorPatternsForm,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    pbehaviorId: {
      type: String,
      required: false,
    },
    form: {
      type: Object,
      required: true,
    },
    noPattern: {
      type: Boolean,
      default: false,
    },
    noEnabled: {
      type: Boolean,
      default: false,
    },
    noComments: {
      type: Boolean,
      default: false,
    },
    noTimezone: {
      type: Boolean,
      default: false,
    },
    withStartOnTrigger: {
      type: Boolean,
      default: false,
    },
    withInherited: {
      type: Boolean,
      default: false,
    },
    nameLabel: {
      type: String,
      required: false,
    },
    nameTooltip: {
      type: String,
      required: false,
    },
    pbehaviorCounterType: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const activeTab = ref(PBEHAVIOR_FORM_TABS.general);
    const hasGeneralError = ref(false);
    const hasPatternsError = ref(false);

    const generalFormElement = ref(null);
    const patternsFormElement = ref(null);

    useAiChatExpand({ activeTab, neededTab: PBEHAVIOR_FORM_TABS.patterns });

    let stopTabErrorWatchers;

    watch(() => props.noPattern, async (noPattern) => {
      stopTabErrorWatchers?.();
      stopTabErrorWatchers = undefined;

      if (!noPattern) {
        await nextTick();

        const stopGeneral = watch(() => generalFormElement.value?.hasAnyError, (value) => {
          hasGeneralError.value = value ?? false;
        });

        const stopPatterns = watch(() => patternsFormElement.value?.hasAnyError, (value) => {
          hasPatternsError.value = value ?? false;
        });

        stopTabErrorWatchers = () => {
          stopGeneral();
          stopPatterns();
        };
      } else {
        hasGeneralError.value = false;
        hasPatternsError.value = false;
      }
    }, { immediate: true });

    return {
      activeTab,
      hasGeneralError,
      hasPatternsError,
      generalFormElement,
      patternsFormElement,
    };
  },
};
</script>
