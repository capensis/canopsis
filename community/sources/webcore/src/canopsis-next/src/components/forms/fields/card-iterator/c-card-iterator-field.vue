<template>
  <c-draggable-list-field
    v-field="data"
    :disabled="disabled"
    :group="draggableGroup"
    class="c-card-iterator-field"
    handle=".item-drag-handler"
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
.c-card-iterator-field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
