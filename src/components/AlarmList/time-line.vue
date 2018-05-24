<template lang="pug">
  ul.timeline(v-if="!pending")
    li.timeline-item(v-for="step in steps")
      time-item.time(:time="step.t")
      div(v-if="step._t !== 'statecounter'")
        state-flag.flag(:val="step.val" :isStatus="isStatus(step._t)")
        .header
          alarm-chips.chips(:val="step.val" :isStatus="isStatus(step._t)")
          p {{ step._t }}
        .content
          p {{ step.m }}
      div(v-else)
        state-flag.flag(isCroppedState)
        .header
          p Cropped State (since last change of status)
        .content
          table
            tr
              td State increased :
              td {{ step.val.stateinc }}
            tr
              td State decreases :
              td {{ step.val.statedec }}
            tr(v-for="(value, state) in stateSteps(step.val)")
              td State {{ stateName(state) }} :
              td {{ value }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import pickBy from 'lodash/pickBy';

import StateFlag from '@/components/BasicComponent/alarm-flag.vue';
import AlarmChips from '@/components/BasicComponent/alarm-chips.vue';
import TimeItem from '@/components/BasicComponent/time-item.vue';

import { STATES_CHIPS_AND_FLAGS_STYLE } from '@/config';

const { mapGetters, mapActions } = createNamespacedHelpers('entities/alarm');

export default {
  name: 'time-line',
  components: { TimeItem, AlarmChips, StateFlag },
  mounted() {
    this.fetchItem({
      params: {
        filter: { d: this.alarmId },
        opened: 'true',
        resolved: 'true',
        sort_key: 't',
        sort_dir: 'DESC',
        limit: '1',
        with_steps: 'true',
      },
    }).then(() => {
      this.alarm = this.item;
      this.pending = false;
    });
  },
  props: {
    alarmId: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      alarm: null,
      pending: true,
    };
  },
  computed: {
    ...mapGetters(['item']),
    steps() {
      return this.alarm.v.steps;
    },
    stateName(state) {
      const stateValue = parseInt(state.replace('state:', ''), 10);
      return STATES_CHIPS_AND_FLAGS_STYLE[stateValue].text;
    },
    stateSteps(steps) {
      return pickBy(steps, (value, key) => key.startsWith('state:'));
    },
  },
  methods: {
    ...mapActions([
      'fetchItem',
    ]),
    isStatus(stepTitle) {
      return stepTitle.startsWith('status');
    },
  },
};
</script>

<style lang="scss" scoped>
  ul {
    list-style: none;
  }
  .timeline {
    margin: 0 auto;
    width: 90%
  }

  $border_line: #DDDDE0;

  .timeline-item {
    padding: 3em 2em 2em;
    position: relative;
    border-left: 2px solid $border_line;

    .time{
      position: absolute;
      color: #686868;
      left: 2em;
      top: 9px;
      display: block;
      font-size: 11px;
    }

    &:last-child  {
     border-image: linear-gradient(
     to bottom,
     $border-line 60%,
     white) 1 100%;
    }

  }

  .flag {
     top: 0;
     position: absolute;
     left: -19px;
     background: white;
   }

  .content {
    color: #858585;
    padding-left: 20px;
    padding-top: 20px;
    overflow-wrap: break-word;
    width: 90%;
  }

  .header {
    color: #686868;
    display: flex;
    align-items: baseline;
    font-weight: bold;

    .chips{
      font-size: 15px;
    }

    p{
      font-size: 20px;
    }

  }

  p{
    overflow-wrap: break-word;
    text-overflow: ellipsis;
    width: 90%;
  }
</style>
