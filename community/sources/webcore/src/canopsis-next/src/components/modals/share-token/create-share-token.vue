<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        share-token-form(v-model="form")
      template(#actions="")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { formToShareToken, shareTokenToForm } from '@/helpers/forms/share-token';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import ShareTokenForm from '@/components/other/share-token/form/share-token-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createShareToken,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { ShareTokenForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: shareTokenToForm(this.modal.config.shareToken),
    };
  },
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.createShareToken.create.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToShareToken(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
