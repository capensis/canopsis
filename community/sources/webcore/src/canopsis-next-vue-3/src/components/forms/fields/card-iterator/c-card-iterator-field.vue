<template>
  <c-draggable-list-field
    :model-value="data"
    :disabled="disabled"
    :group="draggableGroup"
    :handle="handle"
    :item-key="itemKey"
    class="c-card-iterator-field"
    ghost-class="grey"
    @update:model-value="$emit('update:data', $event)"
  >
    <template #item="{ element: item, index }">
      <v-card class="lighten-2">
        <v-card-text>
          <slot
            :item="item"
            :index="index"
            :handle="handle"
            name="item"
          />
        </v-card-text>
      </v-card>
    </template>
  </c-draggable-list-field>
</template>

<script>
export default {
  props: {
    data: {
      type: Array,
      default: () => [],
    },
    itemKey: {
      type: String,
      default: '_id',
    },
    handle: {
      type: String,
      default: '.item-drag-handler',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    draggableGroup: {
      type: Object,
      required: false,
    },
  },
  emits: ['update:data'],
};
</script>

<style lang="scss">
.c-card-iterator-field {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
