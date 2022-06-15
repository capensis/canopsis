<template lang="pug">
  div.timeline
    ul(v-for="(steps, day) in groupedSteps", :key="day")
      li(v-for="(step, index) in steps", :key="index")
        div.timeline-item(v-show="index === 0")
          div.date {{ day }}
        div.timeline-item
          div.time {{ step.t | date('time') }}
          time-line-flag.flag(:step="step")
          time-line-card(:step="step", :is-html-enabled="isHtmlEnabled")
    c-pagination(
      :total="meta.total_count",
      :limit="meta.per_page",
      :page="meta.page",
      @input="updatePage"
    )
</template>

<script>
import { groupAlarmSteps } from '@/helpers/entities';

import TimeLineFlag from './time-line-flag.vue';
import TimeLineCard from './time-line-card.vue';

export default {
  components: { TimeLineFlag, TimeLineCard },
  props: {
    steps: {
      type: Object,
      required: true,
    },
    isHtmlEnabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    meta() {
      return this.steps?.meta ?? {};
    },

    groupedSteps() {
      return groupAlarmSteps(this.steps?.data ?? []);
    },
  },
  methods: {
    updatePage(page) {
      this.$emit('update:page', page);
    },
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
