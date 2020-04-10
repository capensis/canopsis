<template lang="pug">
  v-form(data-test="addInfoPopupModal", @submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.infoPopupSetting.addInfoPopup.title') }}
      template(slot="text")
        info-popup-form(v-model="form", :columns="config.columns")
      template(slot="actions")
        v-btn(
          flat,
          depressed,
          data-test="addInfoCancelButton",
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :loading="submitting",
          :disabled="isDisabled",
          type="submit",
          data-test="addInfoSubmitButton"
        ) {{ $t('common.submit') }}
</template>

<script>
import { get, pick } from 'lodash';

import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import InfoPopupForm from '@/components/other/alarm/forms/info-popup-form.vue';

import ModalWrapper from '../../modal-wrapper.vue';

export default {
  name: MODALS.addInfoPopup,
  $_veeValidate: {
    validator: 'new',
  },
  components: { InfoPopupForm, ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
  data() {
    const { popup, columns } = this.modal.config;
    let form = {
      template: '',
      column: null,
    };

    if (popup) {
      form = pick(popup, ['template', 'column']);
    } else if (columns && columns.length) {
      form.column = get(columns[0], 'value');
    }

    return {
      form,
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
