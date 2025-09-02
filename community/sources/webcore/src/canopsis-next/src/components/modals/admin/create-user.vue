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
          @click="close"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="isDisabled"
          :loading="submitting"
          class="primary white--text"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { ref, computed } from 'vue';

import { MODALS } from '@/constants';

import { userToForm, formToUserRequest } from '@/helpers/entities/user/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import UserForm from '@/components/other/users/form/user-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createUser,

  $_veeValidate: {
    validator: 'new',
  },

  components: { UserForm, ModalWrapper },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const form = ref(userToForm(config.value.user));

    const isNew = computed(() => !config.value.user);
    const title = computed(() => (config.value.title || t('modals.createUser.create.title')));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action(formToUserRequest(form.value));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      isNew,
      title,
      isDisabled,
      submitting,
      submit,
      close,
      config,
    };
  },
};
</script>
