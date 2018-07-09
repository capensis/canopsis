<template lang="pug">
  ul.timeline
    li(v-for="step in steps")
      .timeline-item(v-if="isNewDate(step.t)")
        .date {{ getFormattedDate(step.t) }}
      .timeline-item
        .time {{ $d(step.t, 'time') }}
        div(v-if="step._t !== 'statecounter'")
          alarm-flag.flag(:value="step.val", :isStatus="isStatus(step._t)")
          .header
            alarm-chips.chips(:value="step.val", :isStatus="isStatus(step._t)")
            p {{ step._t | stepTitle(step.a) }}
          .content
            p {{ step.m }}
        div(v-else)
          alarm-flag.flag(isCroppedState)
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
import moment from 'moment';

import AlarmFlag from '@/components/other/alarm/timeline/alarm-flag.vue';
import AlarmChips from '@/components/other/alarm/timeline/alarm-chips.vue';
import { numericSortObject } from '@/helpers/sorting';
import { ENTITIES_STATES_STYLES } from '@/constants';

const { mapGetters, mapActions } = createNamespacedHelpers('alarm');

export default {
  components: { AlarmChips, AlarmFlag },
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
      // This data represent the last step timestamp encountered
      // Its used for group step under the same date
      lastDate: null,
    },
  },
  computed: {
    ...mapGetters(['item']),
    steps() {
      const alarm = this.item(this.alarmProps._id);
      if (alarm && alarm.v.steps) {
        const steps = [...alarm.v.steps];
        return numericSortObject(steps, 't', 'desc');
      }
      return [];
    },
    stateName(state) {
      const stateValue = parseInt(state.replace('state:', ''), 10);
      return ENTITIES_STATES_STYLES[stateValue].text;
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
    this.lastDate = null;
  },
  updated() {
    // Useful like for example when the user change the translation
    this.lastDate = null;
  },
  methods: {
    ...mapActions([
      'fetchItem',
    ]),
    isNewDate(timestamp) {
      const date = new Date(timestamp);
      if (!this.lastDate ||
          (date.getDay() !== this.lastDate.getDay()
            && date.getMonth() !== this.lastDate.getMonth()
            && date.getFullYear() !== this.lastDate.getFullYear())) {
        this.lastDate = date;
        return true;
      }
      return false;
    },
    isStatus(stepTitle) {
      return stepTitle.startsWith('status');
    },
    getFormattedDate(timestamp) {
      return moment.unix(timestamp).format('DD/MM/YYYY');
    },
  },
};
</script>

<style lang="scss" scoped>
  $border_line: #DDDDE0;
  $background: #FAFAFA;
  ul {
    list-style: none;
    color: #858585;
  }

  li:last-child {
    .timeline-item {
      border-image: linear-gradient(
          to bottom,
          $border-line 60%,
          $background) 1 100%;
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

    .chips {
      font-size: 15px;
      height: 25px;
    }

    p {
      font-size: 17px;
    }

  }

  p {
    overflow-wrap: break-word;
    text-overflow: ellipsis;
    width: 90%;
  }
</style>
