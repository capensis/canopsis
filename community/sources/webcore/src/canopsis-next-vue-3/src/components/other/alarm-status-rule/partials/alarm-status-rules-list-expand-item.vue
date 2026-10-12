<template>
  <div>
    <v-tabs
      v-model="activeTab"
      bg-color="secondary-lighten-1"
      color="primary"
      align-tabs="center"
    >
      <v-tab value="description">
        {{ t('common.description') }}
      </v-tab>
      <v-tab value="patterns">
        {{ t('common.pattern', 2) }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab">
      <c-window-expand-item value="description">
        <c-textarea
          :model-value="rule.description"
          class="my-2 mx-4 pa-0"
          readonly
          auto-grow
          hide-details
        />
      </c-window-expand-item>

      <c-window-expand-item value="patterns" cols="8">
        <alarm-status-rule-patterns-form :form="patterns" :flapping="flapping" readonly />
      </c-window-expand-item>
    </v-window>
  </div>
</template>

<script>
import { computed, ref } from 'vue';

import { PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm } from '@/helpers/entities/filter/form';

import { useI18n } from '@/hooks/i18n';

import AlarmStatusRulePatternsForm from '../form/alarm-status-rule-patterns-form.vue';

export default {
  components: {
    AlarmStatusRulePatternsForm,
  },

  props: {
    rule: {
      type: Object,
      required: true,
    },
    flapping: {
      type: Boolean,
      default: false,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const activeTab = ref('description');
    const patterns = computed(() => filterPatternsToForm(props.rule, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]),
    );

    return {
      t,
      activeTab,
      patterns,
    };
  },
};
</script>
