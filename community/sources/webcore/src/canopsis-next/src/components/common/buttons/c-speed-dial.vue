<template>
  <v-speed-dial
    v-bind="$attrs"
    :value="internalValue"
    class="c-speed-dial"
    @input="input"
  >
    <template #activator="">
      <slot
        name="activator"
        :bind="bind"
      />
    </template>
    <slot :value="value" />
  </v-speed-dial>
</template>

<script>
export default {
  inheritAttrs: false,
  props: {
    value: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      internalValue: this.value,
    };
  },
  computed: {
    bind() {
      return {
        inputValue: this.internalValue,
      };
    },
  },
  watch: {
    value(value) {
      this.internalValue = value;
    },
  },
  methods: {
    input(value) {
      this.internalValue = value;

      this.$emit('input', value);
    },
  },
};
</script>

<style lang="scss">
.c-speed-dial {
  .v-speed-dial__list {
    padding: 8px 0;
  }

  &.v-speed-dial--direction-left, &.v-speed-dial--direction-right {
    .v-speed-dial__list {
      padding: 0 8px;
    }
  }

  .v-btn--fab {
    .v-btn__content {
      &, .v-icon {
        height: inherit;
        width: inherit;
      }

      & > :not(:only-child) {
        &:first-child {
          opacity: 1;
          position: absolute;
          left: 0;
          top: 0;
        }

        &:last-child {
          opacity: 0;
          position: absolute;
          left: 0;
          top: 0;
          transform: rotate(-45deg);
        }
      }
    }

    &.v-btn--active .v-btn__content>:not(:only-child) {
      &:first-child {
        opacity: 0;
        transform: rotate(45deg);
      }

      &:last-child {
        opacity: 1;
        transform: rotate(0);
      }
    }
  }
}
</style>
