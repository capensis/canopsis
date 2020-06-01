<template lang="pug">
  .vue-grid-layout(ref="item", :style="mergedStyle")
    slot
    grid-item.vue-grid-placeholder(
      v-show="isDragging",
      :x="placeholder.x",
      :y="placeholder.y",
      :w="placeholder.w",
      :h="placeholder.h",
      :i="placeholder.i"
    )
</template>

<script>
import { GridLayout } from 'vue-grid-layout';
import { uniq } from 'lodash';

import { bottom } from '../helpers/bottom';
import GridItem from './grid-item.vue';

export default {
  components: { GridItem },
  extends: GridLayout,
  methods: {
    getRowCount() {
      return uniq(this.layout.map(({ y }) => y)).length;
    },

    containerHeight() {
      if (!this.autoSize) {
        return '';
      }
      const [, marginY] = this.margin;

      return `${(bottom(this.layout) * this.rowHeight) + (marginY * this.getRowCount())}px`;
    },
  },
};
</script>
