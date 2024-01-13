<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.infoPopupSetting.addInfoPopup.title') }}</span>
      </template>
      <template #text="">
        <info-popup-form
          v-model="form"
          :columns="config.columns"
        />
      </template>
      <template #actions="">
        <v-btn
          text
          depressed
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          class="primary"
          :loading="submitting"
          :disabled="isDisabled"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import InfoPopupForm from '@/components/widgets/alarm/forms/info-popup-form.vue';

import ModalWrapper from '../../modal-wrapper.vue';

export default {
  name: MODALS.addInfoPopup,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { InfoPopupForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { popup = {}, columns = [] } = this.modal.config;
    const [firstColumn] = columns;

    return {
      form: {
        column: popup?.column ?? firstColumn?.value,
        template: popup.template ?? '',
      },
    };
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
