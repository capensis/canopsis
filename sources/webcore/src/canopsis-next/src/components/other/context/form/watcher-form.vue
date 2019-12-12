<template lang="pug">
  v-tabs(slider-color="primary")
    v-tab.validation-header(
      v-for="tab in tabs",
      :key="tab.name",
      :class="tab.class"
    ) {{ tab.name }}
    v-tab-item
      div
        v-layout(wrap, justify-center)
          v-flex(xs11)
            v-text-field(
              v-field="form.name",
              v-validate="'required'",
              :label="$t('modals.createWatcher.displayName')",
              :error-messages="errors.collect('name')",
              name="name"
            )
        v-layout(wrap, justify-center)
          v-flex(xs11)
            template(v-if="stack === $constants.CANOPSIS_STACK.go")
              v-textarea(
                v-field="form.output_template",
                v-validate="'required'",
                :label="$t('modals.createWatcher.outputTemplate')",
                :error-messages="errors.collect('output_template')",
                name="output_template"
              )
    v-tab-item
      v-card
        v-card-text
          patterns-list(v-if="stack === $constants.CANOPSIS_STACK.go", v-field="form.entities")
          filter-editor(v-else, ref="filterEditor", v-field="form.mfilter", required)
    v-tab-item
      manage-infos(v-field="form.infos")
</template>

<script>
import { CANOPSIS_STACK } from '@/constants';

import FilterEditor from '@/components/other/filter/editor/filter-editor.vue';
import PatternsList from '@/components/other/shared/patterns-list/patterns-list.vue';
import ManageInfos from '@/components/other/context/manage-infos.vue';

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
    stack: {
      type: String,
      default: null,
    },
  },
  data() {
    return {
      hasFilterEditorAnyError: false,
    };
  },
  computed: {
    tabs() {
      const patternsTab = { name: this.$t('eventFilter.pattern') };
      const filterEditorTab = {
        name: this.$t('common.filter'),
        class: { 'error--text': this.hasFilterEditorAnyError },
      };

      return [
        { name: this.$t('modals.createEntity.fields.form') },
        this.stack === CANOPSIS_STACK.go ? patternsTab : filterEditorTab,
        { name: this.$t('modals.createEntity.fields.manageInfos') },
      ];
    },
  },
  mounted() {
    if (this.$refs.filterEditor) {
      this.$watch(() => this.$refs.filterEditor.hasAnyError, (value) => {
        this.hasFilterEditorAnyError = value;
      });
    }
  },
};
</script>

<style>

</style>
