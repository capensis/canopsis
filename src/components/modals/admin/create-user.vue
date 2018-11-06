<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ title }}
    v-card-text
      progress-overlay(:pending="pending")
      v-form
        v-layout(row)
          v-text-field(
          :label="$t('modals.createUser.fields.username')",
          v-model="form._id",
          data-vv-name="username",
          v-validate="'required'",
          :error-messages="errors.collect('username')",
          )
        v-layout(row)
          v-text-field(
          :label="$t('modals.createUser.fields.firstName')",
          v-model="form.firstname",
          )
        v-layout(row)
          v-text-field(
          :label="$t('modals.createUser.fields.lastName')",
          v-model="form.lastname",
          )
        v-layout(row)
          v-text-field(
          :label="$t('modals.createUser.fields.email')",
          v-model="form.mail",
          data-vv-name="email",
          v-validate="'required|email'",
          :error-messages="errors.collect('email')",
          )
        v-layout(row)
          v-text-field(
          type="password",
          :label="$t('modals.createUser.fields.password')",
          v-model="form.password",
          data-vv-name="password",
          v-validate="passwordRules",
          :error-messages="errors.collect('password')",
          )
        v-layout(row)
          v-select(
          :label="$t('modals.createUser.fields.role')",
          v-model="form.role",
          :items="roles",
          item-text="_id",
          item-value="_id",
          data-vv-name="role",
          v-validate="'required'",
          :error-messages="errors.collect('role')",
          )
        v-layout(row)
          v-select(
          :label="$t('modals.createUser.fields.language')",
          v-model="form.ui_language",
          :items="languages"
          )
        v-layout(row)
          v-switch(
            color="primary",
          :label="$t('modals.createUser.fields.enabled')",
          v-model="form.enable",
          )
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import sha1 from 'sha1';
import omit from 'lodash/omit';
import pick from 'lodash/pick';

import { MODALS } from '@/constants';
import { generateUser } from '@/helpers/entities';
import modalInnerMixin from '@/mixins/modal/modal-inner';
import entitiesRoleMixin from '@/mixins/entities/role';
import entitiesUserMixin from '@/mixins/entities/user';

import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';

/**
 * Modal to create an entity (watcher, resource, component, connector)
 */
export default {
  name: MODALS.createUser,
  $_veeValidate: {
    validator: 'new',
  },
  components: { ProgressOverlay },
  mixins: [
    modalInnerMixin,
    entitiesRoleMixin,
    entitiesUserMixin,
  ],
  data() {
    return {
      pending: true,
      languages: ['fr', 'en'],
      form: {
        _id: '',
        firstname: '',
        lastname: '',
        mail: '',
        password: '',
        role: null,
        ui_language: 'fr',
        enable: true,
      },
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createUser.title');
    },

    user() {
      return this.config.userId ? this.getUserById(this.config.userId) : null;
    },

    passwordRules() {
      if (this.isNew) {
        return 'required';
      }

      return null;
    },
    isNew() {
      return !this.user;
    },
  },
  async mounted() {
    await this.fetchRolesList({ params: { limit: 0 } });

    if (!this.isNew) {
      this.form = pick(this.user, [
        '_id',
        'firstname',
        'lastname',
        'mail',
        'password',
        'role',
        'ui_language',
        'enable',
      ]);
    }

    this.pending = false;
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const formData = this.isNew ? { ...generateUser() } : { ...this.user };

        if (this.form.password && this.form.password !== '') {
          formData.shadowpasswd = sha1(this.form.password);
        }

        await this.createUser({ data: { ...formData, ...omit(this.form, ['password']) } });
        await this.fetchUsersListWithPreviousParams();

        this.hideModal();
      }
    },
  },
};
</script>
