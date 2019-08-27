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
              pbehavior-form(v-model="pbehaviorParameters", noFilter)
              pbehavior-exdates-form.mt-2(v-show="pbehaviorParameters.rrule", v-model="pbehaviorParameters.exdate")
              pbehavior-comments-form.mt-2(v-model="pbehaviorParameters.comments")
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
import { MODALS, ACTION_TYPES, ACTION_AUTHOR, WEBHOOK_TRIGGERS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import PbehaviorForm from '@/components/other/pbehavior/form/pbehavior-form.vue';
import PbehaviorExdatesForm from '@/components/other/pbehavior/form/pbehavior-exdates-form.vue';
import PbehaviorCommentsForm from '@/components/other/pbehavior/form/pbehavior-comments-form.vue';
import PatternsList from '@/components/other/shared/patterns-list/patterns-list.vue';

export default {
  name: MODALS.createAction,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    PbehaviorForm,
    PbehaviorExdatesForm,
    PbehaviorCommentsForm,
    PatternsList,
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      form: {
        _id: '',
        type: ACTION_TYPES.snooze,
        hook: {
          event_patterns: [],
          triggers: [],
        },
      },
      pbehaviorParameters: {
        name: '',
        tstart: new Date(),
        tstop: new Date(),
        rrule: '',
        reason: '',
        type_: '',
        comments: [],
        exdate: [],
      },
      snoozeParameters: {
        message: '',
        duration: 60,
      },
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
            data.parameters = { ...this.pbehaviorParameters };
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
