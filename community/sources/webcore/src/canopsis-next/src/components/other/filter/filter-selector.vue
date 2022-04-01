<template lang="pug">
  v-select(
    v-field="value",
    :items="preparedFilters",
    :label="label",
    :multiple="isMultiple",
    :disabled="!hasAccessToListFilters && !hasAccessToUserFilter",
    item-text="title",
    item-value="_id",
    return-object,
    clearable
  )
    template(#prepend-item="")
      v-layout.pl-3
        v-flex(v-show="!hideSelect")
          c-enabled-field(
            :value="isMultiple",
            :label="$t('filterSelector.fields.mixFilters')",
            hide-details,
            @input="updateIsMultipleFlag"
          )
        v-flex(v-show="!hideSelect && isMultiple")
          c-operator-field(:value="condition", @input="updateCondition")
      v-divider.mt-3

    template(#item="{ parent, item, tile }")
      v-list-tile-action(v-if="isMultiple", @click.stop="parent.$emit('select', item)")
        v-checkbox(:input-value="tile.props.value", :color="parent.color")
      v-list-tile-content
        v-list-tile-title
          span {{ item.title }}
</template>

<script>
import { isEmpty } from 'lodash';

import { FILTER_DEFAULT_VALUES } from '@/constants';

import { formMixin } from '@/mixins/form';

import FiltersListForm from '@/components/forms/filters/filters-list-form.vue';

export default {
  components: { FiltersListForm },
  mixins: [formMixin],
  props: {
    value: {
      type: [Object, Array],
      default: () => null,
    },
    filters: {
      type: Array,
      default: () => [],
    },
    lockedFilters: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    condition: {
      type: String,
      default: FILTER_DEFAULT_VALUES.condition,
    },
    hideSelect: { // TODO: remove it
      type: Boolean,
      default: false,
    },
    hasAccessToListFilters: {
      type: Boolean,
      default: false,
    },
    hasAccessToUserFilter: {
      type: Boolean,
      default: true,
    },
    withEntity: {
      type: Boolean,
      default: false,
    },
    withPbehavior: {
      type: Boolean,
      default: false,
    },
    withAlarm: {
      type: Boolean,
      default() {
        return !this.withEntity && !this.withPbehavior;
      },
    },
  },
  computed: {
    isMultiple() {
      return Array.isArray(this.value);
    },

    preparedFilters() {
      const preparedFilters = this.hasAccessToUserFilter ? [...this.filters] : [];
      const preparedLockedFilters = this.lockedFilters.map(filter => ({ ...filter, locked: true }));

      if (preparedFilters.length && preparedLockedFilters.length) {
        return preparedFilters.concat({ divider: true }, preparedLockedFilters);
      }

      if (preparedFilters.length) {
        return preparedFilters;
      }

      return preparedLockedFilters;
    },
  },
  methods: {
    updateIsMultipleFlag(checked) {
      const isValueArray = Array.isArray(this.value);

      if (checked && !isValueArray) {
        this.updateModel(!isEmpty(this.value) ? [this.value] : []);
      } else if (!checked && isValueArray) {
        this.updateModel(!isEmpty(this.value[0]) ? this.value[0] : null);
      }
    },

    updateCondition(newCondition) {
      this.$emit('update:condition', newCondition);
    },

    // TODO: remove
    isFilterEqual(firstFilter, secondFilter) {
      return firstFilter.title === secondFilter.title && firstFilter.filter === secondFilter.filter;
    },
  },
};
</script>
