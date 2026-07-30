<template>
  <div>
    <v-tabs
v-model="activeTab"
bg-color="secondary-lighten-1"
color="primary"
align-tabs="center">
      <v-tab value="patterns">
        {{ t('common.pattern', 2) }}
      </v-tab>
    </v-tabs>

    <c-layout class="pa-3 secondary lighten-2">
      <v-card class="pa-3 meta-alarm-rule-list-expand-panel__card">
        <v-window v-model="activeTab" class="pt-2">
          <v-window-item value="patterns">
            <div class="meta-alarm-rule-list-expand-panel__content">
              <meta-alarm-rule-patterns-form
                :form="patterns"
                :with-total-entity="withTotalEntityPattern"
                readonly
              />
            </div>
          </v-window-item>
        </v-window>
      </v-card>
    </c-layout>
  </div>
</template>

<script>
import { computed, ref } from 'vue';

import {
  isMetaAlarmRuleTypeHasTotalEntityPatterns,
  metaAlarmFilterPatternsToForm,
} from '@/helpers/entities/meta-alarm/rule/form';

import { useI18n } from '@/hooks/i18n';

import MetaAlarmRulePatternsForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-patterns-form.vue';

export default {
  components: { MetaAlarmRulePatternsForm },

  props: {
    metaAlarmRule: {
      type: Object,
      default: () => ({}),
    },
  },

  setup(props) {
    const { t } = useI18n();
    const activeTab = ref('patterns');

    const patterns = computed(() => metaAlarmFilterPatternsToForm(props.metaAlarmRule));
    const withTotalEntityPattern = computed(() =>
      isMetaAlarmRuleTypeHasTotalEntityPatterns(props.metaAlarmRule.type),
    );

    return {
      t,
      activeTab,
      patterns,
      withTotalEntityPattern,
    };
  },
};
</script>

<style lang="scss" scoped>
.meta-alarm-rule-list-expand-panel__card {
  width: 100%;
}

.meta-alarm-rule-list-expand-panel__content {
  width: min(100%, 900px);
  margin: 0 auto;
}
</style>
