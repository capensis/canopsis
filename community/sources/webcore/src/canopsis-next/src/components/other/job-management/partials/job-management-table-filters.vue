<template>
  <v-row>
    <v-col
      v-if="withRuleType"
      cols="12"
      sm="6"
      md="3"
    >
      <c-select
        :model-value="options.rule_type"
        :label="t('jobs.filterByRuleType')"
        :items="ruleTypeItems"
        item-title="text"
        item-value="value"
        clearable
        hide-details
        @update:model-value="updateRuleType"
      />
    </v-col>

    <v-col cols="12" sm="6" md="3">
      <c-select
        :model-value="options.status"
        :label="t('jobs.filterByStatus')"
        :items="statusItems"
        item-title="text"
        item-value="value"
        clearable
        hide-details
        @update:model-value="updateStatus"
      />
    </v-col>
  </v-row>
</template>

<script>
import { computed } from 'vue';

import { JOB_EXECUTION_STATUS, JOB_RULE_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    options: {
      type: Object,
      required: true,
    },
    withRuleType: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:options'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const statusItems = computed(() => Object.values(JOB_EXECUTION_STATUS).map(value => ({
        value,
        text: t(`jobs.executionStatus.${value}`),
      })),
    );
    const ruleTypeItems = computed(() => Object.values(JOB_RULE_TYPES).map(value => ({
        value,
        text: t(`jobs.types.${value}`),
      })),
    );
    const updateOptions = newOptions => emit('update:options', newOptions);
    const updateStatus = status => updateOptions({
        ...props.options,
        status: status ?? undefined,
        page: 1,
      });
    const updateRuleType = ruleType => updateOptions({
        ...props.options,
        rule_type: ruleType ?? undefined,
        page: 1,
      });

    return {
      t,
      statusItems,
      ruleTypeItems,
      updateStatus,
      updateRuleType,
    };
  },
};
</script>
