<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.createAction.create.title') }}
    v-card-text
      v-form
        v-text-field(
        v-validate="'required'",
        v-model="form._id",
        label="Id",
        :error-messages="errors.collect('id')",
        name="id",
        :disabled="modal.config.item"
        )
        v-select(
        v-validate="'required'",
        v-model="form.type",
        label="Type",
        :items="actionTypes"
        :error-messages="errors.collect('type')",
        name="type"
        )
        v-tabs(centered, slider-color="primary")
          v-tab {{ $t('modals.createAction.tabs.general') }}
          v-tab-item
            template(v-if="form.type === $constants.ACTION_TYPES.snooze")
              v-textarea(
              v-model="snoozeParameters.message",
              :label="$t('modals.createAction.fields.message')",
              )
              duration-field(:duration="snoozeParameters.duration")
            template(v-if="form.type === $constants.ACTION_TYPES.pbehavior")
              pbehavior-form(v-model="pbehaviorParameters", :author="$constants.ACTION_AUTHOR", :noFilter="true")
          v-tab {{ $t('modals.createAction.tabs.hook') }}
          v-tab-item
            webhook-form-hook-tab(
            v-model="form.hook",
            :operators="$constants.WEBHOOK_EVENT_FILTER_RULE_OPERATORS",
            )
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

import PbehaviorForm from '@/components/other/pbehavior/form/pbehavior-form.vue';
import WebhookFormHookTab from '@/components/other/webhook/form/tabs/webhook-form-hook-tab.vue';
import DurationField from '@/components/forms/fields/duration.vue';

export default {
  name: MODALS.createAction,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    PbehaviorForm,
    WebhookFormHookTab,
    DurationField,
  },
  mixins: [modalInnerMixin],
  data() {
    // Default form
    let form = {
      _id: uuid('action'),
      type: ACTION_TYPES.snooze,
      hook: {
        event_patterns: [],
        alarm_patterns: [],
        entity_patterns: [],
        triggers: [],
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

    // Default 'snooze' action parameters
    let snoozeParameters = {
      message: '',
      duration: {
        duration: 1,
        durationType: DURATION_UNITS.minute.value,
      },
    };

    if (this.modal.config.item) {
      const { item } = this.modal.config;

      form = omit(item, ['parameters']);

      // If editing a 'pbehavior' action, prepare pbehavior's data. If editing a 'snooze' action copy snooze's data
      if (item.type === ACTION_TYPES.pbehavior) {
        pbehaviorParameters.general = omit(this.$options.filters.pbehaviorToForm(item.parameters), ['filter']);
        pbehaviorParameters.comments = this.$options.filters.commentsToPbehaviorComments(item.parameters.comments);
        pbehaviorParameters.exdate = this.$options.filters.exdateToPbehaviorExdate(item.parameters.exdate);
      } else if (item.type === ACTION_TYPES.snooze) {
        snoozeParameters = { ...snoozeParameters, ...item.parameters };
      }
    }

    return {
      form,
      pbehaviorParameters,
      snoozeParameters,
      actionTypes: Object.values(ACTION_TYPES),
      availableTriggers: Object.values(WEBHOOK_TRIGGERS),
    };
  },
  methods: {
    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        if (this.config.action) {
          let data = { ...this.form };

          const patternsCondition = value => !value || !value.length;

          data = unsetInSeveralWithConditions(data, {
            'hook.event_patterns': patternsCondition,
            'hook.alarm_patterns': patternsCondition,
            'hook.entity_patterns': patternsCondition,
          });

          if (this.form.type === ACTION_TYPES.snooze) {
            const duration = moment.duration(
              parseInt(this.snoozeParameters.duration.duration, 10),
              this.snoozeParameters.duration.durationType,
            ).asSeconds();

            data.parameters = { ...this.snoozeParameters, duration };
          } else if (this.form.type === ACTION_TYPES.pbehavior) {
            const pbehavior = this.$options.filters.formToPbehavior(this.pbehaviorParameters.general);

            pbehavior.comments = this.$options.filters.commentsToPbehaviorComments(this.pbehaviorParameters.comments);
            pbehavior.exdate = this.$options.filters.exdateToPbehaviorExdate(this.pbehaviorParameters.exdate);

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
