<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        snmp-rule-form(v-model="form")
      template(#actions="")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.saveChanges') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS, SNMP_STATE_TYPES } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import SnmpRuleForm from '@/components/other/snmp-rule/form/snmp-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createSnmpRule,
  components: { SnmpRuleForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
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
