<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        c-patterns-field(
          v-model="form",
          :alarm="config.alarm",
          :entity="config.entity",
          :event="config.event"
        )
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
import { cloneDeep } from 'lodash';

import { MODALS } from '@/constants';

import { createSubmittableMixin } from '@/mixins/submittable';
import { createConfirmableModalMixin } from '@/mixins/confirmable-modal';
import { createValidationErrorsMixin } from '@/mixins/form/validation-errors';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.patterns,
  $_veeValidate: {
    validator: 'new',
  },
  components: { ModalWrapper },
  mixins: [
    createSubmittableMixin(),
    createValidationErrorsMixin(),
    createConfirmableModalMixin(),
  ],
  data() {
    return {
      form: this.modal.config.patterns ? cloneDeep(this.modal.config.patterns) : {},
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.patterns.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          if (this.config.action) {
            await this.config.action(this.form);
          }

          this.$modals.hide();
        } catch (err) {
          this.setFormErrors(err);
        }
      }
    },
  },
};
</script>
