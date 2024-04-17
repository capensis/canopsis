<template>
  <v-menu
    ref="menu"
    :value="visible"
    :position-x="positionX"
    :position-y="positionY"
    :close-on-content-click="false"
    :ignore-click-outside="ignoreClickOutside"
    max-height="300"
    @input="$emit('close')"
  >
    <variables-list
      :variables="variables"
      :value="value"
      :dense="dense"
      :z-index="submenuZIndex"
      :show-value="showValue"
      @input="$emit('input', $event)"
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
    variables: {
      type: Array,
      default: () => [],
    },
    showValue: {
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
  },
  data() {
    return {
      submenuZIndex: undefined,
    };
  },
  mounted() {
    this.$watch(() => this.$refs.menu.activeZIndex, (zIndex) => {
      this.submenuZIndex = zIndex + 1;
    });
  },
};
</script>
