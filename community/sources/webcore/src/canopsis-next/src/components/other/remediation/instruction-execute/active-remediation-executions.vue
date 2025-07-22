<template>
  <v-card>
    <v-card-title class="primary white--text">
      <span class="text-h6 font-weight-regular">
        {{ $t('remediation.instruction.manualInstructionsProgress') }}
      </span>
    </v-card-title>
    <v-card-text>
      <v-tabs v-model="activeTab" fixed-tabs>
        <template v-for="(type, index) in types">
          <v-tab :key="type">
            <span>{{ $t(`remediation.instruction.types.${type}`) }}</span>
            <c-circle-badge
              :outlined="activeTab !== index"
              class="ml-2"
              color="primary"
            >
              {{ executionsByType[type]?.length ?? 0 }}
            </c-circle-badge>
          </v-tab>
          <v-tab-item :key="`${type}-item`">
            <v-list>
              <div
                v-if="!executionsByType[type]?.length"
                class="text-center grey--text font-italic"
              >
                {{ $t('common.noData') }}
              </div>
              <template v-for="execution in executionsByType[type]">
                <active-remediation-executions-item
                  :key="execution._id"
                  :execution="execution"
                  @refresh="refresh"
                />
                <v-divider :key="`${execution._id}-divider`" />
              </template>
            </v-list>
          </v-tab-item>
        </template>
      </v-tabs>
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

    return {
      activeTab,
      types,
      executionsByType,
      refresh,
    };
  },
};
</script>
