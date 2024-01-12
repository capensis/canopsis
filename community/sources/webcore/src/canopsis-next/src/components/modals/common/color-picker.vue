<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ title }}</span>
    </template>
    <template #text="">
      <v-layout>
        <v-flex>
          <c-color-chrome-picker-field v-model="color" />
        </v-flex>
        <v-flex>
          <c-color-compact-picker-field v-model="color" />
        </v-flex>
      </v-layout>
    </template>
    <template #actions="">
      <v-btn
        data-test="colorPickerCancelButton"
        depressed
        text
        @click="$modals.hide"
      >
        {{ $t('common.cancel') }}
      </v-btn>
      <v-btn
        class="primary"
        :disabled="isDisabled"
        :loading="submitting"
        data-test="colorPickerSubmitButton"
        @click="submit"
      >
        {{ $t('common.submit') }}
      </v-btn>
    </template>
  </modal-wrapper>
</template>

<script>
import { MODALS } from '@/constants';

import { colorToHex, colorToRgb, isValidColor } from '@/helpers/color';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.colorPicker,
  components: {
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
