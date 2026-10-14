<template>
  <div>
    <v-tabs
      v-model="activeTab"
      bg-color="secondary-lighten-1"
      color="primary"
      align-tabs="center"
    >
      <v-tab value="patterns">
        {{ t('common.pattern', 2) }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab" class="pa-3 bg-secondary-lighten-2">
      <c-window-expand-item value="patterns" cols="8">
        <meta-alarm-rule-patterns-form
          :form="patterns"
          :with-total-entity="withTotalEntityPattern"
          readonly
        />
      </c-window-expand-item>
    </v-window>
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
    const withTotalEntityPattern = computed(() => isMetaAlarmRuleTypeHasTotalEntityPatterns(props.metaAlarmRule.type),
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
