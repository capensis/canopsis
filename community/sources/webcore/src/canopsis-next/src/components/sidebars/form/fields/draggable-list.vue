<template lang="pug">
  c-draggable-list-field(v-field="items", :handle="`.${handleClass}`")
    v-layout(v-for="(item, index) in items", :key="item[itemKey]", row, align-center)
      v-flex(xs1)
        v-icon.draggable(:class="handleClass") drag_indicator
      v-flex(xs8)
        slot(name="title", :item="item", :index="index")
          span {{ item[itemText] }}
      v-flex(xs3)
        c-action-btn(type="edit", @click="edit(item, index)")
        c-action-btn(type="delete", @click="remove(item, index)")
</template>

<script>
export default {
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
    itemText: {
      type: String,
      default: 'text',
    },
  },
  computed: {
    handleClass() {
      return 'drag-handler';
    },
  },
  methods: {
    edit(item, index) {
      this.$emit('edit', item, index);
    },

    remove(item, index) {
      this.$emit('remove', item, index);
    },
  },
};
</script>
