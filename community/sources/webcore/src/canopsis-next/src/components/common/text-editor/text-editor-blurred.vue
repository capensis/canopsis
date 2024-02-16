<template>
  <div
    class="v-input v-textarea v-text-field v-text-field--box v-text-field--enclosed v-input--is-label-active"
    :class="['v-input--is-dirty', { 'v-input--is-disabled': disabled }, themeClasses]"
  >
    <div
      class="v-input__control"
      @click="$emit('click', $event)"
    >
      <div class="v-input__slot">
        <div class="v-text-field__slot">
          <label
            class="v-label"
            :class="[{ 'v-label--active': value }, themeClasses]"
          >{{ label }}</label>
          <c-compiled-template
            ref="content"
            :template="value"
            :class="{ 'v-text-field--input__disabled': disabled }"
          />
        </div>
      </div>
      <div
        class="v-text-field__details"
        v-if="!hideDetails"
      >
        <v-messages
          :value="errorMessages"
          color="error"
        />
      </div>
    </div>
  </div>
</template>

<script>
import Themeable from 'vuetify/lib/mixins/themeable';

import { MODALS } from '@/constants';

export default {
  mixins: [Themeable],
  props: {
    value: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    errorMessages: {
      type: Array,
      default: () => [],
    },
  },
  watch: {
    value() {
      this.addImagesListeners();
    },
  },
  mounted() {
    this.addImagesListeners();
  },
  beforeDestroy() {
    this.removeImagesListeners();
  },
  methods: {
    addImagesListeners() {
      this.removeImagesListeners();

      this.imagesElements = this.$refs.content.$el.querySelectorAll('img');
      this.imagesElements.forEach(image => image.addEventListener('click', this.clickHandler));
    },

    removeImagesListeners() {
      if (this.imagesElements) {
        this.imagesElements.forEach(image => image.removeEventListener('click', this.clickHandler));
      }
    },

    clickHandler(e) {
      e.preventDefault();
      e.stopPropagation();

      this.$modals.show({
        name: MODALS.imageViewer,
        config: {
          src: e.target.getAttribute('src'),
        },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.v-label {
  left: 0;
  right: auto;
  position: absolute;
}

.v-text-field {
  &__slot {
    padding-top: 28px !important;
    min-height: 150px;
    max-width: 100%;

    & ::v-deep img {
      cursor: pointer !important;
    }
  }

  &--input {
    &__disabled {
      color: rgba(0, 0, 0, 0.38);

      & ::v-deep img {
        pointer-events: all !important;
      }
    }
  }
}
</style>
