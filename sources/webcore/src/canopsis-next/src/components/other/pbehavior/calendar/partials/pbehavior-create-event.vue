<template lang="pug">
  v-form.pa-3.pbehavior-form(v-click-outside.zIndex="clickOutsideDirective", @submit.prevent="submitHandler")
    pbehavior-form(v-model="form")
    v-layout(row, justify-end)
      v-btn.error(
        v-show="pbehavior",
        @click="remove"
      ) {{ $t('common.delete') }}
      v-btn.mr-0.mb-0(
        depressed,
        flat,
        @click="cancel"
      ) {{ $t('common.cancel') }}
      v-btn.mr-0.mb-0.primary.white--text(type="submit") {{ $t('common.submit') }}
</template>

<script>
import { get } from 'lodash';
import dependentMixin from 'vuetify/es5/mixins/dependent';

import {
  calendarEventToPbehaviorForm,
  formToCalendarEvent,
} from '@/helpers/forms/planning-pbehavior';

import { MODALS } from '@/constants';

import authMixin from '@/mixins/auth';
import { isOmitEqual } from '@/helpers/is-omit-equal';

import PbehaviorForm from '@/components/other/pbehavior/calendar/partials/pbehavior-form.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: { PbehaviorForm },
  mixins: [authMixin, dependentMixin],
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

    clickOutsideDirective() {
      return {
        handler: this.cancel,
        include: () => [this.$el, ...this.getOpenDependentElements()],
      };
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

    cancel() {
      const oldPbehaviorForm = calendarEventToPbehaviorForm(this.calendarEvent);

      if (isOmitEqual(oldPbehaviorForm, this.form, ['_id']) && this.pbehavior) {
        return this.$emit('close');
      }

      return this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => this.$emit('close'),
        },
      });
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
