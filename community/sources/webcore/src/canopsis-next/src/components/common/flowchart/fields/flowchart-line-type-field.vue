<template lang="pug">
  v-layout(row, justify-space-between, align-center)
    v-label {{ label }}
    v-flex(xs3)
      v-select.mt-0.pt-0(v-field="value", :items="types", hide-details)
        template(v-for="slotName in ['selection', 'item']", #[slotName]="{ item }")
          svg(viewBox="0 0 50 40", width="30", height="30")
            points-line-path(:points="points", :type="item", stroke="currentColor", stroke-width="2")
</template>

<script>
import { LINE_TYPES } from '@/constants';

import PointsLinePath from '@/components/common/flowchart/common/points-line-path.vue';

export default {
  components: { PointsLinePath },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      required: false,
    },
    label: {
      type: String,
      required: false,
    },
    averagePoints: {
      type: Array,
      required: false,
    },
  },
  computed: {
    points() {
      const maxY = 39;
      const maxX = 49;

      const startPoint = { x: 1, y: 1 };
      const endPoint = { x: maxX, y: maxY };

      if (this.averagePoints) {
        const [start, end] = this.averagePoints;

        if (start.y > end.y) {
          startPoint.y = maxY;
          endPoint.y = 1;
        }

        if (start.x > end.x) {
          startPoint.x = maxX;
          endPoint.x = 1;
        }
      }

      return [startPoint, endPoint];
    },

    types() {
      return Object.values(LINE_TYPES);
    },
  },
};
</script>
