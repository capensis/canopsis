<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ title }}
      template(slot="text")
        user-form(
          v-model="form",
          :isNew="isNew",
          :user="config.user",
          :onlyUserPrefs="config.onlyUserPrefs"
        )
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          type="submit",
          data-test="submitButton"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { userToForm } from '@/helpers/forms/user';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import UserForm from '@/components/other/user/user-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create an entity (watcher, resource, component, connector)
 */
export default {
  name: MODALS.createUser,
  $_veeValidate: {
    validator: 'new',
  },
  components: { UserForm, ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
  data() {
    const { user } = this.modal.config;

    return {
      form: userToForm(user),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createUser.title');
    },

    isNew() {
      return !this.config.user;
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
