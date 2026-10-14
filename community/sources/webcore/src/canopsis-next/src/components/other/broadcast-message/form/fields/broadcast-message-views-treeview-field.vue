<template>
  <c-treeview
    :model-value="model"
    :items="items"
    :hide-actions="!hasNestedItems"
    class="broadcast-message-views-treeview-field"
    item-key="value"
    item-title="name"
    item-children="children"
    items-registration="props"
    density="compact"
    select-strategy="classic"
    selected-color="primary"
    selectable
    @update:model-value="updateModel"
  />
</template>

<script>
import { computed } from 'vue';

export default {
  name: 'BroadcastMessageViewsTreeviewField',

  props: {
    modelValue: {
      type: Array,
      default: () => [],
    },
    value: {
      type: Array,
      default: undefined,
    },
    items: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const model = computed(() => props.modelValue ?? props.value ?? []);
    const hasNestedItems = computed(() => props.items.some(item => item.children?.length));

    const updateModel = (nextValue) => {
      emit('update:modelValue', nextValue);
    };

    return {
      model,
      hasNestedItems,
      updateModel,
    };
  },
};
</script>

<style lang="scss">
.broadcast-message-views-treeview-field {
  background: transparent;

  .v-list-item__overlay {
    display: block;
  }

  .v-list-item--active > .v-list-item__overlay {
    opacity: 0;
  }

  .v-list-item:hover > .v-list-item__overlay {
    opacity: calc(var(--v-hover-opacity) * var(--v-theme-overlay-multiplier));
  }
}
</style>
