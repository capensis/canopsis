<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.createDynamicInfoInformation.create.title') }}
      template(slot="text")
        div
          v-text-field(
            v-model="form.name",
            v-validate="'required|unique-name'",
            :label="$t('modals.createDynamicInfoInformation.fields.name')",
            :error-messages="errors.collect('name')",
            name="name"
          )
          v-text-field(
            v-validate="'required'",
            v-model="form.value",
            :label="$t('modals.createDynamicInfoInformation.fields.value')",
            :error-messages="errors.collect('value')",
            name="value"
          )
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(:disabled="isDisabled", type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create dynamic info's information
 */
export default {
  name: MODALS.createDynamicInfoInformation,
  $_veeValidate: {
    validator: 'new',
  },
  components: { ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
  data() {
    const { info = {} } = this.modal.config;

    return {
      form: {
        name: info.name || '',
        value: info.value || '',
      },
    };
  },
  computed: {
    initialName() {
      return this.config.info && this.config.info.name;
    },
    existingNames() {
      return this.config.existingNames;
    },
  },
  created() {
    this.$validator.extend('unique-name', {
      getMessage: () => this.$t('validator.unique'),
      validate: value => (this.initialName && this.initialName === value) ||
        !this.existingNames.find(name => name === value),
    });
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
