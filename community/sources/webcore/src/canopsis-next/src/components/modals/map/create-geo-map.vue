<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close, text-class="pa-0")
      template(#title="")
        span {{ title }}
      template(#text="")
        | Create geo map
      template(#actions="")
        v-btn(
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createGeoMap,
  $_veeValidate: {
    validator: 'new',
  },
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: { ...this.modal.config.map },
    };
  },
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.createGeoMap.create.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.form);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
