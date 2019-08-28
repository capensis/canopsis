<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline Create Action
    v-card-text
      v-form
        v-text-field(
        v-validate="'required'",
        v-model="form._id",
        label="Id",
        :error-messages="errors.collect('id')",
        name="id",
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
          v-tab Général
          v-tab-item
            template(v-if="form.type === $constants.ACTION_TYPES.snooze")
              v-textarea(
              v-model="snoozeParameters.message",
              label="Message",
              )
              v-text-field(
              v-model="snoozeParameters.duration",
              label="Duration",
              type="number"
              )
            template(v-if="form.type === $constants.ACTION_TYPES.pbehavior")
              pbehavior-form(v-model="pbehaviorParameters")
          v-tab Hook
          v-tab-item
            v-select(
            v-model="form.hook.triggers",
            :items="availableTriggers",
            :label="$t('webhook.tabs.hook.fields.triggers')",
            multiple,
            chips
            )
            patterns-list(
            v-model="form.hook.event_patterns",
            :operators="$constants.WEBHOOK_EVENT_FILTER_RULE_OPERATORS",
            )
    v-divider
    v-layout.py-1(justify-end)
      v-btn(depressed, flat, @click="hideModal") {{ $t('common.cancel') }}
      v-btn.primary(:disabled="errors.any()", @click="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { omit } from 'lodash';
import { MODALS, ACTION_TYPES, ACTION_AUTHOR, WEBHOOK_TRIGGERS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import PbehaviorForm from '@/components/other/pbehavior/form/pbehavior-form.vue';
import PatternsList from '@/components/other/shared/patterns-list/patterns-list.vue';

export default {
  name: MODALS.createAction,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    PbehaviorForm,
    PatternsList,
  },
  mixins: [modalInnerMixin],
  data() {
    // Default form
    let form = {
      _id: '',
      type: ACTION_TYPES.snooze,
      hook: {
        event_patterns: [],
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
      duration: 60,
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
          const data = { ...this.form };

          if (this.form.type === ACTION_TYPES.snooze) {
            data.parameters = { ...this.snoozeParameters };
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
