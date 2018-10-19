<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.rowGridSize.title') }}
    v-container
      v-combobox(
      v-model="row",
      @blur="blurRow",
      :items="availableRows",
      :label="$t('settings.rowGridSize.fields.row')",
      :search-input.sync="search",
      data-vv-name="row",
      v-validate="'required'",
      :error-messages="errors.collect('row')",
      item-text="title",
      item-value="title"
      )
        template(slot="no-data")
          v-list-tile
            v-list-tile-content
              v-list-tile-title(v-html="$t('settings.rowGridSize.noData')")
      div
        v-slider(
        v-for="(slider, key) in sliders",
        :key="`slider-${key}`",
        v-bind="slider.bind",
        v-on="slider.on",
        ticks="always"
        always-dirty,
        thumb-label
        )
</template>

<script>
import isEmpty from 'lodash/isEmpty';

import { WIDGET_MAX_SIZE, WIDGET_MIN_SIZE } from '@/constants';
import { generateRow } from '@/helpers/entities';

export default {
  inject: ['$validator'],
  props: {
    rowId: {
      type: String,
      default: null,
    },
    size: {
      type: Object,
      default: () => ({ sm: 3, md: 3, lg: 3 }),
    },
    availableRows: {
      type: Array,
      default: () => [],
    },
    rowForCreation: {
      type: Object,
      default: null,
    },
  },
  data() {
    return {
      search: null,
    };
  },
  computed: {
    row: {
      get() {
        return this.availableRows.find(row => row._id === this.rowId) || null;
      },
      set(value) {
        if (isEmpty(value)) {
          this.$emit('update:rowId', null);
        } else if (value !== '' && value !== this.row) {
          if (typeof value === 'string') {
            let selectedRow = this.availableRows.find(row => row.title === value);

            if (!selectedRow) {
              selectedRow = generateRow();

              selectedRow.title = value;

              this.$emit('createRow', selectedRow);
            }

            this.$emit('update:rowId', selectedRow._id);
          } else if (typeof value === 'object') {
            this.$emit('update:rowId', value._id);
          }

          this.$emit('update:size', { sm: WIDGET_MIN_SIZE, md: WIDGET_MIN_SIZE, lg: WIDGET_MIN_SIZE });
        }
      },
    },
    sliders() {
      const keys = ['sm', 'md', 'lg'];
      const icons = {
        sm: 'smartphone',
        md: 'tablet',
        lg: 'desktop_windows',
      };

      if (!this.row) {
        return keys.map(key => ({
          bind: {
            label: key,
            value: 0,
            max: WIDGET_MAX_SIZE,
            disabled: true,
          },
        }));
      }

      return keys.map(key => ({
        bind: {
          prependIcon: icons[key],
          value: this.size[key],
          max: this.row.availableSize[key],
          errorMessages: this.errors.first(key),
          'data-vv-name': key,
          'v-validate': 'min_value:3',
        },
        on: {
          input: value => this.updateSizeField(key, value),
        },
      }));
    },
  },
  methods: {
    blurRow() {
      this.search = this.row ? this.row.title : '';
    },

    updateSizeField(key, value) {
      this.$emit('update:size', { ...this.size, [key]: value });
    },
  },
};
</script>
