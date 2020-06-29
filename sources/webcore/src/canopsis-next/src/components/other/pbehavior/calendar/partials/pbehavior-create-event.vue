<template lang="pug">
  v-form.pa-3.pbehavior-form(@submit.prevent="submitHandler")
    pbehavior-form(v-model="form")
    v-layout(row, justify-end)
      v-btn.mr-0.mb-0(
        depressed,
        flat,
        @click="$emit('close')"
      ) {{ $t('common.cancel') }}
      v-btn.mr-0.mb-0.primary.white--text(type="submit") {{ $t('common.submit') }}
</template>

<script>
import {
  formToPbehaviorCalendarEvent,
  pbehaviorCalendarEventToForm,
} from '@/helpers/forms/pbehavior';

import authMixin from '@/mixins/auth';

import PbehaviorForm from '@/components/other/pbehavior/calendar/partials/pbehavior-form.vue';

export default {
  components: { PbehaviorForm },
  mixins: [authMixin],
  inject: ['$validator'],
  props: {
    calendarEvent: {
      type: Object,
      required: false,
    },
  },
  data() {
    return {
      form: pbehaviorCalendarEventToForm(this.calendarEvent),
    };
  },
  methods: {
    async submitHandler() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        const pbehavior = formToPbehaviorCalendarEvent(this.form);

        pbehavior.author = this.currentUser._id;

        this.$emit('submit', pbehavior);
      }
    },
  },
};
</script>

<style lang="scss" scoped>
  .pbehavior-form {
    overflow: auto;
    width: 500px;
    max-height: 600px;
  }
</style>
