<template>
  <v-form
    v-click-outside.zIndex="clickOutsideDirective"
    class="pbehavior-form"
    @submit.prevent="submitHandler"
  >
    <pbehavior-form
      v-model="form"
      :no-pattern="!!entityPattern"
      class="py-3"
    />
    <v-layout
      class="pbehavior-form__actions"
      justify-end
    >
      <v-btn
        v-show="pbehavior"
        :outlined="$system.dark"
        class="error"
        @click="remove"
      >
        {{ $t('common.delete') }}
      </v-btn>
      <v-btn
        depressed
        text
        @click="cancel"
      >
        {{ $t('common.cancel') }}
      </v-btn>
      <v-btn
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
import { cloneDeep } from 'lodash';
import dependentMixin from 'vuetify/lib/mixins/dependent';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { calendarEventToPbehaviorForm, formToCalendarEvent } from '@/helpers/entities/pbehavior/form';
import { isOmitEqual } from '@/helpers/collection';
import { getMenuClassByCalendarEvent } from '@/helpers/calendar/calendar';

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
    event: {
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
      form: calendarEventToPbehaviorForm(this.event, this.entityPattern, this.$system.timezone),
    };
  },
  computed: {
    pbehavior() {
      return this.event?.data?.pbehavior;
    },

    clickOutsideDirective() {
      const selectorsForInclude = [
        '.c-calendar__today-btn',
        '.c-calendar__pagination',
        '.c-calendar__menu-right',
        '.v-event',
        `.${getMenuClassByCalendarEvent(this.event.id)}`,
      ];

      return {
        handler: this.cancel,
        include: () => [
          ...this.getOpenDependentElements(),
          ...document.querySelectorAll(selectorsForInclude.join(',')),
        ],
        closeConditional: () => true,
      };
    },
  },
  mounted() {
    this.cacheForm();
  },
  beforeDestroy() {
    if (this.manualClose) {
      // eslint-disable-next-line vue/no-mutating-props
      delete this.event.data.cachedForm;
    } else {
      this.cacheForm();
    }
  },
  methods: {
    cacheForm() {
      // eslint-disable-next-line vue/no-mutating-props
      this.event.data.cachedForm = cloneDeep(this.form);
    },

    async submitHandler(event) {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        const calendarEvent = formToCalendarEvent(this.form, this.event, this.$system.timezone);

        this.manualClose = true;

        this.$emit('submit', calendarEvent, event);
      }
    },

    cancel(event) {
      const { cachedForm } = this.event.data;

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
    width: 100%;

    &__actions {
      gap: 6px;
    }
  }
</style>
