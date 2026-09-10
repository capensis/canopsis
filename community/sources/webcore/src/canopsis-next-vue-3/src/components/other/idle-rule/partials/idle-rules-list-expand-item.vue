<template>
  <div>
    <v-tabs
      v-model="activeTab"
      bg-color="secondary-lighten-1"
      color="primary"
      align-tabs="center"
    >
      <v-tab value="summary">
        {{ t('common.summary') }}
      </v-tab>
      <v-tab value="patterns">
        {{ t('common.pattern', 2) }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab" class="pa-3 bg-secondary-lighten-2">
      <c-window-expand-item value="summary">
        <idle-rules-summary-tab :idle-rule="idleRule" />
      </c-window-expand-item>

      <c-window-expand-item value="patterns" cols="8">
        <idle-rule-patterns-form :form="patterns" :is-entity-type="isEntityType" readonly />
      </c-window-expand-item>
    </v-window>
  </div>
</template>

<script>
import { computed, ref } from 'vue';

import { PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm } from '@/helpers/entities/filter/form';
import { isIdleRuleEntityType } from '@/helpers/entities/idle-rule/form';

import { useI18n } from '@/hooks/i18n';

import IdleRulePatternsForm from '@/components/other/idle-rule/form/idle-rule-patterns-form.vue';

import IdleRulesSummaryTab from './idle-rules-summary-tab.vue';

export default {
  components: { IdleRulePatternsForm, IdleRulesSummaryTab },

  props: {
    idleRule: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const activeTab = ref('summary');

    const isEntityType = computed(() => isIdleRuleEntityType(props.idleRule.type));
    const patterns = computed(
      () => filterPatternsToForm(props.idleRule, [PATTERNS_FIELDS.entity, PATTERNS_FIELDS.alarm]),
    );

    return {
      t,
      activeTab,
      isEntityType,
      patterns,
    };
  },
};
</script>
