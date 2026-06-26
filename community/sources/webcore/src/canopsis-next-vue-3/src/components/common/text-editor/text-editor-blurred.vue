<template>
  <div
    :class="['v-input--is-dirty', { 'v-input--is-disabled': disabled }]"
    class="v-input v-textarea v-text-field v-text-field--box v-text-field--enclosed v-input--is-label-active text-editor-blurred"
  >
    <div
      class="v-input__control"
      @click="$emit('click', $event)"
    >
      <div class="v-input__slot">
        <div
          class="v-text-field__slot"
          @click="wrapperClickHandler"
        >
          <label
            :class="{ 'v-label--active': value }"
            class="v-label"
          >
            {{ label }}
          </label>

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
          :messages="errorMessages"
          color="error"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { MODALS } from '@/constants';

export default {
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

  emits: ['click'],

  methods: {
    wrapperClickHandler(event) {
      if (event.target.tagName !== 'IMG') {
        return;
      }

      event.preventDefault();
      event.stopPropagation();

      this.$modals.show({
        name: MODALS.imageViewer,
        config: {
          src: event.target.getAttribute('src'),
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

.text-editor-blurred {
  padding-top: 28px !important;

  .v-text-field__slot {
    min-height: 150px;
    max-width: 100%;

    :deep(img) {
      cursor: pointer !important;
    }
  }

  .v-text-field--input__disabled {
    color: rgba(0, 0, 0, 0.38);

    :deep(img) {
      pointer-events: all !important;
    }
  }
}
</style>
