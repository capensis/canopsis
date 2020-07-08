<template lang="pug">
  modal-wrapper(data-test="colorPickerModal")
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

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.colorPicker,
  components: {
    Chrome,
    Compact,
    ModalWrapper,
  },
  mixins: [modalInnerMixin, submittableMixin()],
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
    type() {
      return this.config.type || 'rgba';
    },
  },
  methods: {
    async submit() {
      const { rgba, hex } = this.color;

      if (this.config.action) {
        const result = this.type === 'hex' ? hex : `rgba(${rgba.r}, ${rgba.g}, ${rgba.b}, ${rgba.a})`;

        await this.config.action(result);
      }

      this.$modals.hide();
    },
  },
};
</script>
