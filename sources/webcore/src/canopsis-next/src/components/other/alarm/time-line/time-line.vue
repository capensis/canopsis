<template lang="pug">
  v-layout.pa-2(v-if="isShowLoader", justify-center)
    v-progress-circular(
      :size="30",
      color="primary",
      indeterminate
    )
  v-expand-transition(v-else)
    div.timeline(:data-test="`alarmTimeLine-${alarm._id}`")
      ul(v-for="(steps, day) in groupedSteps", :key="day")
        li(v-for="(step, index) in steps", :key="`step-${index}`")
          .timeline-item(v-show="index === 0")
            .date {{ day | date('short', true) }}
          .timeline-item
            .time {{ step.t | date('H:mm:SS', true) }}
            time-line-flag.flag(:step="step")
            time-line-card(:step="step", :isHTMLEnabled="isHTMLEnabled")
</template>

<script>
import moment from 'moment';
import { orderBy, groupBy, get, debounce } from 'lodash';

import TimeLineFlag from '@/components/other/alarm/time-line/time-line-flag.vue';
import TimeLineCard from '@/components/other/alarm/time-line/time-line-card.vue';

import entitiesAlarmMixin from '@/mixins/entities/alarm';

const PENDING_SET_DELAY = 500;

export default {
  components: { TimeLineFlag, TimeLineCard },
  mixins: [entitiesAlarmMixin],
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    isHTMLEnabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      pending: true,
    };
  },

  computed: {
    steps() {
      return this.alarm.v.steps || [];
    },

    hasAlarmSteps() {
      return get(this.alarm, 'v.steps');
    },

    isShowLoader() {
      return this.pending || !this.hasAlarmSteps;
    },

    groupedSteps() {
      if (this.hasAlarmSteps) {
        return this.groupSteps(this.alarm.v.steps);
      }

      return {};
    },
  },

  watch: {
    hasAlarmSteps(value) {
      if (!value) {
        this.fetchItem();
      }
    },
  },

  mounted() {
    this.fetchItem();
  },
  methods: {
    async fetchItem() {
      this.pending = true;

      const params = {
        sort_key: 't',
        sort_dir: 'DESC',
        limit: 1,
        with_steps: true,
      };

      if (this.alarm.v.resolved) {
        params.resolved = true;
      }

      await this.fetchAlarmItem({
        id: this.alarm._id,
        params,
      });

      this.setPendingWithDebounce(false);
    },
    groupSteps(steps) {
      const orderedSteps = orderBy(steps, ['t'], 'desc');

      return groupBy(orderedSteps, step => moment.unix(step.t).startOf('day').format());
    },

    setPendingWithDebounce: debounce(function setRequesting(value) {
      this.pending = value;
    }, PENDING_SET_DELAY),
  },
};
</script>

<style lang="scss" scoped>
  $border_line: #DDDDE0;
  $background: white;

  ul {
    list-style: none;
    color: #858585;

    &:last-child {
      li:last-child {
        .timeline-item:last-child {
          border-image: linear-gradient(
              to bottom,
              $border-line 60%,
              $background) 1 100%;
        }
      }
    }
  }

  .timeline {
    margin: 0 auto;
    width: 90%
  }

  .timeline-item {
    padding: 3em 2em 0;
    position: relative;
    border-left: 2px solid $border_line;

    .time {
      position: absolute;
      left: 2em;
      top: 9px;
      display: block;
      font-size: 11px;
    }
  }

  .flag, .date {
    top: 0;
    position: absolute;
    background: $background;
  }

  .flag {
    top: 4px;
    left: -13px;
  }

  .date {
    top: 4px;
    left: -11px;
  }
</style>
