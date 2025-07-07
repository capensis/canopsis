<template>
  <c-draggable-list-field
    v-field="data"
    :disabled="disabled"
    :group="draggableGroup"
    :handle="handle"
    class="c-card-iterator-field"
    ghost-class="grey"
  >
    <template v-for="(item, index) in data">
      <v-card
        :key="item[itemKey]"
        class="lighten-2"
      >
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
};
</script>

<style lang="scss">
.c-card-iterator-field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
