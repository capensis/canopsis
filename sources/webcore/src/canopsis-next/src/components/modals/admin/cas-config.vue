<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('parameters.casAuthentication.title') }}
    v-card-text.pa-0
      v-container
        v-layout(wrap)
          v-flex(xs12)
            v-switch(v-model="form.enabled", :label="$t('common.enabled')")
          v-flex(xs12)
            v-layout(align-center)
              div(slot="activator") {{ $t('parameters.casAuthentication.fields.defaultRole.title') }}
              v-tooltip.ml-1(left)
                v-icon(slot="activator", small) help
                span {{ $t('parameters.casAuthentication.fields.defaultRole.tooltip') }}
            v-select(
            v-model="form.defaultRole",
            :items="roles",
            item-text="crecord_name",
            item-value="crecord_name"
            name="role",
            v-validate="'required'",
            :error-messages="errors.collect('role')",
            )
          v-flex(xs12)
            v-tooltip.ml-1(left)
              v-text-field(
              slot="activator",
              v-model="form.title",
              :label="$t('parameters.casAuthentication.fields.title.title')",
              name="title",
              v-validate="'required'",
              :error-messages="errors.collect('title')",
              )
              span {{ $t('parameters.casAuthentication.fields.title.tooltip') }}
          v-flex(xs12)
            v-tooltip.ml-1(left)
              v-text-field(
              slot="activator",
              v-model="form.server",
              :label="$t('parameters.casAuthentication.fields.server.title')",
              name="server",
              v-validate="'required'",
              :error-messages="errors.collect('server')",
              )
              span {{ $t('parameters.casAuthentication.fields.server.tooltip') }}
          v-flex(xs12)
            v-tooltip.ml-1(left)
              v-text-field(
              slot="activator",
              v-model="form.service",
              :label="$t('parameters.casAuthentication.fields.service.title')",
              name="service",
              v-validate="'required'",
              :error-messages="errors.collect('service')",
              )
              span {{ $t('parameters.casAuthentication.fields.service.tooltip') }}
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
  name: MODALS.casConfiguration,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInner, roleMixin, authProtocolMixin],
  data() {
    return {
      form: {
        enabled: false,
        defaultRole: 'admin',
        title: '',
        server: '',
        service: '',
      },
    };
  },
  async mounted() {
    await this.fetchRolesList();

    const casConfig = await this.fetchCASConfigWithoutStore();

    const {
      enable: enabled,
      default_role: defaultRole,
      title,
      server,
      service,
    } = casConfig[0];

    this.form = {
      ...this.form,
      enabled,
      defaultRole,
      title,
      server,
      service,
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const data = {
          default_role: this.form.defaultRole,
          enable: this.form.enabled,
          id: 'cservice.casconfig',
          server: this.form.server,
          service: this.form.service,
          title: this.form.title,
          crecord_creation_time: moment().unix(),
          crecord_write_time: moment().unix(),
          crecord_name: 'casconfig',
          crecord_type: 'casconfig',
        };

        await this.updateCASConfig({ data });

        this.hideModal();
      }
    },
  },
};
</script>
