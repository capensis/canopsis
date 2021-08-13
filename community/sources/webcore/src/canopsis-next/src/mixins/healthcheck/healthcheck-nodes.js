export const healthcheckNodesMixin = {
  methods: {
    /**
     * Get texts by node values
     *
     * @param {HealthcheckNode} data
     * @return {string}
     */
    getTooltipText(data) {
      const itemWithDefaultFlags = {
        is_running: true,
        is_queue_overflown: false,
        is_too_few_instances: false,
        is_diff_instances_config: false,
        ...data,
      };
      const statusKeys = [];

      if (!itemWithDefaultFlags.is_running) {
        statusKeys.push('healthcheck.statuses.notRunning');
      }

      if (itemWithDefaultFlags.is_queue_overflown) {
        statusKeys.push('healthcheck.statuses.queueOverflow');
      }

      if (itemWithDefaultFlags.is_too_few_instances) {
        statusKeys.push('healthcheck.statuses.tooFewInstances');
      }

      if (itemWithDefaultFlags.is_diff_instances_config) {
        statusKeys.push('healthcheck.statuses.diffInstancesConfig');
      }

      return statusKeys
        .map(stateKey => this.$t(stateKey, { engine: this.getNodeLabel(data.id) }))
        .join('\n');
    },

    /**
     * Get label for edge between nodes
     *
     * @param {string} nodeName
     * @return {string}
     */
    getEdgeLabel(nodeName) {
      const engineEdgeLabelKey = `healthcheck.engines.${nodeName}.edgeLabel`;
      const serviceEdgeLabelKey = `healthcheck.services.${nodeName}.edgeLabel`;

      if (this.$te(engineEdgeLabelKey)) {
        return this.$t(engineEdgeLabelKey);
      }

      return this.$te(serviceEdgeLabelKey) ? this.$t(serviceEdgeLabelKey) : nodeName;
    },


    /**
     * Get label for node
     *
     * @param {string} nodeName
     * @return {string}
     */
    getNodeLabel(nodeName) {
      const engineLabelKey = `healthcheck.engines.${nodeName}.label`;
      const serviceLabelKey = `healthcheck.services.${nodeName}.label`;

      if (this.$te(engineLabelKey)) {
        return this.$t(engineLabelKey);
      }

      return this.$te(serviceLabelKey) ? this.$t(serviceLabelKey) : nodeName;
    },
  },
};
