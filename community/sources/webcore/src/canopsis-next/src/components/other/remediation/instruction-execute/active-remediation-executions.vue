<template>
  <v-card>
    <v-card-text>
      <v-tabs fixed-tabs>
        <template v-for="type in types">
          <v-tab :key="type">
            {{ $t(`remediation.instruction.types.${type}`) }}
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
import { computed } from 'vue';

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
      types,
      executionsByType,
      refresh,
    };
  },
};
</script>
