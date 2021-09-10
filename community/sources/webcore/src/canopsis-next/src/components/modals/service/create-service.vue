<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ config.title }}
      template(slot="text")
        service-form(v-model="form")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled || advancedJsonWasChanged",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { get } from 'lodash';

import { MODALS } from '@/constants';

import { serviceToForm, formToService } from '@/helpers/forms/service';

import { createSubmittableMixin } from '@/mixins/submittable';
import { createConfirmableModalMixin } from '@/mixins/confirmable-modal';
import { createValidationErrorsMixin } from '@/mixins/form/validation-errors';
import { entitiesContextEntityMixin } from '@/mixins/entities/context-entity';

import ServiceForm from '@/components/other/service/form/service-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createService,
  $_veeValidate: {
    validator: 'new',
  },
  components: { ServiceForm, ModalWrapper },
  mixins: [
    entitiesContextEntityMixin,
    createSubmittableMixin(),
    createConfirmableModalMixin(),
    createValidationErrorsMixin(),
  ],
  data() {
    return {
      form: serviceToForm(this.modal.config.item),
    };
  },
  computed: {
    advancedJsonWasChanged() {
      return get(this.fields, ['advancedJson', 'changed']);
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          const data = formToService(this.form);

          await this.config.action(data);

          this.$modals.hide();
        } catch (err) {
          this.setFormErrors(err);
        }
      }
    },
  },
};
</script>

<style>

</style>
