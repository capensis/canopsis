<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.createPbehaviorType.title') }}
      template(slot="text")
        create-type-form(v-model="form")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(type="submit") {{ $t('common.submit') }}
</template>

<script>
import { has } from 'lodash';

import { MODALS } from '@/constants';
import { typeToForm, formToType } from '@/helpers/forms/type';

import modalInnerMixin from '@/mixins/modal/inner';
import ExceptionsDatesLists from '@/components/other/pbehavior/exceptions/exceptions-dates-lists.vue';
import CreateTypeForm from '@/components/other/planning-administration/form/create-pbehavior-type-form.vue';


import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPbehaviorType,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    CreateTypeForm,
    ExceptionsDatesLists,
    ModalWrapper,
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      form: typeToForm(this.modal.config.type),
    };
  },
  methods: {
    setFormError(err = {}) {
      const existFieldErrors = Object.entries(err).filter(([field]) => has(this.form, field));

      if (existFieldErrors.length) {
        this.errors.add(existFieldErrors.map(([field, msg]) => ({ field, msg })));
      }
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          if (this.config.action) {
            await this.config.action(formToType(this.form));
          }

          this.$modals.hide();
        } catch (err) {
          this.setFormError(err);
        }
      }
    },
  },
};
</script>
