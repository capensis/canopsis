<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        entity-form(v-model="form")
      template(slot="actions")
        v-btn(
          :disabled="submitting",
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

import { entityToForm, formToEntity } from '@/helpers/forms/entity';

import { validationErrorsMixin } from '@/mixins/form/validation-errors';
import { submittableMixin } from '@/mixins/submittable';
import { confirmableModalMixin } from '@/mixins/confirmable-modal';

import EntityForm from '@/components/other/entity/form/entity-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createEntity,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    EntityForm,
    ModalWrapper,
  },
  mixins: [
    submittableMixin(),
    confirmableModalMixin(),
    validationErrorsMixin(),
  ],
  data() {
    return {
      form: entityToForm(this.modal.config.entity),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createEntity.create.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          await this.config.action(formToEntity(this.form));

          this.$modals.hide();
        } catch (err) {
          this.setFormErrors(err);
        }
      }
    },

  },
};
</script>