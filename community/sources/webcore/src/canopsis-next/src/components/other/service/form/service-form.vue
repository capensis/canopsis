<template lang="pug">
  div
    div
      v-layout(wrap, justify-center)
        v-flex
          v-text-field(
            v-field="form.name",
            v-validate="'required'",
            :label="$t('service.fields.name')",
            :error-messages="errors.collect('name')",
            name="name"
          )
      v-layout
        v-flex(xs6)
          c-entity-category-field.mr-3(v-field="form.category", addable, required)
        v-flex(xs6)
          c-impact-level-field(v-field="form.impact_level", required)
      v-layout(wrap, justify-center)
        v-flex
          v-textarea(
            v-field="form.output_template",
            v-validate="'required'",
            :label="$t('service.fields.outputTemplate')",
            :error-messages="errors.collect('output_template')",
            name="output_template"
          )
      v-layout(wrap, justify-center)
        v-flex
          c-enabled-field(v-field="form.enabled")
    v-tabs(slider-color="primary", centered)
      v-tab(:class="{ 'error--text': errors.has('entity_patterns') }") {{ $t('eventFilter.pattern') }}
      v-tab-item
        patterns-list(v-field="form.entity_patterns", v-validate="'required'", name="entity_patterns")
      v-tab.validation-header(:disabled="advancedJsonWasChanged") {{ $t('entity.fields.manageInfos') }}
      v-tab-item
        manage-infos(v-field="form.infos")
</template>

<script>
import { get } from 'lodash';

import FilterEditor from '@/components/other/filter/editor/filter-editor.vue';
import PatternsList from '@/components/common/patterns-list/patterns-list.vue';
import ManageInfos from '@/components/widgets/context/manage-infos.vue';

export default {
  inject: ['$validator'],
  components: {
    FilterEditor,
    PatternsList,
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
  },
};
</script>