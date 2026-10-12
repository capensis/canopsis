<template>
  <c-tooltip :max-width="maxWidth" :location="location">
    <template #activator="{ props: activatorProps }">
      <c-icon
        v-bind="activatorProps"
        :color="color"
        :size="preparedSize"
        :class="preparedIconClass"
        :no-fill="isOutlinedIcon"
      >
        {{ icon }}
      </c-icon>
    </template>

    <slot>
      <div v-html="text" class="pre-wrap" />
    </slot>
  </c-tooltip>
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
      type: [String, Number],
      required: false,
    },
    small: {
      type: Boolean,
      required: false,
    },
    iconClass: {
      type: [String, Object],
      default: '',
    },
  },

  computed: {
    location() {
      return {
        [this.left]: 'left',
        [this.bottom]: 'bottom',
        [this.right]: 'right',
        [this.top]: 'top',
      }.true;
    },

    preparedSize() {
      return this.small ? 16 : this.size;
    },

    preparedIconClass() {
      return this.iconClass || (this.color ? '' : 'text-medium-emphasis');
    },

    isOutlinedIcon() {
      return this.icon.endsWith('_outline');
    },
  },
};
</script>
