<template lang="pug">
  v-form(data-test="createRoleModal", @submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ title }}
      template(slot="text")
        role-form(v-model="form")
      template(slot="actions")
        v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary.white--text(
          :disabled="errors.any()",
          type="submit",
          data-test="submitButton"
        ) {{ $t('common.submit') }}
</template>

<script>
import { pick } from 'lodash';

import { MODALS } from '@/constants';

import { generateRole } from '@/helpers/entities';

import modalInnerMixin from '@/mixins/modal/inner';
import entitiesViewMixin from '@/mixins/entities/view';
import entitiesRoleMixin from '@/mixins/entities/role';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';

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
    entitiesViewMixin,
    entitiesRoleMixin,
    entitiesViewGroupMixin,
  ],
  data() {
    const group = this.modal.config.group || { name: '', description: '', defaultView: '' };

    return {
      form: pick(group, ['_id', 'description', 'defaultview']),
      defaultViewMenu: false,
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createRole.title');
    },

    role() {
      return this.config.roleId ? this.getRoleById(this.config.roleId) : null;
    },

    isNew() {
      return !this.role;
    },
  },
  mounted() {
    if (!this.isNew) {
      this.form = pick(this.role, [
        '_id',
        'description',
        'defaultview',
      ]);
    }
  },
  methods: {
    async submit() {
      try {
        const isFormValid = await this.$validator.validateAll();

        if (isFormValid) {
          const formData = this.isNew ? { ...generateRole() } : { ...this.role };
          formData._id = this.form._id;

          await this.createRole({ data: { ...formData, ...this.form } });
          await this.fetchRolesListWithPreviousParams();

          this.$popups.success({ text: this.$t('success.default') });
          this.$modals.hide();
        }
      } catch (err) {
        this.$popups.error({ text: this.$t('errors.default') });
      }
    },
  },
};
</script>

