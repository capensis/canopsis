<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.createAction.create.title') }}
    v-card-text
      action-form(v-model="form")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(depressed, flat, @click="hideModal") {{ $t('common.cancel') }}
      v-btn.primary(:disabled="errors.any()", @click="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import moment from 'moment';
import { omit } from 'lodash';
import { MODALS, ACTION_TYPES, ACTION_AUTHOR, WEBHOOK_TRIGGERS, DURATION_UNITS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import uuid from '@/helpers/uuid';
import { unsetInSeveralWithConditions } from '@/helpers/immutable';

import ActionForm from '@/components/other/action/form/action-form.vue';

export default {
  name: MODALS.createAction,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    ActionForm,
  },
  filters: {
    actionToForm(action = {}) {
      const defaultHook = {
        event_patterns: [],
        alarm_patterns: [],
        entity_patterns: [],
        triggers: [],
      };

      // Default 'snooze' action parameters
      const snoozeParameters = {
        message: '',
        duration: {
          duration: 1,
          durationType: DURATION_UNITS.minute.value,
        },
      };

      // Default 'pbehavior' action parameters
      const pbehaviorParameters = {
        general: {
          name: '',
          tstart: new Date(),
          tstop: new Date(),
          rrule: null,
          reason: '',
          type_: '',
        },
        comments: [],
        exdate: [],
      };

      // Get basic action parameters
      const generalParameters = {
        _id: action._id || uuid('action'),
        type: action.type || ACTION_TYPES.pbehavior,
        hook: action.hook || defaultHook,
      };

      // If action's type is "snooze", get snooze parameters
      if (action.type === ACTION_TYPES.snooze) {
        let duration = {
          duration: 1,
          durationType: DURATION_UNITS.minute.value,
        };

        if (action.parameters && action.parameters.duration) {
          const durationUnits = Object.values(DURATION_UNITS).map(unit => unit.value);

          // Check for the lowest possible unit to convert the duration in.
          const foundUnit = durationUnits.find(unit => moment.duration(action.parameters.duration, 'seconds').as(unit) % 1 === 0);

          duration = {
            duration: moment.duration(action.parameters.duration, 'seconds').as(foundUnit),
            durationType: foundUnit,
          };

          snoozeParameters.duration = duration;
        }

        if (action.parameters && action.parameters.message) {
          snoozeParameters.message = action.parameters.message;
        }
      }

      // If action's type is "pbehavior", get pbehavior parameters
      if (action.type === ACTION_TYPES.pbehavior) {
        if (action.parameters) {
          pbehaviorParameters.general = omit(this.$options.filters.pbehaviorToForm(action.parameters), ['filter']);

          if (action.parameters.comments) {
            pbehaviorParameters.comments =
              this.$options.filters.commentsToPbehaviorComments(action.parameters.comments);
          }

          if (action.parameters.exdate) {
            pbehaviorParameters.exdate = this.$options.filters.exdateToPbehaviorExdate(action.parameters.exdate);
          }
        }
      }

      return {
        generalParameters,
        snoozeParameters,
        pbehaviorParameters,
      };
    },
  },
  mixins: [modalInnerMixin],
  data() {
    const { item } = this.modal.config;

    return {
      form: this.$options.filters.actionToForm(item),
      availableTriggers: Object.values(WEBHOOK_TRIGGERS),
    };
  },
  methods: {
    async submit() {
      const isBaseFormValid = await this.$validator.validateAll();
      const isPbehaviorFormValid = await this.$validator.validateAll('pbehavior');
      const isSnoozeFormValid = await this.$validator.validateAll('snooze');
      const isHookFormValid = await this.$validator.validateAll('hook');

      if (isBaseFormValid && isPbehaviorFormValid && isSnoozeFormValid && isHookFormValid) {
        if (this.config.action) {
          let data = { ...this.form.generalParameters };

          const patternsCondition = value => !value || !value.length;

          data = unsetInSeveralWithConditions(data, {
            'hook.event_patterns': patternsCondition,
            'hook.alarm_patterns': patternsCondition,
            'hook.entity_patterns': patternsCondition,
          });

          if (this.form.generalParameters.type === ACTION_TYPES.snooze) {
            const duration = moment.duration(
              parseInt(this.form.snoozeParameters.duration.duration, 10),
              this.form.snoozeParameters.duration.durationType,
            ).asSeconds();

            data.parameters = { ...this.form.snoozeParameters, duration };
          } else if (this.form.generalParameters.type === ACTION_TYPES.pbehavior) {
            const pbehavior = this.$options.filters.formToPbehavior(this.form.pbehaviorParameters.general);

            pbehavior.comments =
              this.$options.filters.commentsToPbehaviorComments(this.form.pbehaviorParameters.comments);
            pbehavior.exdate = this.$options.filters.exdateToPbehaviorExdate(this.form.pbehaviorParameters.exdate);

            data.parameters = { ...pbehavior };
          }

          data.parameters.author = ACTION_AUTHOR;

          await this.config.action(data);
        }

        this.hideModal();
      }
    },
  },
};
</script>
