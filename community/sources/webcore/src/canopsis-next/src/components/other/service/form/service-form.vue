<template lang="pug">
  div
    v-layout(column)
      c-name-field(v-field="form.name")
      v-layout
        v-flex.pr-3(xs6)
          c-entity-category-field(v-field="form.category", addable, required)
        v-flex.pr-3(xs4)
          c-entity-state-field(
            v-field="form.sli_avail_state",
            :label="$t('service.availabilityState')",
            required
          )
        v-flex(xs2)
          c-impact-level-field(v-field="form.impact_level", required)
      c-coordinates-field(v-field="form.coordinates", row)
      v-textarea(
        v-field="form.output_template",
        v-validate="'required'",
        :label="$t('service.outputTemplate')",
        :error-messages="errors.collect('output_template')",
        name="output_template"
      )
      c-enabled-field(v-field="form.enabled")
    v-tabs(slider-color="primary", centered)
      v-tab(:class="{ 'error--text': errors.has('entity_patterns') }") {{ $t('common.entityPatterns') }}
      v-tab-item
        c-patterns-field.mt-2(v-field="form.patterns", :entity-attributes="entityAttributes", with-entity)
      v-tab.validation-header(:disabled="advancedJsonWasChanged") {{ $t('entity.manageInfos') }}
      v-tab-item
        manage-infos(v-field="form.infos")
</template>

<script>
import { get } from 'lodash';

import { ENTITY_PATTERN_FIELDS } from '@/constants';

import ManageInfos from '@/components/widgets/context/manage-infos.vue';

export default {
  inject: ['$validator'],
  components: {
    ManageInfos,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    hasFilterEditorAnyError() {
      return this.errors.has('advancedJson') || this.errors.has('filter');
    },

    advancedJsonWasChanged() {
      return get(this.fields, ['advancedJson', 'changed']);
    },

    entityAttributes() {
      return [
        {
          value: ENTITY_PATTERN_FIELDS.lastEventDate,
          options: { disabled: true },
        },
        {
          value: ENTITY_PATTERN_FIELDS.connector,
          options: { disabled: true },
        },
        {
          value: ENTITY_PATTERN_FIELDS.componentInfos,
          options: { disabled: true },
        },
      ];
    },
  },
};
</script>
