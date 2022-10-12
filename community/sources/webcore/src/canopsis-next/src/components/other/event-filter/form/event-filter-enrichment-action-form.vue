<template lang="pug">
  v-card
    v-card-text
      v-layout(row, align-start)
        v-icon.draggable.ml-0.mr-3.mt-3.action-drag-handler drag_indicator
        v-layout(column)
          v-layout(row)
            c-select-field(
              v-model="form.type",
              :items="eventFilterActionTypes",
              :label="$t('common.type')"
            )
            v-btn.mr-0(icon, @click="removeAction")
              v-icon(color="error") delete
          v-expand-transition
            event-filter-action-form-type-info(v-if="form.type", :type="form.type")
          v-layout
            v-flex(xs5)
              c-name-field(v-field="form.name", key="name")
            v-flex(xs7)
              v-text-field(
                v-if="isCopyActionType",
                v-field="form.value",
                v-validate="'required'",
                :label="$t('common.value')",
                :error-messages="errors.collect('value')",
                key="from",
                name="value"
              )
                v-tooltip(#append="", left)
                  v-icon(#activator="{ on, bind }", v-on="on", v-bind="bind") help
                  div(v-html="$t('eventFilter.tooltips.copyFromHelp')")
              c-mixed-field.ml-2(
                v-else,
                v-field="form.value",
                :label="$t('common.value')",
                key="value"
              )
</template>

<script>
import { EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES } from '@/constants';

import EventFilterActionFormTypeInfo from './partials/event-filter-action-form-type-info.vue';

export default {
  inject: ['$validator'],
  components: { EventFilterActionFormTypeInfo },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  computed: {
    eventFilterActionTypes() {
      return Object.values(EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES).map(value => ({
        value,

        text: this.$t(`eventFilter.actionsTypes.${value}.text`),
      }));
    },

    isCopyActionType() {
      return this.form.type === EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copy;
    },
  },
  watch: {
    'form.type': function typeWatcher() {
      this.errors.clear();
    },
  },
  methods: {
    removeAction() {
      this.$emit('remove', this.form);
    },
  },
};
</script>
