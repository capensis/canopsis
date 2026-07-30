<template>
  <v-chip
    :class="chipClass"
    :color="color"
    :variant="outlined ? 'outlined' : undefined"
    class="c-chip"
    size="small"
  >
    <span :style="textStyle" class="c-chip__text">
      <slot />
    </span>
    <v-icon
      v-if="closable"
      class="cursor-pointer ml-2"
      size="small"
      @click.stop="$emit('close', $event)"
    >
      cancel
    </v-icon>
  </v-chip>
</template>

<script>
export default {
  emits: ['close'],
  props: {
    color: {
      type: String,
      required: false,
    },
    closable: {
      type: Boolean,
      default: false,
    },
    small: {
      type: Boolean,
      default: false,
    },
    outlined: {
      type: Boolean,
      default: false,
    },
    rounded: {
      type: Boolean,
      default: false,
    },
    textColor: {
      type: String,
      default: 'white',
    },
  },
  computed: {
    chipClass() {
      return {
        'c-chip--closable': this.closable,
        'c-chip--small': this.small,
        'c-chip--rounded': this.rounded,
      };
    },

    textStyle() {
      return {
        color: this.textColor,
      };
    },
  },
};
</script>

<style lang="scss">
.c-chip.v-chip {
  border-radius: 5px;
  font-size: 12px;
  min-height: 24px;
  height: unset !important;
  padding: 0;

  &.c-chip--rounded {
    border-radius: 20px;
  }

  &__text {
    white-space: initial;
    word-wrap: break-word;
    max-width: 100%;
    overflow: hidden;
  }

  .v-chip__content {
    height: unset !important;
    cursor: pointer;
    max-width: 100%;
  }

  &--closable {
    .v-chip__content {
      padding-right: 4px;
    }
  }

  &--small {
    min-height: 20px !important;
    margin: 2px;
  }

  .v-data-table thead th.column.sortable & .v-icon {
    opacity: 0.6;
  }

  .v-icon {
    transform: unset !important;
  }
}
</style>
