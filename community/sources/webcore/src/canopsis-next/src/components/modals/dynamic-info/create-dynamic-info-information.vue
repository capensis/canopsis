<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ $t('modals.createDynamicInfoInformation.create.title') }}
      template(slot="text")
        div
          v-text-field(
            v-model="form.name",
            v-validate="nameRules",
            :label="$t('modals.createDynamicInfoInformation.fields.name')",
            :error-messages="errors.collect('name')",
            name="name"
          )
          c-mixed-field(
            v-model="form.value",
            :label="$t('modals.createDynamicInfoInformation.fields.value')",
            name="value",
            required
          )
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(:disabled="isDisabled", type="submit") {{ $t('common.submit') }}
</template>

<script>
import { get } from 'lodash';

import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

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
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { info } = this.modal.config;

    return {
      form: {
        name: get(info, 'name', ''),
        value: get(info, 'value', ''),
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

    nameRules() {
      return {
        required: true,
        unique: {
          values: this.existingNames,
          initialValue: this.initialName,
        },
      };
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
