<template lang="pug">
  v-form(data-test="createRoleModal", @submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ title }}
      template(slot="text")
        role-form(v-model="form")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary.white--text(
          :disabled="isDisabled",
          type="submit",
          data-test="submitButton"
        ) {{ $t('common.submit') }}
</template>

<script>
import { pick } from 'lodash';

import { MODALS } from '@/constants';

import { generateRole } from '@/helpers/entities';

import modalInnerMixin from '@/mixins/modal/inner';
import entitiesRoleMixin from '@/mixins/entities/role';
import submittableMixin from '@/mixins/submittable';

import RoleForm from '@/components/other/role/role-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createRole,
  $_veeValidate: {
    validator: 'new',
  },
  components: { RoleForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    entitiesRoleMixin,
    submittableMixin(),
  ],
  data() {
    const role = this.modal.config.role || { name: '', description: '', defaultView: '' };

    return {
      form: pick(role, ['_id', 'description', 'defaultview']),
      defaultViewMenu: false,
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createRole.title');
    },

    isNew() {
      return !this.config.role;
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const formData = this.isNew ? generateRole() : { ...this.role };

        await this.createRole({ data: { ...formData, ...this.form } });
        await this.fetchRolesListWithPreviousParams();

        this.$popups.success({ text: this.$t('success.default') });
        this.$modals.hide();
      }
    },
  },
};
</script>

