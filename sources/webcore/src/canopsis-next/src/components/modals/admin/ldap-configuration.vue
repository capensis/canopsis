<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline LDAP Authentification
    v-card-text.pa-0
      v-tabs(v-model="activeTab", color="secondary", slider-color="primary", dark)
        v-tab General
        v-tab Serveur
        v-tab Champs liés
        v-tab-item
          v-container
            v-layout(wrap)
              v-flex(xs12)
                v-switch(v-model="form.enabled", label="Enabled")
              v-flex(xs12)
                div Role par défaut
                v-select(v-model="form.defaultRole", :items="roles")
        v-tab-item
          v-container
            v-layout(wrap)
              v-flex(xs12)
                v-text-field(
                v-model="form.ldapServerHost",
                label="LDAP Server Host",
                name="ldapServerHost",
                v-validate="'required'",
                :error-messages="errors.collect('ldapServerHost')",
                )
              v-flex(xs12)
                v-text-field(
                v-model="form.ldapServerPort",
                label="LDAP Server Port",
                type="number",
                name="ldapServerPort",
                v-validate="'required|numeric|min:0'",
                :error-messages="errors.collect('ldapServerPort')",
                )
              v-flex(xs12)
                v-text-field(
                v-model="form.adminDn",
                label="Admin DN",
                name="adminDn",
                v-validate="'required'",
                :error-messages="errors.collect('adminDn')",
                )
              v-flex(xs12)
                v-text-field(
                v-model="form.adminPassword",
                label="Admin Password",
                name="password",
                v-validate="'required'",
                :error-messages="errors.collect('password')",
                )
              v-flex(xs12)
                v-text-field(
                v-model="form.adminDn",
                label="User filter",
                name="filter",
                v-validate="'required'",
                :error-messages="errors.collect('filter')",
                )
              v-flex(xs12)
                v-text-field(
                v-model="form.adminDn",
                label="User base",
                name="userBase",
                v-validate="'required'",
                :error-messages="errors.collect('userBase')",
                )
        v-tab-item
          v-container
            div(v-for="(attribute, key) in form.ldapAttributes" :key="key")
              v-layout(justify-center, align-center)
                v-flex
                  v-layout(align-center)
                    v-chip {{ attribute }}
                    v-icon arrow_right_alt
                    v-chip {{ key }}
                v-btn(icon, small)
                  v-icon(color="error") clear
            v-layout(align-center)
              v-text-field.mx-1(v-model="attributeForm.base", label="Base")
              v-text-field.mx-1(v-model="attributeForm.target", label="Target  ")
              v-btn(depressed, color="secondary") Add
      v-divider
      v-layout.py-1(justify-end)
        v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInner from '@/mixins/modal/inner';

export default {
  name: MODALS.ldapConfiguration,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInner],
  data() {
    return {
      activeTab: 0,
      form: {
        enabled: false,
        defaultRole: 'admin',
        ldapServerHost: '',
        ldapServerPort: 0,
        adminDn: '',
        adminPassword: '',
        userFilter: '',
        userBase: '',
        ldapAttributes: {
          password: 'test',
        },
      },
      attributeForm: {
        base: '',
        target: '',
      },
      roles: ['admin', 'CSIO'],
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        // TODO SEND CONFIG

        this.hideModal();
      }
    },
  },
};
</script>

