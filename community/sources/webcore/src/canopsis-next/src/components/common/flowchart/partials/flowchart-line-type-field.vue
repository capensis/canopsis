<template lang="pug">
  v-layout(row, justify-space-between, align-center)
    v-label {{ label }}
    v-flex(xs3)
      v-select.mt-0.pt-0(v-field="value", :items="types", hide-details)
        template(v-for="slotName in ['selection', 'item']", #[slotName]="{ item }")
          svg(viewBox="0 0 50 40", width="30", height="30")
            points-line-path(:points="points", :type="item", stroke="black", stroke-width="2")
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
  },
  computed: {
    points() {
      return [{ x: 1, y: 1 }, { x: 49, y: 39 }];
    },

    types() {
      return Object.values(LINE_TYPES);
    },
  },
};
</script>
