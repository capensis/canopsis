<template lang="pug">
  v-card
    v-card-text
      v-layout(row, align-start)
        v-icon.draggable.ml-0.mr-3.mt-3.action-drag-handler drag_indicator
        v-layout(column)
          v-layout(row)
            v-select(
              v-field="form.type",
              :items="eventFilterActionTypes",
              :label="$t('common.type')"
            )
            v-btn.mr-0(icon, @click="remove")
              v-icon(color="error") delete
          v-expand-transition
            event-filter-enrichment-action-form-type-info(v-if="form.type", :type="form.type")
          v-text-field(
            v-field="form.description",
            :label="$t('common.description')",
            key="description"
          )
          v-layout
            v-flex(xs5)
              c-name-field(v-field="form.name", key="name", required)
            v-flex(xs7)
              v-text-field.ml-2(
                v-if="isStringValueType",
                v-field="form.value",
                v-validate="'required'",
                :label="$t('common.value')",
                :error-messages="errors.collect('value')",
                key="from",
                name="value"
              )
                template(#append="")
                  c-help-icon(icon="help", :text="$t('eventFilter.tooltips.copyFromHelp')", left)
              c-mixed-field.ml-2(
                v-else,
                v-field="form.value",
                :label="$t('common.value')",
                key="value"
              )
</template>

<script>
import { EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES } from '@/constants';

import EventFilterEnrichmentActionFormTypeInfo from './event-filter-enrichment-action-form-type-info.vue';

export default {
  inject: ['$validator'],
  components: { EventFilterEnrichmentActionFormTypeInfo },
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

    isStringValueType() {
      return [
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copy,
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copyToEntityInfo,
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setFieldFromTemplate,
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromTemplate,
      ].includes(this.form.type);
    },
  },
  watch: {
    'form.type': function typeWatcher() {
      this.errors.clear();
    },
  },
  methods: {
    remove() {
      this.$emit('remove', this.form);
    },
  },
};
</script>
