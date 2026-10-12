<template>
  <v-menu
    v-if="hasExecutions"
    v-model="isMenuOpen"
    :close-on-content-click="false"
    class="remediation-executions-menu"
    content-class="remediation-executions-menu__content"
    location="top"
    min-width="400"
    z-index="8"
    persistent
  >
    <template #activator="{ props: activatorProps }">
      <v-badge
        class="view-executions-btn__badge"
        color="blue-grey-lighten-1"
        offset-x="4"
        offset-y="4"
      >
        <template #badge>
          <span>{{ badgeValue }}</span>
        </template>
        <v-btn
          :loading="pending"
          :aria-label="$t('remediation.instruction.manualInstructionsProgress')"
          color="grey-darken-2"
          size="x-large"
          icon
          v-bind="activatorProps"
        >
          <v-icon>
            $vuetify.icons.manual_instruction
          </v-icon>
        </v-btn>
      </v-badge>
    </template>
    <active-remediation-executions
      :executions="executions"
      class="fill-height"
      @refresh="fetchList"
      @minimize="minimizeMenu"
    />
  </v-menu>
</template>

<script>
import { computed, onMounted, ref } from 'vue';

import { SOCKET_ROOMS } from '@/config';

import { usePendingHandler } from '@/hooks/query/pending';
import { useSocketRoom } from '@/hooks/socket';
import { useRemediationInstructionExecution } from '@/hooks/store/modules/remediation-instruction-execution';

import ActiveRemediationExecutions from '@/components/other/remediation/instruction-execute/active-remediation-executions.vue';

const MAX_BADGE_VALUE = 9;

export default {
  components: {
    ActiveRemediationExecutions,
  },
  setup() {
    const isMenuOpen = ref(false);
    const executions = ref([]);

    const { fetchRemediationInstructionExecutionStatusesListWithoutStore } = useRemediationInstructionExecution();

    const { pending, handler: fetchList } = usePendingHandler(async () => {
      executions.value = await fetchRemediationInstructionExecutionStatusesListWithoutStore() ?? [];
    });

    const badgeValue = computed(() => (
      executions.value.length > MAX_BADGE_VALUE ? `${MAX_BADGE_VALUE}+` : executions.value.length
    ));

    const hasExecutions = computed(() => executions.value.length > 0);

    const minimizeMenu = () => {
      isMenuOpen.value = false;
    };

    const updateExecutions = (data = []) => {
      executions.value = data;
    };

    useSocketRoom({ room: SOCKET_ROOMS.executions, listener: updateExecutions });

    onMounted(fetchList);

    return {
      isMenuOpen,
      executions,
      pending,
      badgeValue,
      hasExecutions,
      fetchList,
      minimizeMenu,
    };
  },
};
</script>

<style lang="scss">
.view-executions-btn__badge {
  .v-badge__badge {
    align-items: center;
    display: flex;
    font-size: 13px;
    height: 20px;
    justify-content: center;
    min-height: 20px;
    min-width: 20px;
    padding: 0;
    width: 20px;
  }
}

.remediation-executions-menu {
  z-index: 8 !important;
}

.remediation-executions-menu__content {
  max-height: 95vh;
  height: 600px;
}
</style>
