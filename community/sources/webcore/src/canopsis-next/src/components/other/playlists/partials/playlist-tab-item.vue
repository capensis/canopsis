<template lang="pug">
  v-layout.tab-panel-content(xs12)
    v-flex.tab-content-block.secondary.pa-2.white--text(xs4) {{ view.group.title }}
    v-flex.tab-content-block.secondary.lighten-1.pa-2.white--text(xs4) {{ view.title }}
    v-flex.tab-content-block.pa-2.white--text(xs4) {{ tab.title }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { SCHEMA_EMBEDDED_KEY } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

const { mapGetters: entitiesMapGetters } = createNamespacedHelpers('entities');

export default {
  props: {
    tab: {
      type: Object,
      required: true,
    },
  },
  computed: {
    ...entitiesMapGetters({
      getEntityItem: 'getItem',
    }),

    view() {
      const tabWithEmbedded = this.getEntityItem(ENTITIES_TYPES.viewTab, this.tab._id, true);

      const { parents: [parent] } = tabWithEmbedded[SCHEMA_EMBEDDED_KEY];

      return this.getEntityItem(parent.type, parent.id);
    },
  },
};
</script>

<style lang="scss" scoped>
.tab-panel-content {
  height: 100%;
}

.tab-content-block {
  display: flex;
  align-items: center;
}
</style>
