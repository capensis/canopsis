<template lang="pug">
  v-card
    v-card-title.green.darken-3.white--text
      h2 {{ $t(config.title) }}
    v-card-text
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
          :label="$t('modals.createUser.fields.language')",
          v-model="form.ui_language",
          :items="languages"
          )
        v-layout(row)
          v-switch(
          :label="$t('modals.createUser.fields.enabled')",
          v-model="form.enable",
          )
    v-card-actions
      v-btn.green.darken-3.white--text(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import sha1 from 'sha1';
import omit from 'lodash/omit';
import cloneDeep from 'lodash/cloneDeep';

import { MODALS } from '@/constants';
import modalInnerMixin from '@/mixins/modal/modal-inner';
import entitiesUserMixin from '@/mixins/entities/user';

/**
 * Modal to create an entity (watcher, resource, component, connector)
 */
export default {
  name: MODALS.createUser,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [
    modalInnerMixin,
    entitiesUserMixin,
  ],
  data() {
    return {
      languages: ['fr', 'en'],
      form: {
        _id: '',
        firstname: '',
        lastname: '',
        password: '',
        mail: '',
        enable: true,
        ui_language: 'fr',
      },
    };
  },
  computed: {
    passwordRules() {
      if (this.config.item) {
        return null;
      }

      return 'required';
    },
  },
  mounted() {
    if (this.config.item) {
      this.form = cloneDeep(this.config.item);
    }
  },
  methods: {
    updateImpact(entities) {
      this.form.impacts = entities.map(entity => entity._id);
    },
    updateDependencies(entities) {
      this.form.dependencies = entities.map(entity => entity._id);
    },
    async submit() {
      const isFormValid = await this.$validator.validateAll();
      const defaultUserData = {
        crecord_write_time: null,
        crecord_type: 'user',
        crecord_creation_time: null,
        crecord_name: null,
        user_contact: null,
        rights: null,
        user_role: null,
        user_groups: null,
        authkey: null,
        role: null,
        external: false,
        defaultview: null,
        id: null,
      };

      if (isFormValid) {
        const formData = {
          ...defaultUserData,
          ...omit(this.form, ['password']),
        };

        if (this.form.password !== '') {
          formData.shadowpasswd = sha1(this.form.password);
        }

        await this.createUser({ data: formData });

        this.hideModal();
      }
    },
  },
};
</script>
