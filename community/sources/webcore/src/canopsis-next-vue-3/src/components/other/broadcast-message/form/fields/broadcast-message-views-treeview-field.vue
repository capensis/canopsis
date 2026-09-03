<template>
  <c-treeview
    :model-value="model"
    :items="items"
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

    const updateModel = (nextValue) => {
      emit('update:modelValue', nextValue);
    };

    return {
      model,
      updateModel,
    };
  },
};
</script>

<style lang="scss" scoped>
.broadcast-message-views-treeview-field {
  :deep(.v-list-item__overlay) {
    display: none;
  }
}
</style>
