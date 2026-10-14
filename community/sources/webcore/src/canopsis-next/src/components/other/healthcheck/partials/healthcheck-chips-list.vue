<template>
  <c-responsive-list
    :items="preparedEngines"
    class="healthcheck-chips-list ml-4"
    item-key="name"
    item-value="label"
  >
    <template #default="{ item }">
      <c-tooltip :disabled="!item.tooltip" location="bottom">
        <template #activator="{ props }">
          <c-engine-chip
            :color="item.color"
            :clickable="hasAccessToHealthcheck"
            class="ma-1"
            v-bind="props"
            @click="redirectToHealthcheck"
          >
            {{ item.label }}
          </c-engine-chip>
        </template>

        <span>{{ item.tooltip }}</span>
      </c-tooltip>
    </template>
  </c-responsive-list>
</template>

<script>
import { isEqual, sortBy } from 'lodash';
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

import { COLORS, SOCKET_ROOMS } from '@/config';
import { HEALTHCHECK_SERVICES_NAMES, ROUTES_NAMES, USER_PERMISSIONS } from '@/constants';

import { getHealthcheckNodeColor } from '@/helpers/entities/healthcheck/color';

import { useCanPermission } from '@/hooks/auth';
import { useHealthcheckNodes } from '@/hooks/healthcheck-nodes';
import { useI18n } from '@/hooks/i18n';
import { useSocketRoom } from '@/hooks/socket';
import { useHealthcheck } from '@/hooks/store/modules/healthcheck';

export default {
  setup() {
    const router = useRouter();
    const { t } = useI18n();
    const { getNodeName, getTooltipText } = useHealthcheckNodes();
    const { fetchHealthcheckStatusWithoutStore } = useHealthcheck();
    const { hasAccess: hasAccessToHealthcheck } = useCanPermission(
      USER_PERMISSIONS.technical.healthcheck,
    );

    const hasServerError = ref(false);
    const wrongNodes = ref([]);

    const setHealthcheckStatus = ({ services = [], engines = [] } = {}) => {
      const newWrongNodes = sortBy([...services, ...engines], ['name']);

      if (!isEqual(newWrongNodes, wrongNodes.value)) {
        wrongNodes.value = newWrongNodes;
      }
    };

    const preparedEngines = computed(() => {
      if (hasServerError.value) {
        return [
          {
            name: 'error',
            color: COLORS.healthcheck.error,
            tooltip: t('healthcheck.notRunning', {
              name: getNodeName(HEALTHCHECK_SERVICES_NAMES.healthcheck),
            }),
            label: t('common.error'),
          },
        ];
      }

      if (!wrongNodes.value.length) {
        return [
          {
            name: 'ok',
            tooltip: t('healthcheck.systemsOperational'),
            label: t('common.ok'),
          },
        ];
      }

      return wrongNodes.value.map(engine => ({
        ...engine,
        color: getHealthcheckNodeColor(engine),
        tooltip: getTooltipText(engine),
        label: getNodeName(engine.name),
      }));
    });

    const fetchList = async () => {
      try {
        const response = await fetchHealthcheckStatusWithoutStore();

        hasServerError.value = false;
        setHealthcheckStatus(response);
      } catch {
        hasServerError.value = true;
      }
    };

    const redirectToHealthcheck = () => {
      if (hasAccessToHealthcheck.value && router.hasRoute(ROUTES_NAMES.adminHealthcheck)) {
        router.push({ name: ROUTES_NAMES.adminHealthcheck });
      }
    };

    useSocketRoom({ room: SOCKET_ROOMS.healthcheckStatus, listener: setHealthcheckStatus });

    onMounted(fetchList);
    onBeforeUnmount(() => (wrongNodes.value = []));

    return {
      hasAccessToHealthcheck,
      preparedEngines,
      redirectToHealthcheck,
    };
  },
};
</script>

<style lang="scss" scoped>
.healthcheck-chips-list {
  max-width: 420px;
  min-width: 0;
}
</style>
