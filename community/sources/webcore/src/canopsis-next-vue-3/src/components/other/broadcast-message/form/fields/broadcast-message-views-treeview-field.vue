<template>
  <v-list class="broadcast-message-views-treeview-field pa-0" density="compact">
    <template v-for="item in items" :key="item.value">
      <c-layout class="broadcast-message-views-treeview-field__item" align="center">
        <v-checkbox-btn
          :model-value="isSelected(item.value)"
          color="primary"
          density="compact"
          @update:model-value="toggleItem(item.value)"
        />
        <span>{{ item.name }}</span>
      </c-layout>

      <broadcast-message-views-treeview-field
        v-if="item.children?.length"
        :model-value="model"
        :items="item.children"
        class="ml-6"
        @update:model-value="updateModel"
      />
    </template>
  </v-list>
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

    const isSelected = value => model.value.includes(value);

    const toggleItem = (value) => {
      const nextValue = model.value.includes(value)
        ? model.value.filter(item => item !== value)
        : [...model.value, value];

      updateModel(nextValue);
    };

    return {
      model,
      isSelected,
      updateModel,
      toggleItem,
    };
  },
};
</script>

<style lang="scss" scoped>
.broadcast-message-views-treeview-field__item {
  min-height: 32px;
}
</style>
