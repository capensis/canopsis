<template lang="pug">
  v-form(data-test="createRoleModal", @submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        role-form(v-model="form")
      template(#actions="")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary.white--text(
          :disabled="isDisabled",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { roleToForm, formToRole } from '@/helpers/forms/role';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import RoleForm from '@/components/other/role/form/role-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createRole,
  $_veeValidate: {
    validator: 'new',
  },
  components: { RoleForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: roleToForm(this.modal.config.role),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createRole.create.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.config.action(formToRole(this.form));

        this.$modals.hide();
      }
    },
  },
};
</script>
