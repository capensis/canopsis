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

import GridItem from './grid-item.vue';

export function bottom(layout) {
  let max = 0;
  let bottomY;

  for (let i = 0, len = layout.length; i < len; i += 1) {
    bottomY = layout[i].y + layout[i].h;
    if (bottomY > max) max = bottomY;
  }
  return max;
}

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
