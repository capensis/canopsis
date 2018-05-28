<template lang="pug">
  ul.timeline
    li(v-for="step in steps")
      .timeline-item(v-if="isNewDate(step.t)")
        .date {{ $d(step.t, 'short') }}
      .timeline-item
          .time {{ $d(step.t, 'time') }}
          div(v-if="step._t !== 'statecounter'")
            state-flag.flag(:val="step.val" :isStatus="isStatus(step._t)")
            .header
              alarm-chips.chips(:val="step.val" :isStatus="isStatus(step._t)")
              p {{ step._t | stepTitle(step.a) }}
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
import capitalize from 'lodash/capitalize';

import StateFlag from '@/components/BasicComponent/alarm-flag.vue';
import AlarmChips from '@/components/BasicComponent/alarm-chips.vue';

import { STATES_CHIPS_AND_FLAGS_STYLE } from '@/config';

const { mapGetters, mapActions } = createNamespacedHelpers('alarm');

export default {
  name: 'time-line',
  components: { AlarmChips, StateFlag },
  filters: {
    stepTitle(stepTitle, stepAuthor) {
      let formattedStepTitle = '';
      formattedStepTitle = stepTitle.replace(/(status)|(state)/g, '$& ');
      formattedStepTitle = formattedStepTitle.replace(/(inc)|(dec)/g, '$&reased ');
      formattedStepTitle += 'by ';
      if (stepAuthor === 'canopsis.engine') {
        formattedStepTitle += 'system';
      } else {
        formattedStepTitle += stepAuthor;
      }
      return capitalize(formattedStepTitle);
    },
  },
  props: {
    alarmProps: {
      type: Object,
      required: true,
      lastDate: null,
    },
  },
  computed: {
    ...mapGetters(['item']),
    steps() {
      const alarm = this.item(this.alarmProps._id);
      if (alarm && alarm.v.steps) {
        console.log(alarm.v.steps);
        return alarm.v.steps;
      }
      return [];
    },
    stateName(state) {
      const stateValue = parseInt(state.replace('state:', ''), 10);
      return STATES_CHIPS_AND_FLAGS_STYLE[stateValue].text;
    },
    stateSteps(steps) {
      return pickBy(steps, (value, key) => key.startsWith('state:'));
    },
  },
  mounted() {
    this.fetchItem({
      id: this.alarmProps.d,
      params: {
        opened: 'true',
        resolved: 'true',
        sort_key: 't',
        sort_dir: 'DESC',
        limit: '1',
        with_steps: 'true',
      },
    });
  },
  updated() {
    // Useful for example the user change the translation
    this.lastDate = null;
  },
  methods: {
    ...mapActions([
      'fetchItem',
    ]),
    isNewDate(timestamp) {
      if (timestamp !== this.lastDate) {
        this.lastDate = timestamp;
        return true;
      }
      return false;
    },
    isStatus(stepTitle) {
      return stepTitle.startsWith('status');
    },
  },
};
</script>

<style lang="scss" scoped>
  $border_line: #DDDDE0;
  ul {
    list-style: none;
    color: #858585;
  }

  li:last-child  {
    .timeline-item {
      border-image: linear-gradient(
          to bottom,
          $border-line 60%,
          white) 1 100%;
    }
  }
  .timeline {
    margin: 0 auto;
    width: 90%
  }

  .timeline-item {
    padding: 3em 2em 0em;
    position: relative;
    border-left: 2px solid $border_line;

    .time{
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
     background: white;
   }

  .flag {
    left: -19px;
  }

  .date {
    left: -11px;
  }

  .content {
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
      height: 25px;
    }

    p{
      font-size: 17px;
    }

  }

  p{
    overflow-wrap: break-word;
    text-overflow: ellipsis;
    width: 90%;
  }
</style>
