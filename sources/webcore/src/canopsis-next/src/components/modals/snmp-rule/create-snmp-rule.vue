<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ title }}
    v-card-text
      snmp-rule-form(v-model="form")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
      v-btn.primary(@click="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS, SNMP_STATE_TYPES } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import SnmpRuleForm from '@/components/other/snmp-rule/form/snmp-rule-form.vue';

export default {
  name: MODALS.createSnmpRule,
  components: { SnmpRuleForm },
  mixins: [modalInnerMixin],
  data() {
    const defaultModuleMibObjectForm = {
      value: '',
      regex: '',
      formatter: '',
    };

    const defaultSnmpRule = {
      oid: {
        oid: '',
        mibName: '',
        moduleName: '',
      },
      component: { ...defaultModuleMibObjectForm },
      connector_name: { ...defaultModuleMibObjectForm },
      output: { ...defaultModuleMibObjectForm },
      resource: { ...defaultModuleMibObjectForm },
      state: {
        type: SNMP_STATE_TYPES.simple,
      },
    };

    return {
      form: this.modal.config.snmpRule ? cloneDeep(this.modal.config.snmpRule) : defaultSnmpRule,
    };
  },
  computed: {
    title() {
      if (this.config.snmpRule) {
        return this.$t('modals.createSnmpRule.edit.title');
      }

      return this.$t('modals.createSnmpRule.create.title');
    },
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

      this.$modals.hide();
    },
  },
};
</script>
