<template lang="pug">
  v-spacer(v-if="healthcheckPending")
  c-responsive-list.ml-4(v-else, :items="preparedEngines", item-key="name", item-value="label")
    v-tooltip(:disabled="!item.tooltip", slot-scope="{ item }", bottom)
      c-engine-chip.ma-1(slot="activator", :color="item.color") {{ item.label }}
      span {{ item.tooltip }}
</template>

<script>
import { getHealthcheckNodeColor } from '@/helpers/color';

import { healthcheckNodesMixin } from '@/mixins/healthcheck/healthcheck-nodes';
import { entitiesHealthcheckMixin } from '@/mixins/entities/healthcheck';

export default {
  mixins: [healthcheckNodesMixin, entitiesHealthcheckMixin],
  computed: {
    preparedEngines() {
      return this.engines.nodes.map(engine => ({
        ...engine,
        color: getHealthcheckNodeColor(engine),
        tooltip: this.getTooltipText(engine),
        label: this.getNodeLabel(engine.name),
      }));
    },
  },
};
</script>

<style scoped>

</style>
