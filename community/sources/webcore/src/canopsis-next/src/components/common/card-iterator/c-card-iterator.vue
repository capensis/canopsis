<template lang="pug">
  c-draggable-list-field.c-card-iterator(
    v-field="data",
    handle=".item-drag-handler",
    ghost-class="grey",
    :disabled="disabled",
    :group="draggableGroup"
  )
    template(v-for="(item, index) in data")
      v-card.lighten-2(:key="item[itemKey]")
        v-card-text
          slot(name="item", :item="item", :index="index")
</template>

<script>
export default {
  model: {
    prop: 'data',
    event: 'input',
  },
  props: {
    data: {
      type: Array,
      default: () => [],
    },
    itemKey: {
      type: String,
      default: '_id',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    draggableGroup: {
      type: Object,
      required: true,
    },
  },
};
</script>

<style lang="scss">
.c-card-iterator {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
