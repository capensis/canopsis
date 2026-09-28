<template>
  <v-treeview
    v-field="value"
    :items="items"
    selected-color="primary"
    item-key="value"
    item-children="children"
    selectable
    open-on-click
    transition
  >
    <template #label="{ item }">
      <span>{{ item.name }}</span>
    </template>
  </v-treeview>
</template>

<script>
export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
      default: () => [],
    },
    items: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
