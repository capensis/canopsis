<template lang="pug">
  v-card.white--text.weather-item(
    :class="itemClasses",
    :style="{ height: itemHeight + 'em', backgroundColor: color }",
    tile
  )
    v-btn.helpBtn.ma-0(
      v-if="hasVariablesHelpAccess",
      icon,
      small,
      @click.stop="showVariablesHelpModal"
    )
      v-icon help
    div
      v-layout(justify-start)
        v-icon.px-3.py-2.white--text(size="2em") {{ icon }}
        v-runtime-template.weather-item__service-name.pt-3(:template="compiledTemplate")
        v-btn.see-alarms-btn(
          v-if="hasAlarmsListAccess",
          flat,
          @click.stop="showAlarmListModal"
        ) {{ $t('serviceWeather.seeAlarms') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import VRuntimeTemplate from 'v-runtime-template';

import {
  MODALS,
  USERS_PERMISSIONS,
  ENTITIES_STATES_KEYS,
  COUNTER_STATES_ICONS,
} from '@/constants';

import { compile } from '@/helpers/handlebars';
import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities';
import { convertObjectToTreeview } from '@/helpers/treeview';

import { authMixin } from '@/mixins/auth';

const { mapActions } = createNamespacedHelpers('alarm');

export default {
  components: {
    VRuntimeTemplate,
  },
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
  asyncComputed: {
    compiledTemplate: {
      async get() {
        const { blockTemplate, levels } = this.widget.parameters;
        const compiledTemplate = await compile(blockTemplate, {
          levels,

          counter: this.counter,
        });

        return `<div>${compiledTemplate}</div>`;
      },
      default: '',
    },
  },
  computed: {
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
        'v-card__with-see-alarms-btn',
        `mt-${this.widget.parameters.margin.top}`,
        `mr-${this.widget.parameters.margin.right}`,
        `mb-${this.widget.parameters.margin.bottom}`,
        `ml-${this.widget.parameters.margin.left}`,
      ];
    },

    itemHeight() {
      return 4 + this.widget.parameters.heightFactor;
    },
  },
  methods: {
    ...mapActions({
      fetchAlarmsListWithoutStore: 'fetchListWithoutStore',
    }),

    async showAlarmListModal() {
      const widget = generatePreparedDefaultAlarmListWidget();

      widget.parameters = {
        ...widget.parameters,
        ...this.widget.parameters.alarmsList,
      };

      this.$modals.show({
        name: MODALS.alarmsList,
        config: {
          widget,
          title: this.$t('modals.alarmsList.prefixTitle', { prefix: this.counter.filter?.title }),
          fetchList: params => this.fetchAlarmsListWithoutStore({
            params: { ...this.query, ...params, filters: [this.counter.filter?._id] },
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
        },
      });
    },
  },
};
</script>
