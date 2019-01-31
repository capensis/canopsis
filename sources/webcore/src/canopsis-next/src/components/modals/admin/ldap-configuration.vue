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
                v-text-field(v-model="form.ldapServerHost", label="LDAP Server Host")
              v-flex(xs12)
                v-text-field(v-model="form.ldapServerPost", label="LDAP Server Port", type="number")
              v-flex(xs12)
                v-text-field(v-model="form.adminDn", label="Admin DN")
              v-flex(xs12)
                v-text-field(v-model="form.adminPassword", label="Admin Password")
              v-flex(xs12)
                v-text-field(v-model="form.adminDn", label="User filter")
              v-flex(xs12)
                v-text-field(v-model="form.adminDn", label="User base")
        v-tab-item
      v-divider
      v-layout.py-1(justify-end)
        v-btn.primary(@click.prevent="") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

export default {
  name: MODALS.ldapConfiguration,
  data() {
    return {
      activeTab: 0,
      form: {
        enabled: false,
        defaultRole: 'admin',
        ldapServerHost: '',
        ldapServerPost: 0,
        adminDn: '',
        adminPassword: '',
        userFilter: '',
        userBase: '',
      },
      roles: ['admin', 'CSIO'],
    };
  },
};
</script>

