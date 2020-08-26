<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.createRrule.title') }}
      template(slot="text")
        r-rule-form(v-model="form.rrule")
        pbehavior-exception-form(v-if="form.rrule", v-model="form.exdates", :exceptions.sync="form.exceptions")
      template(slot="actions")
        v-btn(
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS } from '@/constants';

import RRuleForm from '@/components/forms/rrule.vue';
import PbehaviorExceptionForm from '@/components/other/pbehavior/calendar/partials/pbehavior-exception-form.vue';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createRRule,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    PbehaviorExceptionForm,
    RRuleForm,
    ModalWrapper,
  },
  mixins: [modalInnerMixin, submittableMixin()],
  data() {
    return {
      form: {
        rrule: this.modal.config.rrule || '',
        exdates: this.modal.config.exdates ? cloneDeep(this.modal.config.exdates) : [],
        exceptions: this.modal.config.exceptions ? cloneDeep(this.modal.config.exceptions) : [],
      },
    };
  },
  methods: {
    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        if (this.config.action) {
          this.config.action(this.form);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
