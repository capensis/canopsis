<template>
  <c-responsive-list
    :items="preparedEngines"
    class="ml-4"
    item-key="name"
    item-value="label"
  >
    <template #default="{ item }">
      <v-tooltip
        :disabled="!item.tooltip"
        bottom
      >
        <template #activator="{ on }">
          <c-engine-chip
            :color="item.color"
            class="ma-1"
            v-on="{ ...chipListeners, ...on }"
          >
            {{ item.label }}
          </c-engine-chip>
        </template>
        <span>{{ item.tooltip }}</span>
      </v-tooltip>
    </template>
  </c-responsive-list>
</template>

<script>
import { isEqual, sortBy } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { COLORS, SOCKET_ROOMS } from '@/config';
import { HEALTHCHECK_SERVICES_NAMES, ROUTES_NAMES, USERS_PERMISSIONS } from '@/constants';

import { getHealthcheckNodeColor } from '@/helpers/entities/healthcheck/color';

import { authMixin } from '@/mixins/auth';
import { healthcheckNodesMixin } from '@/mixins/healthcheck/healthcheck-nodes';

const { mapActions } = createNamespacedHelpers('healthcheck');

export default {
  mixins: [authMixin, healthcheckNodesMixin],
  data() {
    return {
      hasServerError: false,
      data: {
        services: [],
        engines: [],
      },
    };
  },
  computed: {
    chipListeners() {
      if (this.checkAccess(USERS_PERMISSIONS.technical.healthcheck)) {
        return {
          click: this.redirectToHealthcheck,
        };
      }

      return {};
    },

    preparedEngines() {
      const wrongNodes = [...this.data.services, ...this.data.engines];

      if (this.hasServerError) {
        return [{
          name: 'error',
          color: COLORS.healthcheck.error,
          tooltip: this.$t('healthcheck.notRunning', {
            name: this.getNodeName(HEALTHCHECK_SERVICES_NAMES.healthcheck),
          }),
          label: this.$tc('common.error'),
        }];
      }

      if (!wrongNodes.length) {
        return [{
          name: 'ok',
          tooltip: this.$t('healthcheck.systemsOperational'),
          label: this.$t('common.ok'),
        }];
      }

      return sortBy(wrongNodes, ['name']).map(engine => ({
        ...engine,

        color: getHealthcheckNodeColor(engine),
        tooltip: this.getTooltipText(engine),
        label: this.getNodeName(engine.name),
      }));
    },
  },
  mounted() {
    this.fetchList();

    this.$socket
      .join(SOCKET_ROOMS.healthcheckStatus)
      .addListener(this.setHealthcheckStatus);
  },
  beforeDestroy() {
    this.$socket
      .leave(SOCKET_ROOMS.healthcheckStatus)
      .removeListener(this.setHealthcheckStatus);
  },
  methods: {
    ...mapActions({
      fetchHealthcheckStatusWithoutStore: 'fetchStatusWithoutStore',
    }),

    redirectToHealthcheck() {
      this.$router.push({
        name: ROUTES_NAMES.adminHealthcheck,
      });
    },

    setHealthcheckStatus(data) {
      if (!isEqual(data, this.data)) {
        this.data = data;
      }
    },

    async fetchList() {
      try {
        const response = await this.fetchHealthcheckStatusWithoutStore();

        this.setHealthcheckStatus(response);
      } catch (err) {
        this.hasServerError = true;
      }
    },
  },
};
</script>
