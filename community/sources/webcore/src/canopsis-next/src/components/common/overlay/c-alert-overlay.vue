<template>
  <v-fade-transition>
    <v-layout
      class="alert"
      v-show="value"
    >
      <div
        class="overlay"
        :class="backgroundColor"
        :style="{ opacity: opacity }"
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
      default: 'white',
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
    z-index: 2;

    &, .overlay {
      min-height: 100px;
    }

    .content {
      width: 100%;
      display: flex;
      justify-content: center;
      align-items: center;
    }
  }
</style>
