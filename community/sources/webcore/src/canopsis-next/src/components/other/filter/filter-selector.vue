<template lang="pug">
  v-select(
    v-field="value",
    :items="preparedFilters",
    :label="label",
    :multiple="isMultiple",
    :disabled="disabled",
    item-text="title",
    item-value="_id",
    clearable
  )
    template(#prepend-item="", v-if="!hidePrepend")
      v-layout.pl-3
        v-flex
          c-enabled-field(
            :value="isMultiple",
            :label="$t('filterSelector.fields.mixFilters')",
            hide-details,
            @input="updateIsMultipleFlag"
          )
        v-flex(v-show="isMultiple")
          c-operator-field(:value="condition", @input="updateCondition")
      v-divider.mt-3

    template(#item="{ parent, item, tile }")
      v-list-tile-action(v-if="isMultiple", @click.stop="parent.$emit('select', item)")
        v-checkbox(:input-value="tile.props.value", :color="parent.color")
      v-list-tile-content
        v-list-tile-title
          span {{ item.title }}
          v-icon.ml-2(
            v-if="!hideIcon",
            :color="tile.props.value ? parent.color : ''",
            small
          ) {{ item.is_private ? 'person' : 'lock' }}
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
    hidePrepend: {
      type: Boolean,
      default: false,
    },
    hideIcon: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
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
      const preparedFilters = [...this.filters];

      if (!this.lockedFilters.length) {
        return preparedFilters;
      }

      if (preparedFilters.length) {
        preparedFilters.push({ divider: true });
      }

      preparedFilters.push(...this.lockedFilters);

      return preparedFilters;
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
