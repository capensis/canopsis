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
            return-object,
            item-text="id",
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
import { MODALS } from '@/constants';

import modalInner from '@/mixins/modal/inner';
import roleMixin from '@/mixins/entities/role';

export default {
  name: MODALS.casConfiguration,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInner, roleMixin],
  data() {
    return {
      form: {
        enabled: false,
        defaultRole: 'admin',
        title: '',
        server: '',
        service: '',
      },
      roles: [],
    };
  },
  async mounted() {
    const roles = await this.fetchRolesListWithoutStore({ limit: 10000 });

    this.roles = roles.data;
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        // TODO: SEND CONFIG

        this.hideModal();
      }
    },
  },
};
</script>
