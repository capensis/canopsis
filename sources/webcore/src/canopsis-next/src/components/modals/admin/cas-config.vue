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
            v-select(v-model="form.defaultRole", :items="roles")
          v-flex(xs12)
            v-text-field(v-model="form.title", label="Title")
          v-flex(xs12)
            v-text-field(v-model="form.server", label="Server")
          v-flex(xs12)
            v-text-field(v-model="form.service", label="Service")
    v-divider
    v-layout.py-1(justify-end)
      v-btn.primary(@click.prevent="") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

export default {
  name: MODALS.casConfiguration,
  data() {
    return {
      form: {
        enabled: false,
        defaultRole: 'admin',
        title: '',
        server: '',
        service: '',
      },
      roles: ['admin', 'CSIO'],
    };
  },
};
</script>
