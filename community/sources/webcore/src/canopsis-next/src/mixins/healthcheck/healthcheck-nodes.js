export const healthcheckNodesMixin = {
  methods: {
    /**
     * Check is valid healthcheck node
     *
     * @param {HealthcheckNode} engine
     * @return {boolean}
     */
    isWrongEngine(engine) {
      return !engine.is_running
        || engine.is_queue_overflown
        || engine.is_too_few_instances
        || engine.is_diff_instances_config
        || engine.is_unknown;
    },

    /**
     * Get texts by node values
     *
     * @param {HealthcheckNode} data
     * @return {string}
     */
    getTooltipText(data) {
      const statusKeys = [];
      const itemWithDefaultFlags = {
        is_running: true,
        is_queue_overflown: false,
        is_too_few_instances: false,
        is_diff_instances_config: false,
        ...data,
      };

      if (!itemWithDefaultFlags.is_running) {
        statusKeys.push('healthcheck.notRunning');
      }

      if (itemWithDefaultFlags.is_queue_overflown) {
        statusKeys.push('healthcheck.queueOverflow');
      }

      if (itemWithDefaultFlags.is_too_few_instances) {
        statusKeys.push('healthcheck.lackOfInstances');
      }

      if (itemWithDefaultFlags.is_diff_instances_config) {
        statusKeys.push('healthcheck.diffInstancesConfig');
      }

      if (itemWithDefaultFlags.is_unknown) {
        statusKeys.push('healthcheck.unknown');
      }

      return statusKeys
        .map(stateKey => this.$t(stateKey, { name: this.getNodeName(data.name) }))
        .join('\n');
    },

    /**
     * Get label for edge between nodes
     *
     * @param {string} nodeName
     * @return {string}
     */
    getNodeEdgeLabel(nodeName) {
      const nodeEdgeLabelKey = `healthcheck.nodes.${nodeName}.edgeLabel`;

      return this.$te(nodeEdgeLabelKey) ? this.$t(nodeEdgeLabelKey) : nodeName;
    },

    /**
     * Get label for node
     *
     * @param {string} nodeName
     * @return {string}
     */
    getNodeName(nodeName) {
      const nodeLabelKey = `healthcheck.nodes.${nodeName}.name`;

      return this.$te(nodeLabelKey) ? this.$t(nodeLabelKey) : nodeName;
    },
  },
};
