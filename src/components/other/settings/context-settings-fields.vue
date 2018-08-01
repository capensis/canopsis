<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.title")
      v-divider
      field-default-column-sort(v-model="settings.defaultSortColumn")
      v-divider
      field-context-entities-types-filter(v-model="settings.selectedTypes")
      v-divider
    v-btn(@click="submit", color="green darken-4 white--text", depressed, fixed, right) {{ $t('common.save') }}
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';

import FieldTitle from '@/components/other/settings/fields/title.vue';
import FieldDefaultColumnSort from '@/components/other/settings/fields/default-column-sort.vue';
import FieldContextEntitiesTypesFilter from '@/components/other/settings/fields/context-entities-types-filter.vue';
import entitiesWidgetMixin from '@/mixins/entities/widget';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';

/**
 * Component to regroup the entities list settings fields
 *
 * @prop {Object} widget - active widget
 * @prop {bool} isNew - is widget new
 *
 * @event closeSettings#click
 */
export default {
  components: {
    FieldTitle,
    FieldDefaultColumnSort,
    FieldContextEntitiesTypesFilter,
  },
  mixins: [
    entitiesWidgetMixin,
    entitiesUserPreferenceMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    isNew: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      settings: {
        title: this.widget.title,
        defaultSortColumn: cloneDeep(this.widget.default_sort_column),
        selectedTypes: [],
      },
    };
  },
  created() {
    this.settings.selectedTypes = this.userPreference.widget_preferences.selectedTypes || [];
  },
  methods: {
    async submit() {
      const widget = {
        ...this.widget,
        title: this.settings.title,
        default_sort_column: this.settings.defaultSortColumn,
      };

      const userPreference = {
        ...this.userPreference,
        widget_preferences: {
          ...this.userPreference.widget_preferences,
          selectedTypes: this.settings.selectedTypes,
        },
      };

      const actions = [this.createUserPreference({ userPreference })];

      if (this.isNew) {
        actions.push(this.createWidget({ widget }));
      } else {
        actions.push(this.updateWidget({ widget }));
      }

      await this.$store.dispatch('query/startPending', { id: widget.id }); // TODO: fix it
      await Promise.all(actions);
      await this.$store.dispatch('query/stopPending', { id: widget.id }); // TODO: fix it

      this.$emit('closeSettings');
    },
  },
};
</script>

<style scoped>
  .closeIcon:hover {
    cursor: pointer;
  }
</style>
