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
import find from 'lodash/find';
import omit from 'lodash/omit';

import FieldTitle from '@/components/other/settings/fields/title.vue';
import FieldDefaultColumnSort from '@/components/other/settings/fields/default-column-sort.vue';
import FieldContextEntitiesTypesFilter from '@/components/other/settings/fields/context-entities-types-filter.vue';
import entitiesWidgetMixin from '@/mixins/entities/widget';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';
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
    entitiesContextEntityMixin,
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
    const filter = find(this.userPreference.widget_preferences.user_filters, { title: 'default_type_filter' });

    if (filter) {
      const filterData = JSON.parse(filter.filter);

      this.settings.selectedTypes = filterData.$or.map(({ type }) => type);
    }
  },
  methods: {
    async submit() {
      const widget = {
        ...this.widget,
        title: this.settings.title,
        default_sort_column: this.settings.defaultSortColumn,
      };

      const userPreference = omit(this.userPreference, ['crecord_creation_time', 'crecord_write_time', 'enable']);
      const userFilters = userPreference.widget_preferences.user_filters || [];
      const defaultTypeFilterIndex = userFilters.findIndex(filter => filter.title === 'default_type_filter');

      if (defaultTypeFilterIndex >= 0) {
        if (this.settings.selectedTypes.length) {
          userPreference.widget_preferences.user_filters[defaultTypeFilterIndex] = {
            title: 'default_type_filter',
            filter: JSON.stringify({
              $or: this.settings.selectedTypes.map(type => ({ type })),
            }),
          };
        } else {
          delete userPreference.widget_preferences.user_filters[defaultTypeFilterIndex];
        }
      } else if (this.settings.selectedTypes.length) {
        userPreference.widget_preferences.user_filters = [...userFilters, {
          title: 'default_type_filter',
          filter: JSON.stringify({
            $or: this.settings.selectedTypes.map(type => ({ type })),
          }),
        }];
      }

      await this.createUserPreference({ userPreference });

      if (this.isNew) {
        await this.createWidget({ widget });
      } else {
        await this.updateWidget({ widget });
      }

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
