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

    <c-layout class="pa-3 secondary lighten-2">
      <v-card class="pa-3 idle-rules-list-expand-item__card">
        <v-window v-model="activeTab" class="pt-2">
          <v-window-item value="summary">
            <div class="idle-rules-list-expand-item__content">
              <idle-rules-summary-tab :idle-rule="idleRule" />
            </div>
          </v-window-item>

          <v-window-item value="patterns">
            <div class="idle-rules-list-expand-item__content">
              <idle-rule-patterns-form :form="patterns" :is-entity-type="isEntityType" readonly />
            </div>
          </v-window-item>
        </v-window>
      </v-card>
    </c-layout>
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

<style lang="scss" scoped>
.idle-rules-list-expand-item__card {
  width: 100%;
}

.idle-rules-list-expand-item__content {
  width: min(100%, 900px);
  margin: 0 auto;
}
</style>
