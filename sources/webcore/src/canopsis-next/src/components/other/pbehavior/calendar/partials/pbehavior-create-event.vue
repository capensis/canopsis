<template lang="pug">
  v-form.pa-3.pbehavior-form(@submit.prevent="submitHandler")
    pbehavior-form(v-model="form")
    v-layout(row, justify-end)
      v-btn.error(
        v-show="pbehavior",
        @click="remove"
      ) {{ $t('common.delete') }}
      v-btn.mr-0.mb-0(
        depressed,
        flat,
        @click="$emit('close')"
      ) {{ $t('common.cancel') }}
      v-btn.mr-0.mb-0.primary.white--text(type="submit") {{ $t('common.submit') }}
</template>

<script>
import { get } from 'lodash';

import {
  calendarEventToPbehaviorForm,
  formToCalendarEvent,
} from '@/helpers/forms/planning-pbehavior';

import authMixin from '@/mixins/auth';

import PbehaviorForm from '@/components/other/pbehavior/calendar/partials/pbehavior-form.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: { PbehaviorForm },
  mixins: [authMixin],
  props: {
    calendarEvent: {
      type: Object,
      required: false,
    },
  },
  data() {
    return {
      form: calendarEventToPbehaviorForm(this.calendarEvent),
    };
  },
  computed: {
    pbehavior() {
      return get(this.calendarEvent, 'data.pbehavior');
    },
  },
  methods: {
    async submitHandler() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        this.form.author = this.currentUser._id;

        const calendarEvent = formToCalendarEvent(this.form, this.calendarEvent);

        this.$emit('submit', calendarEvent);
      }
    },

    remove() {
      this.$emit('close');
      this.$emit('remove', this.pbehavior);
    },
  },
};
</script>

<style lang="scss" scoped>
  .pbehavior-form {
    overflow: auto;
    width: 500px;
    max-height: 555px; // For validation errors
  }
</style>
