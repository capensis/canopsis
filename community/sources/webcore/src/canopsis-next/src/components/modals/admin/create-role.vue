<template>
  <v-form
    data-test="createRoleModal"
    @submit.prevent="submit"
  >
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <role-form
          v-model="form"
          :with-template="config.withTemplate"
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

import { roleToForm, formToRole } from '@/helpers/entities/role/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import RoleForm from '@/components/other/role/form/role-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createRole,

  $_veeValidate: {
    validator: 'new',
  },

  components: {
    RoleForm,
    ModalWrapper,
  },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const form = ref(roleToForm(config.value.role));

    const title = computed(() => (
      config.value.title || t('modals.createRole.create.title')
    ));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action(formToRole(form.value));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
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
