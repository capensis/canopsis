<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ title }}
    template(#text="")
      v-layout
        v-flex
          chrome(v-model="color")
        v-flex
          compact(v-model="color", :palette="config.palette")
    template(#actions="")
      v-btn(
        data-test="colorPickerCancelButton",
        depressed,
        flat,
        @click="$modals.hide"
      ) {{ $t('common.cancel') }}
      v-btn.primary(
        :disabled="isDisabled",
        :loading="submitting",
        data-test="colorPickerSubmitButton",
        @click="submit"
      ) {{ $t('common.submit') }}
</template>

<script>
import { Chrome, Compact } from 'vue-color';

import { MODALS } from '@/constants';

import { colorToHex, colorToRgb, isValidColor } from '@/helpers/color';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.colorPicker,
  components: {
    Chrome,
    Compact,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
  ],
  data() {
    const { config } = this.modal;
    const color = {};

    if (config.color) {
      if (isValidColor(config.color)) {
        color.hex = colorToHex(config.color);
      }
    }

    return {
      color,
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.colorPicker.title');
    },
    isHexType() {
      return this.config.type === 'hex';
    },
  },
  methods: {
    async submit() {
      if (this.config.action) {
        const { hex } = this.color;
        const result = this.isHexType ? hex : colorToRgb(hex);

        await this.config.action(result);
      }

      this.$modals.hide();
    },
  },
};
</script>
