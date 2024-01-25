<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <user-form
          v-model="form"
          :is-new="isNew"
          :user="config.user"
          :only-user-prefs="config.onlyUserPrefs"
        />
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          class="primary"
          :disabled="isDisabled"
          type="submit"
          data-test="submitButton"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS } from '@/constants';

import { userToForm, formToUserRequest } from '@/helpers/entities/user/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import UserForm from '@/components/other/users/form/user-form.vue';

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
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: userToForm(this.modal.config.user),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createUser.create.title');
    },

    isNew() {
      return !this.config.user;
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.config.action(formToUserRequest(this.form));

        this.$modals.hide();
      }
    },
  },
};
</script>
