<template>
  <div>
    <slot name="prepend" />
    <movable-card-iterator-item
      v-for="(item, index) in items"
      :key="item[itemKey]"
      :disabled-up="index === 0"
      :disabled-down="index === items.length - 1"
      class="my-2"
      @up="up(index)"
      @down="down(index)"
      @remove="removeItemFromArray(index)"
    >
      <template>
        <slot
          :item="item"
          :index="index"
          name="item"
        />
      </template>
    </movable-card-iterator-item>
    <slot name="append" />
    <v-layout wrap>
      <v-btn
        v-if="addable"
        class="mr-2 mx-0"
        color="primary"
        outlined
        @click.prevent="$emit('add')"
      >
        {{ $t('common.add') }}
      </v-btn>
      <slot name="actions" />
    </v-layout>
  </div>
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
    addable: {
      type: Boolean,
      default: false,
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
