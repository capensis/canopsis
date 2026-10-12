<template>
  <component
    :is="component.is"
    v-bind="component.bind"
    v-on="componentOn"
  />
</template>

<script>
import { get } from 'lodash';

import { ENTITY_FIELDS, ENTITY_INFOS_TYPE } from '@/constants';

import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';
import { hasStateSetting } from '@/helpers/entities/entity/entity';

import { widgetColumnsFiltersMixin } from '@/mixins/widget/columns-filters';

import EntityColumnEventStatistics from './entity-column-event-statistics.vue';
import EntityColumnPbehaviorInfo from './entity-column-pbehavior-info.vue';
import EntityColumnState from './entity-column-state.vue';

export default {
  components: {
    EntityColumnEventStatistics,
    EntityColumnPbehaviorInfo,
    EntityColumnState,
  },

  mixins: [widgetColumnsFiltersMixin],

  emits: ['click:state'],

  props: {
    entity: {
      type: Object,
      required: true,
    },
    column: {
      type: Object,
      required: true,
    },
    showRootCauseByStateClick: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    columnsFiltersMap() {
      return this.columnsFilters.reduce((acc, { column, filter, attributes = [] }) => {
        acc[column.replace(/^entity\./, '')] = this.getFilter(filter, attributes);

        return acc;
      }, {});
    },

    columnFilter() {
      const PROPERTIES_FILTERS_MAP = {
        [ENTITY_FIELDS.lastEventDate]: convertDateToStringWithFormatForToday,
        [ENTITY_FIELDS.lastPbehaviorDate]: convertDateToStringWithFormatForToday,
        [ENTITY_FIELDS.imported]: convertDateToStringWithFormatForToday,

        ...this.columnsFiltersMap,
      };

      return PROPERTIES_FILTERS_MAP[this.column.value];
    },

    value() {
      const value = get(this.entity, this.column.value, '');

      return this.columnFilter ? this.columnFilter(value) : value;
    },

    stateCellProperties() {
      const component = {
        is: EntityColumnState,
        bind: {
          type: ENTITY_INFOS_TYPE.state,
          value: this.value,
        },
      };

      if (this.showRootCauseByStateClick && hasStateSetting(this.entity)) {
        component.bind.class = 'cursor-pointer';
        component.bind.appendIconName = '$vuetify.icons.root_cause';
        component.on = {
          click: () => this.$emit('click:state', this.entity),
        };
      }

      return component;
    },

    component() {
      const PROPERTIES_COMPONENTS_MAP = {
        enabled: {
          is: 'c-enabled',
          bind: {
            value: this.value,
          },
        },
        idle_since: {
          is: 'c-no-events-icon',
          bind: {
            value: Number(this.value),
            top: true,
            color: 'error',
          },
        },
        ko_events: {
          is: EntityColumnEventStatistics,
          bind: {
            entity: this.entity,
          },
        },
        ok_events: {
          is: EntityColumnEventStatistics,
          bind: {
            entity: this.entity,
          },
        },
        pbehavior_info: {
          is: EntityColumnPbehaviorInfo,
          bind: {
            pbehaviorInfo: this.entity.pbehavior_info,
          },
        },
        state: this.stateCellProperties,
        links: {
          is: 'c-links-chips',
          bind: {
            links: this.entity.links,
            onlyIcon: this.column.onlyIcon,
            inlineCount: this.column.inlineLinksCount,
          },
        },
      };

      const cell = PROPERTIES_COMPONENTS_MAP[this.column.value];

      if (cell) {
        return cell;
      }

      if (this.column.value.startsWith('links.')) {
        const category = this.column.value.replace('links.', '');

        return {
          is: 'c-links-chips',
          bind: {
            category,
            links: this.entity.links,
            onlyIcon: this.column.onlyIcon,
            inlineCount: this.column.inlineLinksCount,
          },
        };
      }

      return {
        is: 'c-ellipsis',
        bind: {
          text: String(this.value),
        },
      };
    },

    componentOn() {
      return this.component.on ?? {};
    },
  },
};
</script>
