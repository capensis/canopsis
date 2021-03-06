<template lang="pug">
  modal-wrapper(data-test="colorPickerModal", close)
    template(slot="title")
      span {{ title }}
    template(slot="text")
      v-layout
        v-flex
          chrome(data-test="colorPickerChrome", v-model="color")
        v-flex
          compact(data-test="colorPickerCompact", v-model="color")
    template(slot="actions")
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
import tinycolor from 'tinycolor2';
import { Chrome, Compact } from 'vue-color';

import { MODALS } from '@/constants';

import submittableMixin from '@/mixins/submittable';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.colorPicker,
  components: {
    Chrome,
    Compact,
    ModalWrapper,
  },
  mixins: [submittableMixin()],
  data() {
    const { config } = this.modal;
    const color = {};

    if (config.color) {
      const colorObject = tinycolor(config.color);

      if (colorObject.isValid()) {
        color.hex = colorObject.toHexString();
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
        const colorObject = tinycolor(hex);
        const result = this.isHexType ? hex : colorObject.toRgbString();

        await this.config.action(result);
      }

      this.$modals.hide();
    },
  },
};
</script>
