<template>
  <c-layout
    gap="2"
    column
  >
    <c-enabled-field
      v-model="enabledField"
      :disabled="disablable || defaultRule"
      hide-details
      with-background
    />

    <v-tabs
      v-model="activeTab"
      align-tabs="center"
      color="primary"
    >
      <v-tab value="general">
        {{ t('common.general') }}
      </v-tab>
      <v-tab
        v-if="!defaultRule"
        value="patterns"
      >
        {{ t('common.pattern', 2) }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab">
      <v-window-item value="general">
        <alarm-status-rule-general-form
          :flapping="flapping"
          :default-rule="defaultRule"
        />
      </v-window-item>

      <v-window-item
        v-if="!defaultRule"
        value="patterns"
      >
        <alarm-status-rule-patterns-form
          :flapping="flapping"
          class="mt-2"
        />
      </v-window-item>
    </v-window>
  </c-layout>
</template>

<script>
import { ref } from 'vue';

import { useFormModel } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import AlarmStatusRuleGeneralForm from './alarm-status-rule-general-form.vue';
import AlarmStatusRulePatternsForm from './alarm-status-rule-patterns-form.vue';

export default {
  components: {
    AlarmStatusRuleGeneralForm,
    AlarmStatusRulePatternsForm,
  },

  props: {
    flapping: {
      type: Boolean,
      default: false,
    },
    disablable: {
      type: Boolean,
      default: false,
    },
    defaultRule: {
      type: Boolean,
      default: false,
    },
  },

  setup() {
    const { t } = useI18n();
    const activeTab = ref('general');

    const enabledField = useFormModel('enabled');

    return {
      t,
      activeTab,
      enabledField,
    };
  },
};
</script>
