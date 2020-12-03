<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.engines') }}
    engines-list(:loading="pending")
    fab-buttons(@refresh="fetchList")
</template>

<script>
import entitiesEngineRunInfoMixin from '@/mixins/entities/engine-run-info';

import FabButtons from '@/components/other/fab-buttons/fab-buttons.vue';
import EnginesList from '@/components/other/engines/exploitation/engines-list.vue';

export default {
  components: { EnginesList, FabButtons },
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
