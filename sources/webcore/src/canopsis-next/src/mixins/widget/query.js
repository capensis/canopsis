import { omit } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { DATETIME_FORMATS, REMEDIATION_INSTRUCTION_FILTER_ALL } from '@/constants';

import { dateParse } from '@/helpers/date-intervals';

import queryMixin from '@/mixins/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import vuetifyPaginationMixinCreator from '@/mixins/vuetify/pagination-creator';

/**
 * @mixin Add query logic
 */
export default {
  mixins: [
    queryMixin,
    entitiesUserPreferenceMixin,
    vuetifyPaginationMixinCreator({
      field: 'vDataTablePagination',
      mutating: true,
    }),
  ],
  props: {
    tabId: {
      type: String,
      required: true,
    },
    defaultQueryId: {
      type: [Number, String],
      required: false,
    },
  },
  computed: {
    query: {
      get() {
        return this.getQueryById(this.queryId);
      },
      set(query) {
        return this.updateQuery({ id: this.queryId, query });
      },
    },

    queryId() {
      return this.defaultQueryId || this.widget._id;
    },

    tabQueryNonce() {
      return this.getQueryNonceById(this.tabId);
    },
  },
  methods: {
    getQuery() {
      const query = omit(this.query, [
        'sortKey',
        'sortDir',
        'tstart',
        'tstop',
        'remediationInstructionsFilters',
      ]);

      const {
        tstart,
        tstop,
        remediationInstructionsFilters,
        limit = PAGINATION_LIMIT,
      } = this.query;

      if (tstart) {
        const convertedTstart = dateParse(tstart, 'start', DATETIME_FORMATS.dateTimePicker);

        query.tstart = convertedTstart.unix();
      }

      if (tstop) {
        const convertedTstop = dateParse(tstop, 'stop', DATETIME_FORMATS.dateTimePicker);

        query.tstop = convertedTstop.unix();
      }

      if (remediationInstructionsFilters) {
        const result = remediationInstructionsFilters.reduce((acc, filter) => {
          const key = filter.condition ? 'without' : 'with';

          if (filter.all) {
            acc[key] = [REMEDIATION_INSTRUCTION_FILTER_ALL];
          } else if (!acc[key].includes(REMEDIATION_INSTRUCTION_FILTER_ALL)) {
            acc[key].push(...filter.instructions);
          }

          return acc;
        }, { with: [], without: [] });

        if (result.with.length) {
          query.with_instructions = result.with.join(',');
        }

        if (result.without.length) {
          query.without_instructions = result.without.join(',');
        }
      }

      if (this.query.sortKey) {
        query.sort_key = this.query.sortKey;
        query.sort_dir = this.query.sortDir.toLowerCase();
      }

      query.limit = limit;

      return query;
    },

    updateRecordsPerPage(limit) {
      this.updateLockedQuery({
        id: this.queryId,
        query: { limit },
      });
    },
  },
};
