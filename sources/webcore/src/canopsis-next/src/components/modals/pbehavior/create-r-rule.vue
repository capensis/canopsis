<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.createRrule.title') }}
      template(slot="text")
        r-rule-form(v-model="form.rrule")
        pbehavior-exception-dates-form(v-if="form.rrule", v-model="form.exdate")
      template(slot="actions")
        v-btn(
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(type="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { MODALS } from '@/constants';

import RRuleForm from '@/components/forms/rrule.vue';
import PbehaviorExceptionDatesForm from '@/components/other/pbehavior/calendar/partials/pbehavior-exception-dates-form.vue';

import modalInnerMixin from '@/mixins/modal/inner';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createRRule,
  components: {
    PbehaviorExceptionDatesForm,
    RRuleForm,
    ModalWrapper,
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      form: {
        rrule: this.modal.config.rrule || '',
        exdate: this.modal.config.exdate || [],
      },
    };
  },
  methods: {
    submit() {
      if (this.config.action) {
        this.config.action(this.form);
      }

      this.$modals.hide();
    },
  },
};
</script>
