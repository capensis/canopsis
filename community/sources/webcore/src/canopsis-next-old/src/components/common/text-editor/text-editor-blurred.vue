<template>
  <div
    :class="['v-input--is-dirty', { 'v-input--is-disabled': disabled }, themeClasses]"
    class="v-input v-textarea v-text-field v-text-field--box v-text-field--enclosed v-input--is-label-active"
  >
    <div
      class="v-input__control"
      @click="$emit('click', $event)"
    >
      <div class="v-input__slot">
        <div class="v-text-field__slot" @click="wrapperClickHandler">
          <label
            :class="[{ 'v-label--active': value }, themeClasses]"
            class="v-label"
          >{{ label }}</label>
          <c-compiled-template
            :template="value"
            :class="{ 'v-text-field--input__disabled': disabled }"
          />
        </div>
      </div>
      <div
        v-if="!hideDetails"
        class="v-text-field__details"
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
  mixins: [Themeable], // TODO: rewrite it to composition in the future
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
  methods: {
    wrapperClickHandler(e) {
      if (e.target.tagName !== 'IMG') {
        return;
      }

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
  padding-top: 28px !important;

  &__slot {
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
