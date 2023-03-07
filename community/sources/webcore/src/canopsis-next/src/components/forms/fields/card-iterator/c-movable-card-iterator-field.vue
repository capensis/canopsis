<template lang="pug">
  div
    movable-card.my-2(
      v-for="(item, index) in items",
      :key="item[itemKey]",
      :disabled-up="index === 0",
      :disabled-down="index === items.length - 1",
      @up="up(index)",
      @down="down(index)",
      @remove="removeItemFromArray(index)"
    )
      template
        slot(name="item", :item="item", :index="index")
    v-btn.mx-0(color="primary", @click.prevent="$emit('add')") {{ $t('common.add') }}
</template>

<script>
import { formArrayMixin } from '@/mixins/form';

import MovableCard from './movable-card.vue';

export default {
  components: { MovableCard },
  mixins: [formArrayMixin],
  model: {
    prop: 'items',
    event: 'input',
  },
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    itemKey: {
      type: String,
      default: 'key',
    },
  },
  methods: {
    up(index) {
      if (index > 0) {
        const items = [...this.items];
        const temp = items[index];

        items[index] = items[index - 1];
        items[index - 1] = temp;

        this.updateModel(items);
      }
    },

    down(index) {
      if (index < this.items.length - 1) {
        const items = [...this.items];
        const temp = items[index];

        items[index] = items[index + 1];
        items[index + 1] = temp;

        this.updateModel(items);
      }
    },
  },
};
</script>
