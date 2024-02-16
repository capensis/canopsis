<template>
  <c-draggable-list-field
    class="c-card-iterator-field"
    v-field="data"
    handle=".item-drag-handler"
    ghost-class="grey"
    :disabled="disabled"
    :group="draggableGroup"
  >
    <template v-for="(item, index) in data">
      <v-card
        class="lighten-2"
        :key="item[itemKey]"
      >
        <v-card-text>
          <slot
            name="item"
            :item="item"
            :index="index"
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
