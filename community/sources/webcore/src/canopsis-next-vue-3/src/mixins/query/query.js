import { isEqual, omit } from 'lodash'

import { PAGINATION_LIMIT } from '@/config'

import { getQueryForList } from '@/helpers/entities/shared/query'

/**
 * Mixin for local list query state.
 *
 * Provides Vuetify table options, query field updates and request query conversion.
 */
export const localQueryMixin = {
  data() {
    return {
      query: {
        page: 1,
        itemsPerPage: PAGINATION_LIMIT,
        search: '',
        sortBy: [],
        sortDesc: [],
      },
    }
  },

  watch: {
    query(query, oldQuery) {
      if (this.customQueryCondition(query, oldQuery)) {
        this.fetchList()
      }
    },
  },

  computed: {
    options: {
      get() {
        const { page = 1, itemsPerPage = PAGINATION_LIMIT, sortBy = [], sortDesc = [], ...restQuery } = this.query

        return { page, itemsPerPage, sortBy, sortDesc, ...restQuery }
      },
      set(value) {
        this.query = {
          search: value.search || '',
          page: value.itemsPerPage <= this.query.itemsPerPage ? value.page : 1,
          itemsPerPage: value.itemsPerPage || PAGINATION_LIMIT,
          sortBy: value.sortBy || [],
          sortDesc: value.sortDesc || [],

          ...omit(value, [
            'search',
            'page',
            'itemsPerPage',
            'totalItems',
            'groupBy',
            'groupDesc',
            'multiSort',
            'mustSort',
          ]),
        }
      },
    },
  },

  methods: {
    /**
     * Update query field by name.
     *
     * @param {string} field
     * @param {*} value
     */
    updateQueryField(field, value) {
      this.query = {
        ...this.query,
        [field]: value,
      }
    },

    /**
     * Check whether query update should trigger list fetching.
     *
     * @param {Object} query
     * @param {Object} oldQuery
     * @returns {boolean}
     */
    customQueryCondition(query, oldQuery) {
      return !isEqual(query, oldQuery)
    },

    /**
     * Convert local query state to request params.
     *
     * @param {Object} [queryParams=this.query]
     * @returns {Object}
     */
    getQuery(queryParams = this.query) {
      return getQueryForList(queryParams)
    },
  },

  fetchList() {},
}
