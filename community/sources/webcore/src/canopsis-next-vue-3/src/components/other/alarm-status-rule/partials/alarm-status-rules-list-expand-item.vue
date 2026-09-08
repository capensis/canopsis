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
      <v-window-item value="description">
        <c-layout class="py-3">
          <c-textarea
            :model-value="rule.description"
            class="my-2 mx-4 pa-0"
            readonly
            auto-grow
            hide-details
          />
        </c-layout>
      </v-window-item>

      <v-window-item value="patterns">
        <v-row class="py-3" justify="center">
          <v-col cols="8">
            <v-card>
              <v-card-text>
                <alarm-status-rule-patterns-form :form="patterns" :flapping="flapping" readonly />
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-window-item>
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
