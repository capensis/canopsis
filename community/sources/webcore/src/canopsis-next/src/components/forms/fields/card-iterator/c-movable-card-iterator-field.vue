<template lang="pug">
  div
    movable-card-iterator-item.my-2(
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

import MovableCardIteratorItem from './movable-card-iterator-item.vue';

export default {
  components: { MovableCardIteratorItem },
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
    moveItemByOffset(index, offset) {
      const items = [...this.items];
      const temp = items[index];

      items[index] = items[index + offset];
      items[index + offset] = temp;

      this.updateModel(items);
    },

    up(index) {
      if (index > 0) {
        this.moveItemByOffset(index, -1);
      }
    },

    down(index) {
      if (index < this.items.length - 1) {
        this.moveItemByOffset(index, 1);
      }
    },
  },
};
</script>
