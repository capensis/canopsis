<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline Create SNMP rule
    v-card-text
      snmp-rule-form(v-model="form")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(depressed, flat, @click="hideModal") {{ $t('common.cancel') }}
      v-btn.primary(:disabled="errors.any()", @click="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import SnmpRuleForm from '@/components/other/snmp-rule/form/snmp-rule-form.vue';

export default {
  name: MODALS.createSnmpRule,
  $_veeValidate: {
    validator: 'new',
  },
  components: { SnmpRuleForm },
  mixins: [modalInnerMixin],
  data() {
    return {
      form: {
        oid: {
          oid: '',
          mibName: '',
          moduleName: '',
        },
        component: {
          value: '',
          regex: '',
          formatter: '',
        },
        connector_name: {
          value: '',
          regex: '',
          formatter: '',
        },
        output: {
          value: '',
          regex: '',
          formatter: '',
        },
        resource: {
          value: '',
          regex: '',
          formatter: '',
        },
        state: {},
        _id: null,
      },
      items: [],
      loading: false,
      search: '',
    };
  },
  methods: {
    submit() {
      // TO DO SOMETHING
    },
  },
};
</script>
