<template lang="pug">
  v-layout(wrap)
    v-flex(xs3)
      v-select
    v-flex(xs12)
      v-btn.mx-0(:loading="markAsReadPending", color="primary", @click="markNewErrorsAsRead")
        v-icon(left) done_all
        span {{ $t('eventFilter.markAsRead') }}
    v-flex(xs12)
      event-filter-errors-list(
        :errors="eventFilterErrors",
        :pagination.sync="pagination",
        :total-items="eventFilterErrorsMeta.total_count",
        :pending="eventFilterErrorsPending"
      )
</template>

<script>
import { localQueryMixin } from '@/mixins/query-local/query';
import { entitiesEventFilterMixin } from '@/mixins/entities/event-filter';

import EventFilterErrorsList from './event-filter-errors-list.vue';

export default {
  components: { EventFilterErrorsList },
  mixins: [localQueryMixin, entitiesEventFilterMixin],
  props: {
    eventFilter: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      eventFilterErrors: [],
      eventFilterErrorsPending: false,
      markAsReadPending: false,
    };
  },
  computed: {
    eventFilterErrorsMeta() {
      return {
        total_count: this.eventFilterErrors.length,
      };
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    async markNewErrorsAsRead() {
      this.markAsReadPending = true;

      try {
        await this.markNewEventFilterErrorsAsRead({ id: this.eventFilter._id });

        this.fetchList();
      } catch (err) {
        console.error(err);
      } finally {
        this.markAsReadPending = false;
      }
    },

    async fetchList() {
      this.eventFilterErrorsPending = true;

      try {
        const params = this.getQuery();

        const { data } = await this.fetchEventFilterErrorsListWithoutStore({
          id: this.eventFilter._id,
          params,
        });

        this.eventFilterErrors = data;
      } catch (err) {
        /* TODO: Should be removed */
        this.eventFilterErrors = [
          {
            timestamp: 123,
            type: 'type',
            message: 'First message',
            new: true,
          },
          {
            timestamp: 123,
            type: 'type',
            message: 'Second message',
          },
        ];
        console.error(err);
      } finally {
        this.eventFilterErrorsPending = false;
      }
    },
  },
};
</script>
