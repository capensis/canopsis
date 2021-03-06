<template lang="pug">
  v-form.pa-3.pbehavior-form(v-click-outside.zIndex="clickOutsideDirective", @submit.prevent="submitHandler")
    pbehavior-form(v-model="form", :noFilter="!!filter")
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
import { get, cloneDeep } from 'lodash';
import dependentMixin from 'vuetify/es5/mixins/dependent';

import {
  calendarEventToPbehaviorForm,
  formToCalendarEvent,
} from '@/helpers/forms/planning-pbehavior';

import { MODALS } from '@/constants';

import { isOmitEqual } from '@/helpers/validators/is-omit-equal';
import { getMenuClassByCalendarEvent } from '@/helpers/calendar/dayspan';

import authMixin from '@/mixins/auth';

import PbehaviorForm from '@/components/other/pbehavior/calendar/partials/pbehavior-form.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  inject: ['$system'],
  components: { PbehaviorForm },
  mixins: [authMixin, dependentMixin],
  props: {
    calendarEvent: {
      type: Object,
      required: false,
    },
    filter: {
      type: Object,
      required: false,
    },
  },
  data() {
    return {
      manualClose: false,
      form: calendarEventToPbehaviorForm(this.calendarEvent, this.filter, this.$system.timezone),
    };
  },
  computed: {
    pbehavior() {
      return get(this.calendarEvent, 'data.pbehavior');
    },

    clickOutsideDirective() {
      const selectorsForInclude = [
        '.ds-calendar-app-action',
        `.${getMenuClassByCalendarEvent(this.calendarEvent)}`,
      ];

      return {
        handler: this.cancel,
        include: () => [
          ...this.getOpenDependentElements(),
          ...document.querySelectorAll(selectorsForInclude.join(',')),
        ],
      };
    },
  },
  mounted() {
    this.cacheForm();
  },
  beforeDestroy() {
    if (this.manualClose) {
      delete this.calendarEvent.data.cachedForm;
    } else {
      this.cacheForm();
    }
  },
  methods: {
    cacheForm() {
      this.calendarEvent.data.cachedForm = cloneDeep(this.form);
    },

    async submitHandler(event) {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        this.form.author = this.currentUser._id;

        const calendarEvent = formToCalendarEvent(this.form, this.calendarEvent, this.$system.timezone);

        this.manualClose = true;

        this.$emit('submit', calendarEvent, event);
      }
    },

    cancel(event) {
      const { cachedForm } = this.calendarEvent.data;

      if (isOmitEqual(cachedForm, this.form, ['_id'])) {
        return this.close(event, true);
      }

      return this.$modals.show({
        name: MODALS.confirmation,
        config: {
          text: this.$t('modals.createPbehavior.cancelConfirmation'),
          action: () => this.close(event, true),
        },
      });
    },

    remove(event) {
      this.$emit('remove', this.pbehavior);
      this.close(event);
    },

    close(event, manualClose = false) {
      this.manualClose = manualClose;

      this.$emit('close', event);
    },
  },
};
</script>

<style lang="scss" scoped>
  .pbehavior-form {
    overflow: auto;
    width: 500px;
    max-height: 615px;
  }
</style>
