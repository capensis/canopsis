<template lang="pug">
  div.v-input.v-textarea.v-text-field.v-text-field--box.v-text-field--enclosed.v-input--is-label-active(
    class="v-input--is-dirty theme--light"
  )
    div.v-input__control(@click="$emit('click', $event)")
      div.v-input__slot
        div.v-text-field__slot
          label.v-label.theme--light(:class="{ 'v-label--active': value }") {{ label }}
          div(ref="content", v-html="value")
      div.v-text-field__details(v-if="!hideDetails")
        div.v-messages.theme--light.error--text
          div.v-messages__wrapper
            div.v-messages__message(v-for="errorMessage in errorMessages") {{ errorMessage }}
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

      this.imagesElements = this.$refs.content.querySelectorAll('img');
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

.v-text-field__slot {
  padding-top: 28px !important;
  min-height: 150px;
  max-width: 100%;

  & /deep/ img {
    cursor: pointer !important;
  }
}
</style>
