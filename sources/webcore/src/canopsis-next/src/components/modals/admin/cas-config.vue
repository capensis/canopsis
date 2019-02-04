<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline CAS Authentification
    v-card-text.pa-0
      v-container
        v-layout(wrap)
          v-flex(xs12)
            v-switch(v-model="form.enabled", label="Enabled")
          v-flex(xs12)
            div Role par défaut
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
            v-text-field(
            v-model="form.title",
            label="Title",
            name="title",
            v-validate="'required'",
            :error-messages="errors.collect('title')",
            )
          v-flex(xs12)
            v-text-field(
            v-model="form.server",
            label="Server",
            name="server",
            v-validate="'required'",
            :error-messages="errors.collect('server')",
            )
          v-flex(xs12)
            v-text-field(
            v-model="form.service",
            label="Service",
            name="service",
            v-validate="'required'",
            :error-messages="errors.collect('service')",
            )
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
