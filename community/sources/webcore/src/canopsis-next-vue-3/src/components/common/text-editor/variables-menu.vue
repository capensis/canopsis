<template>
  <v-menu
    :model-value="visible"
    :target="[positionX, positionY]"
    :close-on-content-click="false"
    max-height="300"
    location="bottom start"
    @update:model-value="$emit('close')"
  >
    <variables-list
      :items="items"
      :return-object="returnObject"
      :children-key="childrenKey"
      :value="value"
      :dense="dense"
      :show-value="showValue"
      :hide-empty-value="hideEmptyValue"
      :clickable-parent="clickableParent"
      @update:value="$emit('update:value', $event)"
    />
  </v-menu>
</template>

<script>
import VariablesList from './variables-list.vue';

export default {
  components: { VariablesList },

  props: {
    visible: {
      type: Boolean,
      default: false,
    },
    value: {
      type: String,
      default: '',
    },
    positionX: {
      type: Number,
      required: false,
    },
    positionY: {
      type: Number,
      required: false,
    },
    items: {
      type: Array,
      default: () => [],
    },
    showValue: {
      type: Boolean,
      default: false,
    },
    hideEmptyValue: {
      type: Boolean,
      default: false,
    },
    ignoreClickOutside: {
      type: Boolean,
      default: false,
    },
    dense: {
      type: Boolean,
      default: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
    clickableParent: {
      type: Boolean,
      default: false,
    },
    childrenKey: {
      type: String,
      default: 'variables',
    },
  },

  emits: ['update:value', 'close'],
};
</script>
