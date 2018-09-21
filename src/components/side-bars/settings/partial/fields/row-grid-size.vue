<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.rowGridSize.title') }}
    v-container
      v-combobox(
      :value="row",
      @change="changeRow"
      @blur="blurRow"
      :items="items",
      label="Row",
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
        v-for="(slider, key) in sliders"
        :key="`slider-${key}`"
        v-bind="slider.bind",
        v-on="slider.on",
        ticks="always"
        always-dirty,
        thumb-label
        )
</template>

<script>
import omit from 'lodash/omit';

import { WIDGET_MAX_SIZE, WIDGET_MIN_SIZE } from '@/constants';
import { generateRow } from '@/helpers/entities';

import entitiesViewRowMixin from '@/mixins/entities/view/row';

export default {
  inject: ['$validator'],
  mixins: [entitiesViewRowMixin],
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
      items: [...this.availableRows],
    };
  },
  computed: {
    row() {
      return this.items.find(row => row._id === this.rowId) || null;
    },
    sliders() {
      const keys = ['sm', 'md', 'lg'];

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
          label: key,
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
      // this.search = this.row ? this.row.title : '';
    },

    updateSizeField(key, value) {
      this.$emit('update:size', { ...this.size, [key]: value });
    },

    changeRow(value) {
      if (value !== '' && value !== this.row) {
        if (typeof value === 'string') {
          let newRow = this.items.find(v => v.title === value);

          if (!newRow) {
            newRow = generateRow();

            newRow.title = value;
            newRow.availableSize = {
              sm: WIDGET_MAX_SIZE,
              md: WIDGET_MAX_SIZE,
              lg: WIDGET_MAX_SIZE,
            };

            this.items.push(newRow);

            this.$emit('update:rowForCreation', omit(newRow, ['availableSize']));
          }

          this.$emit('update:rowId', newRow._id);
        } else if (typeof value === 'object') {
          this.$emit('update:rowId', value._id);
        } else {
          this.$emit('update:rowId', value);
        }

        this.$emit('update:size', { sm: WIDGET_MIN_SIZE, md: WIDGET_MIN_SIZE, lg: WIDGET_MIN_SIZE });
      }
    },
  },
};
</script>
