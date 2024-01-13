<template>
  <c-lazy-search-field
    v-field="value"
    :label="label"
    :loading="pending"
    :items="availableParameters"
    :name="name"
    :has-more="hasMoreExternalMetrics"
    :required="required"
    :hide-no-data="addable"
    item-text="text"
    item-value="value"
    @fetch="fetchListExternalMetrics"
    @fetch:more="fetchMoreExternalMetrics"
    @update:search="updateSearch"
  >
    <template
      v-if="!addable"
      #selection="{ item }"
    >
      <v-icon class="mr-2">
        language
      </v-icon>
      <span>{{ item }}</span>
    </template>
    <template
      v-if="!addable"
      #icon=""
    >
      <v-icon class="mr-2">
        language
      </v-icon>
    </template>
  </c-lazy-search-field>
</template>

<script>
import { uniq } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

const { mapActions: mapMetricsActions } = createNamespacedHelpers('metrics');

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [String, Array],
      required: true,
    },
    name: {
      type: String,
      default: 'parameters',
    },
    label: {
      type: String,
      required: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    addable: {
      type: Boolean,
      default: false,
    },
    limit: {
      type: Number,
      default: 20,
    },
  },
  data() {
    return {
      externalMetrics: [],
      pending: false,
      pageCount: 1,

      query: {
        page: 1,
        search: null,
      },
    };
  },
  computed: {
    availableParameters() {
      return this.externalMetrics.map(value => ({
        value,
        text: value,
      }));
    },

    hasMoreExternalMetrics() {
      return this.pageCount > this.query.page;
    },
  },
  methods: {
    ...mapMetricsActions({ fetchExternalMetricsWithoutStore: 'fetchExternalMetricsWithoutStore' }),

    getQuery() {
      return {
        limit: this.limit,
        page: this.query.page,
        search: this.query.search,
      };
    },

    async fetchListExternalMetrics() {
      try {
        this.pending = true;

        const { data, meta } = await this.fetchExternalMetricsWithoutStore({
          params: this.getQuery(),
        });

        this.pageCount = meta?.page_count;

        this.externalMetrics = this.query.page !== 1
          ? uniq([...this.externalMetrics, ...data])
          : data;
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },

    fetchMoreExternalMetrics() {
      this.query.page += 1;

      this.fetchListExternalMetrics();
    },

    updateSearch(search) {
      this.query.search = search;
      this.query.page = 1;

      this.fetchListExternalMetrics();
    },
  },
};
</script>
