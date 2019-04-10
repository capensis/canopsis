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
import { cloneDeep } from 'lodash';

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
    const defaultSnmpRule = {
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
      state: {
        type: '',
      },
    };

    return {
      form: this.modal.config.snmpRule ? cloneDeep(this.modal.config.snmpRule) : defaultSnmpRule,
    };
  },
  methods: {
    async submit() {
      if (this.config.action) {
        const preparedData = this.form;

        if (preparedData._id) {
          preparedData.id = preparedData._id;
        }

        await this.config.action({
          document: preparedData,
        });
      }

      this.hideModal();
    },
  },
};
</script>
