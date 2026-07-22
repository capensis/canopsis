<template>
  <div class="anomaly-monitored-connector-history-graph-placeholder text-body-1 text--secondary">
    {{ $t('healthcheck.connectorsBlocks.graphNotEnoughData') }}
  </div>
</template>

<script>
export default {};
</script>

<style lang="scss" scoped>
.anomaly-monitored-connector-history-graph-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 300px;
  text-align: center;
}
</style>
