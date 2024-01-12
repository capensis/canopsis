<template>
  <v-tooltip
    :max-width="maxWidth"
    :top="top"
    :right="right"
    :bottom="bottom"
    :left="left"
  >
    <template #activator="{ bind, on }">
      <v-icon
        v-bind="bind"
        v-on="on"
        :color="color"
        :size="size"
        :class="iconClass"
      >
        {{ icon }}
      </v-icon>
    </template>
    <slot>
      <div
        class="pre-wrap"
        v-html="text"
      />
    </slot>
  </v-tooltip>
</template>

<script>
export default {
  props: {
    icon: {
      type: String,
      default: 'help_outline',
    },
    color: {
      type: String,
      default: '',
    },
    text: {
      type: String,
      required: false,
    },
    maxWidth: {
      type: [String, Number],
      required: false,
    },
    top: {
      type: Boolean,
      required: false,
    },
    right: {
      type: Boolean,
      required: false,
    },
    bottom: {
      type: Boolean,
      required: false,
    },
    left: {
      type: Boolean,
      required: false,
    },
    size: {
      type: String,
      required: false,
    },
    iconClass: {
      type: [String, Object],
      required: false,
    },
  },
};
</script>
