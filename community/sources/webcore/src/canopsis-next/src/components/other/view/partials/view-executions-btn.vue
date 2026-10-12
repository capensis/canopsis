<template>
  <v-menu
    v-if="hasExecutions"
    v-model="isMenuOpen"
    ref="menuEl"
    :close-on-content-click="false"
    content-class="remediation-executions-menu__content"
    min-width="400"
    offset-overflow
    ignore-click-outside
    top
    offset-y
  >
    <template #activator="{ on }">
      <v-btn
        :loading="pending"
        color="grey darken-2"
        fab
        dark
        v-on="on"
      >
        <v-badge color="blue-grey lighten-1">
          <template #badge>
            <span>{{ badgeValue }}</span>
          </template>
          <v-icon>$vuetify.icons.manual_instruction</v-icon>
        </v-badge>
      </v-btn>
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
import { ref, computed, onMounted } from 'vue';

import { SOCKET_ROOMS } from '@/config';

import { useSocketRoom } from '@/hooks/socket';
import { usePendingHandler } from '@/hooks/query/pending';
import { useRemediationInstructionExecution } from '@/hooks/store/modules/remediation-instruction-execution';

import ActiveRemediationExecutions from '@/components/other/remediation/instruction-execute/active-remediation-executions.vue';

const MAX_BADGE_VALUE = 9;

export default {
  components: {
    ActiveRemediationExecutions,
  },
  setup() {
    const menuEl = ref(null);
    const isMenuOpen = ref(false);
    const executions = ref([]);

    const { fetchExecutionsStatusesWithoutStore } = useRemediationInstructionExecution();

    const { pending, handler: fetchList } = usePendingHandler(async () => {
      executions.value = await fetchExecutionsStatusesWithoutStore() ?? [];
    });

    const badgeValue = computed(() => (
      executions.value.length > MAX_BADGE_VALUE ? `${MAX_BADGE_VALUE}+` : executions.value.length
    ));

    const hasExecutions = computed(() => executions.value.length > 0);

    /**
     * Closes the executions menu by setting isMenuOpen to false
     */
    const minimizeMenu = () => isMenuOpen.value = false;

    useSocketRoom({ room: SOCKET_ROOMS.executions, listener: (data = []) => executions.value = data });

    onMounted(fetchList);

    return {
      menuEl,
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
.remediation-executions-menu__content {
  max-height: 95vh;
  height: 600px;
}
</style>
