<template lang="pug">
  c-responsive-list.ml-4(:items="preparedEngines", item-key="name", item-value="label")
    v-tooltip(:disabled="!item.tooltip", slot-scope="{ item }", bottom)
      c-engine-chip.ma-1(slot="activator", :color="item.color") {{ item.label }}
      span {{ item.tooltip }}
</template>

<script>
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
      services: [],
      engines: [],
    };
  },
  computed: {
    preparedEngines() {
      const wrongNodes = [...this.services, ...this.engines];

      if (!wrongNodes.length) {
        return [{
          name: 'ok',
          tooltip: this.$t('healthcheck.systemsOperational'),
          label: this.$t('common.ok'),
        }];
      }

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

      return wrongNodes.map(engine => ({
        ...engine,
        color: getHealthcheckNodeColor(engine),
        tooltip: this.getTooltipText(engine),
        label: this.getNodeName(engine.name),
      }));
    },
  },
  mounted() {
    this.fetchList();

    this.$socket.join(SOCKET_ROOMS.healthcheckStatus);
    this.$socket
      .getRoom(SOCKET_ROOMS.healthcheckStatus)
      .addListener(({ engines, services }) => {
        this.engines = engines;
        this.services = services;
      });
  },
  beforeDestroy() {
    this.$socket.leave(SOCKET_ROOMS.healthcheckStatus);
  },
  methods: {
    ...mapActions({
      fetchHealthcheckStatusWithoutStore: 'fetchStatusWithoutStore',
    }),

    async fetchList() {
      try {
        const { services, engines } = await this.fetchHealthcheckStatusWithoutStore();

        this.services = services.filter(this.isWrongEngine);
        this.engines = engines.graph.nodes.reduce((acc, name) => {
          const parameters = engines.parameters[name];

          if (this.isWrongEngine(parameters)) {
            acc.push({ name, ...parameters });
          }

          return acc;
        }, []);
      } catch (err) {
        this.hasServerError = true;
      }
    },
  },
};
</script>
