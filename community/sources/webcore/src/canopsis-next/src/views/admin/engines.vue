<template lang="pug">
  v-container
    c-page-header {{ $t('common.engines') }}
    engines-list(:loading="pending", :engines="engines")
    c-fab-btn(@refresh="fetchList")
</template>

<script>
import entitiesEngineRunInfoMixin from '@/mixins/entities/engine-run-info';

import EnginesList from '@/components/other/engines/exploitation/engines-list.vue';

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
