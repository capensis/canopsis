<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('parameters.ldapAuthentication.title') }}
    v-card-text.pa-0
      v-tabs(v-model="activeTab", color="secondary", slider-color="primary", dark)
        v-tab {{ $t('parameters.ldapAuthentication.tabs.general') }}
        v-tab {{ $t('parameters.ldapAuthentication.tabs.server') }}
        v-tab {{ $t('parameters.ldapAuthentication.tabs.linkFields') }}
        v-tab-item
          v-container
            v-layout(wrap)
              v-flex(xs12)
                v-switch(v-model="form.enabled", :label="$t('common.enabled')")
              v-flex(xs12)
                v-layout(align-center)
                  div(slot="activator") {{ $t('parameters.ldapAuthentication.fields.defaultRole.title') }}
                  v-tooltip.ml-1(left)
                    v-icon(slot="activator", small) help
                    span {{ $t('parameters.ldapAuthentication.fields.defaultRole.tooltip') }}
                v-select(
                v-model="form.defaultRole",
                :items="roles",
                item-text="crecord_name",
                item-value="crecord_name",
                )
        v-tab-item
          v-container
            v-layout(wrap)
              v-flex(xs12)
                v-tooltip.ml-1(left)
                  v-text-field(
                  slot="activator",
                  v-model="form.ldapServerHost",
                  :label="$t('parameters.ldapAuthentication.fields.ldapServerHost.title')",
                  name="ldapServerHost",
                  v-validate="'required'",
                  :error-messages="errors.collect('ldapServerHost')",
                  )
                  span {{ $t('parameters.ldapAuthentication.fields.ldapServerHost.tooltip') }}
              v-flex(xs12)
                v-tooltip.ml-1(left)
                  v-text-field(
                  slot="activator",
                  v-model="form.ldapServerPort",
                  :label="$t('parameters.ldapAuthentication.fields.ldapServerPort.title')",
                  type="number",
                  name="ldapServerPort",
                  v-validate="'required|numeric|min:0'",
                  :error-messages="errors.collect('ldapServerPort')",
                  )
                  span {{ $t('parameters.ldapAuthentication.fields.ldapServerPort.tooltip') }}
              v-flex(xs12)
                v-tooltip.ml-1(left)
                  v-text-field(
                  slot="activator",
                  v-model="form.adminDn",
                  :label="$t('parameters.ldapAuthentication.fields.adminDn.title')",
                  name="adminDn",
                  v-validate="'required'",
                  :error-messages="errors.collect('adminDn')",
                  )
                  span {{ $t('parameters.ldapAuthentication.fields.adminDn.tooltip') }}
              v-flex(xs12)
                v-tooltip.ml-1(left)
                  v-text-field(
                  slot="activator",
                  v-model="form.adminPassword",
                  :label="$t('parameters.ldapAuthentication.fields.adminPassword.title')",
                  type="password",
                  name="password",
                  v-validate="'required'",
                  :error-messages="errors.collect('password')",
                  )
                  span {{ $t('parameters.ldapAuthentication.fields.adminPassword.tooltip') }}
              v-flex(xs12)
                v-tooltip.ml-1(left)
                  v-text-field(
                  slot="activator",
                  v-model="form.adminDn",
                  :label="$t('parameters.ldapAuthentication.fields.userFilter.title')",
                  name="filter",
                  v-validate="'required'",
                  :error-messages="errors.collect('filter')",
                  )
                  span {{ $t('parameters.ldapAuthentication.fields.userFilter.tooltip') }}
              v-flex(xs12)
                v-tooltip.ml-1(left)
                  v-text-field(
                  slot="activator",
                  v-model="form.adminDn",
                  :label="$t('parameters.ldapAuthentication.fields.userBase.title')",
                  name="userBase",
                  v-validate="'required'",
                  :error-messages="errors.collect('userBase')",
                  )
                  span {{ $t('parameters.ldapAuthentication.fields.userBase.tooltip') }}
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
              v-text-field.mx-1(
              v-model="attributeForm.base",
              :label="$t('parameters.ldapAuthentication.fields.linkFields.base')",
              )
              v-text-field.mx-1(
              v-model="attributeForm.target",
              :label="$t('parameters.ldapAuthentication.fields.linkFields.target')",
              )
              v-btn(depressed, color="secondary") {{ $t('common.add') }}
      v-divider
      v-layout.py-1(justify-end)
        v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import moment from 'moment';

import { MODALS } from '@/constants';

import modalInner from '@/mixins/modal/inner';
import roleMixin from '@/mixins/entities/role';
import authProtocolMixin from '@/mixins/entities/authProtocol';

export default {
  name: MODALS.ldapConfiguration,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInner, roleMixin, authProtocolMixin],
  data() {
    return {
      activeTab: 0,
      form: {
        enabled: false,
        defaultRole: '',
        ldapServerHost: '',
        ldapServerPort: 0,
        adminDn: '',
        adminPassword: '',
        userFilter: '',
        userBase: '',
        ldapAttributes: {},
      },
      attributeForm: {
        base: '',
        target: '',
      },
    };
  },
  async mounted() {
    await this.fetchRolesList();
    const ldapConfiguration = await this.fetchLDAPConfigWithoutStore();

    const {
      enable: enabled,
      default_role: defaultRole,
      host: ldapServerHost,
      port: ldapServerPort,
      admin_dn: adminDn,
      admin_passwd: adminPassword,
      ufilter: userFilter,
      user_dn: userBase,
      attrs: ldapAttributes,
    } = ldapConfiguration[0];

    this.form = {
      ...this.form,
      enabled,
      defaultRole,
      ldapServerHost,
      ldapServerPort,
      adminDn,
      adminPassword,
      userFilter,
      userBase,
      ldapAttributes,
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const data = {
          admin_dn: this.form.adminDn,
          admin_passwd: this.form.adminPassword,
          attrs: this.form.ldapAttributes,
          default_role: this.form.defaultRole,
          enable: this.form.enabled,
          host: this.form.ldapServerHost,
          port: this.form.ldapServerPort,
          ufilter: this.form.userFilter,
          user_dn: this.form.userBase,
          crecord_creation_time: moment().unix(),
          crecord_write_time: moment().unix(),
          crecord_name: 'ldapconfig',
          crecord_type: 'ldapconfig',
        };

        await this.updateLDAPConfig({ data });

        this.hideModal();
      }
    },
  },
};
</script>

