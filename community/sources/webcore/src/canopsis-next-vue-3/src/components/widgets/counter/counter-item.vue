<template>
  <card-with-see-alarms-btn
    :class="itemClasses"
    :show-button="!hideActions && hasAlarmsListAccess"
    :style="itemStyle"
    class="counter-item"
    variant="flat"
    tile
    @show:alarms="showAlarmListModal"
  >
    <v-btn
      v-if="!hideActions && hasVariablesHelpAccess"
      class="counter-item__help-btn ma-0"
      icon="help"
      size="small"
      variant="text"
      @click.stop="showVariablesHelpModal"
    />

    <c-layout class="counter-item__content" align="start">
      <v-icon :icon="icon" class="counter-item__icon px-3 py-2" size="2em" />

      <c-compiled-template
        :template="widget.parameters.blockTemplate"
        :context="templateContext"
        class="counter-item__template pt-3"
      />
    </c-layout>
  </card-with-see-alarms-btn>
</template>

<script>
import { invert } from 'lodash';

import {
  ALARM_STATES,
  COUNTER_EXPORT_FILE_NAME_PREFIX,
  COUNTER_STATES_ICONS,
  MODALS,
  USER_PERMISSIONS,
} from '@/constants';

import { convertObjectToTreeview } from '@/helpers/treeview';

import { authMixin } from '@/mixins/auth';

import CardWithSeeAlarmsBtn from '@/components/common/card/card-with-see-alarms-btn.vue';

const INVERTED_ALARM_STATES = invert(ALARM_STATES);

export default {
  components: { CardWithSeeAlarmsBtn },

  mixins: [authMixin],

  props: {
    counter: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    hideActions: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    templateContext() {
      return {
        levels: this.widget.parameters.levels,
        counter: this.counter,
      };
    },

    state() {
      const { counter, values } = this.widget.parameters.levels;

      const count = this.counter[counter];

      return (
        [ALARM_STATES.critical, ALARM_STATES.major, ALARM_STATES.minor].find(
          state => count >= values?.[INVERTED_ALARM_STATES[state]],
        ) || ALARM_STATES.ok
      );
    },

    hasVariablesHelpAccess() {
      return this.checkAccess(USER_PERMISSIONS.business.counter.actions.variablesHelp);
    },

    hasAlarmsListAccess() {
      return this.checkAccess(USER_PERMISSIONS.business.counter.actions.alarmsList);
    },

    color() {
      const { colors } = this.widget.parameters.levels;

      return colors[INVERTED_ALARM_STATES[this.state]];
    },

    icon() {
      return COUNTER_STATES_ICONS[this.state];
    },

    itemClasses() {
      return [
        `mt-${this.widget.parameters.margin.top}`,
        `mr-${this.widget.parameters.margin.right}`,
        `mb-${this.widget.parameters.margin.bottom}`,
        `ml-${this.widget.parameters.margin.left}`,
      ];
    },

    itemHeight() {
      return 4 + this.widget.parameters.heightFactor;
    },

    itemStyle() {
      return {
        height: `${this.itemHeight}em`,
        backgroundColor: this.color,
      };
    },
  },

  methods: {
    showAlarmListModal() {},

    showVariablesHelpModal() {
      const counterFields = convertObjectToTreeview(this.counter, 'counter');
      const levelsFields = convertObjectToTreeview(this.widget.parameters.levels, 'levels');
      const variables = [counterFields, levelsFields];

      this.$modals.show({
        name: MODALS.variablesHelp,
        config: {
          variables,
          exportEntity: this.counter,
          exportEntityName: COUNTER_EXPORT_FILE_NAME_PREFIX,
        },
      });
    },
  },
};
</script>

<style lang="scss">
.counter-item {
  &__icon {
    display: inline-flex;
    align-self: center;
  }

  &__template {
    width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1.2em;
  }

  &__help-btn {
    position: absolute;
    top: 0;
    right: 0.2em;
    z-index: 1;

    &:hover,
    &:focus {
      position: absolute;
    }
  }
}
</style>
