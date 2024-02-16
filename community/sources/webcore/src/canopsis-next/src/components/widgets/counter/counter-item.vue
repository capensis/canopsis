<template>
  <card-with-see-alarms-btn
    class="counter-item"
    :class="itemClasses"
    :style="itemStyle"
    :show-button="hasAlarmsListAccess"
    tile
    @show:alarms="showAlarmListModal"
  >
    <v-btn
      class="counter-item__help-btn ma-0"
      v-if="hasVariablesHelpAccess"
      icon
      small
      @click.stop="showVariablesHelpModal"
    >
      <v-icon>help</v-icon>
    </v-btn>
    <div>
      <v-layout justify-start>
        <v-icon
          class="px-3 py-2"
          size="2em"
        >
          {{ icon }}
        </v-icon>
        <c-compiled-template
          class="counter-item__template pt-3"
          :template="widget.parameters.blockTemplate"
          :context="templateContext"
        />
      </v-layout>
    </div>
  </card-with-see-alarms-btn>
</template>

<script>
import { omit } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import {
  MODALS,
  USERS_PERMISSIONS,
  ENTITIES_STATES_KEYS,
  COUNTER_STATES_ICONS,
  COUNTER_EXPORT_FILE_NAME_PREFIX,
} from '@/constants';

import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';
import { convertObjectToTreeview } from '@/helpers/treeview';

import { authMixin } from '@/mixins/auth';

import CardWithSeeAlarmsBtn from '@/components/common/card/card-with-see-alarms-btn.vue';

const { mapActions } = createNamespacedHelpers('alarm');

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
    query: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    templateContext() {
      return {
        levels: this.widget.parameters.levels,
        counter: this.counter,
      };
    },

    stateKey() {
      const {
        counter,
        values,
      } = this.widget.parameters.levels;

      const count = this.counter[counter];

      return [
        ENTITIES_STATES_KEYS.critical,
        ENTITIES_STATES_KEYS.major,
        ENTITIES_STATES_KEYS.minor,
      ].find(state => count >= values[state]) || ENTITIES_STATES_KEYS.ok;
    },

    hasVariablesHelpAccess() {
      return this.checkAccess(USERS_PERMISSIONS.business.counter.actions.variablesHelp);
    },

    hasAlarmsListAccess() {
      return this.checkAccess(USERS_PERMISSIONS.business.counter.actions.alarmsList);
    },

    color() {
      const { colors } = this.widget.parameters.levels;

      return colors[this.stateKey];
    },

    icon() {
      return COUNTER_STATES_ICONS[this.stateKey];
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
    ...mapActions({
      fetchAlarmsListWithoutStore: 'fetchListWithoutStore',
    }),

    showAlarmListModal() {
      const widget = generatePreparedDefaultAlarmListWidget();

      widget.parameters = {
        ...widget.parameters,
        ...this.widget.parameters.alarmsList,

        opened: this.widget.parameters.opened,
      };

      this.$modals.show({
        name: MODALS.alarmsList,
        config: {
          widget,
          title: this.$t('modals.alarmsList.prefixTitle', { prefix: this.counter.filter?.title }),
          fetchList: params => this.fetchAlarmsListWithoutStore({
            params: {
              ...this.query,
              ...omit(params, ['tstart', 'tstop']),

              filters: [this.counter.filter?._id],
            },
          }),
        },
      });
    },

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
  &__template {
    width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1.2em;
  }

  &__help-btn {
    position: absolute;
    right: 0.2em;
    top: 0;
    z-index: 1;

    &:hover, &:focus {
      position: absolute;
    }
  }
}
</style>
