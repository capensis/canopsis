<template lang="pug">
  c-responsive-list.ml-4(:items="preparedEngines", item-key="name", item-value="label")
    v-tooltip(:disabled="!item.tooltip", slot-scope="{ item }", bottom)
      c-engine-chip.ma-1.cursor-pointer(
        slot="activator",
        :color="item.color",
        @click="redirectToHealthcheck"
      ) {{ item.label }}
      span {{ item.tooltip }}
</template>

<script>
import { isEqual, sortBy } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { COLORS } from '@/config';
import { HEALTHCHECK_SERVICES_NAMES, SOCKET_ROOMS } from '@/constants';

import { getHealthcheckNodeColor } from '@/helpers/color';

import { healthcheckNodesMixin } from '@/mixins/healthcheck/healthcheck-nodes';

const { mapActions } = createNamespacedHelpers('healthcheck');

export default {
  mixins: [healthcheckNodesMixin],
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
        name: 'admin-healthcheck',
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
