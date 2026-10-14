<template>
  <v-fade-transition>
    <v-layout
      v-show="value"
      class="alert"
    >
      <div
        :style="{ opacity: opacity, backgroundColor }"
        class="overlay"
      />
      <div class="content">
        <slot>
          <v-alert type="error">
            {{ errorMessage }}
          </v-alert>
        </slot>
      </div>
    </v-layout>
  </v-fade-transition>
</template>

<script>
import { CSS_COLORS_VARS } from '@/config';

export default {
  props: {
    value: {
      type: Boolean,
      default: false,
    },
    opacity: {
      type: Number,
      default: 0.5,
    },
    backgroundColor: {
      type: String,
      default: CSS_COLORS_VARS.background,
    },
    message: {
      type: String,
      default: '',
    },
  },

  computed: {
    errorMessage() {
      return this.message || this.$t('errors.default');
    },
  },
};
</script>

<style lang="scss" scoped>
  .alert {
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    right: 0;
    z-index: 2;

    &, .overlay {
      min-height: 100px;
    }

    .overlay {
      position: absolute;
      top: 0;
      left: 0;
      bottom: 0;
      right: 0;
    }

    .content {
      width: 100%;
      display: flex;
      justify-content: center;
      align-items: center;
    }
  }
</style>
