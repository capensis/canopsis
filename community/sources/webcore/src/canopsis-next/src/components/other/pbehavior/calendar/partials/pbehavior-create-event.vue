<template>
  <v-form
    class="pa-3 pbehavior-form"
    v-click-outside.zIndex="clickOutsideDirective"
    @submit.prevent="submitHandler"
  >
    <pbehavior-form
      v-model="form"
      :no-pattern="!!entityPattern"
    />
    <v-layout justify-end="justify-end">
      <v-btn
        class="error"
        v-show="pbehavior"
        :outlined="$system.dark"
        @click="remove"
      >
        {{ $t('common.delete') }}
      </v-btn>
      <v-btn
        class="mr-0 mb-0"
        depressed="depressed"
        text
        @click="cancel"
      >
        {{ $t('common.cancel') }}
      </v-btn>
      <v-btn
        class="mr-0 mb-0"
        :disabled="errors.any()"
        color="primary"
        type="submit"
      >
        {{ $t('common.submit') }}
      </v-btn>
    </v-layout>
  </v-form>
</template>

<script>
import { get, cloneDeep } from 'lodash';
import dependentMixin from 'vuetify/es5/mixins/dependent';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { calendarEventToPbehaviorForm, formToCalendarEvent } from '@/helpers/entities/pbehavior/form';
import { isOmitEqual } from '@/helpers/collection';
import { getMenuClassByCalendarEvent } from '@/helpers/calendar/dayspan';

import PbehaviorForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-form.vue';

export default {
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  inject: ['$system'],
  components: { PbehaviorForm },
  mixins: [dependentMixin],
  props: {
    calendarEvent: {
      type: Object,
      required: false,
    },
    entityPattern: {
      type: Array,
      required: false,
    },
  },
  data() {
    return {
      manualClose: false,
      form: calendarEventToPbehaviorForm(this.calendarEvent, this.entityPattern, this.$system.timezone),
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
      // eslint-disable-next-line vue/no-mutating-props
      delete this.calendarEvent.data.cachedForm;
    } else {
      this.cacheForm();
    }
  },
  methods: {
    cacheForm() {
      // eslint-disable-next-line vue/no-mutating-props
      this.calendarEvent.data.cachedForm = cloneDeep(this.form);
    },

    async submitHandler(event) {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
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
    width: 100%;
  }
</style>
