<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
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
import { MODALS } from '@/constants';

import RRuleForm from '@/components/forms/rrule.vue';
import PbehaviorExceptionForm from '@/components/other/pbehavior/calendar/partials/pbehavior-exception-form.vue';

import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createRRule,
  $_veeValidate: {
    validator: 'new',
  },
  inject: ['$system'],
  components: {
    PbehaviorExceptionForm,
    RRuleForm,
    ModalWrapper,
  },
  mixins: [
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    const { rrule, exdates, exceptions } = this.modal.config;

    return {
      form: {
        rrule: rrule || '',
        exdates: exdates || [],
        exceptions: exceptions || [],
      },
    };
  },
  methods: {
    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        if (this.config.action) {
          const { rrule, exdates, exceptions } = this.form;

          this.config.action({
            rrule,
            exdates,
            exceptions,
          });
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
