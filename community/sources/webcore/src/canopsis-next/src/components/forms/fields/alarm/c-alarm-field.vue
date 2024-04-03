<template>
  <c-lazy-search-field
    v-field="value"
    :label="$t('entity.fields.alarmDisplayName')"
    :loading="pending"
    :items="alarms"
    :name="name"
    :has-more="hasMoreAlarms"
    :required="required"
    :item-text="itemText"
    :item-value="itemValue"
    :disabled="disabled"
    :no-data-text="$t('alarm.noAlarmFound')"
    clearable
    autocomplete
    @fetch="fetchAlarms"
    @fetch:more="fetchMoreAlarms"
    @update:search="updateSearch"
  />
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { isArray, keyBy, pick } from 'lodash';

import { formArrayMixin } from '@/mixins/form';

const { mapActions: mapAlarmActions } = createNamespacedHelpers('alarm');

export default {
  inject: ['$validator'],
  mixins: [formArrayMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Array, String, Object],
      default: '',
    },
    name: {
      type: String,
      default: 'alarm',
    },
    itemText: {
      type: String,
      default: 'display_name',
    },
    itemValue: {
      type: String,
      default: '_id',
    },
    limit: {
      type: Number,
      default: 20,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    params: {
      type: Object,
      default: () => {},
    },
  },
  data() {
    return {
      alarmsById: {},
      pending: false,
      pageCount: 1,

      query: {
        page: 1,
        search: null,
      },
    };
  },
  computed: {
    alarms() {
      return Object.values(this.alarmsById);
    },

    hasMoreAlarms() {
      return this.pageCount > this.query.page;
    },
  },
  watch: {
    params() {
      this.query.page = 1;

      this.fetchAlarms();
    },
  },
  methods: {
    ...mapAlarmActions({ fetchAlarmsDisplayNamesWithoutStore: 'fetchDisplayNamesWithoutStore' }),

    getQuery() {
      return {
        limit: this.limit,
        page: this.query.page,
        search: this.query.search,
        ...this.params,
      };
    },

    async fetchAlarms() {
      try {
        this.pending = true;

        const { data, meta } = await this.fetchAlarmsDisplayNamesWithoutStore({
          params: this.getQuery(),
        });

        this.pageCount = meta.page_count;

        this.alarmsById = {
          ...(this.query.page !== 1 ? this.alarmsById : {}),
          ...keyBy(data, '_id'),
          ...pick(this.alarmsById, isArray(this.value) ? this.value : [this.value]),
        };
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },

    fetchMoreAlarms() {
      this.query.page += 1;

      this.fetchAlarms();
    },

    updateSearch(search) {
      this.query.search = search;
      this.query.page = 1;

      this.fetchAlarms();
    },
  },
};
</script>
