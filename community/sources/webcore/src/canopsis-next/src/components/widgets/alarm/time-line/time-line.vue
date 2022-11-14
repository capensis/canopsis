<template lang="pug">
  div.timeline
    ul(v-for="(steps, day) in groupedSteps", :key="day")
      li(v-for="(step, index) in steps", :key="index")
        div.timeline-date(v-show="index === 0")
          div.date {{ day }}
        div.timeline-card
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
$borderLine: #DDDDE0;

.timeline {
  margin: 0 auto;
  width: 90%;

  ul {
    list-style: none;
    color: #858585;
  }

  &-date {
    padding: 3em 2em 0;
    position: relative;

    .date {
      top: 4px;
      left: -11px;
      position: absolute;
    }

    &:before, &:after {
      content: '';
      position: absolute;
      left: -2px;
      width: 2px;
      background-color: $borderLine;
    }

    &:before {
      top: 0;
      height: 4px;
    }

    &:after {
      top: 24px;
      bottom: 0;
    }
  }

  &-card {
    padding: 3em 2em 0;
    position: relative;

    .time {
      position: absolute;
      left: 2em;
      top: 9px;
      display: block;
      font-size: 11px;
    }

    .flag {
      height: 30px;
      top: 0;
      left: -13px;
      position: absolute;
    }

    &:before, &:after {
      content: '';
      position: absolute;
      left: -2px;
      width: 2px;
      background-color: $borderLine;
    }

    &:after {
      top: 30px;
      bottom: 0;
    }
  }

  ul:last-of-type li:last-child &-card:after {
    background-color: unset;
    background-image: linear-gradient(
        to bottom,
        $borderLine 60%,
        transparent
    );
  }
}
</style>
