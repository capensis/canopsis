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
import { find } from 'lodash';

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
    return {
      form: {
        selectedColumn: {},
        template: '',
      },
    };
  },
  mounted() {
    if (this.config) {
      [this.form.selectedColumn] = this.config.columns;

      if (this.config.popup) {
        const { template, column } = this.config.popup;

        this.form.template = template;
        this.form.selectedColumn = find(this.config.columns, { value: column });
      }
    }
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action({ column: this.form.selectedColumn.value, template: this.form.template });
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
