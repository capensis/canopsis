<template lang="pug">
  div.engines
    c-progress-overlay(:pending="pending")
    c-page-header
    engines-list(:engines="engines")
    c-fab-btn(@refresh="fetchList")
</template>

<script>
import { entitiesEngineRunInfoMixin } from '@/mixins/entities/engine-run-info';

import EnginesList from '@/components/other/engines/engines-list.vue';

export default {
  components: { EnginesList },
  mixins: [entitiesEngineRunInfoMixin],
  data() {
    return {
      pending: true,
      engines: [],
    };
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    async fetchList() {
      this.pending = true;

      this.engines = await this.fetchEnginesListWithoutStore();

      this.pending = false;
    },
  },
};
</script>

<style lang="scss" scoped>
.engines {
  display: flex;
  flex-direction: column;
  height: 100%;
}
</style>
