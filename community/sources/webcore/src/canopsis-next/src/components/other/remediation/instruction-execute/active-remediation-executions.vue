<template>
  <v-card class="fill-height">
    <v-card-title class="bg-primary text-white">
      <c-layout align="center" justify="space-between">
        <span class="text-h6 font-weight-regular">
          {{ $t('remediation.instruction.manualInstructionsProgress') }}
        </span>
        <v-btn
          size="small"
          variant="text"
          icon
          @click="minimize"
        >
          <v-icon color="white">
            minimize
          </v-icon>
        </v-btn>
      </c-layout>
    </v-card-title>
    <v-card-text>
      <v-tabs v-model="activeTab" fixed-tabs>
        <v-tab
          v-for="type in types"
          :key="type"
          :value="type"
        >
          <span>{{ $t(`remediation.instruction.types.${type}`) }}</span>
          <c-circle-badge
            :outlined="activeTab !== type"
            class="ml-2"
            color="primary"
          >
            {{ executionsByType[type]?.length ?? 0 }}
          </c-circle-badge>
        </v-tab>
      </v-tabs>
      <v-window v-model="activeTab">
        <v-window-item
          v-for="type in types"
          :key="`${type}-item`"
          :value="type"
        >
          <v-list>
            <div
              v-if="!executionsByType[type]?.length"
              class="text-center text-grey font-italic"
            >
              {{ $t('common.noData') }}
            </div>
            <template v-for="execution in executionsByType[type]" :key="execution._id">
              <active-remediation-executions-item
                :execution="execution"
                @refresh="refresh"
              />
              <v-divider />
            </template>
          </v-list>
        </v-window-item>
      </v-window>
    </v-card-text>
  </v-card>
</template>

<script>
import { computed, ref } from 'vue';

import { REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

import ActiveRemediationExecutionsItem from './partials/active-remediation-executions-item.vue';

export default {
  components: {
    ActiveRemediationExecutionsItem,
  },
  props: {
    executions: {
      type: Array,
      default: () => [],
    },
  },
  emits: ['refresh', 'minimize'],
  setup(props, { emit }) {
    const activeTab = ref(REMEDIATION_INSTRUCTION_TYPES.manual);

    const executionsByType = computed(() => props.executions.reduce((acc, execution) => {
      if (!acc[execution.type]) {
        acc[execution.type] = [];
      }

      acc[execution.type].push(execution);

      return acc;
    }, {}));

    const types = computed(() => [REMEDIATION_INSTRUCTION_TYPES.manual, REMEDIATION_INSTRUCTION_TYPES.simpleManual]);

    const refresh = () => emit('refresh');

    const minimize = () => emit('minimize');

    return {
      activeTab,
      types,
      executionsByType,
      refresh,
      minimize,
    };
  },
};
</script>
